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

package utils_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Chris-Greaves/stencil/utils"
	"github.com/stretchr/testify/assert"
)

func TestOsExecClient_ExecuteCommandInDir(t *testing.T) {
	e := &utils.OsExecClient{}
	t.Run("Can execute a command in a directory", func(t *testing.T) {
		err := e.ExecuteCommandInDir(".", "ls")
		assert.NoError(t, err)
	})

	t.Run("Error when path doesn't exist", func(t *testing.T) {
		err := e.ExecuteCommandInDir("/non/existent/path", "ls")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, os.ErrNotExist)
		}
	})

	t.Run("Error when command doesn't exist", func(t *testing.T) {
		err := e.ExecuteCommandInDir(".", "nonexistentcommand")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, exec.ErrNotFound)
		}
	})
}

func TestOsExecClient_CheckPathForExecutable(t *testing.T) {
	e := &utils.OsExecClient{}
	t.Run("Check for existing executable", func(t *testing.T) {
		path, err := e.CheckPathForExecutable("ls")
		assert.NoError(t, err)
		assert.NotEmpty(t, path)
	})

	t.Run("Check for non-existing executable", func(t *testing.T) {
		path, err := e.CheckPathForExecutable("nonexistentcommand")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, exec.ErrNotFound)
		}
		assert.Empty(t, path)
	})
}
