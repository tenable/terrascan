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

package dockerv1

import (
	"fmt"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

// LoadIacDir loads the docker file specified in given folder.
func (dc *DockerV1) LoadIacDir(absRootDir string, nonRecursive bool) (output.AllResourceConfigs, error) {
	// set the root directory being scanned
	dc.absRootDir = absRootDir

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	// find all the files in the folder with name `Dockerfile`
	fileMap, err := utils.FindFilesBySuffix(absRootDir, []string{DockerFileName})
	if err != nil {
		zap.S().Errorf("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(dc.errIacLoadDirs, results.DirScanErr{IacType: "docker", Directory: absRootDir, ErrMessage: err.Error()})
	}

	if len(fileMap) == 0 {
		zap.S().Warnf("directory '%s' has no files named Dockerfile. Use -f flag if Dockerfiles follow a different naming convention.", absRootDir)
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])
			var configData output.AllResourceConfigs
			if configData, err = dc.LoadIacFile(file); err != nil {
				errMsg := fmt.Sprintf("error while parsing file %s", file)
				zap.S().Errorf("error while searching for iac files", zap.String("root dir", absRootDir), errMsg)
				dc.errIacLoadDirs = multierror.Append(dc.errIacLoadDirs, results.DirScanErr{IacType: "docker", Directory: absRootDir, ErrMessage: errMsg})
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, dc.errIacLoadDirs

}
