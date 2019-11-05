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

func newDatabase(client *client) *Database {
	return &Database{
		client: client,
	}
}

// Database provides methods for managing the Atomix database
type Database struct {
	*client
}

// Partitions returns a list of partitions in the database
func (s *Database) Partitions(group string) *Partitions {
	return newPartitions(group, s.client)
}
