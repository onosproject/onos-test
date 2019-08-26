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

package atomix

import (
	"context"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

// BenchAtomixMap : benchmark
func BenchAtomixMap(b *testing.B) {
	client, err := env.NewAtomixClient("map")
	assert.NoError(b, err)
	assert.NotNil(b, client)

	group, err := client.GetGroup(context.Background(), "raft")
	assert.NoError(b, err)
	assert.NotNil(b, group)

	m, err := group.GetMap(context.Background(), "bench")
	assert.NoError(b, err)
	assert.NotNil(b, m)

	b.Run("map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = m.Put(context.Background(), "foo", []byte("Hello world!"))
		}
	})
}

func init() {
	test.Registry.RegisterBench("atomix-map", BenchAtomixMap, []*runner.BenchSuite{AtomixBenchmarks})
}
