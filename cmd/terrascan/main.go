package main

import (
	"flag"

	"github.com/accurics/terrascan/pkg/cli"
	httpServer "github.com/accurics/terrascan/pkg/http-server"
)

func main() {
	var (
		server      = flag.Bool("server", false, "run terrascan in server mode")
		iacType     = flag.String("iac", "", "IaC provider (supported values: terraform)")
		iacVersion  = flag.String("iac-version", "default", "IaC version (supported values: 'v12' for terraform)")
		cloudType   = flag.String("cloud", "", "cloud provider (supported values: aws)")
		iacFilePath = flag.String("f", "", "IaC file path")
		iacDirPath  = flag.String("d", "", "IaC directory path")
	)
	flag.Parse()

	// if server mode set, run terrascan as a server, else run it as CLI
	if *server {
		httpServer.Start()
	} else {
		cli.Run(*iacType, *iacVersion, *cloudType, *iacFilePath, *iacDirPath)
	}
}
