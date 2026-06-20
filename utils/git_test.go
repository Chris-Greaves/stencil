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
	"os/exec"
	"testing"

	"github.com/Chris-Greaves/stencil/test/mocks"
	"github.com/Chris-Greaves/stencil/utils"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGitClient(t *testing.T) {
	gitClient := utils.DefaultGitClient()
	assert.NotNil(t, gitClient)
}

func TestGitClient_IsGitInstalled(t *testing.T) {
	t.Run("verify git is installed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().CheckPathForExecutable("git").Return("/usr/bin/git", nil)
		gitClient := utils.NewGitClient(mockExecutor)

		assert.True(t, gitClient.IsGitInstalled(), "Expected git to be installed, but it was not detected.")
	})

	t.Run("verify git is not installed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().CheckPathForExecutable("git").Return("", exec.ErrNotFound)
		gitClient := utils.NewGitClient(mockExecutor)

		assert.False(t, gitClient.IsGitInstalled(), "Expected git to not be installed, but it was detected.")
	})
}

func TestGitClient_RunCommand(t *testing.T) {
	t.Run("can run a command", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"status"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.RunCommand(".", "status")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"status"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.RunCommand(".", "status")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_Initialize(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"init"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Initialize(".")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"init"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Initialize(".")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_AddRemote(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"remote", "add", "origin", "https://github.com/user/repo.git"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.AddRemote(".", "origin", "https://github.com/user/repo.git")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"remote", "add", "origin", "https://github.com/user/repo.git"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.AddRemote(".", "origin", "https://github.com/user/repo.git")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_Fetch(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"fetch"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Fetch(".")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"fetch"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Fetch(".")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_FetchRemote(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"fetch", "origin"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.FetchRemote(".", "origin")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"fetch", "origin"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.FetchRemote(".", "origin")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_Pull(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"pull"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Pull(".")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"pull"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.Pull(".")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func TestGitClient_PullRemote(t *testing.T) {
	t.Run("command can succeed", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"pull", "origin", "main"}).Return(nil)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.PullRemote(".", "origin", "main")
		assert.NoError(t, err)
	})

	t.Run("will return error when command fails", func(t *testing.T) {
		mockExecutor := mocks.NewMockExecutor(t)
		mockExecutor.EXPECT().ExecuteCommandInDir(".", "git", []string{"pull", "origin", "main"}).Return(assert.AnError)
		gitClient := utils.NewGitClient(mockExecutor)

		err := gitClient.PullRemote(".", "origin", "main")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}
