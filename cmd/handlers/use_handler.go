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

package handlers

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/Chris-Greaves/stencil/processor"
	"github.com/Chris-Greaves/stencil/repository"
	"github.com/Chris-Greaves/stencil/utils"
)

type UseHandlerFlags struct {
	Debug bool
	Repo  string
}

type UseHandler struct {
	flags UseHandlerFlags
}

func NewUseHandler() UseHandler {
	return UseHandler{}
}

func (h *UseHandler) ValidateArgs(args []string) error {
	if len(args) != 2 {
		return ErrInvalidNumberOfArguments
	}

	if !utils.PathExistsAndIsDir(args[0]) {
		return ErrUnableToFindTemplate
	}

	if !utils.PathExistsAndIsDir(args[1]) {
		return ErrUnableToFindOutput
	}

	return nil
}

func (h *UseHandler) SetFlags(debug bool, repo string) {
	h.flags.Debug = debug
	h.flags.Repo = repo
}

func (h *UseHandler) Handle(args []string) error {
	templatePath := args[0]
	outputPath := args[1]

	if h.flags.Debug {
		fmt.Printf("Template Path: %s\n", templatePath)
		fmt.Printf("Output Path: %s\n", outputPath)
	}

	// If a repo flag is provided, we need to look up the path to that repo and override the template path with it
	// This allows users to specify templates within repos without needing to know the local path to the repo on their machine
	if h.flags.Repo != "" {
		repoPath, err := repository.GetRepositoryPath(h.flags.Repo)
		if err != nil {
			return err
		}

		templatePath = path.Join(repoPath, templatePath)
		templatePath = filepath.Clean(templatePath)
		fmt.Printf("Using repo '%s' with template path '%s'\n", h.flags.Repo, templatePath)
	}

	proc, err := processor.NewProcessor(templatePath, outputPath)
	if err != nil {
		return err
	}

	if h.flags.Debug {
		fmt.Printf("Config: %s\n", proc.DumpConfig())
	}

	err = proc.PromptUserForInput()
	if err != nil {
		return err
	}

	if h.flags.Debug {
		fmt.Printf("Values: %s\n", proc.DumpValues())
	}

	err = proc.ExecuteTemplate(h.flags.Debug)
	if err != nil {
		return err
	}

	return nil
}
