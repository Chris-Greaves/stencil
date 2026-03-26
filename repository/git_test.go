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
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsGitInstalled(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		pc   pathCheckerFunc
		want bool
	}{
		{
			name: "if error then return false",
			pc:   func(file string) (string, error) { return "", errors.New("bang!") },
			want: false,
		},
		{
			name: "if no error then return true",
			pc:   func(file string) (string, error) { return "something", nil },
			want: true,
		},
		{
			name: "running true path succeeds",
			pc:   nil,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			if tt.pc != nil {
				got = testable_IsGitInstalled(tt.pc)
			} else {
				got = IsGitInstalled()
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("IsGitInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initializeGitRepo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			// Currently, only the happy path is testable... If this method does more than just one command in the future, then more test cases can be added.
			name: "errors are passed back to caller",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir()
			gotErr := initializeGitRepo(path)

			if gotErr != nil {
				t.Errorf("initializeGitRepo() failed: %v", gotErr)
				return
			}

			if _, err := os.Lstat(filepath.Join(path, ".git")); err != nil {
				t.Errorf("initializeGitRepo() getting .git folder returned error: %v", err)
			}
		})
	}
}

func Test_addGitRemote(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		remote string
	}{
		{
			// Currently, only the happy path is testable... If this method does more than just one command in the future, then more test cases can be added.
			name:   "can add a remote to existing git repo",
			remote: "https://github.com/Chris-Greaves/stencil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir()
			createRepoErr := initializeGitRepo(path)
			if createRepoErr != nil {
				t.Errorf("addGitRemote() failed creating required repo to run test")
			}

			gotErr := addGitRemote(path, tt.remote)

			if gotErr != nil {
				t.Errorf("addGitRemote() failed: %v", gotErr)
				return
			}

			listRemotesCmd := exec.Command("git", "remote")
			listRemotesCmd.Dir = path
			output, err := listRemotesCmd.Output()
			if err != nil {
				t.Errorf("failed to assert that remote was added due to error: %v", err)
			}
			outStr := string(output)
			if strings.TrimSpace(outStr) != "origin" {
				t.Errorf("addGitRemote() = %v, want %v", outStr, "origin")
			}
		})
	}
}

func Test_fetchGitRemote(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			// Currently, only the happy path is testable... If this method does more than just one command in the future, then more test cases can be added.
			name: "can fetch from the remote",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir()
			createRepoErr := initializeGitRepo(path)
			if createRepoErr != nil {
				t.Errorf("fetchGitRemote() failed creating required repo to run test")
			}
			addRemoteErr := addGitRemote(path, "https://github.com/Chris-Greaves/stencil")
			if addRemoteErr != nil {
				t.Errorf("fetchGitRemote() failed adding remote to git repo to run test")
			}

			gotErr := fetchGitRemote(path)

			if gotErr != nil {
				t.Errorf("fetchGitRemote() failed: %v", gotErr)
				return
			}
		})
	}
}

func Test_pullGitRemote(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			// Currently, only the happy path is testable... If this method does more than just one command in the future, then more test cases can be added.
			name: "can fetch from the remote",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir()
			createRepoErr := initializeGitRepo(path)
			if createRepoErr != nil {
				t.Errorf("pullGitRemote() failed creating required repo to run test")
			}
			addRemoteErr := addGitRemote(path, "https://github.com/Chris-Greaves/stencil")
			if addRemoteErr != nil {
				t.Errorf("pullGitRemote() failed adding remote to git repo to run test")
			}
			fetchRemoteErr := fetchGitRemote(path)
			if fetchRemoteErr != nil {
				t.Errorf("pullGitRemote() failed fetching from remote")
			}

			gotErr := pullGitRemote(path)

			if gotErr != nil {
				t.Errorf("pullGitRemote() failed: %v", gotErr)
				return
			}
		})
	}
}
