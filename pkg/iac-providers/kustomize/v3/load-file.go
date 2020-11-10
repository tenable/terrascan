package kustomizev3

import (
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

var (
	errLoadIacFileNotSupported = fmt.Errorf("load iac file is not supported for kustomize")
)

// LoadIacFile is not supported for helm. Only loading chart directories are supported
func (k *KustomizeV3) LoadIacFile(absRootPath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Error(errLoadIacFileNotSupported)
	return make(map[string][]output.ResourceConfig), errLoadIacFileNotSupported
}
