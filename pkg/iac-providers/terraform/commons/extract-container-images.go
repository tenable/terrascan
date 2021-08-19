package commons

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
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

var k8sResources = make(map[string]struct{})

// all the type of resources which has container definitaions
var k8sResourcesTypes = []string{"kubernetes_deployment",
	"kubernetes_pod", "kubernetes_stateful_set", "kubernetes_job",
	"kubernetes_cron_job", "kubernetes_daemonset", "kubernetes_replication_controller"}

func init() {
	for _, resource := range k8sResourcesTypes {
		k8sResources[resource] = struct{}{}
	}
}

// isKuberneteResource - verifies resource is k8s type and can we fetch container details from it
func isKuberneteResource(resource *hclConfigs.Resource) bool {
	_, ok := k8sResources[resource.Type]
	return ok
}

//isAzureConatinerResource verifies resource is azure type and can we fetch container details from it
func isAzureConatinerResource(resource *hclConfigs.Resource) bool {
	return resource.Type == azureContainerResource
}

//isAwsConatinerResource verifies resource is aws type and can we fetch container details from it
func isAwsConatinerResource(resource *hclConfigs.Resource) bool {
	return resource.Type == awsContainerResources
}

//fetchConatinersFromAzureResource extracts all the containers from azure resource
func fetchContainersFromAzureResource(resource jsonObj) []output.ContainerNameAndImage {
	results := []output.ContainerNameAndImage{}
	if v, ok := resource[container]; ok {
		if containers, vok := v.([]jsonObj); vok {
			for _, container := range containers {
				tempContainer := output.ContainerNameAndImage{}
				if image, iok := container[image]; iok {
					tempContainer.Image = image.(string)
				}
				if name, nok := container[name]; nok {
					tempContainer.Name = name.(string)
				}
				if tempContainer.Name == "" && tempContainer.Image == "" {
					continue
				}
				results = append(results, tempContainer)
			}
		}
	}
	return results
}

//fetchConatinersFromAwsResource extracts all the containers from aws ecs resource
func fetchContainersFromAwsResource(resource jsonObj, absRootDir string) []output.ContainerNameAndImage {
	results := []output.ContainerNameAndImage{}
	if v, ok := resource[containerDefinitions]; ok {
		def := v.(string)
		if strings.HasPrefix(def, jsonCodeSuffix) {
			def = strings.TrimPrefix(def, jsonCodeSuffix)
			def = strings.TrimSuffix(def, ")}")
		} else if strings.HasPrefix(def, fileSuffix) {
			file := strings.TrimPrefix(def, fileSuffix)
			file = strings.TrimSuffix(file, `")}`)
			dir := filepath.Dir(absRootDir)
			file = filepath.Join(dir, file)
			fileData, err := ioutil.ReadFile(file)
			if err != nil {
				zap.S().Errorf("error reading file: %s : %v", file, err)
				return results
			}
			def = string(fileData)
		}
		containers := []jsonObj{}
		err := json.Unmarshal([]byte(def), &containers)
		if err != nil {
			zap.S().Errorf("error unmarshaling string: %s : %v", def, err)
			return results
		}
		for _, container := range containers {
			tempContainer := output.ContainerNameAndImage{}
			if image, iok := container[image]; iok {
				tempContainer.Image = image.(string)
			}
			if name, nok := container[name]; nok {
				tempContainer.Name = name.(string)
			}
			if tempContainer.Name == "" && tempContainer.Image == "" {
				continue
			}
			results = append(results, tempContainer)
		}
	}
	return results
}

//extractContainerImagesFromk8sResources extracts containers from k8s resource
func extractContainerImagesFromk8sResources(resource *hclConfigs.Resource, body *hclsyntax.Body) (containers, initContainers []output.ContainerNameAndImage) {
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

//getContainerAndInitContainerFromTemplateBlocks extracts container config from template block of resource
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

//getContainerConfigFromContainerBlock creates container config from container block of resource
func getContainerConfigFromContainerBlock(containerBlocks []*hclsyntax.Block) (containerImages []output.ContainerNameAndImage) {
	for _, conatainerBlock := range containerBlocks {
		containerImage := output.ContainerNameAndImage{}
		for _, attr := range conatainerBlock.Body.Attributes {
			if attr.Name == image {
				val, _ := attr.Expr.Value(nil)
				containerImage.Image = val.AsString()
			}
			if attr.Name == name {
				val, _ := attr.Expr.Value(nil)
				containerImage.Name = val.AsString()
			}
		}
		if containerImage.Image == "" && containerImage.Name == "" {
			continue
		}
		containerImages = append(containerImages, containerImage)
	}
	return
}
