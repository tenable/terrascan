package kustomizev4

import (
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

// KustomizeV4 struct
type KustomizeV4 struct {
	errIacLoadDirs *multierror.Error
}

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
