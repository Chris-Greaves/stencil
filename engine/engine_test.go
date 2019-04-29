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

package engine

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validSettings = struct {
	ProjectName string
	Text        string
}{
	"Foobar",
	"Hello World",
}

const exampleFileContents = `Computer says: {{.Text}}`

var defaultEngine = New()

const validPathTemplate = "{{.ProjectName}}.txt"

func TestPathCanBeExecuted(t *testing.T) {
	executedPath, err := defaultEngine.ParseAndExecutePath("{{.ProjectName}}.txt", validSettings)
	require.NoError(t, err, "No error was expected from ParseAndExecutePath")

	assert.Equal(t, "Foobar.txt", executedPath, "Returned path was incorrect.")
}

func TestInvalidPathStringReturnsError(t *testing.T) {
	_, err := defaultEngine.ParseAndExecutePath("{{.Project-Name}}", validSettings)

	assert.Error(t, err, "Error was expected from ParseAndExecutePath")
	assert.Contains(t, err.Error(), "Error parsing path '{{.Project-Name}}' to template", "Incorrect error returned.")
}

func TestMissingPropertyInSettingsReturnsErrorWhenExecutinPath(t *testing.T) {
	_, err := defaultEngine.ParseAndExecutePath("{{.ProjectName}}-{{.NonExistantValue}}", validSettings)

	assert.Error(t, err, "Error was expected from ParseAndExecutePath")
	assert.Contains(t, err.Error(), "Error executing path as template. Path:'{{.ProjectName}}-{{.NonExistantValue}}'", "Incorrect error returned")
}

func TestFileCanBeExecutedCorrectly(t *testing.T) {
	testFilePath := CreateTestTemplateFile(t, exampleFileContents)
	defer os.RemoveAll(testFilePath)
	var b bytes.Buffer

	err := defaultEngine.ParseAndExecuteFile(testFilePath, validSettings, &b)
	require.NoError(t, err, "No error was expected")

	assert.Equal(t, "Computer says: Hello World", b.String())
}

func TestUnparseableTemplateErrors(t *testing.T) {
	testFilePath := CreateTestTemplateFile(t, "{{{{{{..}{\\.}}")
	defer os.RemoveAll(testFilePath)
	var b bytes.Buffer

	err := defaultEngine.ParseAndExecuteFile(testFilePath, validSettings, &b)
	require.Error(t, err, "An error was expected")
	assert.Contains(t, err.Error(), "Error Parsing template for file", "Incorrect error returned")
}

func TestMissingPropertyInSettingsReturnsErrorWhenExecutinFile(t *testing.T) {
	testFilePath := CreateTestTemplateFile(t, "{{ .DoesntExist }}")
	defer os.RemoveAll(testFilePath)
	var b bytes.Buffer

	err := defaultEngine.ParseAndExecuteFile(testFilePath, validSettings, &b)
	require.Error(t, err, "An error was expected")
	assert.Contains(t, err.Error(), "Error executing template file", "Incorrect error returned")
}

func CreateTestTemplateFile(t *testing.T, contents string) string {
	file, err := ioutil.TempFile("", "stencil-test-file-*.txt")
	require.NoError(t, err, "Unable to create temp file for test")

	file.WriteString(contents)

	return file.Name()
}
