package commons

import (
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
)

const (
	hashiCorp     = "hashicorp"
	container     = "container"
	initContainer = "init_container"
	spec          = "spec"
	template      = "template"
	kubernetes    = "kubernetes"
	image         = "image"
	name          = "name"
	jobTemplate   = "job_template"
)

// ResourceMetadata holds details about the provider
type ResourceMetadata struct {
	ProviderType string // 	kubernetes
	ResourceType string // 	kubernetes_service
	Namespace    string // hoshicorp
}

var eligibiltylist []ResourceMetadata = []ResourceMetadata{
	{ProviderType: kubernetes, ResourceType: "kubernetes_deployment"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_pod"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_stateful_set"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_job"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_cron_job"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_daemonset"},
	{ProviderType: kubernetes, ResourceType: "kubernetes_replication_controller"},
}

func isEligibleForContainerImageExtraction(resource *hclConfigs.Resource, reqdProviderNameMapping map[string]ResourceMetadata) bool {
	if resource.Provider.Namespace == "" && resource.Provider.Type == "" {
		if v, ok := reqdProviderNameMapping[kubernetes]; ok {
			for _, item := range eligibiltylist {
				if strings.ToLower(v.Namespace) == hashiCorp &&
					strings.ToLower(resource.Type) == item.ResourceType {
					return true
				}
			}
		} else {
			for _, item := range eligibiltylist {
				if strings.ToLower(resource.Type) == item.ResourceType {
					return true
				}
			}
		}
		return false
	}

	for _, item := range eligibiltylist {
		// only official providers from hashicorp will be eligible for now
		if strings.ToLower(resource.Provider.Namespace) == hashiCorp &&
			strings.ToLower(resource.Provider.Type) == item.ProviderType &&
			strings.ToLower(resource.Type) == item.ResourceType {
			return true
		}
	}
	return false
}

func extractContainerImages(resource *hclConfigs.Resource, body *hclsyntax.Body) (containers, initContainers []output.ContainerNameAndImage) {
	for _, block := range body.Blocks {
		if block.Type == spec {
			containerBlocks, initContainerBlocks := getContainerAndInitContainerFromBlocks(block.Body)
			containers = getContainerConfigFromContainerBlock(containerBlocks)
			initContainers = getContainerConfigFromContainerBlock(initContainerBlocks)

		}
	}
	return
}

func getContainerAndInitContainerFromBlocks(specs *hclsyntax.Body) (containers, initContainers []*hclsyntax.Block) {
	for _, block := range specs.Blocks {
		if block.Type == template {
			return getContainerAndInitContainerFromTemplateBlocks(block.Body.Blocks)
		} else if block.Type == jobTemplate {
			for _, jobTemplateBlock := range block.Body.Blocks {
				if jobTemplateBlock.Type == spec {
					return getContainerAndInitContainerFromBlocks(jobTemplateBlock.Body)
				}
			}
		} else if block.Type == container {
			containers = append(containers, block)
		}
	}
	return
}

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
