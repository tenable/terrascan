package kustomizev3

import (
	"fmt"

	iacloaderror "github.com/accurics/terrascan/pkg/iac-providers/iac-load-error"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

var (
	errLoadIacFileNotSupported = fmt.Errorf("load iac file is not supported for kustomize")
)

// LoadIacFile is not supported for kustomize. Only loading directories that have kustomization.y(a)ml file are supported
func (k *KustomizeV3) LoadIacFile(absRootPath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Debug(errLoadIacFileNotSupported)
	return make(map[string][]output.ResourceConfig), &iacloaderror.LoadError{Err: errLoadIacFileNotSupported}
}
