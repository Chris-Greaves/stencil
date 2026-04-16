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
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_readReposFile(t *testing.T) {
	malformedReposFile := createMalformedReposFile(t)
	validReposFile := createReposFile(t, []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}})

	tests := []struct {
		name    string // description of this test case
		path    string
		want    []Repository
		wantErr bool
	}{
		{
			name:    "return an empty slice when the file does not exist",
			path:    "does/not/exist.yaml",
			want:    []Repository{},
			wantErr: false,
		},
		{
			name:    "return an error when the file is malformed",
			path:    malformedReposFile,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "return the expected repositories when the file is valid",
			path:    validReposFile,
			want:    []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := readReposFile(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("readReposFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("readReposFile() succeeded unexpectedly")
			}

			if !compareTwoSlices(got, tt.want) {
				t.Errorf("readReposFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveReposFile(t *testing.T) {
	tempDir := t.TempDir()
	repos := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
	stencilDir := filepath.Join(tempDir, ".stencil")

	if err := saveReposFile(repos, stencilDir); err != nil {
		t.Fatalf("saveReposFile() failed: %v", err)
	}

	got, err := readReposFile(stencilDir)
	if err != nil {
		t.Fatalf("readReposFile() failed: %v", err)
	}

	if !compareTwoSlices(got, repos) {
		t.Errorf("readReposFile() = %v, want %v", got, repos)
	}
}

func Test_repositoryLifecycle(t *testing.T) {
	homeDir := setupTemporaryHome(t)

	existingRepos := []Repository{{Name: "repo1", URL: "https://example.com/repo1.git"}}
	if err := saveReposFile(existingRepos, filepath.Join(homeDir, ".stencil")); err != nil {
		t.Fatalf("saveReposFile() failed: %v", err)
	}

	t.Run("add duplicate repository returns error", func(t *testing.T) {
		if err := AddRepository("repo1", "https://example.com/repo1.git"); err == nil {
			t.Fatal("AddRepository() succeeded unexpectedly")
		} else if err.Error() != fmt.Sprintf(ErrRepositoryAlreadyExists, "repo1") {
			t.Fatalf("AddRepository() error = %v, want %v", err, fmt.Sprintf(ErrRepositoryAlreadyExists, "repo1"))
		}
	})

	t.Run("add repository and list repositories", func(t *testing.T) {
		if err := AddRepository("repo2", "https://example.com/repo2.git"); err != nil {
			t.Fatalf("AddRepository() failed: %v", err)
		}

		got, err := ListRepositories()
		if err != nil {
			t.Fatalf("ListRepositories() failed: %v", err)
		}

		want := []Repository{
			{Name: "repo1", URL: "https://example.com/repo1.git"},
			{Name: "repo2", URL: "https://example.com/repo2.git"},
		}

		if !compareTwoSlices(got, want) {
			t.Errorf("ListRepositories() = %v, want %v", got, want)
		}
	})

	t.Run("repository exists checks", func(t *testing.T) {
		exists, err := RepositoryExists("repo1")
		if err != nil {
			t.Fatalf("RepositoryExists() failed: %v", err)
		}
		if !exists {
			t.Fatal("RepositoryExists() = false, want true")
		}

		exists, err = RepositoryExists("repo-missing")
		if err != nil {
			t.Fatalf("RepositoryExists() failed: %v", err)
		}
		if exists {
			t.Fatal("RepositoryExists() = true, want false")
		}
	})

	t.Run("remove repository and preserve remaining entries", func(t *testing.T) {
		if err := RemoveRepository("repo1"); err != nil {
			t.Fatalf("RemoveRepository() failed: %v", err)
		}

		got, err := ListRepositories()
		if err != nil {
			t.Fatalf("ListRepositories() failed: %v", err)
		}
		if len(got) != 1 || got[0].Name != "repo2" {
			t.Fatalf("ListRepositories() = %v, want remaining repo2", got)
		}

		removeErr := RemoveRepository("repo-missing")
		if removeErr == nil {
			t.Fatal("RemoveRepository() succeeded unexpectedly")
		} else if removeErr.Error() != fmt.Sprintf(ErrRepositoryNotFound, "repo-missing") {
			t.Fatalf("RemoveRepository() error = %v, want %v", removeErr, fmt.Sprintf(ErrRepositoryNotFound, "repo-missing"))
		}
	})
}

func Test_GetRepositoryPath(t *testing.T) {
	homeDir := setupTemporaryHome(t)

	got, err := GetRepositoryPath("repo1")
	if err != nil {
		t.Fatalf("GetRepositoryPath() failed: %v", err)
	}

	want := filepath.Join(homeDir, ".stencil", "repos", "repo1")
	if got != want {
		t.Fatalf("GetRepositoryPath() = %v, want %v", got, want)
	}
}

func setupTemporaryHome(t *testing.T) string {
	t.Helper()

	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Setenv("USERPROFILE", homeDir)

	return homeDir
}

func createMalformedReposFile(t *testing.T) string {
	tempDir := t.TempDir()
	// Create a malformed repositories file
	malformedContent := []byte("invalid yaml content")
	err := os.WriteFile(filepath.Join(tempDir, "repositories.yaml"), malformedContent, 0644)
	if err != nil {
		t.Fatal(err)
	}

	return tempDir
}

func createReposFile(t *testing.T, repository []Repository) string {
	tempDir := t.TempDir()
	// Create a valid repositories file
	validContent, err := json.Marshal(repository)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(tempDir, "repositories.yaml"), validContent, 0644)
	if err != nil {
		t.Fatal(err)
	}

	return tempDir
}

func compareTwoSlices(got, repository []Repository) bool {
	if len(got) != len(repository) {
		return false
	}

	for i := range got {
		if got[i] != repository[i] {
			return false
		}
	}

	return true
}
