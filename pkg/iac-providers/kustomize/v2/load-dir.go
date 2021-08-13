package kustomizev2

import (
	"github.com/accurics/terrascan/pkg/iac-providers/kustomize/commons"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

const (
	versionSuffix = "V2"
)

// LoadIacDir loads the kustomize directory and returns the ResourceConfig mapping which is evaluated by the policy engine
func (k *KustomizeV2) LoadIacDir(absRootDir string, options map[string]interface{}) (output.AllResourceConfigs, error) {
	return commons.NewKustomizeDirectoryLoader(absRootDir, options, true, versionSuffix).LoadIacDir()
}
