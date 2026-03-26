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
	"os/exec"
)

type pathCheckerFunc func(file string) (string, error)

func IsGitInstalled() bool {
	return testable_IsGitInstalled(exec.LookPath)
}

func testable_IsGitInstalled(pc pathCheckerFunc) bool {
	_, err := pc("git")
	return err == nil
}

// Initialize a new git repo inside a folder
func initializeGitRepo(path string) error {
	return runCommandInDir(path, "git", "init")
}

// Add the origin remote URL to the git repo
func addGitRemote(path string, remoteURL string) error {
	return runCommandInDir(path, "git", "remote", "add", "origin", remoteURL)
}

// Fetch the latest data from the remote URL
func fetchGitRemote(path string) error {
	return runCommandInDir(path, "git", "fetch", "origin")
}

// Pull down the latest changes from the remote URL
func pullGitRemote(path string) error {
	return runCommandInDir(path, "git", "pull", "origin", "main")
}

func runCommandInDir(path, name string, args ...string) error {
	execCmd := exec.Command(name, args...)
	execCmd.Dir = path
	return execCmd.Run()
}
