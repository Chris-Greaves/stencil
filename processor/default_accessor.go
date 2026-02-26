// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package processor

// DefaultAccessor is a copy of EmbeddedAccessor, which allows for the setting of default values
type DefaultAccessor[T any] struct {
	value T
}

// Get gets the value.
func (a *DefaultAccessor[T]) Get() T {
	return a.value
}

// Set sets the value.
func (a *DefaultAccessor[T]) Set(value T) {
	a.value = value
}

func NewDefaultAccessor[T any](value T) *DefaultAccessor[T] {
	return &DefaultAccessor[T]{
		value: value,
	}
}
