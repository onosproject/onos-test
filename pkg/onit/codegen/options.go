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
		Groups: make(map[string]*GroupOptions),
	}

	for _, resource := range config.Resources {
		groupOpts, ok := options.Groups[resource.Group]
		if !ok {
			groupOpts = &GroupOptions{
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
				Versions: make(map[string]*VersionOptions),
			}
			options.Groups[resource.Group] = groupOpts
		}

		versionOpts, ok := groupOpts.Versions[resource.Version]
		if !ok {
			versionOpts = &VersionOptions{
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
				Resources: make(map[string]*ResourceOptions),
			}
			groupOpts.Versions[resource.Version] = versionOpts
		}

		resourceOptions, ok := versionOpts.Resources[resource.Kind]
		if !ok {
			pkg := resource.Package
			if pkg == "" {
				pkg = fmt.Sprintf("k8s.io/api/%s/%s", resource.Group, resource.Version)
			}
			resourceOptions = &ResourceOptions{
				Client: &ResourceClientOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
						File: fmt.Sprintf("%sclient.go", toLowerCase(resource.PluralKind)),
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
				},
				Reader: &ResourceReaderOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
						File: fmt.Sprintf("%s.go", toLowerCase(resource.PluralKind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Group, resource.Version),
						Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
					},
					Types: ResourceReaderTypes{
						Interface: resource.PluralKind,
						Struct:    toLowerCamelCase(resource.PluralKind),
					},
				},
				Resource: &ResourceObjectOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, resource.Group, resource.Version),
						File: fmt.Sprintf("%s.go", toLowerCase(resource.Kind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Group, resource.Version),
						Alias: fmt.Sprintf("%s%s", resource.Group, resource.Version),
					},
					Kind: ResourceObjectKind{
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
					Types: ResourceObjectTypes{
						Kind:     fmt.Sprintf("%sKind", resource.Kind),
						Resource: fmt.Sprintf("%sResource", resource.Kind),
						Struct:   resource.Kind,
					},
					Names: ResourceObjectNames{
						Singular: resource.Kind,
						Plural:   resource.PluralKind,
					},
				},
			}
			versionOpts.Resources[resource.Kind] = resourceOptions
		}
	}

	for _, resource := range config.Resources {
		if resource.SubResources == nil {
			continue
		}
		subResources := make([]*ResourceOptions, 0, len(resource.SubResources))
		for _, ref := range resource.SubResources {
			group, ok := options.Groups[ref.Group]
			if !ok {
				continue
			}
			version, ok := group.Versions[ref.Version]
			if !ok {
				continue
			}
			subResource, ok := version.Resources[ref.Kind]
			if !ok {
				continue
			}
			subResources = append(subResources, subResource)
		}
		resourceOpts := options.Groups[resource.Group].
			Versions[resource.Version].
			Resources[resource.Kind]
		resourceOpts.Resource.SubResources = subResources
	}

	println(fmt.Sprintf("%v", options))
	return options
}
