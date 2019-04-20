package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chris-Greaves/stencil/confighelper"
	"github.com/pkg/errors"
)

// Config is an interface for Stencil's config file
type Config interface {
	GetAllValues() ([]confighelper.Setting, error)
	SetValues(settings []confighelper.Setting) error
	Object() interface{}
}

// Engine is a interface to wrap the functions needed to create a templating engine that Stencil can understand
type Engine interface {
	ParseAndExecutePath(path string, settings interface{}) (string, error)
	ParseAndExecuteFile(sourcePath string, settings interface{}, wr io.Writer) error
}

// IOWrapper is a wrapper around the Input / Output for Stencil
type IOWrapper interface {
	GetOverrides(allSettings []confighelper.Setting) ([]confighelper.Setting, error)
}

// RootHandler is the Handler object for the Root cm
type RootHandler struct {
	confhelper     Config
	TemplateEngine Engine
	IO             IOWrapper
}

// NewRootHandler creates and returns a new RootHandler instance
func NewRootHandler(conf Config, templateEngine Engine, io IOWrapper) RootHandler {
	return RootHandler{confhelper: conf, TemplateEngine: templateEngine, IO: io}
}

// OfferConfigOverrides will take the current configuration and offer the user the ability to override the default values
func (h RootHandler) OfferConfigOverrides() error {
	editableSettings, err := h.confhelper.GetAllValues()
	if err != nil {
		return err
	}

	updatedSets, err := h.IO.GetOverrides(editableSettings)
	if err != nil {
		return err
	}

	h.confhelper.SetValues(updatedSets)

	return nil
}

// ProcessTemplate will walk through the Template and Parse it using the existing configuration
func (h RootHandler) ProcessTemplate(templatePath, outputPath string) error {
	return filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			// Skip if root or part of git
			if path == templatePath || strings.Contains(path, ".git") {
				return nil
			}

			fmt.Printf("Creating %v\n", path)

			tarPath, err := h.GetTargetPath(templatePath, outputPath, path, h.confhelper.Object())
			if err != nil {
				return err
			}

			if info.IsDir() {
				// If its a Directory, create the directory in the target
				if err = os.MkdirAll(tarPath, info.Mode()); err != nil {
					return errors.Wrapf(err, "Error making directory %v", path)
				}
			} else {
				// Open the file to write the contents into.
				destinationFile, err := os.OpenFile(tarPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
				if err != nil {
					return errors.Wrapf(err, "Error creating file at '%v'", tarPath)
				}
				defer destinationFile.Close()

				// If its a file, parse and execute the file and copy the result to the target
				if err = h.TemplateEngine.ParseAndExecuteFile(path, h.confhelper.Object(), destinationFile); err != nil {
					return errors.Wrapf(err, "Error processing file %v", path)
				}
			}

			return nil
		})
}

// GetTargetPath Converts a template path into the output path
func (h RootHandler) GetTargetPath(templatePath, outputPath, path string, settings interface{}) (string, error) {
	relPath, err := filepath.Rel(templatePath, path)
	if err != nil {
		return "", errors.Wrap(err, "Error getting relative path")
	}

	relTarPath, err := h.TemplateEngine.ParseAndExecutePath(relPath, settings)
	if err != nil {
		return "", err
	}
	tarPath := filepath.Join(outputPath, relTarPath)
	return tarPath, nil
}
