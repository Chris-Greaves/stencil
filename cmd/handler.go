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

func offerConfigOverrides(conf *confighelper.Conf) error {
	editableSettings, err := conf.GetAllValues()
	if err != nil {
		return err
	}

	var updatedSets []confighelper.Setting

	for _, setting := range editableSettings {
		if output := offerSettingToUser(setting); output != "" {
			updatedSets = append(updatedSets, confighelper.Setting{Name: setting.Name, Value: output})
		}
	}

	conf.SetValues(updatedSets)

	return nil
}

func offerSettingToUser(setting confighelper.Setting) string {
	fmt.Printf("Conf Override: \"%v\" [%v]: ", setting.Name, setting.Value)
	output := ""
	fmt.Scanln(&output)
	return output
}

func processTemplate(templatePath, outputPath string, conf *confighelper.Conf) {
	if err := filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			// Skip if root or part of git
			if path == templatePath || strings.Contains(path, ".git") {
				return nil
			}

			fmt.Printf("Creating %v\n", path)

			tarPath, err := getTargetPath(templatePath, outputPath, path, conf.Object())
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
				if err = engine.ParseAndExecuteFile(conf.Object(), tarPath, path, info.Mode()); err != nil {
					return errors.Wrapf(err, "Error processing file %v", path)
				}
			}

			return nil
		}); err != nil {
		log.Panicf("Error while creating project from template, %v", err)
	}
}

func getTargetPath(templatePath, outputPath, path string, settings interface{}) (string, error) {
	relPath, err := filepath.Rel(templatePath, path)
	if err != nil {
		return "", errors.Wrap(err, "Error getting relative path")
	}

	relTarPath, err := engine.ParseAndExecutePath(settings, relPath)
	if err != nil {
		return "", err
	}
	tarPath := filepath.Join(outputPath, relTarPath)
	return tarPath, nil
}
