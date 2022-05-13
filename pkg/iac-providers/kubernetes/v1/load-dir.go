/*
    Copyright (C) 2022 Tenable, Inc.

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

package k8sv1

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
)

func (*K8sV1) getFileType(file string) string {
	if strings.HasSuffix(file, YAMLExtension) {
		return YAMLExtension
	} else if strings.HasSuffix(file, YAMLExtension2) {
		return YAMLExtension2
	} else if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}

// LoadIacDir loads all k8s files in the current directory
func (k *K8sV1) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	// set the root directory being scanned
	k.absRootDir = absRootDir

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, K8sFileExtensions())
	if err != nil {
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(k.errIacLoadDirs, results.DirScanErr{IacType: "k8s", Directory: absRootDir, ErrMessage: err.Error()})
	}
	if len(fileMap) == 0 {
		errMsg := fmt.Sprintf("kubernetes files not found in the directory %s", k.absRootDir)
		return allResourcesConfig, multierror.Append(k.errIacLoadDirs, results.DirScanErr{IacType: "k8s", Directory: k.absRootDir, ErrMessage: errMsg})
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])

			var configData output.AllResourceConfigs
			if configData, err = k.LoadIacFile(file, options); err != nil {
				errMsg := fmt.Sprintf("error while loading iac file '%s'. err: %v", file, err)
				zap.S().Debug("error while loading iac files", zap.String("IAC file", file), zap.Error(err))
				k.errIacLoadDirs = multierror.Append(k.errIacLoadDirs, results.DirScanErr{IacType: "k8s", Directory: fileDir, ErrMessage: errMsg})
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, k.errIacLoadDirs
}

// Name returns name of the provider
func (*K8sV1) Name() string {
	return kubernetesTypeNameShort
}
