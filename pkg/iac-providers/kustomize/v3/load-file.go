package kustomizev3

import (
	"github.com/accurics/terrascan/pkg/iac-providers/kustomize/commons"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// LoadIacFile is not supported for kustomize. Only loading directories that have kustomization.y(a)ml file are supported
func (k *KustomizeV3) LoadIacFile(absRootPath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	return commons.LoadIacFile()
}
