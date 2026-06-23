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

/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Chris-Greaves/stencil/repository"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the stencil repositories",
	Long: `By listing your stencil repositories, you can see all the repositories that you have added using stencil repo add.
	This is useful for keeping track of your repositories and ensuring that you have access to the templates you need.`,
	Run: func(cmd *cobra.Command, args []string) {
		rm := repository.NewDefaultRepositoryManager()
		repos, err := rm.ListRepositories()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return
		}

		for _, repo := range repos {
			fmt.Printf("Name: %s, URL: %s\n", repo.Name, repo.URL)
		}
	},
}

func init() {
	repoCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
