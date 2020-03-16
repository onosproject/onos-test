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

// ClientOptions contains options for generating a client
type ClientOptions struct {
	Location Location
	Package  Package
	Types    ClientTypes
	Groups   map[string]*GroupOptions
}

// ClientTypes contains types for generating a client
type ClientTypes struct {
	Interface string
	Struct    string
}

func generateClient(options ClientOptions) error {
	if err := generateTemplate(getTemplate("client.tpl"), path.Join(options.Location.Path, options.Location.File), options); err != nil {
		return err
	}

	for _, group := range options.Groups {
		if err := generateVersionClient(*group); err != nil {
			return err
		}
	}
	return nil
}
