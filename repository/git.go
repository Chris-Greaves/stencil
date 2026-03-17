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

func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Initialize a new git repo inside a folder
func initializeGitRepo(path string) error {
	execCmd := exec.Command("git", "init")
	execCmd.Dir = path
	err := execCmd.Run()
	return err
}

// Add the origin remote URL to the git repo
func addGitRemote(path string, remoteURL string) error {
	execCmd := exec.Command("git", "remote", "add", "origin", remoteURL)
	execCmd.Dir = path
	err := execCmd.Run()
	return err
}

// Fetch the latest data from the remote URL
func fetchGitRemote(path string) error {
	execCmd := exec.Command("git", "fetch", "origin")
	execCmd.Dir = path
	err := execCmd.Run()
	return err
}

// Pull down the latest changes from the remote URL
func pullGitRemote(path string) error {
	execCmd := exec.Command("git", "pull", "origin", "main")
	execCmd.Dir = path
	err := execCmd.Run()
	return err
}
