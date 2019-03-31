package tests

import (
	"os"
	"testing"

	"github.com/chris-greaves/stencil/fetch"
)

func TestIsPathReturnsTrueIfPathExists(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("An error was thrown '%v'", err.Error())
	}

	result := fetch.IsPath(wd)

	if result != true {
		t.Errorf("IsPath result was incorrect. Expected: 'true' Actual: '%v'", result)
	}
}

func TestIsPathReturnsFalseIfPathDoesntExists(t *testing.T) {
	result := fetch.IsPath("thisShouldNotExist12314")

	if result != false {
		t.Errorf("IsPath result was incorrect. Expected: 'false' Actual: '%v'", result)
	}
}
