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

package model

// TemplatesPath is the path containing model templates
const TemplatesPath = "./models"

// StorePath is the path to which models are stored
const StorePath = "/etc/test/models"

// CheckerPort is the model checker port
const CheckerPort = 8000

// Model is a model to check
type Model struct {
	// Name is the name of the model
	Name string
}
