package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/cli"
	httpServer "github.com/accurics/terrascan/pkg/http-server"
	"github.com/accurics/terrascan/pkg/logger"
)

func main() {
	var (
		server      = flag.Bool("server", false, "run terrascan in server mode")
		iacType     = flag.String("iac", "", "IaC provider (supported values: terraform)")
		iacVersion  = flag.String("iac-version", "default", "IaC version (supported values: 'v12' for terraform)")
		cloudType   = flag.String("cloud", "", "cloud provider (supported values: aws)")
		iacFilePath = flag.String("f", "", "IaC file path")
		iacDirPath  = flag.String("d", "", "IaC directory path")

		// logging flags
		logLevel = flag.String("log-level", "info", "logging level (debug, info, warn, error, panic, fatal)")
		logType  = flag.String("log-type", "console", "log type (json, console)")
	)
	flag.Parse()

	// if server mode set, run terrascan as a server, else run it as CLI
	if *server {
		logger.Init(*logType, *logLevel)
		httpServer.Start()
	} else {
		logger.Init(*logType, *logLevel)
		zap.S().Debug("running terrascan in cli mode")
		cli.Run(*iacType, *iacVersion, *cloudType, *iacFilePath, *iacDirPath)
	}
}
