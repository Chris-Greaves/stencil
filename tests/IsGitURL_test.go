package tests

import (
	"testing"

	"github.com/chris-greaves/stencil/fetch"
)

func TestIsGitURLReturnsTrueWhenRepoExists(t *testing.T) {
	result := fetch.IsGitURL("https://github.com/src-d/go-git")

	if result != true {
		t.Errorf("IsGitURL result was incorrect. Expected: 'true' Actual: '%v'", result)
	}
}

func TestIsGitURLReturnsFalseWhenRepoDoesntExists(t *testing.T) {
	result := fetch.IsGitURL("https://christophergreaves.co.uk")

	if result != false {
		t.Errorf("IsGitURL result was incorrect. Expected: 'false' Actual: '%v'", result)
	}
}
