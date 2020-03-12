package codegen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"os"
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
	Package    string  `yaml:"package,omitempty"`
	Group      string  `yaml:"group,omitempty"`
	Version    string  `yaml:"version,omitempty"`
	Kind       string  `yaml:"kind,omitempty"`
	ListKind   string  `yaml:"listKind,omitempty"`
	PluralKind string  `yaml:"pluralKind,omitempty"`
}

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

var quote = func(value string) string {
	return "\"" + value + "\""
}

func getTemplate(name string) *template.Template {
	_, filepath, _, _ := runtime.Caller(0)
	file := path.Join(path.Dir(filepath), name)
	funcs := template.FuncMap{
		"toCamel":      toCamelCase,
		"toLowerCamel": toLowerCamelCase,
		"lower":        toLowerCase,
		"upper":        toUpperCase,
		"quote":        quote,
	}
	return template.Must(template.New(name).Funcs(funcs).ParseFiles(file))
}

func generateTemplate(t *template.Template, outputFile string, options interface{}) error {
	fmt.Println(fmt.Sprintf("Generating file %s from template %s", outputFile, t.Name()))
	file, err := openFile(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.Execute(file, options)
}

func openFile(filename string) (*os.File, error) {
	if !dirExists(path.Dir(filename)) {
		if err := os.MkdirAll(path.Dir(filename), os.ModeDir); err != nil {
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
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
