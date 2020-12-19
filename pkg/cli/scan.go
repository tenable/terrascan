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
	"strings"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var scanCommand = new(ScanCommand)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Detect compliance and security violations across Infrastructure as Code.",
	Long: `Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
`,
	PreRun: initial,
	Run:    scan,
}

func scan(cmd *cobra.Command, args []string) {
	zap.S().Debug("running terrascan in cli mode")
	scanCommand.configFile = ConfigFile
	scanCommand.outputType = OutputType
	scanCommand.StartScan()
}

func init() {
	scanCmd.Flags().StringSliceVarP(&scanCommand.policyType, "policy-type", "t", []string{"all"}, fmt.Sprintf("policy type (%s)", strings.Join(policy.SupportedPolicyTypes(true), ", ")))
	scanCmd.Flags().StringVarP(&scanCommand.iacType, "iac-type", "i", "", fmt.Sprintf("iac type (%v)", strings.Join(iacProvider.SupportedIacProviders(), ", ")))
	scanCmd.Flags().StringVarP(&scanCommand.iacVersion, "iac-version", "", "", fmt.Sprintf("iac version (%v)", strings.Join(iacProvider.SupportedIacVersions(), ", ")))
	scanCmd.Flags().StringVarP(&scanCommand.iacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&scanCommand.iacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringArrayVarP(&scanCommand.policyPath, "policy-path", "p", []string{}, "policy path directory")
	scanCmd.Flags().StringVarP(&scanCommand.remoteType, "remote-type", "r", "", "type of remote backend (git, s3, gcs, http)")
	scanCmd.Flags().StringVarP(&scanCommand.remoteURL, "remote-url", "u", "", "url pointing to remote IaC repository")
	scanCmd.Flags().BoolVarP(&scanCommand.configOnly, "config-only", "", false, "will output resource config (should only be used for debugging purposes)")
	// flag passes a string, but we normalize to bool in PreRun
	scanCmd.Flags().StringVar(&scanCommand.useColors, "use-colors", "auto", "color output (auto, t, f)")
	scanCmd.Flags().BoolVarP(&scanCommand.Verbose, "verbose", "v", false, "will show violations with details (applicable for default output)")
	RegisterCommand(rootCmd, scanCmd)
}
