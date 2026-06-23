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

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Chris-Greaves/stencil/test/mocks"
	"github.com/Chris-Greaves/stencil/utils/fsw"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	defaultDirFileMode  = os.FileMode(0755)
	defaultFileFileMode = os.FileMode(0644)
)

func Test_readReposFile(t *testing.T) {
	defaultRM := NewDefaultRepositoryManager()
	t.Run("return an empty slice when the file does not exist", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)
		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, defaultDirFileMode), "failed to create .stencil directory in temporary home directory")
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().ReadFile(reposFilePath).Return(nil, os.ErrNotExist)

		// Act
		got, gotErr := defaultRM.readReposFile()

		// Assert
		if assert.NoError(t, gotErr, "error not expected") {
			assert.Empty(t, got, "expected an empty slice when the file does not exist")
		}
	})

	t.Run("return an error when the file is malformed", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, defaultDirFileMode), "failed to create .stencil directory in temporary home directory")

		if err := createMalformedReposFile(t, stencilDir); err != nil {
			t.Fatalf("failed to create malformed repos file: %v", err)
		}
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil).Once()
		m.EXPECT().ReadFile(reposFilePath).Passthrough().Once()

		// Act
		got, gotErr := defaultRM.readReposFile()

		// Assert
		m.AssertExpectations(t)
		assert.Error(t, gotErr, "error was expected, but got nil. Returned repositories: %v", got)
		assert.Nil(t, got, "expected nil when the file is malformed")
	})

	t.Run("return the expected repositories when the file is valid", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, defaultDirFileMode), "failed to create .stencil directory in temporary home directory")

		validReposContent := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		if err := createReposFile(t, stencilDir, validReposContent); err != nil {
			t.Fatalf("failed to create valid repos file: %v", err)
		}
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().ReadFile(reposFilePath).Passthrough().Once()

		// Act
		got, gotErr := defaultRM.readReposFile()

		// Assert
		assert.NoError(t, gotErr, "error not expected")
		assert.Equal(t, validReposContent, got)
	})

	t.Run("return the error when home directory cannot be determined", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("failed to determine home directory")
		m.EXPECT().GetHomeDirectory().Return("", returnErr)

		// Act
		got, gotErr := defaultRM.readReposFile()

		// Assert
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, returnErr)
		}
		assert.Nil(t, got, "expected nil when home directory cannot be determined")
	})

	t.Run("return the error when file cannot be read", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("failed to read file")
		m.EXPECT().GetHomeDirectory().Return("/some/path/", nil)
		m.EXPECT().ReadFile(filepath.Join("/some/path/", ".stencil", "repositories.json")).Return(nil, returnErr)

		// Act
		got, gotErr := defaultRM.readReposFile()

		// Assert
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, returnErr)
		}
		assert.Nil(t, got, "expected nil when file cannot be read, but got: %v", got)
	})
}

func Test_saveReposFile(t *testing.T) {
	defaultRM := NewDefaultRepositoryManager()
	t.Run("creates file if it does not exist yet", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		reposToSave := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		fileContents, err := json.Marshal(reposToSave)
		require.NoError(t, err, "failed to marshal repositories for test setup")

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, defaultDirFileMode), "failed to create .stencil directory in temporary home directory")

		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().MkdirAll(stencilDir, mock.Anything).Return(nil)
		m.EXPECT().WriteFile(reposFilePath, fileContents, defaultFileFileMode).Passthrough()

		// Act
		gotErr := defaultRM.saveReposFile(reposToSave)

		// Assert
		m.AssertExpectations(t)
		assert.NoError(t, gotErr, "error not expected")
		assert.FileExists(t, filepath.Join(stencilDir, "repositories.json"))
	})

	t.Run("creates folder and file if they do not exist yet", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		reposToSave := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		fileContents, err := json.Marshal(reposToSave)
		require.NoError(t, err, "failed to marshal repositories for test setup")

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")

		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().MkdirAll(stencilDir, defaultDirFileMode).Passthrough()
		m.EXPECT().WriteFile(reposFilePath, fileContents, defaultFileFileMode).Passthrough()

		// Act
		gotErr := defaultRM.saveReposFile(reposToSave)

		// Assert
		m.AssertExpectations(t)
		assert.NoError(t, gotErr, "error not expected")
		assert.FileExists(t, filepath.Join(stencilDir, "repositories.json"))
	})

	t.Run("returns error when file cannot be written", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		reposToSave := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		errToReturn := errors.New("failed to write file")

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")

		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().MkdirAll(stencilDir, defaultDirFileMode).Passthrough()
		m.EXPECT().WriteFile(reposFilePath, mock.Anything, defaultFileFileMode).
			Return(errToReturn)

		// Act
		gotErr := defaultRM.saveReposFile(reposToSave)

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, errToReturn)
		}
		assert.NoFileExists(t, filepath.Join(stencilDir, "repositories.json"))
	})

	t.Run("returns error when home directory cannot be determined", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		reposToSave := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		errToReturn := errors.New("failed to determine home directory")

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")

		m.EXPECT().GetHomeDirectory().Return("", errToReturn)

		// Act
		gotErr := defaultRM.saveReposFile(reposToSave)

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, errToReturn)
		}
		assert.NoFileExists(t, filepath.Join(stencilDir, "repositories.json"))
	})

	t.Run("returns error when folder cannot be created", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		reposToSave := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		errToReturn := errors.New("failed to create folder")

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")

		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().MkdirAll(stencilDir, mock.Anything).Return(errToReturn)

		// Act
		gotErr := defaultRM.saveReposFile(reposToSave)

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, errToReturn)
		}
		assert.NoFileExists(t, filepath.Join(stencilDir, "repositories.json"))
	})
}

