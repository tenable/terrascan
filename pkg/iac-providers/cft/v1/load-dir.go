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

package cftv1

import (
	"fmt"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

// LoadIacDir loads all CFT template files in the current directory.
func (a *CFTV1) LoadIacDir(absRootDir string, nonRecursive bool) (output.AllResourceConfigs, error) {
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, CFTFileExtensions())
	if err != nil {
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "cft", Directory: absRootDir, ErrMessage: err.Error()})
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])

			var configData output.AllResourceConfigs
			if configData, err = a.LoadIacFile(file); err != nil {
				errMsg := fmt.Sprintf("error while loading iac file '%s', err: %v", file, err)
				zap.S().Debug("error while loading iac files", zap.String("IAC file", file), zap.Error(err))
				a.errIacLoadDirs = multierror.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "cft", Directory: fileDir, ErrMessage: errMsg})
				continue
			}

			for key := range configData {
				resourceConfigs := configData[key]
				makeSourcePathRelative(absRootDir, resourceConfigs)
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, a.errIacLoadDirs
}

// makeSourcePathRelative modifies the source path of each resource from absolute to relative path
func makeSourcePathRelative(absRootDir string, resourceConfigs []output.ResourceConfig) {
	for i := range resourceConfigs {
		r := &resourceConfigs[i]
		var err error
		oldSource := r.Source
		// update the source path
		r.Source, err = filepath.Rel(absRootDir, r.Source)
		// though this error should never occur, but, if occurs for some reason, assign the old value of source back
		if err != nil {
			r.Source = oldSource
			zap.S().Debug("error while getting the relative path for", zap.String("IAC file", oldSource), zap.Error(err))
			continue
		}
	}
}
