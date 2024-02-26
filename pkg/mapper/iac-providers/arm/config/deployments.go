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

package config

import (
	"encoding/json"

	"github.com/tenable/terrascan/pkg/mapper/convert"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
	"go.uber.org/zap"
)

const (
	armTemplateLink   = "templateLink"
	armTemplate       = "template"
	armParametersLink = "parametersLink"
	armParameters     = "parameters"
	armRelativePath   = types.LinkedTemplateRelativePath
	armURI            = "uri"
)

const (
	tfParametersContent = types.LinkedParametersContent
	tfTemplateContent   = types.LinkedTemplateContent
)

const (
	errParameters = "unable to load linked template parameters"
	errTemplate   = "unable to load linked template data"
)

// DeploymentsConfig returns config for azurerm_resource_group_template_deployment
func DeploymentsConfig(r types.Resource, vars, params map[string]interface{}) map[string]interface{} {
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(vars, params, r.Location),
		tfName:     fn.LookUpString(vars, params, r.Name),
		tfTags:     functions.PatchAWSTags(r.Tags),
	}

	// if template is defined directly
	if template := convert.ToMap(r.Properties, armTemplate); template != nil {
		templateContent, err := json.Marshal(template)
		if err != nil {
			zap.S().Debug(errTemplate, zap.String("resource", r.Type),
				zap.String("name", convert.ToString(cf, tfName)))
		}
		cf[tfTemplateContent] = templateContent
	} else if templateLink := convert.ToMap(r.Properties, armTemplateLink); templateLink != nil {
		// else if template linked is defined

		// if templateLink is relative path, resolve it later as it need the absPath for current template
		if relativePath := convert.ToString(templateLink, armRelativePath); relativePath != "" {
			cf[armRelativePath] = relativePath
		} else if templateURI := convert.ToString(templateLink, armURI); templateURI != "" {
			// if templateLink has a uri
			templateURI := fn.LookUpString(vars, params, templateURI)
			templateContent, err := fn.ResolveLinkedTemplate(templateURI)
			if err != nil {
				zap.S().Debug(errTemplate, zap.String("resource", r.Type), zap.String("uri", templateURI))
			}
			cf[tfTemplateContent] = templateContent
		}
	}

	// if parameters are defined directly
	if parameters := convert.ToMap(r.Properties, armParameters); parameters != nil {
		parametersContent, err := json.Marshal(parameters)
		if err != nil {
			zap.S().Debug(errParameters, zap.String("resource", r.Type),
				zap.String("name", convert.ToString(cf, tfName)),
			)
		}
		cf[tfParametersContent] = parametersContent
	} else if parametersLink := convert.ToMap(r.Properties, armParametersLink); parametersLink != nil {
		if parametersURI := convert.ToString(parametersLink, armURI); parametersURI != "" {
			// if parametersLink has a uri
			parametersURI = fn.LookUpString(vars, params, parametersURI)
			parametersContent, err := fn.ResolveLinkedTemplate(parametersURI)
			if err != nil {
				zap.S().Debug(errParameters, zap.String("resource", r.Type), zap.String("uri", parametersURI))
			}
			cf[tfParametersContent] = parametersContent
		}
	}

	return cf
}
