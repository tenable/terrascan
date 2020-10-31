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
	"path/filepath"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/pkg/writer"
	"go.uber.org/zap"
)

// Run executes terrascan in CLI mode
func Run(iacType, iacVersion string, cloudType []string,
	iacFilePath, iacDirPath, configFile string, policyPath []string,
	format, remoteType, remoteURL string, configOnly, useColors bool) {

	// temp dir to download the remote repo
	tempDir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
	defer os.RemoveAll(tempDir)

	// download remote repository
	d := downloader.NewDownloader()
	path, err := d.DownloadWithType(remoteType, remoteURL, tempDir)
	if err == downloader.ErrEmptyURLType {
		// url and type empty, proceed with regular scanning
		zap.S().Debugf("remote url and type not configured, proceeding with regular scanning")
	} else if err != nil {
		// some error while downloading remote repository
		return
	} else {
		// successfully downloaded remote repository
		iacDirPath = path
	}

	// create a new runtime executor for processing IaC
	executor, err := runtime.NewExecutor(iacType, iacVersion, cloudType,
		iacFilePath, iacDirPath, configFile, policyPath)
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
		os.RemoveAll(tempDir)
		os.Exit(3)
	}
}
