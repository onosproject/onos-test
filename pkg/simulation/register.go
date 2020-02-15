// Copyright 2020-present Open Networking Foundation.
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

package simulation

// Register records simulation traces
type Register interface {
	// Trace records a trace to the register
	Trace(values ...interface{})

	// close closes the register
	close()
}

// newChannelRegister returns a new register that records to the given channel
func newChannelRegister(ch chan<- []interface{}) Register {
	return &channelRegister{
		ch: ch,
	}
}

// channelRegister is a register that writes records to a channel
type channelRegister struct {
	ch chan<- []interface{}
}

func (r *channelRegister) Trace(values ...interface{}) {
	r.ch <- values
}

func (r *channelRegister) close() {
	close(r.ch)
}
