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

package raft

import (
	"context"
	"github.com/atomix/atomix-go-client/pkg/client/map"
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/google/uuid"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestRaftHA : integration test
func (s *HATestSuite) TestRaftHA(t *testing.T) {
	partitions := env.Database().Partitions("raft")
	group, err := partitions.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, group)

	m, err := group.GetMap(context.Background(), "TestRaftHA", session.WithTimeout(5*time.Second))
	assert.NoError(t, err)

	ch := make(chan *_map.Event)
	err = m.Watch(context.Background(), ch)
	assert.NoError(t, err)

	key := uuid.New().String()
	entry, err := m.Put(context.Background(), key, []byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)
	assert.Equal(t, "foo", string(entry.Value))
	version := entry.Version

	event := <-ch
	assert.Equal(t, _map.EventInserted, event.Type)
	assert.Equal(t, key, event.Entry.Key)
	assert.Equal(t, "foo", string(event.Entry.Value))
	assert.Equal(t, version, event.Entry.Version)

	entry, err = m.Get(context.Background(), key)
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)
	assert.Equal(t, "foo", string(entry.Value))
	assert.Equal(t, version, entry.Version)

	key = uuid.New().String()
	entry, err = m.Put(context.Background(), key, []byte("bar"))
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)
	assert.Equal(t, "bar", string(entry.Value))

	event = <-ch
	assert.Equal(t, _map.EventInserted, event.Type)
	assert.Equal(t, key, event.Entry.Key)
	assert.Equal(t, "bar", string(event.Entry.Value))

	key = uuid.New().String()
	entry, err = m.Put(context.Background(), key, []byte("baz"))
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)

	event = <-ch
	assert.Equal(t, _map.EventInserted, event.Type)
	assert.Equal(t, key, event.Entry.Key)
	assert.Equal(t, "baz", string(event.Entry.Value))

	i := 0
	for {
		for _, partition := range partitions.List() {
			if len(partition.Nodes()) == 1 || len(partition.Nodes()) <= i {
				return
			}

			key := uuid.New().String()
			entry, err = m.Put(context.Background(), key, []byte(uuid.New().String()))
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)

			t.Logf("Killing Raft node %s", partition.Nodes()[i].Name())
			err = partition.Nodes()[i].Kill()
			assert.NoError(t, err)

			event = <-ch
			assert.Equal(t, key, event.Entry.Key)

			entry, err = m.Get(context.Background(), key)
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)
		}

		t.Log("Sleeping for 15 seconds")
		time.Sleep(10 * time.Second)

		for range partitions.List() {
			key := uuid.New().String()
			entry, err = m.Put(context.Background(), key, []byte(uuid.New().String()))
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)

			event = <-ch
			assert.Equal(t, key, event.Entry.Key)

			entry, err = m.Get(context.Background(), key)
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)
		}

		t.Log("Waiting for pods to recover")
		for _, partition := range partitions.List() {
			err = partition.AwaitReady()
			assert.NoError(t, err)
		}
		i++
	}
}
