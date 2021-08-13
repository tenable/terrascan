package commons

import (
	"fmt"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

var (
	errLoadIacFileNotSupported = fmt.Errorf("load iac file is not supported for kustomize")
)

// LoadIacFile is not supported for kustomize. Only loading directories that have kustomization.y(a)ml file are supported
func LoadIacFile() (allResourcesConfig output.AllResourceConfigs, err error) {
	zap.S().Error(errLoadIacFileNotSupported)
	return nil, errLoadIacFileNotSupported
}
