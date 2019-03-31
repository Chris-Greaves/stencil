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
	"strings"

	"github.com/chris-greaves/stencil/fetch"

	"github.com/chris-greaves/stencil/engine"
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
			if !fetch.IsGitInstalled() || !fetch.IsGitURL(args[0]) {
				return ErrUnableToFindTemplate
			} else {
				repoDirectory, err := fetch.PullTemplate(args[0])
				if err != nil {
					return err
				}

				templatePath = repoDirectory
				usingGit = true
			}
		} else {
			templatePath = args[0]
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		println(templatePath)

		wd, err := os.Getwd()
		if err != nil {
			log.Panicf("Error getting Working Directory, %v", err)
		}
		fmt.Printf("Current working directory = %v\n", wd)

		if usingGit {
			defer os.RemoveAll(templatePath)
		}

		ProcessTemplate(templatePath, wd)
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

func ProcessTemplate(templatePath string, outputPath string) {
	settings, err := engine.GetSettings(templatePath)
	if err != nil {
		log.Panicf("Error getting settings from settings file. Make sure a stencil.json file exists at the root directory of the template.: %v", err)
	}

	if err = filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			// Skip if root or part of git
			if path == templatePath || strings.Contains(path, ".git") {
				return nil
			}

			fmt.Printf("Creating %v\n", path)

			tarPath, err := GetTargetPath(templatePath, outputPath, path, settings)
			if err != nil {
				return err
			}

			// Create target
			if info.IsDir() {
				// If its a Directory, create the directory in the target
				if err = os.MkdirAll(tarPath, info.Mode()); err != nil {
					return errors.Wrapf(err, "Error making directory %v", path)
				}
			} else {
				// If its a file, parse and execute the file and copy the result to the target
				if err = engine.ParseAndExecuteFile(settings, tarPath, path, info.Mode()); err != nil {
					return errors.Wrapf(err, "Error processing file %v", path)
				}
			}

			return nil
		}); err != nil {
		log.Panicf("Error while creating project from template, %v", err)
	}
}

func GetTargetPath(templatePath string, outputPath string, path string, settings interface{}) (string, error) {
	relPath, err := filepath.Rel(templatePath, path)
	if err != nil {
		return "", errors.Wrap(err, "Error getting relative path")
	}

	relTarPath, err := engine.ParseAndExecutePath(settings, relPath)
	if err != nil {
		return "", err
	}
	tarPath := filepath.Join(outputPath, relTarPath)
	return tarPath, nil
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
