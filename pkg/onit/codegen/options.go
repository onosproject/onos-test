package codegen

import (
	"fmt"
	"path"
)

type Location struct {
	Path string
	File string
}

type Package struct {
	Name  string
	Path  string
	Alias string
}

func getOptionsFromConfig(config Config) ClientOptions {
	options := ClientOptions{
		Location: Location{
			Path: config.Path,
			File: "client.go",
		},
		Package: Package{
			Name:  path.Base(config.Package),
			Path:  config.Package,
			Alias: path.Base(config.Package),
		},
		Types: ClientTypes{
			Interface: "Client",
			Struct:    "client",
		},
		Groups: make(map[string]GroupOptions),
	}

	for _, resource := range config.Resources {
		groupOpts, ok := options.Groups[resource.Group]
		if !ok {
			groupOpts = GroupOptions{
				Location: Location{
					Path: fmt.Sprintf("%s/%s", config.Path, resource.Group),
					File: "client.go",
				},
				Package: Package{
					Name:  resource.Group,
					Path:  fmt.Sprintf("%s/%s", config.Package, resource.Group),
					Alias: resource.Group,
				},
				Group: resource.Group,
				Types: GroupTypes{
					Interface: "Client",
					Struct:    "client",
				},
				Names: GroupNames{
					Proper: toCamelCase(resource.Group),
				},
				Versions: make(map[string]VersionOptions),
			}
			options.Groups[resource.Group] = groupOpts
		}

		versionOpts, ok := groupOpts.Versions[resource.Version]
		if !ok {
			versionOpts = VersionOptions{
				Location: Location{
					Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
					File: "client.go",
				},
				Package: Package{
					Name:  resource.Version,
					Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Group, resource.Version),
					Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
				},
				Group:   resource.Group,
				Version: resource.Version,
				Types: VersionTypes{
					Interface: "Client",
					Struct:    "client",
				},
				Names: VersionNames{
					Proper: toCamelCase(resource.Version),
				},
				Resources: make(map[string]ResourceClientOptions),
			}
			groupOpts.Versions[resource.Version] = versionOpts
		}

		resourceClientOpts, ok := versionOpts.Resources[resource.Kind]
		if !ok {
			pkg := resource.Package
			if pkg == "" {
				pkg = fmt.Sprintf("k8s.io/api/%s/%s", resource.Group, resource.Version)
			}
			resourceClientOpts = ResourceClientOptions{
				Location: Location{
					Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
					File: fmt.Sprintf("%s.go", toLowerCase(resource.PluralKind)),
				},
				Package: Package{
					Name:  resource.Version,
					Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Group, resource.Version),
					Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
				},
				Types: ResourceClientTypes{
					Interface: fmt.Sprintf("%sClient", resource.PluralKind),
					Struct:    toLowerCamelCase(fmt.Sprintf("%sClient", resource.PluralKind)),
				},
				Names: ResourceClientNames{
					Singular: resource.Kind,
					Plural:   resource.PluralKind,
				},
				Resource: ResourceOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
						File: fmt.Sprintf("%s.go", toLowerCase(resource.Kind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Group, resource.Version),
						Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
					},
					Kind: ResourceKind{
						Package: Package{
							Name:  path.Base(pkg),
							Path:  pkg,
							Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
						},
						Group:    resource.Group,
						Version:  resource.Version,
						Kind:     resource.Kind,
						ListKind: resource.ListKind,
					},
					Types: ResourceTypes{
						Kind:     fmt.Sprintf("%sKind", resource.Kind),
						Resource: fmt.Sprintf("%sResource", resource.Kind),
						Struct:   resource.Kind,
					},
					Names: ResourceNames{
						Singular: resource.Kind,
						Plural:   resource.PluralKind,
					},
				},
			}
			versionOpts.Resources[resource.Kind] = resourceClientOpts
		}
	}
	return options
}
