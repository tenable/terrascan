package commons

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/addrs"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"go.uber.org/zap"
	"strings"
)

const hashiCorp = "hashicorp"
type ResourceMetadata struct {
	ProviderType string // 	kubernetes
	ResourceType string // 	kubernetes_service
}

var eligibiltylist []ResourceMetadata = []ResourceMetadata{
	{ProviderType: "kubernetes", ResourceType: "kubernetes_deployment"},
	{ProviderType: "kubernetes", ResourceType: "kubernetes_pod"},
	{ProviderType: "kubernetes", ResourceType: "kubernetes_stateful_set"},
	{ProviderType: "kubernetes", ResourceType: "kubernetes_job"},
	{ProviderType: "kubernetes", ResourceType: "kubernetes_cron_job"},
	{ProviderType: "kubernetes", ResourceType: "kubernetes_daemonset"},
}


func isEligibleForContainerImageExtraction(resource *hclConfigs.Resource, reqdProviderNameMapping map[addrs.Provider]string) bool {

	zap.S().Infof("Checking eligibility for resource:")
	zap.S().Infof(resource.Provider.Namespace)
	zap.S().Infof(resource.Provider.Type)
	zap.S().Infof(resource.Type)
	zap.S().Infof(resource.Name)

	for _, item := range eligibiltylist {
		// only official providers from hashicorp will be eligible for now
		if strings.ToLower(resource.Provider.Namespace) == hashiCorp &&
			strings.ToLower(resource.Provider.Type) == item.ProviderType &&
			strings.ToLower(resource.Type) == item.ResourceType {
			zap.S().Info("true")
			return true
		}
	}
	zap.S().Info("false")
	return false
}

func extractContainerImages(resource *hclConfigs.Resource, body *hclsyntax.Body) []error {
	zap.S().Info("Dump Round 1")
//	spew.Dump(body)

	if strings.ToLower(resource.Provider.Type) == "kubernetes" {
		if strings.ToLower(resource.Type) == "kubernetes_deployment" {
			for _, block := range body.Blocks {
				if block.Type == "spec" {
					zap.S().Info("Inside spec")
					for _, block1 := range block.Body.Blocks {
						if block1.Type == "template" {
							zap.S().Info("Inside spec.template")
							for _, block2 := range block1.Body.Blocks {
								zap.S().Infof("block type : %s", block2.Type)
								if block2.Type == "spec" {
									zap.S().Info("Inside spec.template.spec")
									for _, block3 := range block2.Body.Blocks {
										if block3.Type == "container" {
											zap.S().Info("Inside spec.template.spec.container")
											//spew.Dump(block3)
											for _, attr := range block3.Body.Attributes {
												if attr.Name == "image" {
													val, _ := attr.Expr.Value(nil)
													//spew.Dump(val)
													zap.S().Info(val.AsString())
												}
												if attr.Name == "name" {
													val, _ := attr.Expr.Value(nil)
													//spew.Dump(val)
													zap.S().Info(val.AsString())
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			zap.S().Infof("Extraction logic coming soon for resource type %s", resource.Type)
		}
	}
	return nil
}
