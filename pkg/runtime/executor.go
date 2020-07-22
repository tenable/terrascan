package runtime

import (
	"fmt"
	"log"
	"os"

	"github.com/accurics/terrascan/pkg/utils"

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

	// terrascan can accept either a file or a directory, both inputs cannot
	// be processed together
	if r.filePath != "" && r.dirPath != "" {
		errMsg := fmt.Sprintf("cannot accept both '-f %s' and '-d %s' options together", r.filePath, r.dirPath)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
	}

	if r.dirPath != "" {
		// if directory, check if directory exists
		absDirPath, err := utils.GetAbsPath(r.dirPath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absDirPath); err != nil {
			errMsg := fmt.Sprintf("directory '%s' does not exist", absDirPath)
			log.Printf(errMsg)
			return fmt.Errorf(errMsg)
		}
	} else {

		// if file path, check if file exists
		absFilePath, err := utils.GetAbsPath(r.filePath)
		if err != nil {
			return err
		}

		if _, err := os.Stat(absFilePath); err != nil {
			errMsg := fmt.Sprintf("file '%s' does not exist", absFilePath)
			log.Printf(errMsg)
			return fmt.Errorf(errMsg)
		}
	}

	// check if Iac type is supported
	if !IacProvider.IsIacSupported(r.iacType, r.iacVersion) {
		errMsg := fmt.Sprintf("iac type '%s', version '%s' not supported", r.iacType, r.iacVersion)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
	}

	// check if cloud type is supported
	if !CloudProvider.IsCloudSupported(r.cloudType) {
		errMsg := fmt.Sprintf("cloud type '%s' not supported", r.cloudType)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
	}

	// check if policy type is supported

	// successful
	return nil
}

// Process validates the inputs, processes the IaC, creates json output
func (r *Executor) Process() error {

	// validate inputs
	if err := r.ValidateInputs(); err != nil {
		return err
	}

	// create new IacProvider
	iacProvider, err := IacProvider.NewIacProvider(r.iacType, r.iacVersion)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create a new IacProvider for iacType '%s'. error: '%s'", r.iacType, err)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
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
	utils.PrintJSON(iacOut)

	// create new CloudProvider
	cloudProvider, err := CloudProvider.NewCloudProvider(r.cloudType)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create a new CloudProvider for cloudType '%s'. error: '%s'", r.cloudType, err)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
	}
	cloudProvider.CreateNormalizedJson()

	// write output

	// successful
	return nil
}
