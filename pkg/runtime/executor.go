package runtime

import (
	"fmt"
	"log"

	CloudProvider "github.com/accurics/terrascan/pkg/cloud-providers"
	IacProvider "github.com/accurics/terrascan/pkg/iac-providers"
)

// Executor object
type Executor struct {
	filePath   string
	cloudType  string
	iacType    string
	iacVersion string
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath string) *Executor {
	return &Executor{
		filePath:   filePath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
	}
}

// ValidateInputs validates the inputs to the executor object
func (r *Executor) ValidateInputs() error {

	// terrascan can accept either a file or a directory, both inputs cannot
	// be processed together

	// if file path, check if file exists
	// if directory, check if directory exists

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

	// create config from IaC
	_, err = iacProvider.LoadIacFile(r.filePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to load iac file '%s'. error: '%s'", err)
		log.Printf(errMsg)
		return fmt.Errorf(errMsg)
	}

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
