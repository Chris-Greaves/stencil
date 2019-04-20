package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/Chris-Greaves/stencil/confighelper"
	"github.com/chris-greaves/stencil/cmd/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRootHandlerReturnsValidRootHandler(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	assert.NotNil(t, handler, "returned handler should not be nil")
	assert.IsType(t, RootHandler{}, handler)
}

func TestOfferConfigOverridesReturnsErrorIfGetAllValuesReturnsError(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	mockConfig.On("GetAllValues").Return(nil, errors.New("Something happened"))

	err := handler.OfferConfigOverrides()

	assert.Error(t, err)
	mockConfig.AssertExpectations(t)
}

func TestOfferConfigOverridesReturnsErrorFromGettingOverrides(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	mockConfig.On("GetAllValues").Return([]confighelper.Setting{
		confighelper.Setting{Value: "Something", Name: "Name1"},
		confighelper.Setting{Value: "Something", Name: "Name2"},
	}, nil)

	mockIO.On("GetOverrides", mock.Anything).Return(nil, errors.New("Bang"))

	err := handler.OfferConfigOverrides()

	assert.Error(t, err)
	mockConfig.AssertExpectations(t)
	mockIO.AssertExpectations(t)
}

func TestOfferConfigOverridesOnlyUpdatesEditedValues(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	mockConfig.On("GetAllValues").Return([]confighelper.Setting{
		confighelper.Setting{Value: "Something", Name: "Name1"},
		confighelper.Setting{Value: "Something", Name: "Name2"},
	}, nil)

	overrides := []confighelper.Setting{
		confighelper.Setting{Value: "SomethingElse", Name: "Name1"},
	}

	mockIO.On("GetOverrides", mock.Anything).Return(overrides, nil)

	mockConfig.On("SetValues", overrides).Return(nil)

	err := handler.OfferConfigOverrides()

	require.NoError(t, err)
	mockConfig.AssertExpectations(t)
	mockIO.AssertExpectations(t)
}

func TestGetTargetPathReturnsErrorIfRelCallFails(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	_, err := handler.GetTargetPath("../../path/../to/../../../somewhere", "", "../.../../..//..../to/somewhere", "")
	assert.Error(t, err)
	mockEngine.AssertExpectations(t)
}

func TestProcessTemplateIgnoresGitFolders(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	err = os.Mkdir(path.Join(templatePath, ".git"), os.ModeTemporary)
	require.NoError(t, err)

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, "")
	require.NoError(t, err)
	mockEngine.AssertNotCalled(t, "ParseAndExecutePath", mock.Anything, mock.Anything)
}

func TestProcessTemplateReturnsErrorsFromGetTargetPath(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	err = os.Mkdir(path.Join(templatePath, "{{ .Title }}"), os.ModeTemporary)
	require.NoError(t, err)

	mockConfig.On("Object").Return("")
	mockEngine.On("ParseAndExecutePath", mock.Anything, mock.Anything).Return("", errors.New("Bang!"))

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, "")
	assert.Error(t, err)
	mockEngine.AssertExpectations(t)
	mockEngine.AssertNotCalled(t, "ParseAndExecuteFile", mock.Anything, mock.Anything, mock.Anything)
}

func TestProcessTemplateReturnsErrorsFromParseAndExecuteFile(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	f, err := ioutil.TempFile(templatePath, "test-file-")
	require.NoError(t, err)
	f.Close()

	mockConfig.On("Object").Return("")
	mockEngine.On("ParseAndExecutePath", mock.Anything, mock.Anything).Return(f.Name(), nil)
	mockEngine.On("ParseAndExecuteFile", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Bang!"))

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, "")
	assert.Error(t, err)
	mockEngine.AssertExpectations(t)
}

func TestProcessTemplateReturnsErrorIfDirectoryCantBeMade(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	err = os.Mkdir(path.Join(templatePath, "something"), os.ModeTemporary)
	require.NoError(t, err)

	mockConfig.On("Object").Return("")
	mockEngine.On("ParseAndExecutePath", mock.Anything, mock.Anything).Return("*++`\"", nil)

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error making directory")
}

func TestProcessTemplateReturnsErrorWhenFailingToCreateFile(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	f, err := ioutil.TempFile(templatePath, "test-file-")
	require.NoError(t, err)
	f.Close()

	mockConfig.On("Object").Return("")
	mockEngine.On("ParseAndExecutePath", mock.Anything, mock.Anything).Return("*++`\"", nil)

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error creating file")
}

func TestProcessTemplateWorksCorrectly(t *testing.T) {
	mockEngine, mockConfig, mockIO := createMocks()
	outputPath, err := ioutil.TempDir("", "test-output-folder-")
	require.NoError(t, err)
	templatePath, err := ioutil.TempDir("", "test-template-")
	require.NoError(t, err)
	defer os.RemoveAll(templatePath)
	defer os.RemoveAll(outputPath)
	f, err := ioutil.TempFile(templatePath, "test-file-*.txt")
	require.NoError(t, err)
	f.Close()
	_, filename := filepath.Split(f.Name())

	mockConfig.On("Object").Return("")
	mockEngine.On("ParseAndExecutePath", mock.Anything, mock.Anything).Return(filename, nil)
	mockEngine.On("ParseAndExecuteFile", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	handler := NewRootHandler(mockConfig, mockEngine, mockIO)

	err = handler.ProcessTemplate(templatePath, outputPath)
	require.NoError(t, err)
	mockEngine.AssertExpectations(t)
	info, err := os.Stat(filepath.Join(outputPath, filename))
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
}

func createMocks() (*mocks.Engine, *mocks.Config, *mocks.IOWrapper) {
	return new(mocks.Engine), new(mocks.Config), new(mocks.IOWrapper)
}
