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

package utils_test

import (
	"os"
	"testing"

	"github.com/Chris-Greaves/stencil/utils"
)

func TestPathExistsAndIsDir(t *testing.T) {
	existingDir := t.TempDir()
	existingFile := t.TempDir() + "/file.txt"
	file, err := os.Create(existingFile)
	if err != nil {
		t.Fatal("Failed to create test file")
	}
	defer func() { // Putting the close in a deferred function to appease the linter
		err := file.Close()
		if err != nil {
			t.Errorf("Failed to close destination file: %v", err)
		}
	}()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path string
		want bool
	}{
		{
			name: "path does not exist",
			path: "/non/existent/path",
			want: false,
		},
		{
			name: "path exists and is a directory",
			path: existingDir,
			want: true,
		},
		{
			name: "path exists but is not a directory",
			path: existingFile,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.PathExistsAndIsDir(tt.path)
			if got != tt.want {
				t.Errorf("PathExistsAndIsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
