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
