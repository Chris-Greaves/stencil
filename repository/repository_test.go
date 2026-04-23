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
	"os"
	"path/filepath"
	"testing"

	"github.com/Chris-Greaves/stencil/utils/fsw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// type MockedFsWrapper struct {
// 	mock.Mock
// }

// func (m MockedFsWrapper) GetHomeDirectory() (string, error) {
// 	args := m.Called()
// 	return args.String(0), args.Error(1)
// }

// func (m MockedFsWrapper) ReadFile(name string) ([]byte, error) {
// 	args := m.Called(name)
// 	return args.Get(0).([]byte), args.Error(1)
// }

// func (m MockedFsWrapper) MkdirAll(path string, perm os.FileMode) error {
// 	args := m.Called(path, perm)
// 	return args.Error(0)
// }

// func (m MockedFsWrapper) WriteFile(name string, data []byte, perm os.FileMode) error {
// 	args := m.Called(name, data, perm)
// 	return args.Error(0)
// }

// func (m MockedFsWrapper) Lstat(name string) (os.FileInfo, error) {
// 	args := m.Called(name)
// 	return args.Get(0).(os.FileInfo), args.Error(1)
// }

func Test_readReposFile(t *testing.T) {
	t.Run("return an empty slice when the file does not exist", func(t *testing.T) {
		// Arrange
		m := NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)
		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, 0755), "failed to create .stencil directory in temporary home directory")
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().ReadFile(reposFilePath).Return(nil, os.ErrNotExist)

		// Act
		got, gotErr := readReposFile()

		// Assert
		if assert.NoError(t, gotErr, "error not expected") {
			assert.Empty(t, got, "expected an empty slice when the file does not exist")
		}
	})

	t.Run("return an error when the file is malformed", func(t *testing.T) {
		// Arrange
		m := NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, 0755), "failed to create .stencil directory in temporary home directory")

		if err := createMalformedReposFile(t, stencilDir); err != nil {
			t.Fatalf("failed to create malformed repos file: %v", err)
		}
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil).Once()
		m.EXPECT().ReadFile(reposFilePath).Return(os.ReadFile(reposFilePath)).Once()

		// Act
		got, gotErr := readReposFile()

		// Assert
		m.AssertExpectations(t)
		assert.Error(t, gotErr, "error was expected, but got nil. Returned repositories: %v", got)
		assert.Nil(t, got, "expected nil when the file is malformed")
	})

	t.Run("return the expected repositories when the file is valid", func(t *testing.T) {
		// Arrange
		m := NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		customHomeDir := t.TempDir()
		stencilDir := filepath.Join(customHomeDir, ".stencil")
		reposFilePath := filepath.Join(stencilDir, "repositories.json")
		require.NoError(t, os.MkdirAll(stencilDir, 0755), "failed to create .stencil directory in temporary home directory")

		validReposContent := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
		if err := createReposFile(t, stencilDir, validReposContent); err != nil {
			t.Fatalf("failed to create valid repos file: %v", err)
		}
		m.EXPECT().GetHomeDirectory().Return(customHomeDir, nil)
		m.EXPECT().ReadFile(reposFilePath).Return(os.ReadFile(reposFilePath))

		// Act
		got, gotErr := readReposFile()

		// Assert
		assert.NoError(t, gotErr, "error not expected")
		assert.Equal(t, validReposContent, got)
	})

	t.Run("return the error when home directory cannot be determined", func(t *testing.T) {
		// Arrange
		m := NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("failed to determine home directory")
		m.EXPECT().GetHomeDirectory().Return("", returnErr)

		// Act
		got, gotErr := readReposFile()

		// Assert
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, returnErr)
		}
		assert.Nil(t, got, "expected nil when home directory cannot be determined")
	})

	t.Run("return the error when file cannot be read", func(t *testing.T) {
		// Arrange
		m := NewMockFsWrapper(t)
		fsw.UseCustomWrapper(m)

		returnErr := errors.New("failed to read file")
		m.EXPECT().GetHomeDirectory().Return("/some/path/", nil)
		m.EXPECT().ReadFile(filepath.Join("/some/path/", ".stencil", "repositories.json")).Return(nil, returnErr)

		// Act
		got, gotErr := readReposFile()

		// Assert
		if assert.Error(t, gotErr) {
			assert.ErrorIs(t, gotErr, returnErr)
		}
		assert.Nil(t, got, "expected nil when file cannot be read, but got: %v", got)
	})
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
