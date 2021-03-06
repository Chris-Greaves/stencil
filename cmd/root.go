// Copyright © 2018 Christopher Greaves <cjgreaves97@hotmail.co.uk>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Chris-Greaves/stencil/IO"
	"github.com/Chris-Greaves/stencil/cmd/handlers"
	"github.com/Chris-Greaves/stencil/confighelper"
	"github.com/Chris-Greaves/stencil/engine"
	"github.com/Chris-Greaves/stencil/fetch"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile                 string
	templatePath            string
	usingGit                = false
	ErrNoArguments          = errors.New("You must provide the path to the template")
	ErrUnableToFindTemplate = errors.New("stencil was unable to find a local path or git repository using the path provided")
)

var rootCmd = &cobra.Command{
	Use:   "stencil [path]",
	Short: "stencil is a tool to parse and execute project templates, using Go's built in template package",
	Long: `stencil is designed to be a very customisable and user friendly tool, allowing you to execute templates using Go's text/template package.

Note: Only the first argument passed in will be processed.

By utilising the Go's template package we have opened the ability to create unique and complex templates, easily.

View the documentation on http://christophergreaves.co.uk/projects/stencil/documentation`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return ErrNoArguments
		}

		if !fetch.IsPath(args[0]) {
			if !fetch.IsGitURL(args[0]) {
				return ErrUnableToFindTemplate
			} else {
				usingGit = true
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		templatePath := args[0]
		println(args[0])

		wd, err := os.Getwd()
		if err != nil {
			log.Panicf("Error getting Working Directory, %v", err)
		}
		fmt.Printf("Current working directory = %v\n", wd)

		if usingGit {
			dir, err := fetch.PullTemplate(templatePath)
			if err != nil {
				log.Panicf("Error retrieving git repo: %v", err.Error())
			}
			defer os.RemoveAll(dir)
			templatePath = dir
		}

		config, err := confighelper.New(filepath.Join(templatePath, ".stencil/.stencil.json"))
		if err != nil {
			log.Panicf("Error parsing config file: %v", err.Error())
		}

		templateEngine := engine.New()

		handler := handlers.NewRootHandler(config, templateEngine, new(IO.CLI))

		handler.OfferConfigOverrides()

		err = handler.ProcessTemplate(templatePath, wd)
		if err != nil {
			log.Panicf("Error while creating project from template, %v", err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stencil.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".stencil" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".stencil")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
