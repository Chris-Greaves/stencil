package engine

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

// DefaultEngine is a default implementation of the Template Engine needed for Stencil
type DefaultEngine struct {
}

// New Creates a new instance of the Default Engine
func New() DefaultEngine {
	return DefaultEngine{}
}

// ParseAndExecutePath will parse the path as a template and execute it using the settings provided
func (e DefaultEngine) ParseAndExecutePath(path string, settings interface{}) (string, error) {
	mainTemplate := template.New("main")

	tmpl, err := mainTemplate.Parse(path)
	if err != nil {
		return "", errors.Wrapf(err, "Error parsing path '%v' to template", path)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, settings)
	if err != nil {
		return "", errors.Wrapf(err, "Error executing path as template. Path:'%v'", path)
	}

	return buf.String(), nil
}

// ParseAndExecuteFile will parse a file as a template and execute it using the settings provided. it will write out to the destinationPath using the FileMode supplied.
func (e DefaultEngine) ParseAndExecuteFile(sourcePath string, destinationPath string, settings interface{}, fileMode os.FileMode) error {
	fileTemplate, err := template.ParseFiles(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "Error Parsing template for file '%v'", sourcePath)
	}

	destinationFile, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fileMode)
	if err != nil {
		return errors.Wrapf(err, "Error creating file at '%v'", destinationPath)
	}
	defer destinationFile.Close()

	if err = fileTemplate.ExecuteTemplate(destinationFile, filepath.Base(sourcePath), settings); err != nil {
		return errors.Wrapf(err, "Error executing template file '%v'", sourcePath)
	}

	return nil
}
