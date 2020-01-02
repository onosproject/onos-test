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

package cluster

import (
	"os"
	"strconv"
	"strings"
)

var args = make(map[string]string)

const argPrefix = "ONIT_ARG_"

func init() {
	for _, keyval := range os.Environ() {
		key := keyval[:strings.Index(keyval, "=")]
		if strings.HasPrefix(key, argPrefix) {
			value := keyval[strings.Index(keyval, "=")+1:]
			SetArg(strings.ReplaceAll(strings.ToLower(key[len(argPrefix):]), "_", "."), value)
		}
	}
}

// SetArgs sets the cluster arguments
func SetArgs(args map[string]string) {
	for name, value := range args {
		SetArg(name, value)
	}
}

// SetArg sets the value of an argument
func SetArg(name, value string) {
	args[name] = value
}

// GetArgs gets the cluster arguments
func GetArgs() map[string]string {
	return args
}

// GetArg gets the value of an argument
func GetArg(names ...string) *Arg {
	return &Arg{
		value: args[strings.Join(names, ".")],
	}
}

// Arg is a cluster argument
type Arg struct {
	value string
}

// Bool returns the argument as a bool
func (a *Arg) Bool(def bool) bool {
	if a.value == "" {
		return def
	}
	b, err := strconv.ParseBool(a.value)
	if err != nil {
		panic(err)
	}
	return b
}

// Int returns the argument as an int
func (a *Arg) Int(def int) int {
	if a.value == "" {
		return def
	}
	i, err := strconv.Atoi(a.value)
	if err != nil {
		panic(err)
	}
	return i
}

// String returns the argument as a string
func (a *Arg) String(def string) string {
	if a.value == "" {
		return def
	}
	return a.value
}
