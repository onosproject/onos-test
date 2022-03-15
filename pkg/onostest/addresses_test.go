// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package onostest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestPlaceHolder(t *testing.T) {
	assert.Equal(t, "A-B-atomix", AtomixName("A", "B"))
}
