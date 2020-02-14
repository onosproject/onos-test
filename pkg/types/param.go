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

package types

import (
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// Param is an interface for parameters
type Param interface {
	// Reset resets the parameter
	Reset()

	// Next returns the Next parameter value
	Next() Value
}

// RandomChoice returns a parameter that chooses a random value from a set
func RandomChoice(set Param) Param {
	return &RandomChoiceParam{
		set: set,
	}
}

// RandomChoiceParam is a parameter that chooses a random value from a set
type RandomChoiceParam struct {
	set    Param
	values []Value
}

// Reset resets the parameter
func (p *RandomChoiceParam) Reset() {
	p.set.Reset()
	p.values = p.set.Next().Slice()
}

// Next returns the next random value from the set
func (p *RandomChoiceParam) Next() Value {
	return p.values[rand.Intn(len(p.values))]
}

// SetOf returns a parameter that constructs a set of values of the given parameter type
func SetOf(valueType Param, size int) Param {
	return &SetParam{
		valueType: valueType,
		size:      size,
	}
}

// SetParam is a parameter that returns a random set of values
type SetParam struct {
	valueType Param
	size      int
	set       Value
}

// Reset resets the parameter
func (p *SetParam) Reset() {
	p.valueType.Reset()
	set := make([]Value, p.size)
	for i := 0; i < p.size; i++ {
		set[i] = p.valueType.Next()
	}
	p.set = NewValue(set)
}

// Next returns the random set
func (p *SetParam) Next() Value {
	return p.set
}

// RandomString returns a random string parameter
func RandomString(length int) Param {
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
func (a *RandomStringParam) Next() Value {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return NewValue(string(bytes))
}

// RandomBytes returns a random bytes parameter
func RandomBytes(length int) Param {
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
func (a *RandomBytesParam) Next() Value {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return NewValue(bytes)
}
