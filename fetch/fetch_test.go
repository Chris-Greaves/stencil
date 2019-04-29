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
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestIsPathReturnsTrueIfPathExists(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err, "Error when getting Working Directory")

	result := IsPath(wd)

	assert.True(t, result, "Path should exist")
}

func TestIsPathReturnsFalseIfPathDoesntExists(t *testing.T) {
	result := IsPath("thisShouldNotExist12314")

	assert.False(t, result, "Path doesn't exist")
}

func TestIsGitURLReturnsTrueWhenRepoExists(t *testing.T) {
	result := IsGitURL("https://github.com/src-d/go-git")

	assert.True(t, result, "Git Url is Valid and exists")
}

func TestIsGitURLReturnsFalseWhenRepoDoesntExists(t *testing.T) {
	result := IsGitURL("https://christophergreaves.co.uk")

	assert.False(t, result, "Url isn't a valid git url")
}
