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
	"go.uber.org/zap"
)

var (
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
	//ConfigOnly will output resource config (should only be used for debugging purposes)
	ConfigOnly bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Detect compliance and security violations across Infrastructure as Code.",
	Long: `Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initial(cmd, args)
	},
	Run: scan,
}

func scan(cmd *cobra.Command, args []string) {
	zap.S().Debug("running terrascan in cli mode")
	Run(IacType, IacVersion, PolicyType, IacFilePath, IacDirPath, ConfigFile, PolicyPath, OutputType, ConfigOnly)
}

func init() {
	scanCmd.Flags().StringVarP(&PolicyType, "policy-type", "t", "", "<required> policy type (aws, azure, gcp, k8s)")
	scanCmd.Flags().StringVarP(&IacType, "iac-type", "i", "", "iac type (terraform, k8s)")
	scanCmd.Flags().StringVarP(&IacVersion, "iac-version", "", "", "iac version terraform:(v12) k8s:(v1)")
	scanCmd.Flags().StringVarP(&IacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&IacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringVarP(&PolicyPath, "policy-path", "p", "", "policy path directory")
	scanCmd.Flags().BoolVarP(&ConfigOnly, "config-only", "", false, "will output resource config (should only be used for debugging purposes)")
	scanCmd.MarkFlagRequired("policy-type")
	RegisterCommand(rootCmd, scanCmd)
}
