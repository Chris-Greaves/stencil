// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package processor

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chris-Greaves/stencil/utils"
	"github.com/charmbracelet/huh"
)

type Processor struct {
	templatePath string
	outputPath   string
	cfg          Config
	values       map[string]interface{}
}

func NewProcessor(templatePath string, outputPath string) (Processor, error) {
	var p = Processor{
		templatePath: templatePath,
		outputPath:   outputPath,
		values:       make(map[string]interface{}),
	}

	err := p.parseStencilConfigFolder()
	if err != nil {
		return p, err
	}

	return p, nil
}

// Parse Stencil Config
func (p *Processor) parseStencilConfigFolder() error {
	var configDir = filepath.Join(p.templatePath, ".stencil")
	if !utils.PathExistsAndIsDir(configDir) {
		return errors.New("stencil configuration folder not found at path: " + configDir)
	}

	cfgLoaded := false
	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if filepath.Base(path) == "config.yaml" || filepath.Base(path) == "config.yml" || filepath.Base(path) == "config.json" {
				if cfgLoaded {
					return errors.New("multiple stencil configuration files found in: " + configDir)
				}
				config, err := loadConfig(path)
				if err != nil {
					return err
				}
				p.cfg = config
				cfgLoaded = true
				return nil
			}
		}

		// TODO: Add logic for scripts

		return nil
	})
	if err != nil {
		return err
	}

	if !cfgLoaded {
		return errors.New("no stencil config files found in: " + configDir)
	}

	return nil
}

func (p *Processor) DumpConfig() string {
	return fmt.Sprintf("%v", p.cfg)
}

// Prompt user for input
func (p *Processor) PromptUserForInput() error {
	var fields []huh.Field
	fields, err := getFieldsForPrompt(p, fields)
	if err != nil {
		return err
	}
	// Ask user to provide values requested
	form := huh.NewForm(
		huh.NewGroup(
			fields...,
		),
	)
	err = form.Run()
	if err != nil {
		return err
	}

	// Collect the values provided by the user
	for _, field := range fields {
		p.values[field.GetKey()] = field.GetValue()
	}

	// Add static values
	for key, value := range p.cfg.Vars.Static {
		parsedValue, err := parseTemplateString(value, p.values)
		if err != nil {
			return errors.Join(fmt.Errorf("error parsing static value for key: %s", key), err)
		}
		p.values[key] = parsedValue
	}

	return nil
}

func (p *Processor) DumpValues() string {
	return fmt.Sprintf("%v", p.values)
}

// Execute Template
func (p *Processor) ExecuteTemplate() error {
	return filepath.WalkDir(p.templatePath, func(path string, d fs.DirEntry, err error) (retErr error) {
		if path == p.templatePath {
			// Just ignore and return nil, as there is nothing to do and we want to continue walking
			return nil
		}
		if strings.Contains(path, ".stencil") {
			// Skip the .stencil config folder and all its contents
			return filepath.SkipDir
		}
		if err != nil {
			return errors.Join(errors.New("error while walking into the directory"), err)
		}

		targetPath := getTargetPath(p.outputPath, p.templatePath, path)

		parsedTargetPath, err := parseTemplateString(targetPath, p.values)
		if err != nil {
			return errors.Join(fmt.Errorf("error while parsing template string for path %v", targetPath), err)
		}

		if d.IsDir() {
			err = os.MkdirAll(parsedTargetPath, os.ModePerm)
			if err != nil {
				return errors.Join(fmt.Errorf("error creating directory at %v", parsedTargetPath), err)
			}
		} else {
			// Open the file to write the contents into.
			destinationFile, err := os.OpenFile(parsedTargetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return errors.Join(fmt.Errorf("error creating file at %v", parsedTargetPath), err)
			}
			defer func() { // Putting the close in a deferred function to appease the linter
				retErr = errors.Join(destinationFile.Close())
			}()

			err = parseTemplateFile(path, p.values, destinationFile)
			if err != nil {
				return errors.Join(fmt.Errorf("error while parsing template file for path %v", path), err)
			}
		}

		fmt.Printf("%s -> \t%s\n", path, parsedTargetPath)
		return nil
	})
}
