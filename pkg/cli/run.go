package cli

import (
	"github.com/accurics/terrascan/pkg/runtime"
)

// Run executes terrascan in CLI mode
func Run(iacType, iacVersion, cloudType, iacFilePath string) {

	// create a new runtime executor for processing IaC
	executor := runtime.NewExecutor(iacType, iacVersion, cloudType, iacFilePath)
	executor.Process()
}
