/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/logging"
	"github.com/spf13/cobra"
)

// RegisterCommand Registers a new command under the base command
func RegisterCommand(baseCommand *cobra.Command, command *cobra.Command) {
	baseCommand.AddCommand(command)
}

func subCommands() (commandNames []string) {
	for _, command := range rootCmd.Commands() {
		commandNames = append(commandNames, append(command.Aliases, command.Name())...)
	}
	return
}

// setDefaultCommand sets `scan` as default command if no other command is specified
func setDefaultCommandIfNonePresent() {
	if len(os.Args) > 1 {
		potentialCommand := os.Args[1]
		for _, command := range subCommands() {
			if command == potentialCommand {
				return
			}
		}
		os.Args = append([]string{os.Args[0], "scan"}, os.Args[1:]...)
	}

}

// Execute the entrypoint called by main
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "log level (debug, info, warn, error, panic, fatal)")
	rootCmd.PersistentFlags().StringVarP(&LogType, "log-type", "x", "console", "log output type (console, json)")
	rootCmd.PersistentFlags().StringVarP(&OutputType, "output", "o", "yaml", "output type (json, yaml, xml)")
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config-path", "c", "", "config file path")

	// Function to execute before processing commands
	cobra.OnInitialize(func() {
		// Set up the logger
		logging.Init(LogType, LogLevel)
		// Make sure we load the global config from the specified config file
		config.LoadGlobalConfig(ConfigFile)
	})

	// parse the flags but hack around to avoid exiting with error code 2 on help
	// override usage so that flag.Parse uses root command's usage instead of default one when invoked with -h
	flag.Usage = func() {
		_ = rootCmd.Help()
	}

	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	args := os.Args[1:]
	if err := flag.CommandLine.Parse(args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
	}

	setDefaultCommandIfNonePresent()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
