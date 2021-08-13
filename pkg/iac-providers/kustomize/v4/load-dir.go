package kustomizev4

import (
	"github.com/accurics/terrascan/pkg/iac-providers/kustomize/commons"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

const (
	versionSuffix = "V4"
)

// LoadIacDir loads the kustomize directory and returns the ResourceConfig mapping which is evaluated by the policy engine
func (k *KustomizeV4) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return commons.NewKustomizeDirectoryLoader(absRootDir, options, false, versionSuffix).LoadIacDir()
}
