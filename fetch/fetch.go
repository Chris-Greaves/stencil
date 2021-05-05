// Copyright © 2018 Christopher Greaves <cjgreaves97@hotmail.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fetch

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

// IsPath is used to determin if the input string is a valid path on the local system
func IsPath(input string) bool {
	info, err := os.Stat(input)

	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

// IsGitURL is used to determin if the input string is a valid Git Url
func IsGitURL(input string) bool {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{input},
	})

	_, err := rem.List(&git.ListOptions{}) // List tags to prove that git repo exists
	return err == nil
}

// PullTemplate clones the template from its git repo
func PullTemplate(repo string) (string, error) {
	dir, err := ioutil.TempDir("", "template-")
	if err != nil {
		return "", err
	}

	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               repo,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	ref, err := r.Head()
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	log.Printf("Git repo cloned at hash %v", ref.Hash().String())

	return dir, nil
}
