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

import (
	"fmt"
	"path"
)

// ModelsPath is the path to which models are stored
const ModelsPath = "/etc/model-checker/models"

// DataPath is the path to which model data is stored
const DataPath = "/etc/model-checker/data"

// CheckerPort is the model checker port
const CheckerPort = 6000

// NewModel gets a new model
func NewModel(name string) *Model {
	specFile := fmt.Sprintf("%s.tla", name)
	specPath := path.Join(ModelsPath, specFile)
	dataFile := fmt.Sprintf("%s.json", name)
	dataPath := path.Join(DataPath, dataFile)
	return &Model{
		Name:     name,
		specFile: specFile,
		specPath: specPath,
		dataFile: dataFile,
		dataPath: dataPath,
	}
}

// Model is a model to check
type Model struct {
	// Name is the name of the model
	Name     string
	specFile string
	specPath string
	dataFile string
	dataPath string
}
