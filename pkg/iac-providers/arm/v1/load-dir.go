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

package armv1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

const iacFile = "IAC file"

// LoadIacDir loads all ARM template files in the current directory.
func (a *ARMV1) LoadIacDir(absRootDir string, nonRecursive bool) (output.AllResourceConfigs, error) {
	// set the root directory being scanned
	a.absRootDir = absRootDir

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, ARMFileExtensions())
	if err != nil {
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "arm", Directory: absRootDir, ErrMessage: err.Error()})
	}

	for fileDir, files := range fileMap {
		for i := range files {
			// continue if file is a *.parameters.json or metadata.json
			if files[i] != nil && isParametersFile(*files[i]) || isMetadataFile(*files[i]) {
				continue
			}

			file := filepath.Join(fileDir, *files[i])

			// check if the template has a supporting .parameters.json file or not
			// yes: extract parameter values; no: continue with the default values set in the template
			a.templateParameters = make(map[string]interface{})
			a.tryGetParameters(*files[i], fileDir, files)

			var configData output.AllResourceConfigs
			if configData, err = a.LoadIacFile(file); err != nil {
				errMsg := fmt.Sprintf("error while loading iac file '%s'. err: %v", file, err)
				zap.S().Debug("error while loading iac files", zap.String("IAC file", file), zap.Error(err))
				a.errIacLoadDirs = multierror.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "arm", Directory: fileDir, ErrMessage: errMsg})
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}
	return allResourcesConfig, a.errIacLoadDirs
}

func isParametersFile(file string) bool {
	return strings.Contains(file, ParametersFileExtension)
}

func isMetadataFile(file string) bool {
	return strings.Contains(file, MetadataFileExtension)
}

const errFileLoad = "error while loading iac files"

func (a *ARMV1) tryGetParameters(fileName string, fileDir string, files []*string) {
	pf := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ParametersFileExtension

	file := filepath.Join(fileDir, pf)
	f, err := os.Open(file)
	if err != nil {
		zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
		return
	}

	var params map[string]interface{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
		return
	}

	err = a.extractParameterValues(params)
	if err != nil {
		zap.S().Debug("error extracting parameter values", zap.String(iacFile, file), zap.Error(err))
		return
	}
}

func (a *ARMV1) extractParameterValues(params map[string]interface{}) error {
	data, err := json.Marshal(params["parameters"])
	if err != nil {
		return err
	}
	var npm map[string]struct {
		Value interface{} `json:"value"`
	}
	err = json.Unmarshal(data, &npm)
	if err != nil {
		return err
	}

	for key, value := range npm {
		a.templateParameters[key] = value.Value
	}
	return nil
}
