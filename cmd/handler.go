package cmd

import (
	"fmt"
	"io"
	"log"
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

// RootHandler is the Handler object for the Root cm
type RootHandler struct {
	confhelper     Config
	TemplateEngine Engine
}

// NewRootHandler creates and returns a new RootHandler instance
func NewRootHandler(conf Config, templateEngine Engine) RootHandler {
	return RootHandler{confhelper: conf, TemplateEngine: templateEngine}
}

// OfferConfigOverrides will take the current configuration and offer the user the ability to override the default values
func (h RootHandler) OfferConfigOverrides() error {
	editableSettings, err := h.confhelper.GetAllValues()
	if err != nil {
		return err
	}

	var updatedSets []confighelper.Setting

	for _, setting := range editableSettings {
		if output := offerSettingToUser(setting); output != "" {
			updatedSets = append(updatedSets, confighelper.Setting{Name: setting.Name, Value: output})
		}
	}

	h.confhelper.SetValues(updatedSets)

	return nil
}

func offerSettingToUser(setting confighelper.Setting) string {
	fmt.Printf("Conf Override: \"%v\" [%v]: ", setting.Name, setting.Value)
	output := ""
	fmt.Scanln(&output)
	return output
}

// ProcessTemplate will walk through the Template and Parse it using the existing configuration
func (h RootHandler) ProcessTemplate(templatePath, outputPath string) {
	if err := filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			// Skip if root or part of git
			if path == templatePath || strings.Contains(path, ".git") {
				return nil
			}

			fmt.Printf("Creating %v\n", path)

			tarPath, err := h.getTargetPath(templatePath, outputPath, path, h.confhelper.Object())
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
		}); err != nil {
		log.Panicf("Error while creating project from template, %v", err)
	}
}

func (h RootHandler) getTargetPath(templatePath, outputPath, path string, settings interface{}) (string, error) {
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
