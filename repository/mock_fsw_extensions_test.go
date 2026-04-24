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

package repository

import "os"

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_GetHomeDirectory_Call) Passthrough() *MockFsWrapper_GetHomeDirectory_Call {
	_c.Call.Return(os.UserHomeDir())
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_MkdirAll_Call) Passthrough() *MockFsWrapper_MkdirAll_Call {
	args := _c.Arguments
	_c.Call.Return(os.MkdirAll(args[0].(string), args[1].(os.FileMode)))
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_ReadFile_Call) Passthrough() *MockFsWrapper_ReadFile_Call {
	args := _c.Arguments
	_c.Call.Return(os.ReadFile(args[0].(string)))
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_WriteFile_Call) Passthrough() *MockFsWrapper_WriteFile_Call {
	args := _c.Arguments
	_c.Call.Return(os.WriteFile(args[0].(string), args[1].([]byte), args[2].(os.FileMode)))
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_Lstat_Call) Passthrough() *MockFsWrapper_Lstat_Call {
	args := _c.Arguments
	_c.Call.Return(os.Lstat(args[0].(string)))
	return _c
}
