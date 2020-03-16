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

// GroupOptions contains options for generating a group client
type GroupOptions struct {
	Location Location
	Package  Package
	Group    string
	Types    GroupTypes
	Names    GroupNames
	Versions map[string]*VersionOptions
}

// GroupTypes contains types for generating a group client
type GroupTypes struct {
	Interface string
	Struct    string
}

// GroupNames contains names for generating a group client
type GroupNames struct {
	Natural string
	Proper  string
}

func generateGroupClient(options GroupOptions) error {
	if err := generateTemplate(getTemplate("groupclient.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}

	for _, version := range options.Versions {
		if err := generateVersionClient(*version); err != nil {
			return err
		}
	}
	return nil
}
