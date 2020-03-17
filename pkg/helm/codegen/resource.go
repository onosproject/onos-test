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

import "path"

// ResourceOptions contains options for generating a resource
type ResourceOptions struct {
	Client    *ResourceClientOptions
	Reader    *ResourceReaderOptions
	Reference *ResourceReferenceOptions
	Resource  *ResourceObjectOptions
	Group     *GroupOptions
}

// ResourceObjectOptions contains options for generating a resource object
type ResourceObjectOptions struct {
	Location   Location
	Package    Package
	Kind       ResourceObjectKind
	Types      ResourceObjectTypes
	Names      ResourceObjectNames
	References []*ResourceOptions
}

// ResourceObjectKind contains kinds for generating a resource kind
type ResourceObjectKind struct {
	Package  Package
	Group    string
	Version  string
	Kind     string
	ListKind string
	Scoped   bool
}

// ResourceObjectTypes contains types for generating a resource object
type ResourceObjectTypes struct {
	Kind     string
	Resource string
	Struct   string
}

// ResourceObjectNames contains names for generating a resource object
type ResourceObjectNames struct {
	Singular string
	Plural   string
}

func generateResource(options ResourceOptions) error {
	return generateTemplate(getTemplate("resource.tpl"), path.Join(options.Resource.Location.Path, options.Resource.Location.File), options)
}
