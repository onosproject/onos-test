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

func newStorage(cluster *Cluster) *Storage {
	return &Storage{
		client:    cluster.client,
		cluster:   cluster,
		databases: make(map[string]*Database),
	}
}

// Storage provides methods for managing Atomix storage
type Storage struct {
	*client
	cluster   *Cluster
	databases map[string]*Database
}

// Database returns a database by name
func (s *Storage) Database(name string) *Database {
	if database, ok := s.databases[name]; ok {
		return database
	}
	database := newDatabase(s.cluster, name)
	s.databases[name] = database
	return database
}
