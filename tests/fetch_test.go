package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chris-greaves/stencil/fetch"
)

func TestIsPathReturnsTrueIfPathExists(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("An error was thrown '%v'", err.Error())
	}

	result := fetch.IsPath(wd)

	assert.True(t, result, "Path should exist")
}

func TestIsPathReturnsFalseIfPathDoesntExists(t *testing.T) {
	result := fetch.IsPath("thisShouldNotExist12314")

	assert.False(t, result, "Path doesn't exist")
}

func TestIsGitURLReturnsTrueWhenRepoExists(t *testing.T) {
	result := fetch.IsGitURL("https://github.com/src-d/go-git")

	assert.True(t, result, "Git Url is Valid and exists")
}

func TestIsGitURLReturnsFalseWhenRepoDoesntExists(t *testing.T) {
	result := fetch.IsGitURL("https://christophergreaves.co.uk")

	assert.False(t, result, "Url isn't a valid git url")
}
