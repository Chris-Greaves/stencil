package confighelper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
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
	require.NoError(t, err, "Unable to create temp file for test")

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	_, err = New(file.Name())

	assert.Nil(t, err, "No error was expected")
}

func TestNewConfErrorsWhenFileIsntJson(t *testing.T) {
	file, err := ioutil.TempFile("", "fakefile-*.ext")
	require.NoError(t, err, "Unable to create temp file for test")

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	_, result := New(file.Name())

	assert.NotNil(t, result, "Error was expected")
	assert.Equal(t, "Error ocurred validating settings path. Error: Extension must be '.json'", result.Error(), "Extension error should have been returned")
}

func TestNewConfErrorsWhenFileDoesntExist(t *testing.T) {
	_, result := New("path/to/nowhere")

	assert.NotNil(t, result, "Error was expected")
	assert.Equal(t, "Error ocurred validating settings path. Error: Path to file does not exist", result.Error(), "File path error should have been returned")
}

func TestYouCanGetAllValuesFromConf(t *testing.T) {
	assert := assert.New(t)

	conf := createNewConf(t)
	sets, err := conf.GetAllValues()

	require.NoError(t, err, "Unexpected error when Getting All Values")
	require.Greater(t, len(sets), 1, "GetAllValuesFromFile should have returned an array of settings")

	// mapping to make asserts more readable
	settingMap := map[string]string{}
	for i := range sets {
		settingMap[sets[i].Name] = sets[i].Value
	}

	value, exists := settingMap["Project.Name"]
	assert.True(exists, "Returned settings should contain a setting with Name 'Project.Name'")
	assert.Equal("DefaultProjectName", value, "'Project.Name' value was incorrect")

	value, exists = settingMap["Project.User.Name"]
	assert.True(exists, "Returned settings should contain a setting with Name 'Project.User.Name'")
	assert.Equal("Chris", value, "'Project.User.Name' value was incorrect")

	value, exists = settingMap["Directories.TextFiles"]
	assert.True(exists, "Returned settings should contain a setting with Name 'Directories.TextFiles'")
	assert.Equal("TestFileDirectoryName", value, "'Directories.TextFiles' value was incorrect")

	value, exists = settingMap["Directories.YamlFiles"]
	assert.True(exists, "Returned settings should contain a setting with Name 'Directories.YamlFiles'")
	assert.Equal("YamlFilesHere", value, "'Directories.YamlFiles' value was incorrect")

	value, exists = settingMap["Database.ConnectionString"]
	assert.True(exists, "Returned settings should contain a setting with Name 'Database.ConnectionString'")
	assert.Equal("DefaultConnectionString", value, "'Database.ConnectionString' value was incorrect")
}

func createNewConf(t *testing.T) *Conf {
	file, err := ioutil.TempFile("", "fakefile-*.json")
	require.NoError(t, err, "Unable to create temp file for test")

	defer os.RemoveAll(file.Name())

	file.WriteString(exampleFileContents)

	conf, err := New(file.Name())
	require.NoError(t, err, "Unable to create new conf for test")

	return conf
}
