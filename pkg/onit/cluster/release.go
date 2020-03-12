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

package cluster

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/onosproject/onos-test/pkg/kube"
	metav1 "github.com/onosproject/onos-test/pkg/onit/cluster/meta/v1"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	helm "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/releaseutil"
	"log"
	"os"
	"reflect"
	"strings"
)

const helmDriverEnv = "HELM_DRIVER"

var settings = cli.New()

func newRelease(name string, chart *Chart) *Release {
	var release *Release
	filter := func(object metav1.Object) (bool, error) {
		resources, err := release.getResources()
		if err != nil {
			return false, err
		}
		for _, resource := range resources {
			kind := resource.Object.GetObjectKind().GroupVersionKind()
			if kind.Group == object.Kind.Group &&
				kind.Version == object.Kind.Version &&
				kind.Kind == object.Kind.Kind &&
				resource.Namespace == object.Namespace &&
				resource.Name == object.Name {
				return true, nil
			}
		}
		return false, nil
	}
	release = &Release{
		Client: newClient(metav1.NewObjectsClient(chart.API, filter)),
		API:    chart.API,
		name:   name,
		values: make(map[string]interface{}),
	}

	args := GetArgs(name)
	for key, value := range args {
		release.SetValue(key, value)
	}
	return release
}

var _ Client = &Release{}

// Release is a Helm chart release
type Release struct {
	Client
	kube.API
	chart  *Chart
	name   string
	wait   bool
	values map[string]interface{}
}

// Name returns the release name
func (r *Release) Name() string {
	return r.name
}

// SetWait sets whether to wait for installation to complete
func (r *Release) SetWait(wait bool) *Release {
	r.wait = wait
	return r
}

// Wait returns whether the chart waits for installation to complete
func (r *Release) Wait() bool {
	return r.wait
}

// SetValues sets the chart's values
func (r *Release) SetValues(values map[string]interface{}) *Release {
	r.values = values
	return r
}

// SetValue sets a field in the configuration
func (r *Release) SetValue(path string, value interface{}) *Release {
	setKey(r.values, getPathNames(path), value)
}

// Values returns the chart's values
func (r *Release) Values() map[string]interface{} {
	return r.values
	return r
}

// Value gets a value in the configuration
func (r *Release) Value(path string) interface{} {
	return getValue(r.values, getPathNames(path))
}

// AddValue adds a value to a list in the configuration
func (r *Release) AddValue(path string, value interface{}) *Release {
	addValue(r.values, getPathNames(path), value)
	return r
}

// getConfig gets the Helm configuration
func (r *Release) getConfig() (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), r.API.Namespace(), os.Getenv(helmDriverEnv), log.Printf); err != nil {
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
	releases, err := config.Releases.History(r.Name())
	if err != nil {
		return nil, err
	}
	if len(releases) < 1 {
		return nil, nil
	}

	releaseutil.SortByRevision(releases)
	release := releases[0]

	resources, err := config.KubeClient.Build(bytes.NewBufferString(release.Manifest), true)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// Install installs the Helm chart
func (r *Release) Install() error {
	config, err := r.getConfig()
	if err != nil {
		return err
	}

	install := action.NewInstall(config)
	install.Namespace = r.API.Namespace()
	install.IncludeCRDs = true
	install.RepoURL = r.chart.Repository()
	install.ReleaseName = r.Name()
	install.Wait = r.wait

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

	_, err = install.Run(chart, normalize(r.values).(map[string]interface{}))
	if err != nil {
		return err
	}
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
