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

package commons

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

const (
	hashiCorp              = "hashicorp"
	container              = "container"
	initContainer          = "init_container"
	spec                   = "spec"
	template               = "template"
	kubernetes             = "kubernetes"
	image                  = "image"
	name                   = "name"
	jobTemplate            = "job_template"
	azureContainerResource = "azurerm_container_group"
	awsContainerResources  = "aws_ecs_task_definition"
	jsonCodeSuffix         = "${jsonencode("
	fileSuffix             = `${file("`
	containerDefinitions   = "container_definitions"
)

// all the type of resources which has container definitions
var k8sResources = map[string]struct{}{
	"kubernetes_deployment":             {},
	"kubernetes_pod":                    {},
	"kubernetes_stateful_set":           {},
	"kubernetes_job":                    {},
	"kubernetes_cron_job":               {},
	"kubernetes_daemonset":              {},
	"kubernetes_replication_controller": {},
}

// isKubernetesResource - verifies resource is k8s type
func isKubernetesResource(resource *hclConfigs.Resource) bool {
	_, ok := k8sResources[resource.Type]
	return ok
}

// isAzureContainerResource verifies resource is azure type
func isAzureContainerResource(resource *hclConfigs.Resource) bool {
	return resource.Type == azureContainerResource
}

// isAwsContainerResource verifies resource is aws type
func isAwsContainerResource(resource *hclConfigs.Resource) bool {
	return resource.Type == awsContainerResources
}

// fetchContainersFromAzureResource extracts all the containers from azure resource
func fetchContainersFromAzureResource(resource jsonObj) []output.ContainerDetails {
	results := []output.ContainerDetails{}
	if v, ok := resource[container]; ok {
		if containers, vok := v.([]jsonObj); vok {
			results = getContainers(containers)
		}
	}
	return results
}

// fetchContainersFromAwsResource extracts all the containers from aws ecs resource
func fetchContainersFromAwsResource(resource jsonObj, hclBody *hclsyntax.Body, resourcePath string) []output.ContainerDetails {
	results := []output.ContainerDetails{}
	if v, ok := resource[containerDefinitions]; ok {
		def := v.(string)
		if strings.HasPrefix(def, jsonCodeSuffix) {
			return getContainersFromhclBody(hclBody)
		} else if strings.HasPrefix(def, fileSuffix) {
			fileLocation := strings.TrimSpace(def)
			fileLocation = strings.TrimPrefix(fileLocation, fileSuffix)
			fileLocation = strings.TrimSuffix(fileLocation, `")}`)
			dir := filepath.Dir(resourcePath)
			if !filepath.IsAbs(fileLocation) {
				fileLocation = filepath.Join(dir, fileLocation)
			}
			fileData, err := os.ReadFile(fileLocation)
			if err != nil {
				zap.S().Warnf("failed to fetch containers from aws resource: %v", err)
				return results
			}
			def = string(fileData)
		}
		containers := []jsonObj{}
		err := json.Unmarshal([]byte(def), &containers)
		if err != nil {
			zap.S().Warnf("failed to fetch containers from aws resource: %v", err)
			return results
		}
		results = getContainers(containers)
	}
	return results
}

// getContainersFromhclBody parses the attribute and creates container object
func getContainersFromhclBody(hclBody *hclsyntax.Body) (results []output.ContainerDetails) {
	for _, v := range hclBody.Attributes {
		if v.Name == containerDefinitions {
			switch v.Expr.(type) {
			case *hclsyntax.FunctionCallExpr:
				funcExp := v.Expr.(*hclsyntax.FunctionCallExpr)
				for _, arg := range funcExp.Args {
					re, diags := arg.Value(nil)
					if diags.HasErrors() {
						zap.S().Warnf("failed to fetch the container from aws resource: %v", getErrorMessagesFromDiagnostics(diags))
						return
					}
					if !re.CanIterateElements() {
						return
					}
					it := re.ElementIterator()
					for it.Next() {
						_, val := it.Element()
						containerTemp, err := convertCtyToGoNative(val)
						if err != nil {
							zap.S().Warnf("failed to fetch the container from aws resource: %v", err)
							return
						}
						var (
							containerMap map[string]interface{}
							isMap        bool
						)

						if containerMap, isMap = containerTemp.(map[string]interface{}); !isMap {
							break
						}
						tempContainer := output.ContainerDetails{}
						if image, iok := containerMap[image]; iok {
							if imageName, ok := image.(string); ok {
								tempContainer.Image = imageName
							}
						}
						if name, nok := containerMap[name]; nok {
							if containerName, ok := name.(string); ok {
								tempContainer.Name = containerName
							}
						}
						if tempContainer.Name == "" && tempContainer.Image == "" {
							continue
						}
						results = append(results, tempContainer)
					}
				}
			}
			break
		}
	}
	return
}

