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

var scanOptions = NewScanOptions()

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
	scanOptions.configFile = ConfigFile
	scanOptions.outputType = OutputType
	scanOptions.Scan()
}

func init() {
	scanCmd.Flags().StringSliceVarP(&scanOptions.policyType, "policy-type", "t", []string{"all"}, fmt.Sprintf("policy type (%s)", strings.Join(policy.SupportedPolicyTypes(true), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacType, "iac-type", "i", "", fmt.Sprintf("iac type (%v)", strings.Join(iacProvider.SupportedIacProviders(), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacVersion, "iac-version", "", "", fmt.Sprintf("iac version (%v)", strings.Join(iacProvider.SupportedIacVersions(), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&scanOptions.iacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringArrayVarP(&scanOptions.policyPath, "policy-path", "p", []string{}, "policy path directory")
	scanCmd.Flags().StringVarP(&scanOptions.remoteType, "remote-type", "r", "", "type of remote backend (git, s3, gcs, http)")
	scanCmd.Flags().StringVarP(&scanOptions.remoteURL, "remote-url", "u", "", "url pointing to remote IaC repository")
	scanCmd.Flags().BoolVarP(&scanOptions.configOnly, "config-only", "", false, "will output resource config (should only be used for debugging purposes)")
	// flag passes a string, but we normalize to bool in PreRun
	scanCmd.Flags().StringVar(&scanOptions.useColors, "use-colors", "auto", "color output (auto, t, f)")
	scanCmd.Flags().BoolVarP(&scanOptions.Verbose, "verbose", "v", false, "will show violations with details (applicable for default output)")
	scanCmd.Flags().StringSliceVarP(&scanOptions.scanRules, "scan-rules", "", []string{}, "one or more rules to scan (example: --scan-rules=\"ruleID1,ruleID2\")")
	scanCmd.Flags().StringSliceVarP(&scanOptions.skipRules, "skip-rules", "", []string{}, "one or more rules to skip while scanning (example: --skip-rules=\"ruleID1,ruleID2\")")
	RegisterCommand(rootCmd, scanCmd)
}
