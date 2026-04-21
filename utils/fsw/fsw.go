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

package fsw

import (
	"os"

	"github.com/Chris-Greaves/stencil/utils"
)

var wrapper FsWrapper

type FsWrapper interface {
	GetHomeDirectory() (string, error)
	ReadFile(name string) ([]byte, error)
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(name string, data []byte, perm os.FileMode) error
	Lstat(name string) (os.FileInfo, error)
}

func UseDefaultWrapper() {
	wrapper = utils.FsHelper{}
}

func UseCustomeWrapper(w FsWrapper) {
	wrapper = w
}

func GetHomeDirectory() (string, error) {
	return wrapper.GetHomeDirectory()
}
func ReadFile(name string) ([]byte, error) {
	return wrapper.ReadFile(name)
}
func MkdirAll(path string, perm os.FileMode) error {
	return wrapper.MkdirAll(path, perm)
}
func WriteFile(name string, data []byte, perm os.FileMode) error {
	return wrapper.WriteFile(name, data, perm)
}
func Lstat(name string) (os.FileInfo, error) {
	return wrapper.Lstat(name)
}
