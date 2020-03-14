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

package main

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/helm/codegen"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use:  "onit-generate",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			config := codegen.Config{}
			if err := yaml.Unmarshal(bytes, &config); err != nil {
				return err
			}

			if len(args) > 1 {
				config.Path = args[1]
			}
			pkg, _ := cmd.Flags().GetString("package")
			if pkg != "" {
				config.Package = pkg
			}
			return codegen.Generate(config)
		},
	}
	cmd.Flags().StringP("package", "p", "", "the package in which to generate the code")
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
