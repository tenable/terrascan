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
	"github.com/accurics/terrascan/pkg/initialize"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Terrascan",
	Long: `Terrascan

Initializes Terrascan and clones policies from the Terrascan GitHub repository.
`,
	Run: initial,
}

func initial(cmd *cobra.Command, args []string) {
	// initialize terrascan
	if err := initialize.Run(); err != nil {
		zap.S().Error("failed to initialize terrascan")
		return
	}
}

func init() {
	RegisterCommand(rootCmd, initCmd)
}
