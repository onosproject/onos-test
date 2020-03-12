package main

import (
	"fmt"
	"github.com/onosproject/onos-test/pkg/onit/codegen"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use:  "onit-generate",
		Args: cobra.ExactArgs(1),
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
			return codegen.Generate(config)
		},
	}
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
