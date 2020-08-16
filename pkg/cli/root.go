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
	"github.com/spf13/cobra"
)

var (
	// LogLevel Logging level (debug, info, warn, error, panic, fatal)
	LogLevel string
	// LogType Logging output type (console, json)
	LogType string
	// OutputType Violation output type (text, json, yaml, xml)
	OutputType string
	// ConfigFile Config file path
	ConfigFile string
)

var rootCmd = &cobra.Command{
	Use:   "terrascan",
	Short: "Detect compliance and security violations across Infrastructure as Code.",
	Long: `Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
For more information, please visit https://docs.accurics.com
`,
}
