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

package utils

type GitClient struct {
	executor Executor
}

func NewGitClient(executor Executor) *GitClient {
	return &GitClient{executor: executor}
}

func DefaultGitClient() *GitClient {
	return &GitClient{executor: &OsExecClient{}}
}

func (g *GitClient) IsGitInstalled() bool {
	_, err := g.executor.CheckPathForExecutable("git")
	return err == nil
}

func (g *GitClient) RunCommand(dir string, args ...string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", args...)
	return err
}

func (g *GitClient) Initialize(dir string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "init")
	return err
}

func (g *GitClient) AddRemote(dir string, name string, url string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "remote", "add", name, url)
	return err
}

func (g *GitClient) Fetch(dir string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "fetch")
	return err
}

func (g *GitClient) FetchRemote(dir string, remote string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "fetch", remote)
	return err
}

func (g *GitClient) Pull(dir string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "pull")
	return err
}

func (g *GitClient) PullRemote(dir, remote, branch string) error {
	err := g.executor.ExecuteCommandInDir(dir, "git", "pull", remote, branch)
	return err
}
