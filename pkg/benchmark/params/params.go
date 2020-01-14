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

package params

import (
	"github.com/onosproject/onos-test/pkg/benchmark"
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// RandomChoice returns a parameter that chooses a random value from a set
func RandomChoice(set benchmark.Param) benchmark.Param {
	return &RandomChoiceParam{
		set: set,
	}
}

// RandomChoiceParam is a parameter that chooses a random value from a set
type RandomChoiceParam struct {
	set    benchmark.Param
	values []interface{}
}

// Reset resets the parameter
func (p *RandomChoiceParam) Reset() {
	p.set.Reset()
	p.values = p.set.Next().([]interface{})
}

// Next returns the next random value from the set
func (p *RandomChoiceParam) Next() interface{} {
	return p.values[rand.Intn(len(p.values))]
}

// SetOf returns a parameter that constructs a set of values of the given parameter type
func SetOf(valueType benchmark.Param, size int) benchmark.Param {
	return &SetParam{
		valueType: valueType,
		size:      size,
	}
}

// SetParam is a parameter that returns a random set of values
type SetParam struct {
	valueType benchmark.Param
	size      int
	set       []interface{}
}

// Reset resets the parameter
func (p *SetParam) Reset() {
	p.valueType.Reset()
	p.set = make([]interface{}, p.size)
	for i := 0; i < p.size; i++ {
		p.set[i] = p.valueType.Next()
	}
}

// Next returns the random set
func (p *SetParam) Next() interface{} {
	return p.set
}

// RandomString returns a random string parameter
func RandomString(length int) benchmark.Param {
	return &RandomStringParam{
		length: length,
	}
}

// RandomStringParam is a random string parameter
type RandomStringParam struct {
	length int
}

// Reset resets the parameter
func (a *RandomStringParam) Reset() {}

// Next returns a random string
func (a *RandomStringParam) Next() interface{} {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return string(bytes)
}

// RandomBytes returns a random bytes parameter
func RandomBytes(length int) benchmark.Param {
	return &RandomBytesParam{
		length: length,
	}
}

// RandomBytesParam is a random string parameter
type RandomBytesParam struct {
	length int
}

// Reset resets the parameter
func (a *RandomBytesParam) Reset() {}

// Next returns a random byte slice
func (a *RandomBytesParam) Next() interface{} {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return bytes
}
