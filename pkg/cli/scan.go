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
	"strings"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	// PolicyPath Policy path directory
	PolicyPath []string

	// PolicyType Cloud type (aws, azure, gcp, github)
	PolicyType []string

	// IacType IaC type (terraform)
	IacType string

	// IacVersion IaC version (for terraform:v12)
	IacVersion string

	// IacFilePath Path to a single IaC file
	IacFilePath string

	// IacDirPath Path to a directory containing one or more IaC files
	IacDirPath string

	// RemoteType indicates the type of remote backend. Supported backends are
	// git s3, gcs, http.
	RemoteType string

	// RemoteURL points to the remote Iac repository on git, s3, gcs, http
	RemoteURL string

	// ConfigOnly will output resource config (should only be used for debugging purposes)
	ConfigOnly bool

	// UseColors indicates whether to use color output
	UseColors bool
	useColors string // used for flag processing
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Detect compliance and security violations across Infrastructure as Code.",
	Long: `Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		switch strings.ToLower(useColors) {
		case "auto":
			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				UseColors = true
			} else {
				UseColors = false
			}

		case "true":
			fallthrough
		case "t":
			fallthrough
		case "y":
			fallthrough
		case "1":
			fallthrough
		case "force":
			UseColors = true

		default:
			UseColors = false
		}
		initial(cmd, args)
	},
	Run: scan,
}

func scan(cmd *cobra.Command, args []string) {
	zap.S().Debug("running terrascan in cli mode")
	Run(IacType, IacVersion, PolicyType, IacFilePath, IacDirPath, ConfigFile,
		PolicyPath, OutputType, RemoteType, RemoteURL, ConfigOnly, UseColors)
}

func init() {
	scanCmd.Flags().StringSliceVarP(&PolicyType, "policy-type", "t", []string{"all"}, fmt.Sprintf("policy type (%s)", strings.Join(policy.SupportedPolicyTypes(true), ", ")))
	scanCmd.Flags().StringVarP(&IacType, "iac-type", "i", "", fmt.Sprintf("iac type (%v)", strings.Join(iacProvider.SupportedIacProviders(), ", ")))
	scanCmd.Flags().StringVarP(&IacVersion, "iac-version", "", "", fmt.Sprintf("iac version (%v)", strings.Join(iacProvider.SupportedIacVersions(), ", ")))
	scanCmd.Flags().StringVarP(&IacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&IacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringArrayVarP(&PolicyPath, "policy-path", "p", []string{}, "policy path directory")
	scanCmd.Flags().StringVarP(&RemoteType, "remote-type", "r", "", "type of remote backend (git, s3, gcs, http)")
	scanCmd.Flags().StringVarP(&RemoteURL, "remote-url", "u", "", "url pointing to remote IaC repository")
	scanCmd.Flags().BoolVarP(&ConfigOnly, "config-only", "", false, "will output resource config (should only be used for debugging purposes)")
	// flag passes a string, but we normalize to bool in PreRun
	scanCmd.Flags().StringVar(&useColors, "use-colors", "auto", "color output (auto, t, f)")
	RegisterCommand(rootCmd, scanCmd)
}
