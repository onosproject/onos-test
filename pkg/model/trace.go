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
)

// NewTrace returns a new model trace for the given object
func NewTrace(object interface{}) (*Trace, error) {
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return &Trace{
		Bytes: bytes,
	}, nil
}

// NewTraceFields returns a new model Trace for the given fields and values
func NewTraceFields(fieldsAndValues ...interface{}) (*Trace, error) {
	values := make(map[string]interface{})
	for i := 0; i < len(fieldsAndValues); i += 2 {
		values[fieldsAndValues[i].(string)] = fieldsAndValues[i+1]
	}
	return NewTrace(values)
}

// NewTraceValues returns a new model Trace from the given values
func NewTraceValues(values ...interface{}) (*Trace, error) {
	return NewTrace(values)
}
