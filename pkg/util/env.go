// Copyright 2019-present Open Networking Foundation.
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

package util

import (
	"fmt"
	"strings"
)

const keyValueSep = "="
const entrySep = ","

// SplitMap splits the given string into key-value pairs
func SplitMap(value string) map[string]string {
	values := strings.Split(value, entrySep)
	pairs := make(map[string]string)
	for _, pair := range values {
		if strings.Contains(pair, keyValueSep) {
			key := pair[:strings.Index(pair, keyValueSep)]
			value := pair[strings.Index(pair, keyValueSep)+1:]
			pairs[key] = value
		}
	}
	return pairs
}

// JoinMap joins the given map of key-value pairs into a single string
func JoinMap(pairs map[string]string) string {
	values := make([]string, 0, len(pairs))
	for key, value := range pairs {
		values = append(values, fmt.Sprintf("%s%s%s", key, keyValueSep, value))
	}
	return strings.Join(values, entrySep)
}
