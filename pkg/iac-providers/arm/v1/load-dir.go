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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

const iacFile = "IAC file"

// LoadIacDir loads all ARM template files in the current directory.
func (a *ARMV1) LoadIacDir(absRootDir string, nonRecursive bool) (output.AllResourceConfigs, error) {
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, ARMFileExtensions())
	if err != nil {
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	for fileDir, files := range fileMap {
		for i := range files {
			// continue if file is a *.parameters.json or metadata.json
			if isParametersFile(*files[i]) || isMetadataFile(*files[i]) {
				continue
			}

			file := filepath.Join(fileDir, *files[i])

			// validate if the ARM template has a supporting valid
			// *.parameters.json file or not
			if !a.hasValidParametersFile(i, fileDir, files) {
				zap.S().Debug(errFileLoad,
					zap.String(iacFile, file),
					zap.Error(errors.New("does not have required valid parameters file")))
				continue
			}

			var configData output.AllResourceConfigs
			if configData, err = a.LoadIacFile(file); err != nil {
				zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, nil
}

func isParametersFile(file string) bool {
	return strings.Contains(file, "parameters.json")
}

func isMetadataFile(file string) bool {
	return strings.Contains(file, "metadata.json")
}

const errFileLoad = "error while loading iac files"

func (a *ARMV1) hasValidParametersFile(i int, fileDir string, files []*string) bool {
	f := strings.TrimSuffix(*files[i], filepath.Ext(*files[i]))
	for n := range files {
		if n == i {
			continue
		}

		if strings.EqualFold(*files[n], f+".parameters.json") {
			file := filepath.Join(fileDir, *files[n])
			f, err := os.Open(file)
			if err != nil {
				zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
				return false
			}
			defer f.Close()

			data, err := ioutil.ReadAll(f)
			if err != nil {
				zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
				return false
			}

			var params map[string]interface{}
			err = json.Unmarshal(data, &params)
			if err != nil {
				zap.S().Debug(errFileLoad, zap.String(iacFile, file), zap.Error(err))
				return false
			}
			npm, err := extractParameterValues(params)
			if err != nil {
				zap.S().Debug("error extracting parameter values", zap.String(iacFile, file), zap.Error(err))
				return false
			}
			a.templateParameters = npm
			return true
		}
	}
	return false
}

func extractParameterValues(params map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(params["parameters"])
	if err != nil {
		return nil, err
	}
	var npm map[string]struct {
		Value interface{} `json:"value"`
	}
	err = json.Unmarshal(data, &npm)
	if err != nil {
		return nil, err
	}

	finalParams := map[string]interface{}{}
	for key, value := range npm {
		finalParams[key] = value.Value
	}
	return finalParams, nil
}
