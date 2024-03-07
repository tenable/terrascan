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

package cftv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws-cloudformation/rain/cft/parse"
	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation"
	multierr "github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/mapper"
	cftRes "github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/config"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/store"
	"github.com/tenable/terrascan/pkg/results"
	"go.uber.org/zap"
)

// LoadIacFile loads the specified CFT template file.
// Note that a single CFT template json file may contain multiple resource definitions.
func (a *CFTV1) LoadIacFile(absFilePath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	fileData, err := os.ReadFile(absFilePath)
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

	template, err := parse.File(absFilePath)
	if err != nil {
		zap.S().Debug("unable to parse template for getting line numbers of resources", zap.Error(err), zap.String("file", absFilePath))
	}

	// fill AllResourceConfigs
	allResourcesConfig = make(map[string][]output.ResourceConfig)
	var config *output.ResourceConfig
	for _, resource := range configs {
		config = &resource

		// Fetch line number
		resNode, err := template.GetResource(resource.Name)
		if err != nil {
			zap.S().Debug("unable to get line number of resource", zap.Error(err), zap.String("file", absFilePath),
				zap.String("resource", resource.Name))
		}
		if resNode != nil && resNode.Line > 1 {
			config.Line = resNode.Line

			// If yaml, adjust line number
			// resNode.Line points to first line within the resource
			extension := filepath.Ext(absFilePath)
			if extension == ".yaml" || extension == ".yml" {
				config.Line--
			}
		} else {
			// default
			config.Line = 1
		}

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
	var isYaml bool

	switch fileExt {
	case YAMLExtension, YAMLExtension2:
		isYaml = true
	case JSONExtension:
		isYaml = false
	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return nil, fmt.Errorf("unsupported extension for file %s", file)
	}

	zap.S().Debug("sanitizing cft template file", zap.String("file", file))
	sanitized, err := a.sanitizeCftTemplate(file, *data, isYaml)
	if err != nil {
		zap.S().Debug("failed to sanitize cft template file", zap.String("file", file), zap.Error(err))
		return nil, err
	}

	return a.cleanTemplate(sanitized, file)
}

func (a *CFTV1) cleanTemplate(templateMap map[string]interface{}, absFilePath string) (*cloudformation.Template, error) {
	var onetemplate cloudformation.Template

	resourceMap, ok := templateMap["Resources"].(map[string]interface{})
	if !ok {
		zap.S().Debug("failed to find valid Resources key", zap.String("file", absFilePath))
		return nil, errors.New("failed to find valid Resources key in file: " + absFilePath)
	}

	onetemplate.Resources = make(cloudformation.Resources, len(resourceMap))

	for resourceName := range resourceMap {
		var resourceInfo cftResource

		resourceInfo.Resources = make(map[string]interface{}, 1)
		resourceInfo.Resources[resourceName] = resourceMap[resourceName]

		resourceData, err := json.Marshal(resourceInfo)
		if err != nil {
			zap.S().Debug("failed to marshal json for resource", zap.String("resource", resourceName), zap.Error(err))
			a.errIacLoadDirs = multierr.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "cft", Directory: filepath.Dir(absFilePath), ErrMessage: err.Error()})
			continue
		}

		template, err := goformation.ParseJSON(resourceData)
		if err != nil {
			zap.S().Debug("failed to generate template for resource", zap.String("resource", resourceName), zap.Error(err))
			a.errIacLoadDirs = multierr.Append(a.errIacLoadDirs, results.DirScanErr{IacType: "cft", Directory: filepath.Dir(absFilePath), ErrMessage: err.Error()})
			continue
		}

		for key := range template.Resources {
			onetemplate.Resources[key] = template.Resources[key]
		}
	}

	return &onetemplate, nil
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
		return TXTExtension
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
