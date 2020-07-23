package runtime

import (
	"fmt"
	"os"

	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"

	CloudProvider "github.com/accurics/terrascan/pkg/cloud-providers"
	IacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// Executor object
type Executor struct {
	filePath   string
	dirPath    string
	cloudType  string
	iacType    string
	iacVersion string
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath, dirPath string) *Executor {
	return &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
	}
}

// ValidateInputs validates the inputs to the executor object
func (r *Executor) ValidateInputs() error {

	// error message
	errMsg := "input validation failed"

	// terrascan can accept either a file or a directory, both inputs cannot
	// be processed together
	if r.filePath != "" && r.dirPath != "" {
		zap.S().Errorf("cannot accept both '-f %s' and '-d %s' options together", r.filePath, r.dirPath)
		return fmt.Errorf(errMsg)
	}

	if r.dirPath != "" {
		// if directory, check if directory exists
		absDirPath, err := utils.GetAbsPath(r.dirPath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absDirPath); err != nil {
			zap.S().Errorf("directory '%s' does not exist", absDirPath)
			return fmt.Errorf(errMsg)
		}
		zap.S().Debugf("directory '%s' exists", absDirPath)
	} else {

		// if file path, check if file exists
		absFilePath, err := utils.GetAbsPath(r.filePath)
		if err != nil {
			return fmt.Errorf(errMsg)
		}

		if _, err := os.Stat(absFilePath); err != nil {
			zap.S().Errorf("file '%s' does not exist", absFilePath)
			return fmt.Errorf(errMsg)
		}
		zap.S().Debugf("file '%s' exists", absFilePath)
	}

	// check if Iac type is supported
	if !IacProvider.IsIacSupported(r.iacType, r.iacVersion) {
		zap.S().Errorf("iac type '%s', version '%s' not supported", r.iacType, r.iacVersion)
		return fmt.Errorf(errMsg)
	}
	zap.S().Debugf("iac type '%s', version '%s' is supported", r.iacType, r.iacVersion)

	// check if cloud type is supported
	if !CloudProvider.IsCloudSupported(r.cloudType) {
		zap.S().Errorf("cloud type '%s' not supported", r.cloudType)
		return fmt.Errorf(errMsg)
	}
	zap.S().Debugf("cloud type '%s' supported", r.cloudType)

	// check if policy type is supported

	// successful
	zap.S().Debug("input validation successful")
	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (r *Executor) Execute() error {

	// validate inputs
	if err := r.ValidateInputs(); err != nil {
		return err
	}

	// create new IacProvider
	iacProvider, err := IacProvider.NewIacProvider(r.iacType, r.iacVersion)
	if err != nil {
		zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", r.iacType, err)
		return err
	}

	var iacOut output.AllResourceConfigs
	if r.dirPath != "" {
		iacOut, err = iacProvider.LoadIacDir(r.dirPath)
	} else {
		// create config from IaC
		iacOut, err = iacProvider.LoadIacFile(r.filePath)
	}
	if err != nil {
		return err
	}

	// create new CloudProvider
	cloudProvider, err := CloudProvider.NewCloudProvider(r.cloudType)
	if err != nil {
		zap.S().Errorf("failed to create a new CloudProvider for cloudType '%s'. error: '%s'", r.cloudType, err)
		return err
	}
	normalized, err := cloudProvider.CreateNormalizedJSON(iacOut)
	if err != nil {
		return err
	}
	utils.PrintJSON(normalized)

	// write output

	// successful
	return nil
}
