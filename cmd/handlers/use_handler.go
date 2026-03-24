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
	"path/filepath"

	"github.com/Chris-Greaves/stencil/action"
	"github.com/Chris-Greaves/stencil/processor"
	"github.com/Chris-Greaves/stencil/repository"
	"github.com/Chris-Greaves/stencil/utils"
)

type UseHandlerFlags struct {
	Debug  bool
	Repo   string
	Action string
}

type UseHandler struct {
	flags         UseHandlerFlags
	isUsingRepo   bool
	isUsingAction bool
}

func NewUseHandler() UseHandler {
	return UseHandler{}
}

func (h *UseHandler) ValidateArgs(args []string) error {
	if len(args) > 2 || len(args) == 0 {
		return ErrInvalidNumberOfArguments
	}

	if h.flags.Repo != "" {
		h.isUsingRepo = true
		exists, err := repository.RepositoryExists(h.flags.Repo)
		if err != nil {
			return err
		}
		if !exists {
			return ErrRepositoryNotFound
		}
	}

	if len(args) == 1 {
		if h.flags.Action != "" {
			h.isUsingAction = true
			// TODO: Check action exists
			if !utils.PathExistsAndIsDir(args[0]) {
				return ErrUnableToFindOutput
			}
			return nil // using action, no need to continue validations
		} else {
			// If the action flag is provided, we can assume the single argument is the output path, and the template path will be looked up based on the action name
			return ErrMissingActionFlag
		}
	}

	if len(args) == 2 && h.isUsingAction {
		return ErrDuplicateTemplatePathArgs
	}

	if !utils.PathExistsAndIsDir(args[0]) {
		return ErrUnableToFindTemplate
	}

	if !utils.PathExistsAndIsDir(args[1]) {
		return ErrUnableToFindOutput
	}

	return nil
}

func (h *UseHandler) SetFlags(debug bool, repo string, action string) {
	h.flags.Debug = debug
	h.flags.Repo = repo
	h.flags.Action = action
}

func (h *UseHandler) Handle(args []string) error {
	if h.flags.Debug {
		fmt.Printf("Flags: %v\n", h.flags)
	}

	var (
		templatePath string
		outputPath   string
	)

	if len(args) == 1 {
		outputPath = args[0]
	} else {
		templatePath = args[0]
		outputPath = args[1]
	}

	// If a repo flag is provided, we need to look up the path to that repo and override the template path with it
	// This allows users to specify templates within repos without needing to know the local path to the repo on their machine
	if h.isUsingRepo {
		repoPath, err := repository.GetRepositoryPath(h.flags.Repo)
		if err != nil {
			return err
		}

		templatePath = filepath.Join(repoPath, templatePath)
		templatePath = filepath.Clean(templatePath)
		fmt.Printf("Using repo '%s' with template path '%s'\n", h.flags.Repo, templatePath)
	}

	if h.isUsingAction {
		var (
			actionPath string
			actions    []action.Action
			err        error
		)

		if h.isUsingRepo {
			if h.flags.Debug {
				fmt.Printf("Loading actions from repo path: %s\n", templatePath)
			}

			actions, err = action.LoadActionsFromPath(templatePath)
		} else {
			actions, err = action.LoadActionsFromPath(".")
		}
		if err != nil {
			return err
		}

		if h.flags.Debug {
			fmt.Printf("Loaded actions: %v\n", actions)
		}

		for _, action := range actions {
			if action.Name == h.flags.Action {
				// Found the action, use its template path
				actionPath = action.Path
				break
			}
		}

		if actionPath == "" {
			return ErrActionNotFound
		}

		templatePath = filepath.Clean(filepath.Join(templatePath, actionPath))
	}

	if h.flags.Debug {
		fmt.Printf("Template Path: %s\n", templatePath)
		fmt.Printf("Output Path: %s\n", outputPath)
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
