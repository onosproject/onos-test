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
	"github.com/stretchr/testify/assert"
	"testing"
)

// BenchmarkAtomixMap : benchmark
func (s *BenchmarkSuite) BenchmarkAtomixMap(b *testing.B) {
	env := s.Env()

	group, err := env.Database().Partitions("raft").Connect()
	assert.NoError(b, err)
	assert.NotNil(b, group)

	m, err := group.GetMap(context.Background(), "BenchmarkAtomixMap")
	assert.NoError(b, err)
	assert.NotNil(b, m)

	b.Run("map", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = m.Put(context.Background(), "foo", []byte("Hello world!"))
		}
	})
}
