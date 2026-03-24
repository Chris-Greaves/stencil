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

package cmd

import (
	"fmt"

	"github.com/Chris-Greaves/stencil/cmd/handlers"
	"github.com/spf13/cobra"
)

var useHandler = handlers.NewUseHandler()

var (
	debug  *bool
	repo   *string
	action *string
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use {<template_path> | --action <action>} <output_path>",
	Short: "Use a template",
	Example: `# Basic usage, path to template on fs and where to output it to
stencil use ./templates/create-new ./out

# Same as before, but this time the path is from within a repo
stencil use --repo my_stencils ./create-new ./out

# Using an action defined in a .actions.stencil.yaml|yml|json file in the working directory
stencil use --action create-new ./out

# Calling an action defined in a repo
stencil use --repo my_stencils --action create-new ./out`,
	Long: `Use a template, and output the result to the specified output path.
Check the examples for more details on how to use this command with repos and actions.

The output path must exist and be writable.`,
	Args: func(cmd *cobra.Command, args []string) error {
		useHandler.SetFlags(*debug, *repo, *action)
		return useHandler.ValidateArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := useHandler.Handle(args)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	repo = useCmd.Flags().StringP("repo", "r", "", "Specify a repo containing the template you want")
	action = useCmd.Flags().StringP("action", "a", "", "Specify an action to perform")
	debug = useCmd.Flags().Bool("debug", false, "Enable debug mode")
}
