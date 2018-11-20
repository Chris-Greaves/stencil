package tests

import (
	"strings"
	"testing"

	"github.com/chris-greaves/stencil/engine"
)

var validSettings = struct {
	ProjectName string
	Text        string
}{
	"Foobar",
	"Hello World",
}

const validPathTemplate = "{{.ProjectName}}.txt"

func TestPathCanBeExecuted(t *testing.T) {
	executedPath, err := engine.ParseAndExecutePath(validSettings, "{{.ProjectName}}.txt")
	if err != nil {
		t.Errorf("An error was thrown '%v'", err.Error())
	}
	if executedPath != "Foobar.txt" {
		t.Errorf("returned path was incorrect. Expected: \"Foobar.txt\" Actual: \"%v\"", executedPath)
	}
}

func TestInvalidPathStringReturnsError(t *testing.T) {
	_, err := engine.ParseAndExecutePath(validSettings, "{{.Project-Name}}")
	if err == nil {
		t.Error("No error was return")
	} else {
		if !strings.Contains(err.Error(), "Error parsing path '{{.Project-Name}}' to template") {
			t.Errorf("Incorrect error returned. Expected error to contain: \"Error parsing path '{{.Project-Name}}' to template\" have: \"%v\"", err.Error())
		}
	}
}

func TestMissingPropertyInSettingsReturnsError(t *testing.T) {
	_, err := engine.ParseAndExecutePath(validSettings, "{{.ProjectName}}-{{.NonExistantValue}}")
	if err == nil {
		t.Error("No error was return")
	} else {
		if !strings.Contains(err.Error(), "Error executing path as template. Path:'{{.ProjectName}}-{{.NonExistantValue}}'") {
			t.Errorf("Incorrect error returned. Expected error to contain: \"Error executing path as template. Path:'{{.ProjectName}}-{{.NonExistantValue}}'\" have: \"%v\"", err.Error())
		}
	}
}
