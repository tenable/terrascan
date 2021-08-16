package kustomizev4

import (
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

var (
	errLoadIacFileNotSupported = fmt.Errorf("load iac file is not supported for kustomize")
)

// LoadIacFile is not supported for kustomize. Only loading directories that have kustomization.y(a)ml file are supported
func (k *KustomizeV4) LoadIacFile(absRootPath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Error(errLoadIacFileNotSupported)
	return make(map[string][]output.ResourceConfig), errLoadIacFileNotSupported
}
