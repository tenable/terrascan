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
	"github.com/accurics/terrascan/pkg/config"
	//"github.com/accurics/terrascan/pkg/version"
	//"github.com/accurics/terrascan/pkg/initialize"
	//"gopkg.in/src-d/go-git.v4"

	"io"
	"os"
	"net/http"
	"io/ioutil"



)

var scanOptions = NewScanOptions()

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Detect compliance and security violations across Infrastructure as Code.",
	Long: `Terrascan

Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if scanOptions.configOnly {
			return nil
		}

		 err :=  initial(cmd, args, true)
		return err
	},
	RunE:          scan,
	SilenceUsage:  true,
	SilenceErrors: true,
}


func scan(cmd *cobra.Command, args []string) error {
	zap.S().Debug("running terrascan in cli mode")
	scanOptions.configFile = ConfigFile
	scanOptions.outputType = OutputType
	//fmt.Println("from scan" + config.Tag)

	//test reading from tagversion file 
	basePath := config.GetPolicyBasePath()
	filename := basePath + "/TagVersion"
	f, err := os.Open(filename)
    if err != nil {
        fmt.Printf("error opening %s: %s", filename, err)
    }
    defer f.Close()

    buf := make([]byte, 12)
    if _, err := io.ReadFull(f, buf); err != nil {
        if err == io.EOF {
            err = io.ErrUnexpectedEOF
        }
	}
	tagUsed := string(buf)

	if isLatest(tagUsed) == false {
		fmt.Printf("Using an old release of policy repo (%s). Enter 'C'  to proceed with outdated scan, or enter 'Q' to exit scan. To use the latest release of your policy repo, run terrascan init and scan again. \n", tagUsed)
		var input string 
		fmt.Scanln(&input)
		if input == "Q" || input == "q" {
			return nil
		} else if input == "C" || input == "c" {
			return scanOptions.Scan()
		} else {
			fmt.Println("Input not recognized, please enter Continue or press enter")
		}
	} else {
		fmt.Println("You are using the latest policy release!")
	}
	//end of reading from tag version file 

	return scanOptions.Scan()
}

func init() {
	scanCmd.Flags().StringSliceVarP(&scanOptions.policyType, "policy-type", "t", []string{"all"}, fmt.Sprintf("policy type (%s)", strings.Join(policy.SupportedPolicyTypes(true), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacType, "iac-type", "i", "", fmt.Sprintf("iac type (%v)", strings.Join(iacProvider.SupportedIacProviders(), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacVersion, "iac-version", "", "", fmt.Sprintf("iac version (%v)", strings.Join(iacProvider.SupportedIacVersions(), ", ")))
	scanCmd.Flags().StringVarP(&scanOptions.iacFilePath, "iac-file", "f", "", "path to a single IaC file")
	scanCmd.Flags().StringVarP(&scanOptions.iacDirPath, "iac-dir", "d", ".", "path to a directory containing one or more IaC files")
	scanCmd.Flags().StringArrayVarP(&scanOptions.policyPath, "policy-path", "p", []string{}, "policy path directory")
	scanCmd.Flags().StringVarP(&scanOptions.remoteType, "remote-type", "r", "", "type of remote backend (git, s3, gcs, http, terraform-registry)")
	scanCmd.Flags().StringVarP(&scanOptions.remoteURL, "remote-url", "u", "", "url pointing to remote IaC repository")
	scanCmd.Flags().BoolVarP(&scanOptions.configOnly, "config-only", "", false, "will output resource config (should only be used for debugging purposes)")
	// flag passes a string, but we normalize to bool in PreRun
	scanCmd.Flags().StringVar(&scanOptions.useColors, "use-colors", "auto", "color output (auto, t, f)")
	scanCmd.Flags().BoolVarP(&scanOptions.verbose, "verbose", "v", false, "will show violations with details (applicable for default output)")
	scanCmd.Flags().StringSliceVarP(&scanOptions.scanRules, "scan-rules", "", []string{}, "one or more rules to scan (example: --scan-rules=\"ruleID1,ruleID2\")")
	scanCmd.Flags().StringSliceVarP(&scanOptions.skipRules, "skip-rules", "", []string{}, "one or more rules to skip while scanning (example: --skip-rules=\"ruleID1,ruleID2\")")
	scanCmd.Flags().StringVar(&scanOptions.severity, "severity", "", "minimum severity level of the policy violations to be reported by terrascan")
	scanCmd.Flags().StringSliceVarP(&scanOptions.categories, "categories", "", []string{}, "list of categories of violations to be reported by terrascan (example: --categories=\"category1,category2\")")
	scanCmd.Flags().BoolVarP(&scanOptions.showPassedRules, "show-passed", "", false, "display passed rules, along with violations")
	scanCmd.Flags().BoolVarP(&scanOptions.nonRecursive, "non-recursive", "", false, "do not scan directories and modules recursively")
	scanCmd.Flags().BoolVarP(&scanOptions.useTerraformCache, "use-terraform-cache", "", false, "use terraform init cache for remote modules (when used directory scan will be non recursive, flag applicable only with terraform IaC provider)")
	RegisterCommand(rootCmd, scanCmd)
}

type tag_name struct {
    tagName string
}
func isLatest(initTag string) bool { 
	address := config.GetAPI()
	resp, err := http.Get(address)
	if err != nil {
		fmt.Errorf("Couldn't access latest tag '%v'", err)

	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Couldn't read request response '%v'", err)
	}
	tar := "tag_name"
	sb := string(body)
	si := strings.Index(sb, tar) 
	starting := si + len(tar) + 3 
	indent := 0 
	for string(sb[starting + indent]) != "\"" {
		indent = indent + 1 
	}
	latestRelease := sb[starting:starting+indent]
	fmt.Printf(latestRelease)

	if (initTag != latestRelease) { 
		return false
	} else {
		return true 
	}

}