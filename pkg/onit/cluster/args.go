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
	"github.com/onosproject/onos-test/pkg/util"
	"os"
	"strconv"
	"strings"
)

var args = make(map[string]string)

const argsEnv = "ONIT_ARGS"

func init() {
	for key, value := range util.SplitMap(os.Getenv(argsEnv)) {
		SetArg(key, value)
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

// GetArgs gets the cluster arguments with the given key
func GetArgs(names ...string) map[string]string {
	if len(names) == 0 {
		return args
	}

	prefix := strings.Join(names, ".") + "."
	results := make(map[string]string)
	for key, value := range args {
		if strings.HasPrefix(key, prefix) {
			results[key[len(prefix):]] = value
		}
	}
	return results
}

// GetArg gets the value of an argument
func GetArg(names ...string) *Arg {
	return &Arg{
		value: args[strings.Join(names, ".")],
	}
}

// GetArgsAsEnv returns the given arguments as an environment variable map
func GetArgsAsEnv(args map[string]string) map[string]string {
	env := make(map[string]string)
	env[argsEnv] = util.JoinMap(args)
	return env
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
