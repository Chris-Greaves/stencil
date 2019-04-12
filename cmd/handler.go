package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chris-Greaves/stencil/confighelper"
	"github.com/Chris-Greaves/stencil/engine"
	"github.com/pkg/errors"
)

// RootHandler is the Handler object for the Root cm
type RootHandler struct {
	confhelper     confighelper.Config
	TemplateEngine engine.Engine
}

// NewRootHandler creates and returns a new RootHandler instance
func NewRootHandler(conf confighelper.Config, templateEngine engine.Engine) RootHandler {
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
		if output := h.offerSettingToUser(setting); output != "" {
			updatedSets = append(updatedSets, confighelper.Setting{Name: setting.Name, Value: output})
		}
	}

	h.confhelper.SetValues(updatedSets)

	return nil
}

func (h RootHandler) offerSettingToUser(setting confighelper.Setting) string {
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
				// If its a file, parse and execute the file and copy the result to the target
				if err = h.TemplateEngine.ParseAndExecuteFile(h.confhelper.Object(), tarPath, path, info.Mode()); err != nil {
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

	relTarPath, err := h.TemplateEngine.ParseAndExecutePath(settings, relPath)
	if err != nil {
		return "", err
	}
	tarPath := filepath.Join(outputPath, relTarPath)
	return tarPath, nil
}
