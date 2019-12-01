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
	"testing"
)

type MapBenchmarkSuite struct {
	onit.ScriptSuite
}

func (b *MapBenchmarkSuite) SetupScriptSuite() {
	setup.Partitions("nopaxos").
		NOPaxos().
		SetPartitions(1).
		SetReplicasPerPartition(3)
	setup.SetupOrDie()
}

func (b *MapBenchmarkSuite) getHandler(name string) func() benchmark.Handler {
	return func() benchmark.Handler {
		group, err := env.Database().Partitions("nopaxos").Connect()
		if err != nil {
			panic(err)
		}
		m, err := group.GetMap(context.Background(), name)
		if err != nil {
			panic(err)
		}
		return &MapBenchmarkPutHandler{
			m: m,
		}
	}
}

func (b *MapBenchmarkSuite) RunBenchmarkMapPut() {
	benchmark.New().
		SetHandlerFactory(b.getHandler("RunBenchmarkMapPut")).
		SetParallelism(1).
		SetIterations(1000).
		AddHandlerArg(benchmark.RandomString(1000, 8)).
		AddHandlerArg(benchmark.RandomBytes(1000, 128)).
		Run()
}

type MapBenchmarkPutHandler struct {
	m _map.Map
}

func (h *MapBenchmarkPutHandler) Run(args ...interface{}) error {
	_, err := h.m.Put(context.Background(), args[0].(string), args[1].([]byte))
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
