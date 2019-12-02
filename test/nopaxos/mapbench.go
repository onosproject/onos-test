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
	"github.com/onosproject/onos-test/pkg/onit"
	"github.com/onosproject/onos-test/pkg/onit/benchmark"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	"testing"
)

type MapBenchmarkSuite struct {
	onit.ScriptSuite
}

func (b *MapBenchmarkSuite) SetupScriptSuite() {
	setup.Atomix().
		SetImage("192.168.1.11:30000/atomix/atomix-k8s-controller:latest").
		SetPullPolicy(v1.PullAlways)
	setup.Partitions("nopaxos").
		NOPaxos().
		SetPartitions(1).
		SetReplicasPerPartition(3).
		SetReplicaImage("192.168.1.11:30000/atomix/atomix-nopaxos-node:latest").
		SetSequencerImage("192.168.1.11:30000/atomix/atomix-nopaxos-sequencer:latest").
		SetPullPolicy(v1.PullAlways)
	setup.SetupOrDie()
}

func (b *MapBenchmarkSuite) getMap(name string) _map.Map {
	group, err := env.Database().Partitions("nopaxos").Connect()
	if err != nil {
		panic(err)
	}
	m, err := group.GetMap(context.Background(), name)
	if err != nil {
		panic(err)
	}
	return m
}

func (b *MapBenchmarkSuite) RunBenchmarkMapPut() {
	benchmark.New().
		SetHandlerFactory(func() benchmark.Handler {
			return &MapBenchmarkPutHandler{
				m: b.getMap("RunBenchmarkMapPut"),
			}
		}).
		SetClients(20).
		SetParallelism(10).
		SetRequests(100000).
		AddHandlerArg(benchmark.RandomString(1000, 8)).
		AddHandlerArg(benchmark.RandomBytes(1000, 128)).
		Run()
}

func (b *MapBenchmarkSuite) RunBenchmarkMapGet() {
	benchmark.New().
		SetHandlerFactory(func() benchmark.Handler {
			return &MapBenchmarkGetHandler{
				m: b.getMap("RunBenchmarkMapGet"),
			}
		}).
		SetClients(20).
		SetParallelism(10).
		SetRequests(100000).
		AddHandlerArg(benchmark.RandomString(1000, 8)).
		Run()
}

type MapBenchmarkPutHandler struct {
	m _map.Map
}

func (h *MapBenchmarkPutHandler) Run(args ...interface{}) error {
	_, err := h.m.Put(context.Background(), args[0].(string), args[1].([]byte))
	return err
}

type MapBenchmarkGetHandler struct {
	m _map.Map
}

func (h *MapBenchmarkGetHandler) Run(args ...interface{}) error {
	_, err := h.m.Get(context.Background(), args[0].(string))
	return err
}

// BenchmarkNOPaxosMap : benchmark
func (s *BenchmarkSuite) BenchmarkNOPaxosMap(b *testing.B) {
	group, err := env.Database().Partitions("nopaxos").Connect()
	assert.NoError(b, err)
	assert.NotNil(b, group)

	m, err := group.GetMap(context.Background(), "BenchmarkNOPaxosMap")
	assert.NoError(b, err)
	assert.NotNil(b, m)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Put(context.Background(), "foo", []byte("bar"))
	}
}