// getContainers reads and creates container config
func getContainers(containers []jsonObj) (results []output.ContainerDetails) {
	for _, container := range containers {
		tempContainer := output.ContainerDetails{}
		if image, iok := container[image]; iok {
			if imageName, ok := image.(string); ok {
				tempContainer.Image = imageName
			}
		}
		if name, nok := container[name]; nok {
			if containerName, ok := name.(string); ok {
				tempContainer.Name = containerName
			}
		}
		if tempContainer.Name == "" && tempContainer.Image == "" {
			continue
		}
		results = append(results, tempContainer)
	}
	return
}

// extractContainerImagesFromk8sResources extracts containers from k8s resource
func extractContainerImagesFromk8sResources(resource *hclConfigs.Resource, body *hclsyntax.Body) (containers, initContainers []output.ContainerDetails) {
	for _, block := range body.Blocks {
		if block.Type == spec {
			containerBlocks, initContainerBlocks := getContainerAndInitContainerFromSpecBlocks(block.Body)
			containers = getContainerConfigFromContainerBlock(containerBlocks)
			initContainers = getContainerConfigFromContainerBlock(initContainerBlocks)

		}
	}
	return
}

// getContainerAndInitContainerFromSpecBlocks extracts container config from spec block of resource
func getContainerAndInitContainerFromSpecBlocks(specs *hclsyntax.Body) (containers, initContainers []*hclsyntax.Block) {
	for _, block := range specs.Blocks {
		if block.Type == template {
			return getContainerAndInitContainerFromTemplateBlocks(block.Body.Blocks)
		} else if block.Type == jobTemplate {
			for _, jobTemplateBlock := range block.Body.Blocks {
				if jobTemplateBlock.Type == spec {
					return getContainerAndInitContainerFromSpecBlocks(jobTemplateBlock.Body)
				}
			}
		} else if block.Type == container {
			containers = append(containers, block)
		}
	}
	return
}

// getContainerAndInitContainerFromTemplateBlocks extracts container config from template block of resource
func getContainerAndInitContainerFromTemplateBlocks(templateBlocks []*hclsyntax.Block) (containers, initContainers []*hclsyntax.Block) {
	for _, templateBlocks := range templateBlocks {
		if templateBlocks.Type == spec {
			for _, specBlocks := range templateBlocks.Body.Blocks {
				if specBlocks.Type == container {
					containers = append(containers, specBlocks)
				} else if specBlocks.Type == initContainer {
					initContainers = append(initContainers, specBlocks)
				}
			}
		}
	}
	return
}

// getContainerConfigFromContainerBlock creates container config from container block of resource
func getContainerConfigFromContainerBlock(containerBlocks []*hclsyntax.Block) (containerImages []output.ContainerDetails) {
	for _, containerBlock := range containerBlocks {
		containerImage := output.ContainerDetails{}
		for _, attr := range containerBlock.Body.Attributes {
			if attr.Name == image {
				containerImage.Image = getValueFromCtyExpr(attr.Expr)
			}
			if attr.Name == name {
				containerImage.Name = getValueFromCtyExpr(attr.Expr)
			}
		}
		if containerImage.Image == "" && containerImage.Name == "" {
			continue
		}
		containerImages = append(containerImages, containerImage)
	}
	return
}

// getValueFromCtyExpr get value string from hcl expression
func getValueFromCtyExpr(expr hclsyntax.Expression) (value string) {
	val, diags := expr.Value(nil)
	if diags.HasErrors() {
		zap.S().Errorf("error fetching containers from k8s resource: %v", getErrorMessagesFromDiagnostics(diags))
		return
	}
	valInterface, err := convertCtyToGoNative(val)
	if err != nil {
		zap.S().Errorf("error fetching containers from k8s resource: %v", err)
		return
	}
	if containerName, ok := valInterface.(string); ok {
		value = containerName
	}
	return
}
