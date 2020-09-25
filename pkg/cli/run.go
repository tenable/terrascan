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
	"flag"
	"os"

	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/writer"
)

// Run executes terrascan in CLI mode
func Run(iacType, iacVersion, cloudType, iacFilePath, iacDirPath, configFile,
	policyPath, format string, configOnly bool, useColors bool) {

	// create a new runtime executor for processing IaC
	executor, err := runtime.NewExecutor(iacType, iacVersion, cloudType, iacFilePath,
		iacDirPath, configFile, policyPath)
	if err != nil {
		return
	}

	// executor output
	results, err := executor.Execute()
	if err != nil {
		return
	}

    outputWriter := NewOutputWriter(useColors)

	if configOnly {
		writer.Write(format, results.ResourceConfig, outputWriter)
	} else {
		writer.Write(format, results.Violations, outputWriter)
	}

	if results.Violations.ViolationStore.Count.TotalCount != 0 && flag.Lookup("test.v") == nil {
		os.Exit(3)
	}
}
