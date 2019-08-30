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
	"bytes"
	"context"
	"github.com/atomix/atomix-go-client/pkg/client/map"
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestAtomixMap : integration test
func TestAtomixMap(t *testing.T) {
	client, err := env.NewAtomixClient("map")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetGroup(context.Background(), "raft")
	assert.NoError(t, err)

	m, err := group.GetMap(context.Background(), "test", session.WithTimeout(5*time.Second))
	assert.NoError(t, err)

	ch := make(chan *_map.KeyValue)
	err = m.Entries(context.Background(), ch)
	assert.NoError(t, err)
	for range ch {
		assert.Fail(t, "entries found in map")
	}

	size, err := m.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	value, err := m.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.Nil(t, value)

	value, err = m.Put(context.Background(), "foo", []byte("Hello world!"))
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "foo", value.Key)
	assert.True(t, bytes.Equal([]byte("Hello world!"), value.Value))
	assert.NotEqual(t, int64(0), value.Version)
	version := value.Version

	value, err = m.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "foo", value.Key)
	assert.True(t, bytes.Equal([]byte("Hello world!"), value.Value))
	assert.Equal(t, version, value.Version)

	size, err = m.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, size)

	ch = make(chan *_map.KeyValue)
	err = m.Entries(context.Background(), ch)
	assert.NoError(t, err)
	i := 0
	for kv := range ch {
		assert.Equal(t, "foo", kv.Key)
		assert.Equal(t, "Hello world!", string(kv.Value))
		i++
	}
	assert.Equal(t, 1, i)

	allEvents := make(chan *_map.Event)
	err = m.Watch(context.Background(), allEvents, _map.WithReplay())
	assert.NoError(t, err)

	event := <-allEvents
	assert.NotNil(t, event)
	assert.Equal(t, "foo", event.Key)
	assert.Equal(t, []byte("Hello world!"), event.Value)
	assert.Equal(t, value.Version, event.Version)

	futureEvents := make(chan *_map.Event)
	err = m.Watch(context.Background(), futureEvents)
	assert.NoError(t, err)

	value, err = m.Put(context.Background(), "bar", []byte("Hello world!"))
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "bar", value.Key)
	assert.Equal(t, []byte("Hello world!"), value.Value)
	assert.NotEqual(t, int64(0), value.Version)

	event = <-allEvents
	assert.NotNil(t, event)
	assert.Equal(t, "bar", event.Key)
	assert.Equal(t, []byte("Hello world!"), event.Value)
	assert.Equal(t, value.Version, event.Version)

	event = <-futureEvents
	assert.NotNil(t, event)
	assert.Equal(t, "bar", event.Key)
	assert.Equal(t, []byte("Hello world!"), event.Value)
	assert.Equal(t, value.Version, event.Version)
}

func init() {
	test.Registry.RegisterTest("atomix-map", TestAtomixMap, []*runner.TestSuite{AtomixTests})
}
