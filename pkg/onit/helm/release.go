// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helm

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/onosproject/onos-test/pkg/cluster"
	"github.com/onosproject/onos-test/pkg/onit/helm/api"
	"github.com/onosproject/onos-test/pkg/onit/helm/api/resource"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	helm "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"reflect"
	"strings"
)

var settings = cli.New()

func newRelease(name string, chart *Chart) *Release {
	var release *Release
	var filter resource.Filter = func(kind metav1.GroupVersionKind, meta metav1.ObjectMeta) (bool, error) {
		resources, err := release.getResources()
		if err != nil {
			return false, err
		}
		for _, resource := range resources {
			resourceKind := resource.Object.GetObjectKind().GroupVersionKind()
			if resourceKind.Group == kind.Group &&
				resourceKind.Version == kind.Version &&
				resourceKind.Kind == kind.Kind &&
				resource.Namespace == meta.Namespace &&
				resource.Name == meta.Name {
				return true, nil
			}
		}
		return false, nil
	}
	release = &Release{
		Client: api.NewClient(chart.Client, filter),
		chart:  chart,
		name:   name,
		values: make(map[string]interface{}),
	}

	args := cluster.GetArgs(name)
	for key, value := range args {
		release.Set(key, value)
	}
	return release
}

// Release is a Helm chart release
type Release struct {
	api.Client
	chart    *Chart
	name     string
	values   map[string]interface{}
	skipCRDs bool
	release  *release.Release
}

// Name returns the release name
func (r *Release) Name() string {
	return r.name
}

// Set sets a value
func (r *Release) Set(path string, value interface{}) *Release {
	setKey(r.values, getPathNames(path), value)
	return r
}

// Get gets a value
func (r *Release) Get(path string) interface{} {
	return getValue(r.values, getPathNames(path))
}

// Values is the release's values
func (r *Release) Values() map[string]interface{} {
	return r.values
}

// SetSkipCRDs sets whether to skip CRDs
func (r *Release) SetSkipCRDs(skipCRDs bool) *Release {
	r.skipCRDs = skipCRDs
	return r
}

// SkipCRDs returns whether CRDs are skipped in the release
func (r *Release) SkipCRDs() bool {
	return r.skipCRDs
}

// getConfig gets the Helm configuration
func (r *Release) getConfig() (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), r.Namespace(), "memory", log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

// getResources returns a list of chart resources
func (r *Release) getResources() (helm.ResourceList, error) {
	config, err := r.getConfig()
	if err != nil {
		return nil, err
	}
	resources, err := config.KubeClient.Build(bytes.NewBufferString(r.release.Manifest), true)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// Install installs the Helm chart
func (r *Release) Install(wait bool) error {
	config, err := r.getConfig()
	if err != nil {
		return err
	}

	install := action.NewInstall(config)
	install.Namespace = r.Namespace()
	install.SkipCRDs = r.SkipCRDs()
	install.RepoURL = r.chart.Repository()
	install.ReleaseName = r.Name()
	install.Wait = wait

	// Locate the chart path
	path, err := install.ChartPathOptions.LocateChart(r.chart.Name(), settings)
	if err != nil {
		return err
	}

	// Check chart dependencies to make sure all are present in /charts
	chart, err := loader.Load(path)
	if err != nil {
		return err
	}

	valid, err := isChartInstallable(chart)
	if !valid {
		return err
	}

	if req := chart.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chart, req); err != nil {
			if install.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        path,
					Keyring:          install.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(cli.New()),
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	release, err := install.Run(chart, normalize(r.values).(map[string]interface{}))
	if err != nil {
		return err
	}
	r.release = release
	return nil
}

// Uninstall uninstalls the Helm chart
func (r *Release) Uninstall() error {
	config, err := r.getConfig()
	if err != nil {
		return err
	}
	uninstall := action.NewUninstall(config)
	_, err = uninstall.Run(r.Name())
	return err
}

// getValue gets the value for the given path
func getValue(config map[string]interface{}, path []string) interface{} {
	names, key := getPathAndKey(path)
	parent := getMap(config, names)
	return parent[key]
}

// getMap gets the map at the given path
func getMap(parent map[string]interface{}, path []string) map[string]interface{} {
	if len(path) == 0 {
		return parent
	}
	child, ok := parent[path[0]]
	if !ok {
		return make(map[string]interface{})
	}
	return getMap(child.(map[string]interface{}), path[1:])
}

// getSlice gets the slice at the given path
func getSlice(config map[string]interface{}, path []string) []interface{} {
	names, key := getPathAndKey(path)
	parent := getMap(config, names)
	child, ok := parent[key]
	if !ok {
		return make([]interface{}, 0)
	}
	return child.([]interface{})
}

// setKey sets a key in a map
func setKey(config map[string]interface{}, path []string, value interface{}) {
	names, key := getPathAndKey(path)
	parent := getMapRef(config, names)
	parent[key] = value
}

// addValue adds a value to a slice
func addValue(config map[string]interface{}, path []string, value interface{}) {
	names, key := getPathAndKey(path)
	parent := getMapRef(config, names)
	values, ok := parent[key]
	if !ok {
		values = make([]interface{}, 0)
		parent[key] = value
	}
	parent[key] = append(values.([]interface{}), value)
}

// getMapRef gets the given map reference
func getMapRef(parent map[string]interface{}, path []string) map[string]interface{} {
	if len(path) == 0 {
		return parent
	}
	child, ok := parent[path[0]]
	if !ok {
		child = make(map[string]interface{})
		parent[path[0]] = child
	}
	return getMapRef(child.(map[string]interface{}), path[1:])
}

func getPathNames(path string) []string {
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	names, err := r.Read()
	if err != nil {
		panic(err)
	}
	return names
}

func getPathAndKey(path []string) ([]string, string) {
	return path[:len(path)-1], path[len(path)-1]
}

func normalize(value interface{}) interface{} {
	kind := reflect.ValueOf(value).Kind()
	if kind == reflect.Struct {
		return normalizeStruct(value.(struct{}))
	} else if kind == reflect.Map {
		return normalizeMap(value.(map[string]interface{}))
	} else if kind == reflect.Slice {
		return normalizeSlice(value.([]interface{}))
	}
	return value
}

func normalizeStruct(value struct{}) interface{} {
	elem := reflect.ValueOf(value).Elem()
	elemType := elem.Type()
	normalized := make(map[string]interface{})
	for i := 0; i < elem.NumField(); i++ {
		key := normalizeField(elemType.Field(i))
		value := normalize(elem.Field(i).Interface())
		normalized[key] = value
	}
	return normalized
}

func normalizeMap(values map[string]interface{}) interface{} {
	normalized := make(map[string]interface{})
	for key, value := range values {
		normalized[key] = normalize(value)
	}
	return normalized
}

func normalizeSlice(values []interface{}) interface{} {
	normalized := make([]interface{}, len(values))
	for i, value := range values {
		normalized[i] = normalize(value)
	}
	return normalized
}

func normalizeField(field reflect.StructField) string {
	tag := field.Tag.Get("yaml")
	if tag != "" {
		return strings.Split(tag, ",")[0]
	}
	return strcase.ToLowerCamel(field.Name)
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}
