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

func newChart(name string, parent metav1.ObjectsClient) *Chart {
	var chart *Chart
	filter := func(object metav1.Object) (bool, error) {
		resources, err := chart.getResources()
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
	chart = &Chart{
		Client:  newClient(metav1.NewObjectsClient(parent, filter)),
		API:     parent,
		release: name,
		values:  make(map[string]interface{}),
	}

	args := GetArgs(name)
	for key, value := range args {
		chart.SetValue(key, value)
	}
	return chart
}

var _ Client = &Chart{}

// Chart is a Helm chart
type Chart struct {
	Client
	kube.API
	release    string
	chart      string
	repository string
	wait       bool
	values     map[string]interface{}
}

// SetChart sets the chart name
func (c *Chart) SetChart(name string) {
	c.chart = name
}

// Chart returns the chart name
func (c *Chart) Chart() string {
	return c.chart
}

// SetRelease sets the chart's release name
func (c *Chart) SetRelease(name string) {
	c.release = name
}

// Release returns the release name
func (c *Chart) Release() string {
	return c.release
}

// SetRepository sets the chart's repository URL
func (c *Chart) SetRepository(url string) {
	c.repository = url
}

// Repository returns the chart's repository URL
func (c *Chart) Repository() string {
	return c.repository
}

// SetWait sets whether to wait for installation to complete
func (c *Chart) SetWait(wait bool) {
	c.wait = wait
}

// Wait returns whether the chart waits for installation to complete
func (c *Chart) Wait() bool {
	return c.wait
}

// SetValues sets the chart's values
func (c *Chart) SetValues(values map[string]interface{}) {
	c.values = values
}

// SetValue sets a field in the configuration
func (c *Chart) SetValue(path string, value interface{}) {
	setKey(c.values, getPathNames(path), value)
}

// Values returns the chart's values
func (c *Chart) Values() map[string]interface{} {
	return c.values
}

// Value gets a value in the configuration
func (c *Chart) Value(path string) interface{} {
	return getValue(c.values, getPathNames(path))
}

// AddValue adds a value to a list in the configuration
func (c *Chart) AddValue(path string, value interface{}) {
	addValue(c.values, getPathNames(path), value)
}

// getConfig gets the Helm configuration
func (c *Chart) getConfig() (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), c.API.Namespace(), os.Getenv(helmDriverEnv), log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

// getResources returns a list of chart resources
func (c *Chart) getResources() (helm.ResourceList, error) {
	config, err := c.getConfig()
	if err != nil {
		return nil, err
	}
	releases, err := config.Releases.History(c.release)
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

// Setup installs the Helm chart
func (c *Chart) Setup() error {
	config, err := c.getConfig()
	if err != nil {
		return err
	}

	install := action.NewInstall(config)
	install.Namespace = c.API.Namespace()
	install.IncludeCRDs = true
	install.RepoURL = c.repository
	install.ReleaseName = c.release
	install.Wait = c.wait

	// Locate the chart path
	path, err := install.ChartPathOptions.LocateChart(c.chart, settings)
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

	_, err = install.Run(chart, normalize(c.values).(map[string]interface{}))
	if err != nil {
		return err
	}
	return nil
}

// TearDown uninstalls the Helm chart
func (c *Chart) TearDown() error {
	config, err := c.getConfig()
	if err != nil {
		return err
	}
	uninstall := action.NewUninstall(config)
	_, err = uninstall.Run(c.release)
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
