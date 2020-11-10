package kustomizev3

import (
	"fmt"
)

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
		getFullFileName(KustomizeFileName, YAMLExtension),
		getFullFileName(KustomizeFileName, YAMLExtension2),
	}
}

func getFullFileName(file, ext string) string {
	return fmt.Sprintf("%v.%v", file, ext)
}
