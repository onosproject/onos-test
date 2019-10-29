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

package setup

import (
	"fmt"
	"os"

	"github.com/onosproject/onos-test/pkg/onit/k8s"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/onosproject/onos-test/pkg/onit/console"
)

// exitStatus prints the errors from the given status and exits
func exitStatus(status console.ErrorStatus) {
	for _, err := range status.Errors() {
		fmt.Println(err)
	}
	os.Exit(1)
}

// exitError prints the given errors to stdout and exits with exit code 1
func exitError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

// SetDefaultCluster sets the default cluster
func SetDefaultCluster(clusterID string) error {
	if err := initConfig(); err != nil {
		return err
	}
	viper.Set("cluster", clusterID)
	return viper.WriteConfig()
}

// GetDefaultCluster returns the default cluster
func GetDefaultCluster() string {
	return viper.GetString("cluster")
}

func initConfig() error {
	// If the configuration file is not found, initialize a configuration in the home dir.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*viper.ConfigFileNotFoundError); !ok {
			home, err := homedir.Dir()
			if err != nil {
				return err
			}

			err = os.MkdirAll(home+"/.onos", 0777)
			if err != nil {
				return err
			}

			f, err := os.Create(home + "/.onos/onit.yaml")
			if err != nil {
				return err
			}
			f.Close()
		} else {
			return err
		}
	}
	return nil
}

// NewUUIDString returns a new string UUID
func NewUUIDString() string {
	id, err := uuid.NewUUID()
	if err != nil {
		exitError(err)
	}
	return id.String()
}

// InitImageTags initialize the default values of image tags
func InitImageTags(imageTags map[string]string) {
	if imageTags["config"] == "" {
		imageTags["config"] = string(k8s.Debug)
	}
	if imageTags["topo"] == "" {
		imageTags["topo"] = string(k8s.Debug)
	}
	if imageTags["gui"] == "" {
		imageTags["gui"] = string(k8s.Latest)
	}
	if imageTags["cli"] == "" {
		imageTags["cli"] = string(k8s.Latest)
	}
	if imageTags["atomix"] == "" {
		imageTags["atomix"] = string(k8s.Latest)
	}
	if imageTags["raft"] == "" {
		imageTags["raft"] = string(k8s.Latest)
	}
	if imageTags["simulator"] == "" {
		imageTags["simulator"] = string(k8s.Latest)
	}
	if imageTags["stratum"] == "" {
		imageTags["stratum"] = string(k8s.Latest)
	}
	if imageTags["test"] == "" {
		imageTags["test"] = string(k8s.Latest)
	}

}

// Contains tells whether array contains x.
func Contains(array []string, elem string) bool {
	for _, n := range array {
		if elem == n {
			return true
		}
	}
	return false
}
