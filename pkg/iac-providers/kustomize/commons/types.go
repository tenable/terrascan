package commons

import (
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

const (
	// YAMLExtension yaml
	YAMLExtension = "yaml"
	// YAMLExtension2 yml
	YAMLExtension2 = "yml"
	// KustomizeFileName kustomization
	KustomizeFileName = "kustomization"
)

// KustomizeDirectoryLoader implements kustomize directory/file loading
type KustomizeDirectoryLoader struct {
	absRootDir         string
	options            map[string]interface{}
	errIacLoadDirs     *multierror.Error
	useKustomizeBinary bool
	version            string
}

// KustomizeFileNames returns the valid extensions for kustomize (yaml, yml)
func KustomizeFileNames() []string {
	return []string{
		utils.AddFileExtension(KustomizeFileName, YAMLExtension),
		utils.AddFileExtension(KustomizeFileName, YAMLExtension2),
	}
}
