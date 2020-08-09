/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package runtime

import (
	"github.com/accurics/terrascan/pkg/policy"
	opa "github.com/accurics/terrascan/pkg/policy/opa"

	"go.uber.org/zap"

	cloudProvider "github.com/accurics/terrascan/pkg/cloud-providers"
	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// Executor object
type Executor struct {
	filePath      string
	dirPath       string
	policyPath    string
	cloudType     string
	iacType       string
	iacVersion    string
	iacProvider   iacProvider.IacProvider
	cloudProvider cloudProvider.CloudProvider
	policyEngine  []policy.Engine
	//	policyEngine
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath, dirPath, policyPath string) (e *Executor, err error) {
	e = &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		policyPath: policyPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
	}

	// initialized executor
	if err = e.Init(); err != nil {
		return e, err
	}

	return e, nil
}

// Init validates input and initializes iac and cloud providers
func (e *Executor) Init() error {

	// validate inputs
	err := e.ValidateInputs()
	if err != nil {
		return err
	}

	// create new IacProvider
	e.iacProvider, err = iacProvider.NewIacProvider(e.iacType, e.iacVersion)
	if err != nil {
		zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", e.iacType, err)
		return err
	}

	// create new CloudProvider
	e.cloudProvider, err = cloudProvider.NewCloudProvider(e.cloudType)
	if err != nil {
		zap.S().Errorf("failed to create a new CloudProvider for cloudType '%s'. error: '%s'", e.cloudType, err)
		return err
	}

	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() error {

	// load iac config
	var (
		iacOut output.AllResourceConfigs
		err    error
	)
	if e.dirPath != "" {
		iacOut, err = e.iacProvider.LoadIacDir(e.dirPath)
	} else {
		// create config from IaC
		iacOut, err = e.iacProvider.LoadIacFile(e.filePath)
	}
	if err != nil {
		return err
	}

	// create normalized json
	normalized, err := e.cloudProvider.CreateNormalizedJSON(iacOut)
	if err != nil {
		return err
	}

	// create a new policy engine based on IaC type
	if e.iacType == "terraform" {
		var engine policy.Engine
		engine = &opa.OpaEngine{}

		err = engine.Initialize(e.policyPath)
		if err != nil {
			return err
		}

		engine.Evaluate(&normalized)
	}

	// successful
	return nil
}
