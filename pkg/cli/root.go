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

var (
	// Root Persistent

	// LogLevel Logging level (debug, info, warn, error, panic, fatal)
	LogLevel string
	// LogType Logging output type (console, json)
	LogType string
	// OutputType Violation output type (text, json, yaml, xml)
	OutputType string
	// ConfigFile Config file path
	ConfigFile string

	// Scan Local

	// PolicyPath Policy path directory
	PolicyPath string
	// PolicyType Cloud type (aws, azure, gcp)
	PolicyType string
	// IacType IaC type (terraform)
	IacType string
	// IacVersion IaC version (for terraform:v12)
	IacVersion string
	// IacFilePath Path to a single IaC file
	IacFilePath string
	// IacDirPath Path to a directory containing one or more IaC files
	IacDirPath string
)

var rootCmd = &cobra.Command{
	Use:   "terrascan",
	Short: "Terrascan is an IaC (Infrastructure-as-Code) file scanner",
	Long: `Terrascan: An advanced IaC (Infrastructure-as-Code) file scanner written in Go.
           Secure your cloud deployments at design time.
           For more information, please visit https://www.accurics.com`,
	Version: "1.0.0", //version.Get(), // @TODO: Move to a separate file
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan IaC (Infrastructure-as-Code) files for vulnerabilities.",
	Long:  "Scan IaC (Infrastructure-as-Code) files for vulnerabilities.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initial(cmd, args)
	},
	Run: scan,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Terrascan as an API server",
	Long:  "An API server that inspects incoming IaC (Infrastructure-as-Code) files and returns the scan results.",
	PreRun: func(cmd *cobra.Command, args []string) {
		initial(cmd, args)
	},
	Run: server,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Terrascan",
	Long:  "Initializes Terrascan and clones policies from the Terrascan GitHub repository.",
	Run:   initial,
}

// Execute the entrypoint called by main
func Execute() {
	scanCmd.Flags().StringVarP(&PolicyType, "policy-type", "t", "", "<required> policy type (aws, azure, gcp)")
	scanCmd.Flags().StringVarP(&IacType, "iac-type", "i", "terraform", "iac type (terraform)")
	scanCmd.Flags().StringVarP(&IacVersion, "iac-version", "", "v12", "iac version (v12)")
	scanCmd.Flags().StringVarP(&IacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&IacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringVarP(&PolicyPath, "policy-path", "p", "", "policy path directory")
	scanCmd.MarkFlagRequired("policy-type")

	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "log level (debug, info, warn, error, panic, fatal)")
	rootCmd.PersistentFlags().StringVarP(&LogType, "log-type", "x", "console", "log output type (console, json)")
	rootCmd.PersistentFlags().StringVarP(&OutputType, "output-type", "o", "text", "output type (text, json, yaml, xml)")
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config-path", "c", "", "config file path")
	rootCmd.AddCommand(scanCmd, serverCmd, initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
