/*
    Copyright (C) 2022 Tenable, Inc.

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
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/logging"
	"github.com/tenable/terrascan/pkg/utils"
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
	rootCmd.PersistentFlags().StringVarP(&CustomTempDir, "temp-dir", "", "", "temporary directory path to download remote repository,module and templates")
	rootCmd.PersistentFlags().StringVarP(&LogOutputDir, "log-output-dir", "", "", "directory path to write the log and output files")

	//Added init here in case flag parsing failed we should log which flag was incorrect.
	logging.Init(LogType, LogLevel, LogOutputDir)

	// Function to execute before processing commands
	cobra.OnInitialize(func() {
		// making sure the LogOutputDir Exist
		if LogOutputDir != "" {
			err := os.MkdirAll(LogOutputDir, 0755)
			if err != nil {
				zap.S().Warnf("failed to resolve the log output directory: %s", LogOutputDir)
				LogOutputDir = ""
			}
		}

		// Set up the logger
		logging.Init(LogType, LogLevel, LogOutputDir)

		if len(ConfigFile) == 0 {
			ConfigFile = os.Getenv(config.ConfigEnvvarName)
			zap.S().Debugf("%s:%s", config.ConfigEnvvarName, os.Getenv(config.ConfigEnvvarName))
		}

		// Make sure we load the global config from the specified config file
		if err := config.LoadGlobalConfig(ConfigFile); err != nil {
			zap.S().Error("error while loading global config", zap.Error(err))
			os.Exit(1)
		}
		if CustomTempDir != "" {
			utils.CustomTempDir = CustomTempDir
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
		log.SetOutput(io.Discard)
	}

	if err := rootCmd.Execute(); err != nil {
		// check if the error is related to flag argument missing and
		// log it before terminating the process so user gets idea about the incorrect flag value
		if strings.Contains(err.Error(), "flag needs an argument") {
			zap.S().Error("error while executing command ", zap.Error(err))
		}
		os.Exit(1)
	}
}
