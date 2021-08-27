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
	//"fmt"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes Terrascan and clones policies from the Terrascan GitHub repository.",
	Long: `Terrascan

Initializes Terrascan and clones policies from the Terrascan GitHub repository.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("test elsewhere" + initialize.InitializeTag)

		err := initial(cmd, args, false)
		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func initial(cmd *cobra.Command, args []string, nonInitCmd bool)  error {
	// initialize terrascan
	//fmt.Println("test initial" + initialize.InitializeTag)

	err := initialize.Run(nonInitCmd);
	//fmt.Println("dummy" + release)
	//fmt.Println("Release from initial " + Release )
	if err != nil {
		zap.S().Errorf("failed to initialize terrascan. error : %v", err)
		return err
	}
	return nil
}

func init() {
	RegisterCommand(rootCmd, initCmd)
	
}
