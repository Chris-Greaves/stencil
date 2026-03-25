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

package action

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

var (
	ErrActionsFileNotFound = errors.New("actions file could not be found")
)

type ActionConfig struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
}

type staterFunc func(name string) (os.FileInfo, error)

func LoadActionsFromPath(path string) ([]Action, error) {
	return testableLoadActionsFromPath(os.Lstat, path)
}
func testableLoadActionsFromPath(stater staterFunc, path string) ([]Action, error) {
	//var actionConfig ActionConfig

	if _, err := stater(filepath.Join(path, ".actions.stencil.yaml")); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// File doesn't exist, continue to check for .actions.stencil.yml
	} else {
		return loadActionsFromYAML(filepath.Join(path, ".actions.stencil.yaml"))
	}

	if _, err := stater(filepath.Join(path, ".actions.stencil.yml")); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// File doesn't exist, continue to check for .actions.stencil.json
	} else {
		return loadActionsFromYAML(filepath.Join(path, ".actions.stencil.yml"))
	}

	if _, err := stater(filepath.Join(path, ".actions.stencil.json")); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// File doesn't exist, continue to return error
	} else {
		return loadActionsFromJSON(filepath.Join(path, ".actions.stencil.json"))
	}

	return nil, ErrActionsFileNotFound
}

func loadActionsFromJSON(actionsPath string) ([]Action, error) {
	var actionConfig ActionConfig

	data, err := os.ReadFile(actionsPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &actionConfig)
	return actionConfig.Actions, err
}

func loadActionsFromYAML(actionsPath string) ([]Action, error) {
	var actionConfig ActionConfig

	data, err := os.ReadFile(actionsPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &actionConfig)
	return actionConfig.Actions, err
}
