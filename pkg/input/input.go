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

package input

import (
	"math/rand"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNNOPQRSTUVWXYZ1234567890"

// Source is an interface for input sources
type Source interface {
	// Reset resets the input
	Reset()

	// Next returns the Next input value
	Next() Value
}

// initSource initializes and returns the given input
func initSource(gen Source) Source {
	gen.Reset()
	return gen
}

// RandomChoice returns a input that chooses a random value from a set
func RandomChoice(set Source) Source {
	return initSource(&RandomChoiceSource{
		set: set,
	})
}

// RandomChoiceSource is a input that chooses a random value from a set
type RandomChoiceSource struct {
	set    Source
	values []Value
}

// Reset resets the input
func (p *RandomChoiceSource) Reset() {
	p.set.Reset()
	p.values = p.set.Next().Slice()
}

// Next returns the next random value from the set
func (p *RandomChoiceSource) Next() Value {
	return p.values[rand.Intn(len(p.values))]
}

// SetOf returns a input that constructs a set of values of the given input type
func SetOf(valueType Source, size int) Source {
	return initSource(&SetSource{
		valueType: valueType,
		size:      size,
	})
}

// SetSource is a input that returns a random set of values
type SetSource struct {
	valueType Source
	size      int
	set       Value
}

// Reset resets the input
func (p *SetSource) Reset() {
	p.valueType.Reset()
	set := make([]Value, p.size)
	for i := 0; i < p.size; i++ {
		set[i] = p.valueType.Next()
	}
	p.set = NewValue(set)
}

// Next returns the random set
func (p *SetSource) Next() Value {
	return p.set
}

// RandomString returns a random string input
func RandomString(length int) Source {
	return initSource(&RandomStringSource{
		length: length,
	})
}

// RandomStringSource is a random string input
type RandomStringSource struct {
	length int
}

// Reset resets the input
func (a *RandomStringSource) Reset() {}

// Next returns a random string
func (a *RandomStringSource) Next() Value {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return NewValue(string(bytes))
}

// RandomBytes returns a random bytes input
func RandomBytes(length int) Source {
	return initSource(&RandomBytesSource{
		length: length,
	})
}

// RandomBytesSource is a random string input
type RandomBytesSource struct {
	length int
}

// Reset resets the input
func (a *RandomBytesSource) Reset() {}

// Next returns a random byte slice
func (a *RandomBytesSource) Next() Value {
	bytes := make([]byte, a.length)
	for j := 0; j < a.length; j++ {
		bytes[j] = chars[rand.Intn(len(chars))]
	}
	return NewValue(bytes)
}
