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
	"encoding/json"
	"path"
)

var traceBinaryFile = path.Join(DataPath, "trace.bin")

// NewTrace returns a new model Trace from the given values
func NewTrace(values ...interface{}) (*Trace, error) {
	bytes, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	return &Trace{
		Bytes: bytes,
	}, nil
}