func Test_AddRepository(t *testing.T) {
	defaultRM := NewDefaultRepositoryManager()
	t.Run("repository is added correctly", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		file := setupReposFile(t, m)
		fsw.UseCustomWrapper(m)

		m.EXPECT().MkdirAll(filepath.Dir(file), defaultDirFileMode).Passthrough()
		m.EXPECT().WriteFile(file, mock.Anything, mock.Anything).
			RunAndReturn(func(name string, data []byte, perm os.FileMode) error {
				return os.WriteFile(name, data, perm)
			})

		// Act
		err := defaultRM.AddRepository("test-repo", "example.org/test/stencils")

		// Assert
		m.AssertExpectations(t)
		assert.NoError(t, err)
		assert.FileExists(t, file)
		fileContents, readErr := os.ReadFile(file)
		assert.NoError(t, readErr, "failed to read the file that was created")
		assert.JSONEq(t, "[{\"name\": \"test-repo\", \"url\": \"example.org/test/stencils\"}]", string(fileContents))
	})
	t.Run("error is returned if reading the repos file fails", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("Bang!")
		m.EXPECT().GetHomeDirectory().Return("", returnErr)

		// Act
		err := defaultRM.AddRepository("test-repo", "example.org/test/stencils")

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, returnErr)
		}
	})
	t.Run("error is returned if repo already exists", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		var existingRepos = []Repository{
			{Name: "exists", URL: "https://example.org/stencil"},
		}
		_ = setupReposFileWithContent(t, m, existingRepos)
		fsw.UseCustomWrapper(m)

		// Act
		err := defaultRM.AddRepository("exists", "example.org/test/stencils")

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, fmt.Sprintf(ErrRepositoryAlreadyExists, "exists"))
		}
	})
	t.Run("error is returned if the file cannot be written", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		returnErr := errors.New("Bang!")
		file := setupReposFile(t, m)
		fsw.UseCustomWrapper(m)
		m.EXPECT().MkdirAll(filepath.Dir(file), defaultDirFileMode).Passthrough()
		m.EXPECT().WriteFile(file, mock.Anything, mock.Anything).
			RunAndReturn(func(name string, data []byte, perm os.FileMode) error {
				return returnErr
			})

		// Act
		err := defaultRM.AddRepository("exists", "example.org/test/stencils")

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, returnErr)
		}
	})
}

