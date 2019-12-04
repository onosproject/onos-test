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

package nopaxos

import (
	"context"
	"github.com/atomix/atomix-go-client/pkg/client/map"
	"github.com/onosproject/onos-test/pkg/benchmark"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
)

// MapBenchmarkSuite :: benchmark
type MapBenchmarkSuite struct {
	benchmark.Suite
	m _map.Map
}

// SetupBenchmarkSuite :: benchmark
func (s *MapBenchmarkSuite) SetupBenchmarkSuite(c *benchmark.Context) {
	setup.Partitions("nopaxos").
		NOPaxos().
		SetPartitions(c.GetArg("partitions").Int(1)).
		SetReplicasPerPartition(c.GetArg("replicas").Int(1))
	setup.SetupOrDie()
}

// SetupBenchmark :: benchmark
func (s *MapBenchmarkSuite) SetupBenchmark(b *benchmark.Benchmark) {
	group, err := env.Database().Partitions("nopaxos").Connect()
	if err != nil {
		panic(err)
	}
	m, err := group.GetMap(context.Background(), b.Name)
	if err != nil {
		panic(err)
	}
	s.m = m
}

// BenchmarkMapPut :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapPut(b *benchmark.Benchmark) {
	params := []benchmark.Param{
		benchmark.RandomString(b.GetArg("key-count").Int(1000), b.GetArg("key-length").Int(8)),
		benchmark.RandomBytes(b.GetArg("value-count").Int(1), b.GetArg("value-length").Int(128)),
	}
	b.Run(func(client _map.Map, key string, value []byte) error {
		_, err := client.Put(context.Background(), key, value)
		return err
	}, params...)
}

// BenchmarkMapGet :: benchmark
func (s *MapBenchmarkSuite) BenchmarkMapGet(b *benchmark.Benchmark) {
	params := []benchmark.Param{
		benchmark.RandomString(b.GetArg("key-count").Int(1000), b.GetArg("key-length").Int(8)),
	}
	b.Run(func(client _map.Map, key string) error {
		_, err := client.Get(context.Background(), key)
		return err
	}, params...)
}
