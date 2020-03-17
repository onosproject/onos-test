// Copyright 2020-present Open Networking Foundation.
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

package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/joncalhoun/pipe"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"
)

// Config is the code generator configuration
type Config struct {
	Path      string     `yaml:"path,omitempty"`
	Package   string     `yaml:"package,omitempty"`
	Resources []Resource `yaml:"resources"`
}

// Resource is a code generator resource
type Resource struct {
	Package      string     `yaml:"package,omitempty"`
	Group        string     `yaml:"group,omitempty"`
	Version      string     `yaml:"version,omitempty"`
	Kind         string     `yaml:"kind,omitempty"`
	ListKind     string     `yaml:"listKind,omitempty"`
	PluralKind   string     `yaml:"pluralKind,omitempty"`
	Scope        string     `yaml:"scope,omitempty"`
	SubResources []Resource `yaml:"subResources"`
}

// Generate generates a Helm API for the given configuration
func Generate(config Config) error {
	options := getOptionsFromConfig(config)
	return generateClient(options)
}

var toCamelCase = func(value string) string {
	return strcase.ToCamel(value)
}

var toLowerCamelCase = func(value string) string {
	return strcase.ToLowerCamel(value)
}

var toLowerCase = func(value string) string {
	return strings.ToLower(value)
}

var toUpperCase = func(value string) string {
	return strings.ToUpper(value)
}

var upperFirst = func(value string) string {
	bytes := []byte(value)
	first := strings.ToUpper(string([]byte{bytes[0]}))
	return string(append([]byte(first), bytes[1:]...))
}

var quote = func(value string) string {
	return "\"" + value + "\""
}

func getTemplate(name string) *template.Template {
	_, filepath, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get template path")
	}
	file := path.Join(path.Dir(filepath), name)
	funcs := template.FuncMap{
		"toCamel":      toCamelCase,
		"toLowerCamel": toLowerCamelCase,
		"lower":        toLowerCase,
		"upper":        toUpperCase,
		"upperFirst":   upperFirst,
		"quote":        quote,
	}
	return template.Must(template.New(name).Funcs(funcs).ParseFiles(file))
}

func generateTemplate(t *template.Template, outputFile string, options interface{}) error {
	fmt.Println(fmt.Sprintf("Generating file %s from template %s", outputFile, t.Name()))
	rc, wc, _ := pipe.Commands(
		exec.Command("gofmt"),
	)
	if err := t.Execute(wc, options); err != nil {
		return err
	}
	wc.Close()
	file, err := openFile(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, rc)
	return err
}

func openFile(filename string) (*os.File, error) {
	if !dirExists(path.Dir(filename)) {
		if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
			return nil, err
		}
	}
	if fileExists(filename) {
		if err := os.Remove(filename); err != nil {
			return nil, err
		}
	}
	return os.Create(filename)
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && !info.IsDir()
}
