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
	"fmt"
	"github.com/atomix/atomix-go-client/pkg/client/map"
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/google/uuid"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestRaftHA : integration test
func TestRaftHA(t *testing.T) {
	nodes := env.GetRaftNodes()

	client, err := env.NewAtomixClient("map")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetGroup(context.Background(), "raft")
	assert.NoError(t, err)

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
	assert.Equal(t, key, event.Key)
	assert.Equal(t, "foo", string(event.Value))
	assert.Equal(t, version, event.Version)

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
	assert.Equal(t, key, event.Key)
	assert.Equal(t, "bar", string(event.Value))

	key = uuid.New().String()
	entry, err = m.Put(context.Background(), key, []byte("baz"))
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)

	event = <-ch
	assert.Equal(t, _map.EventInserted, event.Type)
	assert.Equal(t, key, event.Key)
	assert.Equal(t, "baz", string(event.Value))

	i := 0
	for {
		for _, partition := range nodes {
			if len(partition) == 1 || len(partition) <= i {
				return
			}

			key := uuid.New().String()
			entry, err = m.Put(context.Background(), key, []byte(uuid.New().String()))
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)

			println(fmt.Sprintf("Killing Raft node %s", partition[i]))
			err = env.KillNode(partition[i])
			assert.NoError(t, err)

			event = <-ch
			assert.Equal(t, key, event.Key)

			entry, err = m.Get(context.Background(), key)
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)
		}

		println(fmt.Sprintf("Sleeping for 15 seconds"))
		time.Sleep(10 * time.Second)

		for range nodes {
			key := uuid.New().String()
			entry, err = m.Put(context.Background(), key, []byte(uuid.New().String()))
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)

			event = <-ch
			assert.Equal(t, key, event.Key)

			entry, err = m.Get(context.Background(), key)
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)
		}

		println("Waiting for pods to recover")
		for _, partition := range nodes {
			env.AwaitReady(partition[i])
		}
		i++
	}
}

func init() {
	test.Registry.RegisterTest("atomix-ha", TestRaftHA, []*runner.TestSuite{AtomixTests})
}
