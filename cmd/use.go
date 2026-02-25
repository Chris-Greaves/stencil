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

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [Path | URL]",
	Short: "Use a template",
	Long: `Use a template located at the specified path.

The path can be a local file system path or a remote URL. If the path is a remote URL, it should point to a Git repository containing the template.`,
	Args: func(cmd *cobra.Command, args []string) error {
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
}
