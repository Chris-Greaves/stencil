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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

var actionConfig = ActionConfig{
	Actions: []Action{
		Action{
			Name:        "new-project",
			Description: "Some description",
			Path:        "some/path",
		},
	},
}

func TestLoadActionsFromPath(t *testing.T) {
	var ErrBang = errors.New("bang!")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path          string
		stater        staterFunc
		want          []Action
		wantErr       bool
		validateErrAs error
	}{
		{
			name:          "path does not exist",
			path:          "path/does/not/exist",
			stater:        nil,
			want:          nil,
			wantErr:       true,
			validateErrAs: ErrActionsFileNotFound,
		},
		{
			name:          "Lstat throws unexpected error on first call",
			path:          "some/path",
			stater:        func(name string) (os.FileInfo, error) { return nil, ErrBang },
			want:          nil,
			wantErr:       true,
			validateErrAs: ErrBang,
		},
		{
			name: "Lstat throws unexpected error on second call",
			path: "some/path",
			stater: func(name string) (os.FileInfo, error) {
				if strings.Contains(name, ".yml") {
					return nil, ErrBang
				} else {
					return nil, os.ErrNotExist
				}
			},
			want:          nil,
			wantErr:       true,
			validateErrAs: ErrBang,
		},
		{
			name: "Lstat throws unexpected error on third call",
			path: "some/path",
			stater: func(name string) (os.FileInfo, error) {
				if strings.Contains(name, ".json") {
					return nil, ErrBang
				} else {
					return nil, os.ErrNotExist
				}
			},
			want:          nil,
			wantErr:       true,
			validateErrAs: ErrBang,
		},
		{
			name:          ".yaml files are processed correctly",
			path:          createActionsFileForTest(t, "yaml"),
			stater:        nil,
			want:          actionConfig.Actions,
			wantErr:       false,
			validateErrAs: nil,
		},
		{
			name:          ".yml files are processed correctly",
			path:          createActionsFileForTest(t, "yml"),
			stater:        nil,
			want:          actionConfig.Actions,
			wantErr:       false,
			validateErrAs: nil,
		},
		{
			name:          ".json files are processed correctly",
			path:          createActionsFileForTest(t, "json"),
			stater:        nil,
			want:          actionConfig.Actions,
			wantErr:       false,
			validateErrAs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []Action
			var gotErr error

			if tt.stater != nil {
				got, gotErr = testableLoadActionsFromPath(tt.stater, tt.path)
			} else {
				got, gotErr = LoadActionsFromPath(tt.path)
			}

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadActionsFromPath() failed: %v", gotErr)
				}
				if tt.validateErrAs != nil {
					if !errors.Is(gotErr, tt.validateErrAs) {
						t.Errorf("LoadActionsFromPath() return incorrect error: %v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadActionsFromPath() succeeded unexpectedly")
			}

			if stringifyObject(got) != stringifyObject(tt.want) {
				t.Errorf("LoadActionsFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadActionsFromJson(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path          string
		want          []Action
		wantErr       bool
		validateErrAs error
	}{
		{
			name:          "path does not exist",
			path:          "path/does/not/exist",
			want:          nil,
			wantErr:       true,
			validateErrAs: os.ErrNotExist,
		},
		{
			name:          "file cannot be unmarshalled",
			path:          createActionsFileForTest(t, "yaml"),
			want:          nil,
			wantErr:       true,
			validateErrAs: nil,
		},
		{
			name:          "json file is loaded correctly",
			path:          filepath.Join(createActionsFileForTest(t, "json"), ".actions.stencil.json"),
			want:          actionConfig.Actions,
			wantErr:       false,
			validateErrAs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := loadActionsFromJSON(tt.path)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadActionsFromPath() failed: %v", gotErr)
				}
				if tt.validateErrAs != nil {
					if !errors.Is(gotErr, tt.validateErrAs) {
						t.Errorf("LoadActionsFromPath() return incorrect error: %v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadActionsFromPath() succeeded unexpectedly")
			}

			if stringifyObject(got) != stringifyObject(tt.want) {
				t.Errorf("LoadActionsFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadActionsFromYaml(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path          string
		want          []Action
		wantErr       bool
		validateErrAs error
	}{
		{
			name:          "path does not exist",
			path:          "path/does/not/exist",
			want:          nil,
			wantErr:       true,
			validateErrAs: os.ErrNotExist,
		},
		{
			name:          "file cannot be unmarshalled",
			path:          createActionsFileForTest(t, "json"),
			want:          nil,
			wantErr:       true,
			validateErrAs: nil,
		},
		{
			name:          "yaml file is loaded correctly",
			path:          filepath.Join(createActionsFileForTest(t, "yaml"), ".actions.stencil.yaml"),
			want:          actionConfig.Actions,
			wantErr:       false,
			validateErrAs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := loadActionsFromYAML(tt.path)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadActionsFromPath() failed: %v", gotErr)
				}
				if tt.validateErrAs != nil {
					if !errors.Is(gotErr, tt.validateErrAs) {
						t.Errorf("LoadActionsFromPath() return incorrect error: %v", gotErr)
					}
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadActionsFromPath() succeeded unexpectedly")
			}

			if stringifyObject(got) != stringifyObject(tt.want) {
				t.Errorf("LoadActionsFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createActionsFileForTest(t *testing.T, ext string) (tempPath string) {
	tempPath = t.TempDir()
	var (
		contents []byte
		err      error
	)
	if ext == "json" {
		contents, err = json.Marshal(actionConfig)
	} else {
		contents, err = yaml.Marshal(actionConfig)
	}
	if err != nil {
		t.Fatal("could not create file for testing")
	}
	_ = os.WriteFile(filepath.Join(tempPath, ".actions.stencil."+ext), contents, 0655)
	return
}

func stringifyObject(v any) string {
	return fmt.Sprintf("%v", v)
}
