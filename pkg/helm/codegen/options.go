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
	"strings"
)

// Location is the location of a code file
type Location struct {
	Path string
	File string
}

// Package is the package for a code file
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
		group := resource.Group
		index := strings.Index(group, ".")
		if index != -1 {
			group = group[:index]
		}

		versionOpts, ok := options.Groups[fmt.Sprintf("%s%s", resource.Group, resource.Version)]
		if !ok {
			versionOpts = &GroupOptions{
				Location: Location{
					Path: fmt.Sprintf("%s/%s/%s", config.Path, group, resource.Version),
					File: "client.go",
				},
				Package: Package{
					Name:  resource.Version,
					Path:  fmt.Sprintf("%s/%s/%s", config.Package, group, resource.Version),
					Alias: fmt.Sprintf("%s%s", group, resource.Version),
				},
				Group:   resource.Group,
				Version: resource.Version,
				Types: GroupTypes{
					Interface: "Client",
					Struct:    "client",
				},
				Names: GroupNames{
					Proper: fmt.Sprintf("%s%s", upperFirst(group), upperFirst(resource.Version)),
				},
				Resources: make(map[string]*ResourceOptions),
			}
			options.Groups[fmt.Sprintf("%s%s", resource.Group, resource.Version)] = versionOpts
		}

		_, ok = versionOpts.Resources[resource.Kind]
		if !ok {
			pkg := resource.Package
			if pkg == "" {
				pkg = fmt.Sprintf("k8s.io/api/%s/%s", group, resource.Version)
			}
			resourceOpts := &ResourceOptions{
				Client: &ResourceClientOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, group, resource.Version),
						File: fmt.Sprintf("%sclient.go", toLowerCase(resource.PluralKind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, group, resource.Version),
						Alias: fmt.Sprintf("%s%s", group, resource.Version),
					},
					Types: ResourceClientTypes{
						Interface: fmt.Sprintf("%sClient", resource.PluralKind),
						Struct:    toLowerCamelCase(fmt.Sprintf("%sClient", resource.PluralKind)),
					},
				},
				Reader: &ResourceReaderOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, group, resource.Version),
						File: fmt.Sprintf("%sreader.go", toLowerCase(resource.PluralKind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, group, resource.Version),
						Alias: fmt.Sprintf("%s%s", group, resource.Version),
					},
					Types: ResourceReaderTypes{
						Interface: fmt.Sprintf("%sReader", resource.PluralKind),
						Struct:    toLowerCamelCase(fmt.Sprintf("%sReader", resource.PluralKind)),
					},
				},
				Resource: &ResourceObjectOptions{
					Location: Location{
						Path: fmt.Sprintf("%s/%s/%s", config.Path, group, resource.Version),
						File: fmt.Sprintf("%s.go", toLowerCase(resource.Kind)),
					},
					Package: Package{
						Name:  resource.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, group, resource.Version),
						Alias: fmt.Sprintf("%s%s", group, resource.Version),
					},
					Kind: ResourceObjectKind{
						Package: Package{
							Name:  path.Base(pkg),
							Path:  pkg,
							Alias: fmt.Sprintf("%s%s", group, resource.Version),
						},
						Group:    resource.Group,
						Version:  resource.Version,
						Kind:     resource.Kind,
						ListKind: resource.ListKind,
						Scoped:   resource.Scope != "Cluster",
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
				Group: versionOpts,
			}
			versionOpts.Resources[resource.Kind] = resourceOpts
		}
	}

	for _, resource := range config.Resources {
		if resource.SubResources == nil {
			continue
		}
		references := make([]*ResourceOptions, 0, len(resource.SubResources))
		for _, ref := range resource.SubResources {
			group, ok := options.Groups[fmt.Sprintf("%s%s", ref.Group, ref.Version)]
			if !ok {
				continue
			}
			resource, ok := group.Resources[ref.Kind]
			if !ok {
				continue
			}
			if resource.Reference == nil {
				resource.Reference = &ResourceReferenceOptions{
					Location: Location{
						Path: resource.Resource.Location.Path,
						File: fmt.Sprintf("%sreference.go", toLowerCase(resource.Resource.Names.Plural)),
					},
					Package: Package{
						Name:  resource.Resource.Kind.Version,
						Path:  fmt.Sprintf("%s/%s/%s", config.Package, resource.Resource.Kind.Group, resource.Resource.Kind.Version),
						Alias: fmt.Sprintf("%s%s", resource.Resource.Kind.Group, resource.Resource.Kind.Version),
					},
					Types: ResourceReaderTypes{
						Interface: fmt.Sprintf("%sReference", resource.Resource.Names.Plural),
						Struct:    toLowerCamelCase(fmt.Sprintf("%sReference", resource.Resource.Names.Plural)),
					},
				}
			}
			references = append(references, resource)
		}
		resourceOpts := options.Groups[fmt.Sprintf("%s%s", resource.Group, resource.Version)].Resources[resource.Kind]
		resourceOpts.Resource.References = references
	}

	println(fmt.Sprintf("%v", options))
	return options
}
