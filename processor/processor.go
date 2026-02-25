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
	"path/filepath"

	"github.com/Chris-Greaves/stencil/utils"
	"github.com/charmbracelet/huh"
)

type Processor struct {
	path   string
	cfg    Config
	values map[string]interface{}
}

func NewProcessor(path string) (Processor, error) {
	var p = Processor{
		path:   path,
		values: make(map[string]interface{}),
	}

	err := p.parseStencilConfigFolder()
	if err != nil {
		return p, err
	}

	return p, nil
}

// Parse Stencil Config
func (p *Processor) parseStencilConfigFolder() error {
	var configDir = filepath.Join(p.path, ".stencil")
	if !utils.PathExistsAndIsDir(configDir) {
		return errors.New("stencil configuration folder not found at path: " + configDir)
	}

	cfgLoaded := false
	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if filepath.Base(path) == ".stencil.yaml" || filepath.Base(path) == ".stencil.yml" || filepath.Base(path) == ".stencil.json" {
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
	for key, prompt := range p.cfg.Vars.Prompt {

		switch prompt.Type {
		case "string":
			var value string
			if err := promptUserForString(key, prompt, &value); err != nil {
				return err
			}
			p.values[key] = value
		// case "select":
		// 	promptUserForSelect(prompt)
		default:
			return errors.New("unsupported prompt type: " + prompt.Type)
		}

	}
	// Ask user to provide values requested

	// Add static values

	// Put the values somewhere ready for template execution
	return nil
}

func (p *Processor) DumpValues() string {
	return fmt.Sprintf("%v", p.values)
}

func promptUserForString(key string, prompt ConfigPrompt, value *string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(key).
				Description(prompt.Description).
				Value(value),
		),
	)

	return form.Run()
}

// Execute Template
