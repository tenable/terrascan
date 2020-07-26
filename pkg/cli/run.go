package cli

import (
	"github.com/accurics/terrascan/pkg/runtime"
)

// Run executes terrascan in CLI mode
func Run(iacType, iacVersion, cloudType, iacFilePath, iacDirPath string) {

	// create a new runtime executor for processing IaC
	executor, err := runtime.NewExecutor(iacType, iacVersion, cloudType, iacFilePath,
		iacDirPath)
	if err != nil {
		return
	}
	executor.Execute()
}
