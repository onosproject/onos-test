package main

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/helm/codegen"
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
