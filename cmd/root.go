// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/chris-greaves/stencil/engine"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "stencil [path]",
	Short: "stencil is a tool to parse and execute project templates, using Go's built in template package",
	Long: `stencil is designed to be a very customisable and user friendly tool, allowing you to execute templates using Go's text/template package.

By utilising the Go's template package we have opened the ability to create unique and complex templates, easily.

View the documentation on http://christophergreaves.co.uk/projects/stencil/documentation`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("You must provide at least 1 argument")
		}

		info, err := os.Stat(args[0])
		if err != nil {
			return fmt.Errorf("Error finding path: %v", err)

		}

		if !info.IsDir() {
			return errors.New("Path specified must be a directory")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		println(args[0])

		settings, err := engine.GetSettings(args[0])
		if err != nil {
			log.Panicf("Error getting settings from settings file. Make sure a stencil.json file exists at the root directory of the template.: %v", err)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Panicf("Error getting Working Directory, %v", err)
		}

		fmt.Printf("Current working directory = %v\n", wd)

		if err = filepath.Walk(args[0],
			func(path string, info os.FileInfo, err error) error {
				fmt.Printf("Creating %v\n", path)

				relPath, err := filepath.Rel(args[0], path)
				if err != nil {
					return errors.Wrap(err, "Error getting relative path")
				}

				if relPath == "." {
					return nil
				}

				relDestPath, err := engine.ParseAndExecutePath(settings, relPath)
				if err != nil {
					return errors.Wrap(err, "Error while creating relative destination path")
				}

				fmt.Printf("Relative destination path: %v\n", relDestPath)

				destinationPath := filepath.Join(wd, relDestPath)

				if info.IsDir() {
					if err = os.MkdirAll(destinationPath, info.Mode()); err != nil {
						return errors.Wrapf(err, "Error making directory %v", path)
					}
					return nil
				}

				if err = engine.ParseAndExecuteFile(settings, destinationPath, path, info.Mode()); err != nil {
					return errors.Wrapf(err, "Error processing file %v", path)
				}

				return nil
			}); err != nil {
			log.Panicf("Error while creating project from tempalate, %v", err)
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