func Test_RemoveRepository(t *testing.T) {
	defaultRM := NewDefaultRepositoryManager()
	t.Run("Can successfully remove a repository", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		var existingRepos = []Repository{
			{Name: "exists", URL: "https://example.org/stencil"},
		}
		file := setupReposFileWithContent(t, m, existingRepos)
		fsw.UseCustomWrapper(m)

		m.EXPECT().MkdirAll(filepath.Dir(file), defaultDirFileMode).Passthrough()
		m.EXPECT().WriteFile(file, mock.Anything, defaultFileFileMode).
			RunAndReturn(func(name string, data []byte, perm os.FileMode) error {
				return os.WriteFile(name, data, perm)
			})

		// Act
		err := defaultRM.RemoveRepository("exists")

		// Assert
		m.AssertExpectations(t)
		assert.NoError(t, err)
		fileContents, readErr := os.ReadFile(file)
		assert.NoError(t, readErr)
		assert.JSONEq(t, "[]", string(fileContents))
	})
	t.Run("return error when repo doesn't exist", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		var existingRepos = []Repository{
			{Name: "exists", URL: "https://example.org/stencil"},
		}
		_ = setupReposFileWithContent(t, m, existingRepos)
		fsw.UseCustomWrapper(m)

		// Act
		err := defaultRM.RemoveRepository("does-not-exist")

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorContains(t, err, fmt.Sprintf(ErrRepositoryNotFound, "does-not-exist"))
		}
	})
	t.Run("return error when file cannot be read", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("Bang!")
		m.EXPECT().GetHomeDirectory().Return("", returnErr)

		// Act
		err := defaultRM.RemoveRepository("does-not-exist")

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, returnErr)
		}
	})
}

func Test_ListRepositories(t *testing.T) {
	defaultRM := NewDefaultRepositoryManager()
	t.Run("can read the file successfully", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		var existingRepos = []Repository{
			{Name: "exists", URL: "https://example.org/stencil"},
		}
		_ = setupReposFileWithContent(t, m, existingRepos)
		fsw.UseCustomWrapper(m)

		// Act
		repos, err := defaultRM.ListRepositories()

		// Assert
		m.AssertExpectations(t)
		assert.NoError(t, err)
		assert.NotNil(t, repos)
		assert.Equal(t, 1, len(repos))
		assert.Equal(t, existingRepos[0].Name, repos[0].Name)
		assert.Equal(t, existingRepos[0].URL, repos[0].URL)
	})
	t.Run("return error when fail to read", func(t *testing.T) {
		// Arrange
		m := mocks.NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("Bang!")
		m.EXPECT().GetHomeDirectory().Return("", returnErr)

		// Act
		repos, err := defaultRM.ListRepositories()

		// Assert
		m.AssertExpectations(t)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, returnErr)
		}
		assert.Nil(t, repos)
	})
}

func Test_Repository_Update(t *testing.T) {
	t.Run("Creates the repository directory and initializes a git repo on first update", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(mock.Anything, defaultDirFileMode).Passthrough()
		mockFS.EXPECT().Lstat(mock.Anything).Passthrough()
		fsw.UseCustomWrapper(mockFS)
		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		assert.NoError(t, err)
	})
	t.Run("Successfully updates the repo", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Passthrough()
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, nil)
		fsw.UseCustomWrapper(mockFS)

		cloneRepoToDirectory(t, "https://github.com/Chris-Greaves/stencil.git", repoPath)
		setRepoBackOneCommit(t, repoPath)
		startingHash := getRepoCommit(t, repoPath)

		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		assert.NoError(t, err)
		finalHash := getRepoCommit(t, repoPath)
		assert.NotEqual(t, startingHash, finalHash)
	})
	t.Run("Returns error when home directory cannot be found", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		expectedErr := errors.New("Bang!")
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, expectedErr)
		fsw.UseCustomWrapper(mockFS)

		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, expectedErr)
		}
	})
	t.Run("Returns error when repos directory cannot be created", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		expectedErr := errors.New("Bang!")
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(expectedErr)
		fsw.UseCustomWrapper(mockFS)

		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, expectedErr)
		}
	})
	t.Run("Returns error when repo directory cannot be accessed", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		expectedErr := errors.New("Bang!")
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, expectedErr)
		fsw.UseCustomWrapper(mockFS)

		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, expectedErr)
		}
	})
	t.Run("Returns error when repo directory cannot be accessed", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		rm := NewDefaultRepositoryManager()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		expectedErr := errors.New("Bang!")
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().MkdirAll(repoPath, defaultDirFileMode).Return(expectedErr)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, os.ErrNotExist)
		fsw.UseCustomWrapper(mockFS)

		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, expectedErr)
		}
	})
	t.Run("Returns an error when initialize fails", func(t *testing.T) {
		// Arrange
		newHome := t.TempDir()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		mockGitRunner := mocks.NewMockGitRunner(t)
		mockGitRunner.EXPECT().Initialize(repoPath).Return(assert.AnError)
		rm := NewRepositoryManager(mockGitRunner)
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().MkdirAll(repoPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, os.ErrNotExist)
		fsw.UseCustomWrapper(mockFS)
		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
	t.Run("Returns an error when AddRemote fails", func(t *testing.T) {
		// Arrange
		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}
		newHome := t.TempDir()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		mockGitRunner := mocks.NewMockGitRunner(t)
		mockGitRunner.EXPECT().Initialize(repoPath).Return(nil)
		mockGitRunner.EXPECT().AddRemote(repoPath, "origin", repo.URL).Return(assert.AnError)
		rm := NewRepositoryManager(mockGitRunner)
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().MkdirAll(repoPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, os.ErrNotExist)
		fsw.UseCustomWrapper(mockFS)

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
	t.Run("Returns an error when FetchRemote fails", func(t *testing.T) {
		// Arrange
		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}
		newHome := t.TempDir()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		mockGitRunner := mocks.NewMockGitRunner(t)
		mockGitRunner.EXPECT().FetchRemote(repoPath, "origin").Return(assert.AnError)
		rm := NewRepositoryManager(mockGitRunner)
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, nil)
		fsw.UseCustomWrapper(mockFS)

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
	t.Run("Returns an error when PullRemote fails", func(t *testing.T) {
		// Arrange
		repo := Repository{Name: "test-repo", URL: "https://github.com/Chris-Greaves/stencil.git"}
		newHome := t.TempDir()
		reposPath := filepath.Join(newHome, ".stencil", "repos")
		repoPath := filepath.Join(reposPath, "test-repo")
		mockGitRunner := mocks.NewMockGitRunner(t)
		mockGitRunner.EXPECT().FetchRemote(repoPath, "origin").Return(nil)
		mockGitRunner.EXPECT().PullRemote(repoPath, "origin", "main").Return(assert.AnError)
		rm := NewRepositoryManager(mockGitRunner)
		mockFS := mocks.NewMockFsWrapper(t)
		mockFS.EXPECT().GetHomeDirectory().Return(newHome, nil)
		mockFS.EXPECT().MkdirAll(reposPath, defaultDirFileMode).Return(nil)
		mockFS.EXPECT().Lstat(mock.Anything).Return(nil, nil)
		fsw.UseCustomWrapper(mockFS)

		// Act
		err := rm.runGitUpdate(repo)

		// Assert
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, assert.AnError)
		}
	})
}

