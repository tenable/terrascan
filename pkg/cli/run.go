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
	"fmt"
	"os"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/pkg/writer"
	"go.uber.org/zap"
)

// Run executes terrascan in CLI mode
func Run(iacType, iacVersion, cloudType, iacFilePath, iacDirPath, configFile,
	policyPath, format, remoteType, remoteURL string, configOnly, useColors bool) {

	// download remote repository
	var tempDir string
	// temp dir to download the remote repo
	tempDir = filepath.Join(os.TempDir(), utils.GenRandomString(6))
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

// validateRemoteOpts validate remote repository options
func validateRemoteOpts(remoteType, remoteURL string) (bool, error) {

	// 1. remoteType and remoteURL both are empty
	if remoteType == "" && remoteURL == "" {
		return false, nil
	}

	// 2. remoteType and remoteURL both are not empty
	if remoteType != "" && remoteURL != "" {
		zap.S().Debugf("remoteType: %q, remoteURL: %q", remoteType, remoteURL)
		return true, nil
	}

	// 3. remoteType is empty and remoteURL is not
	if remoteType != "" || remoteURL != "" {
		zap.S().Errorf("remote type and remote url both options should be specified")
		return false, fmt.Errorf("incorrect remote options")
	}

	return false, nil
}

// downloadRemoteRepo downloads the remote repo in the temp directory and
// returns the path of the dir where the remote repository has been downloaded
func downloadRemoteRepo(remoteType, remoteURL, destDir string) (string, error) {

	// new downloader
	d := downloader.NewDownloader()
	url := fmt.Sprintf("%s::%s", remoteType, remoteURL)
	path, err := d.Download(url, destDir)
	if err != nil {
		zap.S().Errorf("failed to download remote repo url: %q, type: %q. error: '%v'",
			remoteURL, remoteType, err)
		return "", err
	}

	// successful
	return path, nil
}
