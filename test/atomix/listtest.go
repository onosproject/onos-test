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
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/onosproject/onos-test/pkg/runner"
	"github.com/onosproject/onos-test/test"
	"github.com/onosproject/onos-test/test/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestAtomixList : integration test
func TestAtomixList(t *testing.T) {
	client, err := env.NewAtomixClient("list")
	assert.NoError(t, err)
	assert.NotNil(t, client)

	group, err := client.GetGroup(context.Background(), "raft")
	assert.NoError(t, err)

	list, err := group.GetList(context.Background(), "test", session.WithTimeout(5*time.Second))
	assert.NoError(t, err)

	size, err := list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = list.Append(context.Background(), "Hello world!")
	assert.NoError(t, err)

	size, err = list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, size)

	value, err := list.Get(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!", value)

	value, err = list.Remove(context.Background(), 0)
	assert.NoError(t, err)
	assert.Equal(t, "Hello world!", value)

	size, err = list.Len(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, size)

	err = list.Append(context.Background(), "Hello world!")
	assert.NoError(t, err)

	err = list.Append(context.Background(), "Hello world again!")
	assert.NoError(t, err)

	ch := make(chan string)
	err = list.Items(context.Background(), ch)
	i := 0
	for value := range ch {
		if i == 0 {
			assert.Equal(t, "Hello world!", value)
			i++
		} else if i == 1 {
			assert.Equal(t, "Hello world again!", value)
			i++
		} else {
			assert.Fail(t, "Too many values")
		}
	}
	assert.NoError(t, err)
}

func init() {
	test.Registry.RegisterTest("atomix-list", TestAtomixList, []*runner.TestSuite{AtomixTests})
}
