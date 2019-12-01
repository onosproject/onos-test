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

package kubetest

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
)

// ScriptingSuite is a suite of scripts
type ScriptingSuite interface{}

// ScriptSuite is an identifier interface for script suites
type ScriptSuite struct{}

// SetupScriptSuite is an interface for setting up a suite of scripts
type SetupScriptSuite interface {
	SetupScriptSuite()
}

// SetupScript is an interface for setting up individual scripts
type SetupScript interface {
	SetupScript()
}

// TearDownScriptSuite is an interface for tearing down a suite of scripts
type TearDownScriptSuite interface {
	TearDownScriptSuite()
}

// TearDownScript is an interface for tearing down individual scripts
type TearDownScript interface {
	TearDownScript()
}

// BeforeScript is an interface for executing code before every script
type BeforeScript interface {
	BeforeScript(testName string)
}

// AfterScript is an interface for executing code after every script
type AfterScript interface {
	AfterScript(testName string)
}

type InternalScript struct {
	Name string
	F    func()
}

// RunScripts runs a script suite
func RunScripts(suite ScriptingSuite, config *TestConfig) {
	suiteSetupDone := false

	methodFinder := reflect.TypeOf(suite)
	scripts := []InternalScript{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := scriptFilter(method.Name, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if !ok {
			continue
		}
		if !suiteSetupDone {
			if setupScriptSuite, ok := suite.(SetupScriptSuite); ok {
				setupScriptSuite.SetupScriptSuite()
			}
			defer func() {
				if tearDownScriptSuite, ok := suite.(TearDownScriptSuite); ok {
					tearDownScriptSuite.TearDownScriptSuite()
				}
			}()
			suiteSetupDone = true
		}
		script := InternalScript{
			Name: method.Name,
			F: func() {
				if setupScriptSuite, ok := suite.(SetupScript); ok {
					setupScriptSuite.SetupScript()
				}
				if beforeScriptSuite, ok := suite.(BeforeScript); ok {
					beforeScriptSuite.BeforeScript(method.Name)
				}
				defer func() {
					if afterScriptSuite, ok := suite.(AfterScript); ok {
						afterScriptSuite.AfterScript(method.Name)
					}
					if tearDownScriptSuite, ok := suite.(TearDownScript); ok {
						tearDownScriptSuite.TearDownScript()
					}
				}()
				method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			},
		}
		scripts = append(scripts, script)
	}
	runScripts(scripts)
}

// runScript runs a script
func runScripts(scripts []InternalScript) {
	for _, script := range scripts {
		println(script.Name)
		script.F()
	}
}

// scriptFilter filters script method names
func scriptFilter(name string, config *TestConfig) (bool, error) {
	if ok, _ := regexp.MatchString("^Run", name); !ok {
		return false, nil
	}
	if config.Test != "" {
		return config.Test == name, nil
	}
	return true, nil
}
