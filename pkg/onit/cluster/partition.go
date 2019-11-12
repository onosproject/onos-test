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

import "strings"

func newPartition(name string, client *client) *Partition {
	group := name[:strings.LastIndex(name, "-")]
	partition := name[strings.LastIndex(name, "-")+1:]
	labels := map[string]string{
		typeLabel:   partitionType.name(),
		"group":     group,
		"partition": partition,
	}
	return &Partition{
		Service: newService(name, 5678, labels, raftImage, client),
	}
}

// Partition provides methods for querying a database partition
type Partition struct {
	*Service
}
