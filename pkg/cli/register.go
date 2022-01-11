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
	"io/ioutil"
	"log"
	"os"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// RegisterCommand Registers a new command under the base command
func RegisterCommand(baseCommand *cobra.Command, command *cobra.Command) {
	baseCommand.AddCommand(command)
}

// Execute the entrypoint called by main
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "log level (debug, info, warn, error, panic, fatal)")
	rootCmd.PersistentFlags().StringVarP(&LogType, "log-type", "x", "console", "log output type (console, json)")
	rootCmd.PersistentFlags().StringVarP(&OutputType, "output", "o", "human", "output type (human, json, yaml, xml, junit-xml, sarif, github-sarif)")
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config-path", "c", "", "config file path")

	//Added init here in case flag parsing failed we should log which flag was incorrect.
	logging.Init(LogType, LogLevel)

	// Function to execute before processing commands
	cobra.OnInitialize(func() {
		// Set up the logger
		logging.Init(LogType, LogLevel)

		if len(ConfigFile) == 0 {
			ConfigFile = os.Getenv(config.ConfigEnvvarName)
			zap.S().Debugf("%s:%s", config.ConfigEnvvarName, os.Getenv(config.ConfigEnvvarName))
		}

		// Make sure we load the global config from the specified config file
		if err := config.LoadGlobalConfig(ConfigFile); err != nil {
			zap.S().Error("error while loading global config", zap.Error(err))
			os.Exit(1)
		}
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

	// disable terraform logs when TF_LOG env variable is not set
	if os.Getenv("TF_LOG") == "" {
		log.SetOutput(ioutil.Discard)
	}

	if err := rootCmd.Execute(); err != nil {
		// log error before terminating the process so user gets idea about the error
		zap.S().Error("error while executing command ", zap.Error(err))
		os.Exit(1)
	}
}
