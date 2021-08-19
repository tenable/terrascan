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

package commons

import (
	"fmt"
	"io/ioutil"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"go.uber.org/zap"
)

// CreateResourceConfig creates output.ResourceConfig
func CreateResourceConfig(managedResource *hclConfigs.Resource, reqdProviderNameMapping map[string]ResourceMetadata) (resourceConfig output.ResourceConfig, err error) {

	// read source file
	fileBytes, err := ioutil.ReadFile(managedResource.DeclRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to read terraform file")
	}

	// convert resource config from hcl.Body to map[string]interface{}
	c := converter{bytes: fileBytes}
	var hclBody *hclsyntax.Body
	var ok bool
	if hclBody, ok = managedResource.Config.(*hclsyntax.Body); !ok {
		return resourceConfig, fmt.Errorf("failed type assertion for hcl.Body in *hclConfigs.Resource. error: expected hcl.Body type is *hclsyntax.Body, but got %T", managedResource.Config)
	}

	conatiners := []output.ContainerNameAndImage{}
	initContainers := []output.ContainerNameAndImage{}

	if isEligibleForContainerImageExtraction(managedResource, reqdProviderNameMapping) {
		conatiners, initContainers = extractContainerImages(managedResource, hclBody)
	}

	goOut, err := c.convertBody(hclBody)
	if err != nil {
		zap.S().Errorf("failed to convert hcl.Body to go struct; resource '%s', file: '%s'. error: '%v'",
			managedResource.Name, managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to convert hcl.Body to go struct")
	}

	minSeverity, maxSeverity := utils.GetMinMaxSeverity(c.rangeSource(hclBody.Range()))
	// create a resource config
	resourceConfig = output.ResourceConfig{
		ID:                  fmt.Sprintf("%s.%s", managedResource.Type, managedResource.Name),
		Name:                managedResource.Name,
		Type:                managedResource.Type,
		Source:              managedResource.DeclRange.Filename,
		Line:                managedResource.DeclRange.Start.Line,
		Config:              goOut,
		SkipRules:           utils.GetSkipRules(c.rangeSource(hclBody.Range())),
		MaxSeverity:         maxSeverity,
		MinSeverity:         minSeverity,
		ContainerImages:     conatiners,
		InitContainerImages: initContainers,
	}

	// successful
	zap.S().Debugf("created resource config for resource '%s', file: '%s'", resourceConfig.Name, resourceConfig.Source)
	return resourceConfig, nil
}
