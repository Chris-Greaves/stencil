package engine

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/chris-greaves/stencil/models"
)

func ParseAndExecutePath(settings models.Settings, path string) (string, error) {
	mainTemplate := template.New("main")

	tmpl, err := mainTemplate.Parse(path)
	if err != nil {
		return "", fmt.Errorf("Error parsing path to template: %v", err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, settings)
	if err != nil {
		return "", fmt.Errorf("Error applying template to path: %v", err)
	}

	return buf.String(), nil
}

func ParseAndExecuteFile(settings models.Settings, destinationPath string, sourcePath string, fileMode os.FileMode) error {
	fileTemplate, err := template.ParseFiles(sourcePath)
	if err != nil {
		return fmt.Errorf("Error Parsing template for file: %v", err)
	}

	destinationFile, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fileMode)
	if err != nil {
		return fmt.Errorf("Error creating destination file: %v", err)
	}
	defer destinationFile.Close()

	if err = fileTemplate.ExecuteTemplate(destinationFile, filepath.Base(sourcePath), settings); err != nil {
		return fmt.Errorf("Error executing template: %v", err)
	}

	return nil
}
