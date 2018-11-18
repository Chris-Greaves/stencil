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
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/chris-greaves/stencil/models"

	"github.com/chris-greaves/stencil/engine"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "stencil",
	Short: "stencil is a template engine, with a twist!",
	Long: `stencil is designed to be a very customisable
and user friendly templating engine.

Using the combined power of Go's built in template 
renderer and the user friendly features provided, 
this tool can allow for the utilisation of complex 
and highly customisable templates`,
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

		settings := models.Settings{"foo", "bar", "Hello World"}

		wd, err := os.Getwd()
		if err != nil {
			log.Panicf("Error getting Working Directory, %v", err)
		}

		fmt.Printf("Current working directory = %v\n", wd)

		filepath.Walk(args[0],
			func(path string, info os.FileInfo, err error) error {
				fmt.Printf("Creating %v\n", path)

				relPath, err := filepath.Rel(args[0], path)
				if err != nil {
					return fmt.Errorf("Error getting relative path: %v", err)
				}

				if relPath == "." {
					return nil
				}

				relDestPath, err := engine.ParseAndExecutePath(settings, relPath)
				if err != nil {
					return fmt.Errorf("Error while creating relative destination path: %v", err)
				}

				fmt.Printf("Relative destination path: %v\n", relDestPath)

				destinationPath := filepath.Join(wd, relDestPath)

				if info.IsDir() {
					if err = os.MkdirAll(destinationPath, info.Mode()); err != nil {
						return fmt.Errorf("Error making directory %v: %v", path, err)
					}
					return nil
				}

				if err = engine.ParseAndExecuteFile(settings, destinationPath, path, info.Mode()); err != nil {
					return fmt.Errorf("Error processing file %v: %v", path, err)
				}

				return nil
			})
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
