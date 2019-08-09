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

package k8s

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

// GetFreePort asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	debugPort := DebugPort

	for i := 0; i < count; {
		host := "localhost:" + strconv.Itoa(debugPort)
		addr, err := net.ResolveTCPAddr("tcp", host)
		if err != nil {
			return nil, err
		}
		_, err = net.ListenTCP("tcp", addr)
		if err == nil {
			ports = append(ports, debugPort)
			i++
		}
		debugPort++
	}
	return ports, nil
}

// imageName returns a fully qualified name for the given image
func (c *ClusterController) imageName(image string, tag string) string {
	imageName := bytes.Buffer{}
	imageName.WriteString(c.imagePrefix())
	imageName.WriteString(image)
	imageName.WriteString(":")
	imageName.WriteString(tag)
	return imageName.String()
}

// imagePrefix returns a prefix for images
func (c *ClusterController) imagePrefix() string {
	if c.config.Registry != "" {
		return fmt.Sprintf("%s/", c.config.Registry)
	}
	return ""
}
