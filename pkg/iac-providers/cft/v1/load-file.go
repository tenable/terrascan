/*
    Copyright (C) 2021 Accurics, Inc.

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/mapper"
	cftRes "github.com/accurics/terrascan/pkg/mapper/iac-providers/cft/config"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/cft/store"
	"github.com/awslabs/goformation/v4"
	"github.com/awslabs/goformation/v4/cloudformation"
	"go.uber.org/zap"
)

// LoadIacFile loads the specified CFT template file.
// Note that a single CFT template json file may contain multiple resource definitions.
func (a *CFTV1) LoadIacFile(absFilePath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	fileData, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		zap.S().Debug("unable to read file", zap.Error(err), zap.String("file", absFilePath))
		return allResourcesConfig, fmt.Errorf("unable to read file %s", absFilePath)
	}

	// parse the file as cloudformation.Template
	configs, err := a.getConfig(absFilePath, &fileData, nil)
	if err != nil {
		zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
		return allResourcesConfig, err
	}

	// fill AllResourceConfigs
	allResourcesConfig = make(map[string][]output.ResourceConfig)
	var config *output.ResourceConfig
	for _, resource := range configs {
		config = &resource
		config.Line = 1
		if config.Source == "" {
			config.Source = a.getSourceRelativePath(absFilePath)
		}
		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}
	return allResourcesConfig, nil
}

func (a *CFTV1) getConfig(absFilePath string, fileData *[]byte, parameters *map[string]string) ([]output.ResourceConfig, error) {
	// parse the file as cloudformation.Template
	template, err := a.extractTemplate(absFilePath, fileData)
	if err != nil {
		return nil, err
	}

	// replace template parameter values
	if parameters != nil {
		for key, value := range *parameters {
			if parameter, ok := template.Parameters[key]; ok {
				parameter.Default = value
				template.Parameters[key] = parameter
			}
		}
	}

	// map resource to a terrascan type
	configs, err := a.translateResources(template, absFilePath)
	if err != nil {
		zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
		return nil, err
	}
	return configs, nil
}

func (a *CFTV1) extractTemplate(file string, data *[]byte) (*cloudformation.Template, error) {
	fileExt := a.getFileType(file, data)

	switch fileExt {
	case YAMLExtension, YAMLExtension2:
		zap.S().Debug("sanitizing cft template file", zap.String("file", file))
		sanitized, err := a.sanitizeCftTemplate(*data, true)
		if err != nil {
			zap.S().Debug("failed to sanitize cft template file", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		template, err := goformation.ParseYAML(sanitized)
		if err != nil {
			zap.S().Debug("failed to parse file", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		return template, nil
	case JSONExtension:
		zap.S().Debug("sanitizing cft template file", zap.String("file", file))
		sanitized, err := a.sanitizeCftTemplate(*data, false)
		if err != nil {
			zap.S().Debug("failed to sanitize cft template file", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		template, err := goformation.ParseJSON(sanitized)
		if err != nil {
			zap.S().Debug("failed to parse file", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		return template, nil
	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return nil, fmt.Errorf("unsupported extension for file %s", file)
	}
}

func (a *CFTV1) translateResources(template *cloudformation.Template, absFilePath string) ([]output.ResourceConfig, error) {
	m := mapper.NewMapper("cft")
	configs, err := m.Map(template)
	if err != nil {
		zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
		return nil, err
	}

	for _, config := range configs {
		if config.Type == store.AwsCloudFormationStack {
			if stackConfig, ok := config.Config.(cftRes.CloudFormationStackConfig); ok {
				if stackConfig.TemplateData != nil {
					stackResourceConfigs, err := a.getConfig(stackConfig.TemplateURL, &stackConfig.TemplateData, &stackConfig.Parameters)
					if err == nil {
						for i := range stackResourceConfigs {
							// Add template url as source for the nested resources
							stackResourceConfigs[i].Source = stackConfig.TemplateURL
						}
						configs = append(configs, stackResourceConfigs...)
					}
				}
			}
		}
	}

	return configs, nil
}

func (*CFTV1) getFileType(file string, data *[]byte) string {
	if strings.HasSuffix(file, YAMLExtension) {
		return YAMLExtension
	} else if strings.HasSuffix(file, YAMLExtension2) {
		return YAMLExtension2
	} else if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	} else if strings.HasSuffix(file, TXTExtension) || strings.HasSuffix(file, TemplateExtension) {
		if isJSON(string(*data)) {
			return JSONExtension
		}
		return YAMLExtension
	}
	return UnknownExtension
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// getSourceRelativePath fetches the relative path of file being loaded
func (a *CFTV1) getSourceRelativePath(sourceFile string) string {

	// rootDir should be empty when file scan was initiated by user
	if a.absRootDir == "" {
		return filepath.Base(sourceFile)
	}
	relPath, err := filepath.Rel(a.absRootDir, sourceFile)
	if err != nil {
		zap.S().Debug("error while getting the relative path for", zap.String("IAC file", sourceFile), zap.Error(err))
		return sourceFile
	}
	return relPath
}
