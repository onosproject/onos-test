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
	"github.com/onosproject/onos-test/pkg/helm/api"
	"github.com/onosproject/onos-test/pkg/helm/api/resource"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	helm "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path"
	"reflect"
	"sigs.k8s.io/yaml"
	"strings"
)

const ValuesEnv = "HELM_VALUES"
const ValuesPath = "/etc/onit"
const ValuesFile = "values.yaml"

var settings = cli.New()

func getValues(name string) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	file, err := os.Open(path.Join(ValuesPath, ValuesFile))
	if err == nil {
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		rawValues := make(map[string]interface{})
		if err := yaml.Unmarshal(bytes, &rawValues); err != nil {
			return nil, err
		}
		if releaseValues, ok := rawValues[name]; ok {
			for _, value := range releaseValues.([]interface{}) {
				values = mergeMaps(values, value.(map[string]interface{}))
			}
		}
	}

	valuesEnv := os.Getenv("HELM_VALUES")
	if valuesEnv != "" {
		rawValues := make(map[string]interface{})
		if err := yaml.Unmarshal([]byte(valuesEnv), &rawValues); err != nil {
			return nil, err
		}
		if releaseValues, ok := rawValues[name]; ok {
			for _, value := range releaseValues.([]string) {
				if err := strvals.ParseInto(value, values); err != nil {
					return nil, err
				}
			}
		}
	}
	return values, nil
}

func newRelease(name string, chart *Chart, config *action.Configuration) *Release {
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

	values, err := getValues(name)
	if err != nil {
		panic(err)
	}

	release = &Release{
		Client:    api.NewClient(chart.Client, filter),
		chart:     chart,
		config:    config,
		name:      name,
		values:    make(map[string]interface{}),
		overrides: values,
	}
	return release
}

// Release is a Helm chart release
type Release struct {
	api.Client
	chart     *Chart
	config    *action.Configuration
	name      string
	values    map[string]interface{}
	overrides map[string]interface{}
	skipCRDs  bool
	release   *release.Release
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

// getResources returns a list of chart resources
func (r *Release) getResources() (helm.ResourceList, error) {
	resources, err := r.config.KubeClient.Build(bytes.NewBufferString(r.release.Manifest), true)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// Install installs the Helm chart
func (r *Release) Install(wait bool) error {
	install := action.NewInstall(r.config)
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

	values := mergeMaps(r.overrides, normalize(r.values).(map[string]interface{}))
	release, err := install.Run(chart, values)
	if err != nil {
		return err
	}
	r.release = release
	return nil
}

// Uninstall uninstalls the Helm chart
func (r *Release) Uninstall() error {
	uninstall := action.NewUninstall(r.config)
	_, err := uninstall.Run(r.Name())
	return err
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
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

// setKey sets a key in a map
func setKey(config map[string]interface{}, path []string, value interface{}) {
	names, key := getPathAndKey(path)
	parent := getMapRef(config, names)
	parent[key] = value
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
