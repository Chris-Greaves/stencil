package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/chris-greaves/stencil/confighelper"
)

const exampleFileContents = `{
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
}`

func TestNewConfCanBeCreated(t *testing.T) {
	file, err := ioutil.TempFile("", "fakefile-*.json")
	if err != nil {
		t.Fatalf("Unable to create temp file for test, Error: %v", err.Error())
	}

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	_, err = confighelper.New(file.Name())

	if err != nil {
		t.Errorf("Unexpected occured: %v", err)
		return
	}
}

func TestNewConfErrorsWhenFileIsntJson(t *testing.T) {
	file, err := ioutil.TempFile("", "fakefile-*.ext")
	if err != nil {
		t.Fatalf("Unable to create temp file for test, Error: %v", err.Error())
	}

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	_, result := confighelper.New(file.Name())

	if result == nil {
		t.Errorf("Expected error but found nil")
		return
	}

	if result.Error() != "Error ocurred validating settings path. Error: Extension must be '.json'" {
		t.Errorf("GetAllValuesFromFile error message was incorrect. Expected: 'Error ocurred validating settings path. Error: Extension must be '.json'', Actual: '%v'", result.Error())
	}
}

func TestNewConfErrorsWhenFileDoesntExist(t *testing.T) {
	_, result := confighelper.New("path/to/nowhere")

	if result == nil {
		t.Errorf("Expected error but found nil")
		return
	}

	if result.Error() != "Error ocurred validating settings path. Error: Path to file does not exist" {
		t.Errorf("GetAllValuesFromFile error message was incorrect. Expected: 'Error ocurred validating settings path. Error: Path to file does not exist', Actual: '%v'", result.Error())
	}
}

func TestYouCanGetAllValuesFromConf(t *testing.T) {
	conf := createNewConf(t)

	sets, err := conf.GetAllValues()

	if err != nil {
		t.Errorf("GetAllValuesFromFile should not have returned an error")
		return
	}

	if len(sets) < 1 {
		t.Errorf("GetAllValuesFromFile should have returned an array of settings")
	}

	settingMap := map[string]string{}
	for i := range sets {
		settingMap[sets[i].Name] = sets[i].Value
	}

	value, exists := settingMap["Project.Name"]
	if !exists {
		t.Error("Returned settings should contain a setting with Name 'Project.Name'")
	}
	if value != "DefaultProjectName" {
		t.Errorf("Expected 'Project.Name' to have Value: 'DefaultProjectName', Actual: %v", value)
	}

	value, exists = settingMap["Project.User.Name"]
	if !exists {
		t.Error("Returned settings should contain a setting with Name 'Project.User.Name'")
	}
	if value != "Chris" {
		t.Errorf("Expected 'Project.User.Name' to have Value: 'Chris', Actual: %v", value)
	}

	value, exists = settingMap["Directories.TextFiles"]
	if !exists {
		t.Error("Returned settings should contain a setting with Name 'Directories.TextFiles'")
	}
	if value != "TestFileDirectoryName" {
		t.Errorf("Expected 'Directories.TextFiles' to have Value: 'TestFileDirectoryName', Actual: %v", value)
	}

	value, exists = settingMap["Directories.YamlFiles"]
	if !exists {
		t.Error("Returned settings should contain a setting with Name 'Directories.YamlFiles'")
	}
	if value != "YamlFilesHere" {
		t.Errorf("Expected 'Directories.YamlFiles' to have Value: 'YamlFilesHere', Actual: %v", value)
	}

	value, exists = settingMap["Database.ConnectionString"]
	if !exists {
		t.Error("Returned settings should contain a setting with Name 'Database.ConnectionString'")
	}
	if value != "DefaultConnectionString" {
		t.Errorf("Expected 'Database.ConnectionString' to have Value: 'DefaultConnectionString', Actual: %v", value)
	}
}

func createNewConf(t *testing.T) *confighelper.Conf {
	file, err := ioutil.TempFile("", "fakefile-*.json")
	if err != nil {
		t.Fatalf("Unable to create temp file for test, Error: %v", err.Error())
	}

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	conf, err := confighelper.New(file.Name())
	if err != nil {
		t.Fatalf("Unable to create new conf for test: %v", err.Error())
	}

	return conf
}
