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
	httpserver "github.com/accurics/terrascan/pkg/http-server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Terrascan as an API server",
	Long: `Terrascan

Run Terrascan as an API server that inspects incoming IaC (Infrastructure-as-Code) files and returns the scan results.
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initial(cmd, args)
	},
	Run: server,
}

func server(cmd *cobra.Command, args []string) {
	httpserver.Start()
}

func init() {
	RegisterCommand(rootCmd, serverCmd)
}
