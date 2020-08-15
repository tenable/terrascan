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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RegisterCommand Registers a new command under the base command
func RegisterCommand(baseCommand *cobra.Command, command *cobra.Command) {
	baseCommand.AddCommand(command)
}

// Execute the entrypoint called by main
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "log level (debug, info, warn, error, panic, fatal)")
	rootCmd.PersistentFlags().StringVarP(&LogType, "log-type", "x", "console", "log output type (console, json)")
	rootCmd.PersistentFlags().StringVarP(&OutputType, "output-type", "o", "text", "output type (text, json, yaml, xml)")
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config-path", "c", "", "config file path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
