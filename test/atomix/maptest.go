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
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAtomixMap(t *testing.T) {
	client, err := env.NewAtomixClient("map")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetGroup(context.Background(), "raft")
	assert.NoError(t, err)

	map_, err := group.GetMap(context.Background(), "test", session.WithTimeout(5 * time.Second))
	assert.NoError(t, err)

	size, err := map_.Size(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	value, err := map_.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.Nil(t, value)

	value, err = map_.Put(context.Background(), "foo", []byte("Hello world!"))
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "foo", value.Key)
	assert.True(t, bytes.Equal([]byte("Hello world!"), value.Value))
	assert.NotEqual(t, int64(0), value.Version)
	version := value.Version

	value, err = map_.Get(context.Background(), "foo")
	assert.NoError(t, err)
	assert.NotNil(t, value)
	assert.Equal(t, "foo", value.Key)
	assert.True(t, bytes.Equal([]byte("Hello world!"), value.Value))
	assert.Equal(t, version, value.Version)

	size, err = map_.Size(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, size)
}

func init() {
	test.Registry.RegisterTest("atomix-map", TestAtomixMap, []*runner.TestSuite{AtomixTests})
}
