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

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/cli"
	httpServer "github.com/accurics/terrascan/pkg/http-server"
	"github.com/accurics/terrascan/pkg/logging"
)

func main() {

	// command line flags
	var (
		// server mode
		server = flag.Bool("server", false, "run terrascan in server mode")

		// IaC flags
		iacType     = flag.String("iac", "", "IaC provider (supported values: terraform)")
		iacVersion  = flag.String("iac-version", "default", "IaC version (supported values: 'v12' for terraform)")
		iacFilePath = flag.String("f", "", "IaC file path")
		iacDirPath  = flag.String("d", "", "IaC directory path")
		policyPath  = flag.String("p", "", "Policy directory path")

		// cloud flags
		cloudType = flag.String("cloud", "", "cloud provider (supported values: aws)")

		// logging flags
		logLevel = flag.String("log-level", "info", "logging level (debug, info, warn, error, panic, fatal)")
		logType  = flag.String("log-type", "console", "log type (json, console)")

		// config file
		configFile = flag.String("config", "", "config file path")
	)
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
		cli.Run(*iacType, *iacVersion, *cloudType, *iacFilePath, *iacDirPath, *configFile, *policyPath)
	}
}
