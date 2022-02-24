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
	"github.com/awslabs/goformation/v5"
	"github.com/awslabs/goformation/v5/cloudformation"
	"github.com/ghodss/yaml"
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
	var err error
	var sanitized []byte

	switch fileExt {
	case YAMLExtension, YAMLExtension2:
		zap.S().Debug("sanitizing cft template file", zap.String("file", file))
		sanitizedYaml, err := a.sanitizeCftTemplate(*data, true)
		if err != nil {
			zap.S().Debug("failed to sanitize cft template file", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		sanitized, err = yaml.YAMLToJSON(sanitizedYaml)
		if err != nil {
			zap.S().Debug("invalid yaml", zap.String("file", file), zap.Error(err))
			return nil, err
		}

	case JSONExtension:
		zap.S().Debug("sanitizing cft template file", zap.String("file", file))
		sanitized, err = a.sanitizeCftTemplate(*data, false)
		if err != nil {
			zap.S().Debug("failed to sanitize cft template file", zap.String("file", file), zap.Error(err))
			return nil, err
		}

	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return nil, fmt.Errorf("unsupported extension for file %s", file)
	}

	resourcesList, err := getResourcesList(sanitized)
	if err != nil {
		zap.S().Debug("failed to unmarshal sanitized json", zap.String("file", file), zap.Error(err))
		return nil, err
	}

	var onetemplate cloudformation.Template
	onetemplate.Resources = make(map[string]cloudformation.Resource, 1)

	for i := range resourcesList {
		var resourceName string
		for key := range resourcesList[i].Resources {
			resourceName = key
		}

		resourceData, err := json.Marshal(resourcesList[i])
		if err != nil {
			zap.S().Debug("failed to marshal json for resource", zap.String("resource", resourceName), zap.Error(err))
			continue
		}

		template, err := goformation.ParseJSON(resourceData)
		if err != nil {
			zap.S().Debug("failed to generate template for resource", zap.String("resource", resourceName), zap.Error(err))
			continue
		}

		onetemplate.AWSTemplateFormatVersion = template.AWSTemplateFormatVersion
		for key := range template.Resources {
			onetemplate.Resources[key] = template.Resources[key]
		}
	}

	return &onetemplate, nil
}

const (
	AWSTemplateFormatVersion = "AWSTemplateFormatVersion"
	Resources                = "Resources"
)

type cftResource struct {
	AWSTemplateFormatVersion string                 `json:"AWSTemplateFormatVersion"`
	Resources                map[string]interface{} `json:"Resources"`
}

func getResourcesList(sanitized []byte) ([]cftResource, error) {
	var err error
	var jsonMap map[string]interface{}
	var resourcesList []cftResource

	err = json.Unmarshal(sanitized, &jsonMap)
	if err != nil {
		return nil, err
	}

	resourceMap := jsonMap[Resources].(map[string]interface{})
	for key := range resourceMap {
		var resourceInfo cftResource
		resourceInfo.AWSTemplateFormatVersion = jsonMap[AWSTemplateFormatVersion].(string)
		resourceInfo.Resources = make(map[string]interface{}, 1)
		resourceInfo.Resources[key] = resourceMap[key]

		resourcesList = append(resourcesList, resourceInfo)
	}

	return resourcesList, nil
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