func cloneRepoToDirectory(t *testing.T, url, path string) {
	cmd := exec.Command("git", "clone", url, path)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run command %v: %v", cmd.Args, err)
	}
}

func setRepoBackOneCommit(t *testing.T, repoPath string) {
	cmds := []*exec.Cmd{
		exec.Command("git", "reset", "--hard", "HEAD~1"),
		exec.Command("git", "branch", "--set-upstream-to=origin/main", "main"),
	}
	for _, cmd := range cmds {
		cmd.Dir = repoPath
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run command %v: %v", cmd.Args, err)
		}
	}
}

func getRepoCommit(t *testing.T, path string) string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = path
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run command %v: %v", cmd.Args, err)
	}
	return string(outBytes)
}

// Create a malformed repositories file
func createMalformedReposFile(t *testing.T, dir string) error {
	filePath := filepath.Join(dir, "repositories.json")
	t.Logf("Creating malformed file at %s", filePath)
	return os.WriteFile(filePath, []byte("invalid json content"), 0644)
}

// Create a valid repositories file with the given content
func createReposFile(t *testing.T, dir string, repository []Repository) error {
	// Create a valid repositories file
	validContent, err := json.Marshal(repository)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "repositories.json"), validContent, 0644)
}

func setupReposFileWithContent(t *testing.T, m *mocks.MockFsWrapper, repos []Repository) string {
	tempDir := t.TempDir()
	stencilDir := filepath.Join(tempDir, ".stencil")
	reposFile := filepath.Join(stencilDir, "repositories.json")
	m.EXPECT().GetHomeDirectory().Return(tempDir, nil)
	m.EXPECT().ReadFile(reposFile).RunAndReturn(func(name string) ([]byte, error) {
		return os.ReadFile(name)
	})

	content, err := json.Marshal(repos)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(stencilDir, defaultDirFileMode)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(reposFile, content, defaultFileFileMode)
	if err != nil {
		t.Fatal(err)
	}

	return reposFile
}

func setupReposFile(t *testing.T, m *mocks.MockFsWrapper) string {
	return setupReposFileWithContent(t, m, make([]Repository, 0))
}
