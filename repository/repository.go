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
	"path/filepath"
)

var (
	ErrRepositoryAlreadyExists = "repository with name '%s' already exists"
	ErrRepositoryNotFound      = "repository with name '%s' not found"
)

type Repository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func readReposFile() ([]Repository, error) {
	var repos []Repository
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fileContents, err := os.ReadFile(filepath.Join(homeDir, ".stencil", "repositories.yaml"))
	if err != nil {
		// If the file doesn't exist, we can assume there are no repositories yet and return an empty list
		if errors.Is(err, os.ErrNotExist) {
			return repos, nil
		}
		return nil, err
	}

	err = json.Unmarshal(fileContents, &repos)
	return repos, err
}

func saveReposFile(repos []Repository) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Ensure the .stencil directory exists
	err = os.MkdirAll(filepath.Join(homeDir, ".stencil"), 0755)
	if err != nil {
		return errors.Join(errors.New("failed to create $HOME/.stencil directory"), err)
	}

	fileContents, err := json.Marshal(repos)
	if err != nil {
		return errors.Join(errors.New("failed to marshal repositories"), err)
	}

	return os.WriteFile(filepath.Join(homeDir, ".stencil", "repositories.yaml"), fileContents, 0644)
}

func AddRepository(name string, url string) error {
	repos, err := readReposFile()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if repo.Name == name {
			return fmt.Errorf(ErrRepositoryAlreadyExists, name)
		}
	}

	repos = append(repos, Repository{Name: name, URL: url})
	return saveReposFile(repos)
}

func RemoveRepository(name string) error {
	repos, err := readReposFile()
	if err != nil {
		return err
	}

	for i, repo := range repos {
		if repo.Name == name {
			repos = append(repos[:i], repos[i+1:]...)
			return saveReposFile(repos)
		}
	}

	return fmt.Errorf(ErrRepositoryNotFound, name)
}

func ListRepositories() ([]Repository, error) {
	return readReposFile()
}

func UpdateRepositories() error {
	repos, err := readReposFile()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Printf("Updating repository '%s' from URL '%s'\n", repo.Name, repo.URL)

		err := repo.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateRepository(name string) error {
	repos, err := readReposFile()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if repo.Name == name {
			fmt.Printf("Updating repository '%s' from URL '%s'\n", repo.Name, repo.URL)

			err := repo.Update()
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf(ErrRepositoryNotFound, name)
}

func (r *Repository) Update() error {
	reposPath, err := getRepositoriesDirectory()
	if err != nil {
		return err
	}
	var repoPath = filepath.Join(reposPath, r.Name)

	// Ensure the $HOME/.stencil/repos directory exists
	err = os.MkdirAll(reposPath, 0755)
	if err != nil {
		return errors.Join(errors.New("failed to create $HOME/.stencil/repos directory"), err)
	}

	_, err = os.Lstat(repoPath)
	if err != nil && !os.IsNotExist(err) {
		println(err.Error())
	}

	if _, err := os.Lstat(repoPath); os.IsNotExist(err) {
		fmt.Printf("first time setting up %s", r.Name)
		err = os.MkdirAll(repoPath, 0755)
		if err != nil {
			return errors.Join(errors.New("failed to create $HOME/.stencil/repos/"+r.Name+" directory"), err)
		}

		// Run git init in the new repository directory
		if err := initializeGitRepo(repoPath); err != nil {
			return errors.Join(errors.New("failed to initialize git repository"), err)
		}

		// Run the git command to add the remote URL
		if err := addGitRemote(repoPath, r.URL); err != nil {
			return errors.Join(errors.New("failed to add git remote"), err)
		}
	}

	if err := fetchGitRemote(repoPath); err != nil {
		return errors.Join(errors.New("failed to fetch git remote"), err)
	}

	if err := pullGitRemote(repoPath); err != nil {
		return errors.Join(errors.New("failed to pull git remote"), err)
	}

	return nil
}

func getRepositoriesDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var reposPath = filepath.Join(homeDir, ".stencil", "repos")
	return reposPath, nil
}

func GetRepositoryPath(name string) (string, error) {
	reposPath, err := getRepositoriesDirectory()
	if err != nil {
		return "", err
	}

	var repoPath = filepath.Join(reposPath, name)
	return repoPath, nil
}

func RepositoryExists(name string) (bool, error) {
	repos, err := readReposFile()
	if err != nil {
		return false, err
	}

	for _, repo := range repos {
		if repo.Name == name {
			return true, nil
		}
	}

	return false, nil
}
