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
	atomixlock "github.com/atomix/atomix-go-client/pkg/client/lock"
	"github.com/atomix/atomix-go-client/pkg/client/session"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

// TestAtomixLock : integration test
func (s *SmokeTestSuite) TestAtomixLock(t *testing.T) {
	group, err := env.Database().Partitions("raft").Connect()
	assert.NoError(t, err)
	assert.NotNil(t, group)

	lock1, err := group.GetLock(context.Background(), "TestAtomixLock", session.WithTimeout(5*time.Second))
	assert.NoError(t, err)

	lock2, err := group.GetLock(context.Background(), "TestAtomixLock", session.WithTimeout(5*time.Second))
	assert.NoError(t, err)

	id, err := lock1.Lock(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), id)

	var lock uint64
	wait := make(chan struct{})
	go func() {
		id, err := lock2.Lock(context.Background())
		assert.NoError(t, err)
		assert.NotEqual(t, uint64(0), id)
		atomic.StoreUint64(&lock, id)
		wait <- struct{}{}
	}()

	isLocked, err := lock1.IsLocked(context.Background())
	assert.NoError(t, err)
	assert.True(t, isLocked)

	isLocked, err = lock1.IsLocked(context.Background(), atomixlock.IfVersion(id))
	assert.NoError(t, err)
	assert.True(t, isLocked)

	isLocked, err = lock1.IsLocked(context.Background(), atomixlock.IfVersion(id+1))
	assert.NoError(t, err)
	assert.False(t, isLocked)

	unlocked, err := lock1.Unlock(context.Background())
	assert.NoError(t, err)
	assert.True(t, unlocked)

	<-wait

	id = atomic.LoadUint64(&lock)
	assert.NotEqual(t, uint64(0), id)

	isLocked, err = lock2.IsLocked(context.Background())
	assert.NoError(t, err)
	assert.True(t, isLocked)

	unlocked, err = lock1.Unlock(context.Background(), atomixlock.IfVersion(id))
	assert.NoError(t, err)
	assert.True(t, unlocked)

	isLocked, err = lock2.IsLocked(context.Background())
	assert.NoError(t, err)
	assert.False(t, isLocked)

	id, err = lock1.Lock(context.Background())
	assert.NoError(t, err)
	assert.NotEqual(t, uint64(0), id)

	lock = 0
	wait = make(chan struct{})
	go func() {
		id, err := lock2.Lock(context.Background(), atomixlock.WithTimeout(100*time.Millisecond))
		assert.NoError(t, err)
		atomic.StoreUint64(&lock, id)
		wait <- struct{}{}
	}()

	<-wait

	id = atomic.LoadUint64(&lock)
	assert.Equal(t, uint64(0), id)
}
