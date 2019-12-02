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
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/onosproject/onos-test/pkg/test"
)

type MapBenchmarkSuite struct {
	test.BenchmarkSuite
	m _map.Map
}

func (s *MapBenchmarkSuite) SetupBenchmarkSuite(c *test.Context) {
	setup.Partitions("nopaxos").
		NOPaxos().
		SetPartitions(c.GetArg("partitions").Int()).
		SetReplicasPerPartition(c.GetArg("replicas").Int())
	setup.SetupOrDie()
}

func (s *MapBenchmarkSuite) SetupBenchmark(b *test.Benchmark) {
	b.Init(func() (_map.Map, error) {
		group, err := env.Database().Partitions("nopaxos").Connect()
		if err != nil {
			panic(err)
		}
		m, err := group.GetMap(context.Background(), b.Name)
		if err != nil {
			return nil, err
		}
		return m, nil
	})
}

func (s *MapBenchmarkSuite) BenchmarkMapPut(b *test.Benchmark) {
	params := []test.Param{
		test.RandomString(b.GetArg("key-count").Int(), b.GetArg("key-length").Int()),
		test.RandomBytes(b.GetArg("value-count").Int(), b.GetArg("value-length").Int()),
	}
	b.Run(func(client _map.Map, key string, value []byte) error {
		_, err := client.Put(context.Background(), key, value)
		return err
	}, params...)
}

func (s *MapBenchmarkSuite) BenchmarkMapGet(b *test.Benchmark) {
	params := []test.Param{
		test.RandomString(b.GetArg("key-count").Int(), b.GetArg("key-length").Int()),
	}
	b.Run(func(client _map.Map, key string) error {
		_, err := client.Get(context.Background(), key)
		return err
	}, params...)
}
