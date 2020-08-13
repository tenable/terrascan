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

package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/accurics/terrascan/pkg/cli"
	httpServer "github.com/accurics/terrascan/pkg/http-server"
	"github.com/accurics/terrascan/pkg/logging"
	"go.uber.org/zap"
)

const (
	// CmdArgServer server command
	CmdArgServer = "server"
	// CmdArgIacType iac command
	CmdArgIacType = "iac"
	// CmdArgIacVersion iac-version command
	CmdArgIacVersion = "iac-version"
	// CmdArgFilePath file path command
	CmdArgFilePath = "f"
	// CmdArgDirPath dir path command
	CmdArgDirPath = "d"
	// CmdArgPolicyPath policy path command
	CmdArgPolicyPath = "p"
	// CmdArgCloudType cloud type command
	CmdArgCloudType = "cloud"
	// CmdArgVersion version command
	CmdArgVersion = "version"
	// CmdArgLogLevel log-level command
	CmdArgLogLevel = "log-level"
	// CmdArgLogType log-type command
	CmdArgLogType = "log-type"
	// CmdArgConfigFile config command
	CmdArgConfigFile = "config"
	// SectionMode mode section
	SectionMode = "Mode"
	// SectionIac iac section
	SectionIac = "IaC (Infrastructure as Code)"
	// SectionCloud cloud section
	SectionCloud = "Cloud"
	// SectionLogging logging section
	SectionLogging = "Logging"
	// SectionMisc miscellaneous commands
	SectionMisc = "Miscellaneous"
	// TerrascanDataDir Terrascan data directory
	TerrascanDataDir = ".terrascan"
)

// UsageSections Usage sections
var UsageSections = []string{SectionCloud, SectionIac, SectionMode, SectionLogging, SectionMisc}

// UsageSectionMap Sets the section and print order for each command line arg
var UsageSectionMap = map[string][]string{
	SectionCloud:   {CmdArgCloudType},
	SectionIac:     {CmdArgDirPath, CmdArgFilePath, CmdArgIacType, CmdArgIacVersion, CmdArgPolicyPath},
	SectionLogging: {CmdArgLogLevel, CmdArgLogType},
	SectionMode:    {CmdArgServer},
	SectionMisc:    {CmdArgConfigFile, CmdArgVersion},
}

// Usage This overrides the default flags usage information
var Usage = func() {
	fmt.Printf(`
Terrascan

Scan IaC files for security violations

Usage

    terrascan -cloud [aws|azure|gcp] [options...]

Options
`)

	flagMap := make(map[string]*flag.Flag)
	flag.VisitAll(func(f *flag.Flag) {
		flagMap[f.Name] = f
	})

	for i := range UsageSections {
		fmt.Printf("\n%s\n", UsageSections[i])
		for _, arg := range UsageSectionMap[UsageSections[i]] {
			fmt.Printf("    -%-20s %s\n", flagMap[arg].Name, flagMap[arg].Usage)
		}
	}

	fmt.Printf(`
Examples

    Scan Terraform v12 IaC files for the AWS provider under directory /home/user/iac_folder
        terrascan -cloud aws -d /home/user/iac_folder

    Scan Terraform v12 IaC files for the GCP provider under directory /home/user/iac_folder using custom policies under /home/user/policies
        terrascan -cloud gcp -d /home/user/iac_folder -p /home/user/policies

    Launch the API server
        terrascan -server
`)

}

func main() {

	// get home directory for the default policy path
	homeDir, err := homedir.Dir()
	if err != nil {
		zap.S().Error("error obtaining home directory", zap.Error(err))
		return
	}
	homePolicyPath := filepath.Join(homeDir, TerrascanDataDir, "policies")
	// command line flags
	var (
		// server mode
		server = flag.Bool("server", false, "Run Terrascan in server mode")

		// IaC flags
		iacType     = flag.String("iac", "terraform", "IaC provider (supported values: terraform, default: terraform)")
		iacVersion  = flag.String("iac-version", "v12", "IaC version (supported values: 'v12' for Terraform, default: v12)")
		iacFilePath = flag.String("f", "", "IaC file path")
		iacDirPath  = flag.String("d", ".", "IaC directory path (default: current working directory)")
		policyPath  = flag.String("p", homePolicyPath, "Policy directory path")

		// cloud flags
		cloudType = flag.String("cloud", "", "Required. Cloud provider (supported values: aws, azure, gcp)")

		// logging flags
		logLevel = flag.String("log-level", "info", "Logging level (supported values: debug, info, warn, error, panic, fatal)")
		logType  = flag.String("log-type", "console", "Logging type (supported values: json, yaml, console, default: console)")

		// config file
		configFile = flag.String("config", "", "Configuration file path")

		// misc
		version = flag.String("version", "", "Print the Terrascan version")

		// output type
    output = flag.String("output", "yaml", "Output format (supported values: json, xml, yaml, console)")
	)
	// override usage
	flag.Usage = Usage

	flag.Parse()

	// if no flags are passed, print usage
	if flag.NFlag() < 1 {
		flag.Usage()
		return
	}

	// if server mode set, run terrascan as a server, else run it as CLI
	if *server {
		logging.Init(*logType, *logLevel)
		httpServer.Start()
	} else {
		logging.Init(*logType, *logLevel)
		zap.S().Debug("running terrascan in cli mode")
		cli.Run(*iacType, *iacVersion, *cloudType, *iacFilePath, *iacDirPath, *configFile, *policyPath, *output, *version)
	}
}
