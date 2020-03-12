package codegen

import (
	"fmt"
	"os"
	"path"
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
	Group      string `yaml:"group,omitempty"`
	Version    string `yaml:"version,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
	PluralKind string `yaml:"pluralKind,omitempty"`
	ListKind   string `yaml:"listKind,omitempty"`
	Singular   string `yaml:"singular,omitempty"`
	Plural     string `yaml:"plural,omitempty"`
}

func Generate(config Config) error {
	options := getOptionsFromConfig(config)
	return generateClient(options)
}

func getTemplate(name, file string) *template.Template {
	return template.Must(template.New(name).ParseFiles(file))
}

func getOptionsFromConfig(config Config) ClientOptions {
	options := ClientOptions{
		Package:    path.Base(config.Package),
		ImportPath: config.Package,
		FilePath:   fmt.Sprintf("%s/client.go", config.Path),
		Groups:     make(map[string]GroupOptions),
	}
	for _, resource := range config.Resources {
		groupOpts, ok := options.Groups[resource.Group]
		if !ok {
			groupOpts = GroupOptions{
				Package:    resource.Group,
				ImportPath: fmt.Sprintf("%s/%s", options.ImportPath, resource.Group),
				FilePath:   fmt.Sprintf("%s/%s/client.go", config.Path, resource.Group),
				Group:      resource.Group,
				Versions:   make(map[string]VersionOptions),
			}
			options.Groups[resource.Group] = groupOpts
		}

		versionOpts, ok := groupOpts.Versions[resource.Version]
		if !ok {
			versionOpts = VersionOptions{
				Package:    resource.Version,
				ImportPath: fmt.Sprintf("%s/%s/%s", options.ImportPath, resource.Group, resource.Version),
				FilePath:   fmt.Sprintf("%s/%s/%s/client.go", config.Path, resource.Group, resource.Version),
				Group:      resource.Group,
				Version:    resource.Version,
				Resources:  make(map[string]ResourceOptions),
			}
			groupOpts.Versions[resource.Version] = versionOpts
		}

		resourceOpts, ok := versionOpts.Resources[resource.Singular]
		if !ok {
			resourceOpts = ResourceOptions{
				Package:    resource.Version,
				ImportPath: fmt.Sprintf("%s/%s/%s", options.ImportPath, resource.Group, resource.Version),
				FilePath:   fmt.Sprintf("%s/%s/%s/%s.go", config.Path, resource.Group, resource.Version, resource.Singular),
				Group:      resource.Group,
				Version:    resource.Version,
				Kind:       resource.Kind,
				PluralKind: resource.PluralKind,
				ListKind:   resource.ListKind,
				Singular:   resource.Singular,
				Plural:     resource.Plural,
			}
			versionOpts.Resources[resource.Singular] = resourceOpts
		}
	}
	return options
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
