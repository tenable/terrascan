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

package armv1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/mapper"
	"github.com/tenable/terrascan/pkg/mapper/convert"
	"github.com/tenable/terrascan/pkg/mapper/core"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"go.uber.org/zap"
)

// LoadIacFile loads the specified ARM template file.
// Note that a single ARM template json file may contain multiple resource definitions.
func (a *ARMV1) LoadIacFile(absFilePath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(output.AllResourceConfigs)
	if fileExt := a.getFileType(absFilePath); fileExt != JSONExtension {
		return allResourcesConfig, fmt.Errorf("unsupported file %s", absFilePath)
	}

	fileData, err := os.ReadFile(absFilePath)
	if err != nil {
		zap.S().Debug("unable to read file", zap.Error(err), zap.String("file", absFilePath))
		return allResourcesConfig, fmt.Errorf("unable to read file %s", absFilePath)
	}

	template, err := a.extractTemplate(fileData)
	if err != nil {
		zap.S().Debug("unable to parse template", zap.Error(err), zap.String("file", absFilePath))
		return allResourcesConfig, fmt.Errorf("unable to parse file %s", absFilePath)
	}
	if resConfs := a.translateResources(template, absFilePath); resConfs != nil {
		a.addConfig(allResourcesConfig, resConfs)
	}

	return allResourcesConfig, nil
}

func (a *ARMV1) translateResources(template *types.Template, absFilePath string) []output.ResourceConfig {
	mapper := mapper.NewMapper("arm")
	var allResourcesConfig = make([]output.ResourceConfig, 0)

	// set template parameters with default values if not found
	if a.templateParameters == nil {
		a.templateParameters = make(map[string]interface{})
	}
	for key, param := range template.Parameters {
		if _, ok := a.templateParameters[key]; !ok {
			a.templateParameters[key] = param.DefaultValue
		}
	}

	for _, r := range template.Resources {
		configs := a.getConfig(absFilePath, mapper, r, template.Variables)
		for _, config := range configs {
			_, ok := config.Config.(map[string]interface{})
			if !ok {
				zap.S().Debug("unable to parse config.Config data",
					zap.String("resource", r.Type), zap.String("file", absFilePath),
				)
				continue
			}

			for _, nr := range r.Resources {
				if !strings.HasPrefix(nr.Type, "Microsoft.") {
					nr.Type = r.Type + "/" + nr.Type
				}
				resourceConfigs := a.getConfig(absFilePath, mapper, nr, template.Variables)
				allResourcesConfig = append(allResourcesConfig, resourceConfigs...)
			}
		}
		allResourcesConfig = append(allResourcesConfig, configs...)
	}
	return allResourcesConfig
}

func (ARMV1) getFileType(file string) string {
	if ext := filepath.Ext(file); strings.EqualFold(ext, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}

func (ARMV1) extractTemplate(data []byte) (*types.Template, error) {
	var t types.Template
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (ARMV1) addConfig(a output.AllResourceConfigs, configs []output.ResourceConfig) {
	for _, config := range configs {
		if _, present := a[config.Type]; !present {
			a[config.Type] = []output.ResourceConfig{config}
		} else {
			resources := a[config.Type]
			if !output.IsConfigPresent(resources, config) {
				a[config.Type] = append(a[config.Type], config)
			}
		}
	}
}

// getSourceRelativePath fetches the relative path of file being loaded
func (a *ARMV1) getSourceRelativePath(sourceFile string) string {
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

func (a *ARMV1) getConfig(path string, mapper core.Mapper, r types.Resource,
	vars map[string]interface{}) []output.ResourceConfig {

	if _, ok := types.ResourceTypes[r.Type]; !ok {
		return nil
	}

	configs, err := mapper.Map(r, vars, a.templateParameters)
	for i := range configs {
		configs[i].Source = a.getSourceRelativePath(path)
		configs[i].Line = 1
	}

	if err != nil {
		zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", path))
		return nil
	}

	// parse linked templates and translate resources
	for _, config := range configs {
		if linkedTemplate, templatePath := a.getLinkedTemplate(config, path, mapper, vars); linkedTemplate != nil {
			if templatePath != "" {
				return a.translateResources(linkedTemplate, templatePath)
			}
			return a.translateResources(linkedTemplate, path)
		}
	}

	return configs
}

func (a *ARMV1) getLinkedTemplate(config output.ResourceConfig, path string, mapper core.Mapper, vars map[string]interface{}) (*types.Template, string) {

	if config.Type == types.AzureRMDeployments {

		var templateData []byte
		var templateSource string
		var templateParameters map[string]struct {
			Value interface{} `json:"value"`
		}

		// get templateData from config
		if resourceConfig, ok := config.Config.(map[string]interface{}); ok {
			// if linked template is relative path
			if relativePath := convert.ToString(resourceConfig, types.LinkedTemplateRelativePath); relativePath != "" {
				templatePath := filepath.Join(filepath.Dir(path), relativePath)
				data, err := os.ReadFile(templatePath)
				if err != nil {
					zap.S().Debug("error loading linked template", zap.String("path", relativePath), zap.Error(err))
				}
				templateSource = a.getSourceRelativePath(templatePath)
				templateData = data
			} else if templateContent, ok := resourceConfig[types.LinkedTemplateContent]; ok {
				data, ok := templateContent.([]byte)
				if !ok {
					zap.S().Debug("error loading linked template", zap.String("resource", config.ID))
				}
				templateSource = a.getSourceRelativePath(path)
				templateData = data
			}

			// get parameters
			if parametersContent, ok := resourceConfig[types.LinkedParametersContent]; ok {
				parameters, ok := parametersContent.([]byte)
				if ok {
					err := json.Unmarshal(parameters, &templateParameters)
					if err != nil {
						zap.S().Debug("error loading linked template parameters", zap.String("resource", config.ID))
					}
				}
			}
		}

		if len(templateData) != 0 {
			// parse linked template
			linkedTemplate, err := a.extractTemplate(templateData)
			if err != nil {
				zap.S().Debug("unable to parse template", zap.Error(err), zap.String("file", path))
				return nil, path
			}

			// propagate parameters
			for key, param := range linkedTemplate.Parameters {
				if _, ok := a.templateParameters[key]; !ok {
					a.templateParameters[key] = param.DefaultValue
				}
			}

			// add values provided for linked templates
			for key, value := range templateParameters {
				if parameterValue, ok := value.Value.(string); ok {
					val := fn.LookUp(vars, a.templateParameters, parameterValue)
					switch val := val.(type) {
					case string, float64, bool:
						a.templateParameters[key] = val
					default:
					}
				} else {
					a.templateParameters[key] = value.Value
				}
			}

			// propagate template variables
			if linkedTemplate.Variables == nil {
				linkedTemplate.Variables = make(map[string]interface{})
			}
			for key, value := range vars {
				if varValue, ok := value.(string); ok {
					val := fn.LookUp(vars, a.templateParameters, varValue)
					switch val := val.(type) {
					case string, float64, bool:
						linkedTemplate.Variables[key] = val
					default:
					}
				} else {
					linkedTemplate.Variables[key] = value
				}
			}

			return linkedTemplate, templateSource
		}
	}
	return nil, ""
}
