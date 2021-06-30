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

package config

import (
	"encoding/json"

	"github.com/accurics/terrascan/pkg/mapper/convert"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
	"go.uber.org/zap"
)

const (
	armTemplateLink   = "templateLink"
	armTemplate       = "template"
	armParametersLink = "parametersLink"
	armParameters     = "parameters"
	armRelativePath   = types.LinkedTemplateRelativePath
	armUri            = "uri"
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
	// TODO: check if vars are handled correctly
	cf := map[string]interface{}{
		tfLocation: fn.LookUpString(vars, params, r.Location),
		tfName:     fn.LookUpString(vars, params, r.Name),
		tfTags:     r.Tags,
	}

	// if template is defiened directly
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
		} else if templateUri := convert.ToString(templateLink, armUri); templateUri != "" {
			// if templateLink has a uri
			templateUri := fn.LookUpString(vars, params, templateUri)
			templateContent, err := fn.ResolveLinkedTemplate(templateUri)
			if err != nil {
				zap.S().Debug(errTemplate, zap.String("resource", r.Type), zap.String("uri", templateUri))
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
	} else if paramtersLink := convert.ToMap(r.Properties, armParametersLink); paramtersLink != nil {
		if parametersUri := convert.ToString(paramtersLink, armUri); parametersUri != "" {
			// if paramtersLink has a uri
			parametersUri = fn.LookUpString(vars, params, parametersUri)
			parametersContent, err := fn.ResolveLinkedTemplate(parametersUri)
			if err != nil {
				zap.S().Debug(errParameters, zap.String("resource", r.Type), zap.String("uri", parametersUri))
			}
			cf[tfParametersContent] = parametersContent
		}
	}

	return cf
}
