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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

var (
	ErrConfigFileExtensionInvalid = errors.New("config file extension is not supported, supported extensions are .yaml, .yml and .json")
)

type Config struct {
	Metadata ConfigMetadata `yaml:"metadata"`
	Vars     ConfigVars     `yaml:"vars"`
}

type ConfigMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

type ConfigVars struct {
	Prompt map[string]ConfigPrompt `yaml:"prompt"`
	Static map[string]string       `yaml:"static"`
}

type ConfigPrompt struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Default     string   `yaml:"default"`
	Options     []string `yaml:"options"`
	Required    bool     `yaml:"required"`
}

// loadConfig loads the stencil configuration from the specified path
func loadConfig(path string) (Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	fileExtension := strings.ToLower(filepath.Ext(path))
	switch fileExtension {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
		return config, err
	case ".json":
		err = json.Unmarshal(data, &config)
		return config, err
	}

	return config, ErrConfigFileExtensionInvalid
}
