// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import "time"

func Retry(count int, delay time.Duration, check func() bool) bool {
	for attempt := 1; attempt != count; attempt++ {
		if check() {
			return true
		}
		time.Sleep(delay)
	}
	return false
}
