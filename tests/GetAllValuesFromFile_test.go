package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chris-greaves/stencil/settings"
)

func TestGetAllValuesFromFileHappyPath(t *testing.T) {
	file, err := ioutil.TempFile("", "fakefile-*.json")
	if err != nil {
		t.Fatalf("Unable to create temp file for test, Error: %v", err.Error())
	}

	defer os.RemoveAll(file.Name())

	file.WriteString(`{
		"Project": {
			"Name": "DefaultProjectName",
			"User": {
				"Name": "Chris"
			}
		},
		"Directories": {
			"TextFiles": "TestFileDirectoryName",
			"YamlFiles": "YamlFilesHere"
		},
		"Database": {
			"ConnectionString": "DefaultConnectionString"
		}
	}`)

	sets, err := settings.GetAllValuesFromFile(file.Name())

	if err != nil {
		t.Errorf("GetAllValuesFromFile should not have returned an error")
		return
	}

	if len(sets) < 1 {
		t.Errorf("GetAllValuesFromFile should have returned an array of settings")
	}

	if sets[0].Name != "Project.Name" || sets[0].Value != "DefaultProjectName" {
		t.Errorf("Expected sets[0] to have Name: 'Project.Name', Actual: %v", sets[0].Name)
		t.Errorf("Expected sets[0] to have Value: 'DefaultProjectName', Actual: %v", sets[0].Value)
	}
	if sets[1].Name != "Project.User.Name" || sets[1].Value != "Chris" {
		t.Errorf("Expected sets[1] to have Name: 'Project.User.Name', Actual: %v", sets[1].Name)
		t.Errorf("Expected sets[1] to have Value: 'Chris', Actual: %v", sets[1].Value)
	}
	if sets[2].Name != "Directories.TextFiles" || sets[2].Value != "TestFileDirectoryName" {
		t.Errorf("Expected sets[2] to have Name: 'Directories.TextFiles', Actual: %v", sets[2].Name)
		t.Errorf("Expected sets[2] to have Value: 'TestFileDirectoryName', Actual: %v", sets[2].Value)
	}
	if sets[3].Name != "Directories.YamlFiles" || sets[3].Value != "YamlFilesHere" {
		t.Errorf("Expected sets[3] to have Name: 'Directories.YamlFiles', Actual: %v", sets[3].Name)
		t.Errorf("Expected sets[3] to have Value: 'YamlFilesHere', Actual: %v", sets[3].Value)
	}
	if sets[4].Name != "Database.ConnectionString" || sets[4].Value != "DefaultConnectionString" {
		t.Errorf("Expected sets[4] to have Name: 'Database.ConnectionString', Actual: %v", sets[4].Name)
		t.Errorf("Expected sets[4] to have Value: 'DefaultConnectionString', Actual: %v", sets[4].Value)
	}
}

func TestGetAllValuesFromFileReturnsErrorWhenPathDoesntExist(t *testing.T) {
	_, result := settings.GetAllValuesFromFile("path/does/not/exist")

	if result == nil {
		t.Errorf("GetAllValuesFromFile should return an error")
		return
	}

	if result.Error() != "Error ocurred validating settings path. Error: Path to file does not exist" {
		t.Errorf("GetAllValuesFromFile error message was incorrect. Expected: 'Error ocurred validating settings path. Error: Path to file does not exist', Actual: '%v'", result.Error())
	}
}

func TestGetAllValuesFromFileReturnsErrorWhenExtensionIsntJson(t *testing.T) {
	file, err := ioutil.TempFile("", "fakefile-*.ext")
	if err != nil {
		t.Fatalf("Unable to create temp file for test, Error: %v", err.Error())
	}

	defer os.RemoveAll(file.Name())

	_, result := settings.GetAllValuesFromFile(file.Name())

	if result == nil {
		t.Errorf("GetAllValuesFromFile should return an error")
		return
	}

	if result.Error() != "Error ocurred validating settings path. Error: Extension must be '.json'" {
		t.Errorf("GetAllValuesFromFile error message was incorrect. Expected: 'Error ocurred validating settings path. Error: Extension must be '.json'', Actual: '%v'", result.Error())
	}
}
