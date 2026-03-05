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
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

func getTargetPath(targetBase string, sourceBase string, sourcePath string) string {
	relativePath := strings.TrimPrefix(sourcePath, filepath.Clean(sourceBase))
	return filepath.Join(targetBase, relativePath)
}

func parseTemplateString(templateString string, values map[string]interface{}) (string, error) {
	var engine = template.New(templateString)
	tmpl, err := engine.Parse(templateString)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func parseTemplateFile(path string, values map[string]interface{}, wr io.Writer) error {
	_, filename := filepath.Split(path)
	var engine = template.New(filename) // Name template after the filename
	tmpl, err := engine.ParseFiles(path)
	if err != nil {
		return err
	}
	err = tmpl.Execute(wr, values)
	if err != nil {
		return errors.Join(fmt.Errorf("error executing template for file '%v'", path), err)
	}
	return nil
}
