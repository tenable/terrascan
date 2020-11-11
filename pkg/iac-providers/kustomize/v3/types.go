package kustomizev3

import "github.com/accurics/terrascan/pkg/utils"

// KustomizeV3 struct
type KustomizeV3 struct{}

const (
	// YAMLExtension yaml
	YAMLExtension = "yaml"
	// YAMLExtension2 yml
	YAMLExtension2 = "yml"
	// KustomizeFileName kustomization
	KustomizeFileName = "kustomization"
)

// KustomizeFileNames returns the valid extensions for k8s (yaml, yml, json)
func KustomizeFileNames() []string {
	return []string{
		utils.AddFileExtension(KustomizeFileName, YAMLExtension),
		utils.AddFileExtension(KustomizeFileName, YAMLExtension2),
	}
}
