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
	"go.uber.org/zap"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/policy"
	opa "github.com/accurics/terrascan/pkg/policy/opa"
)

// Executor object
type Executor struct {
	filePath     string
	dirPath      string
	policyPath   string
	cloudType    string
	iacType      string
	iacVersion   string
	configFile   string
	iacProvider  iacProvider.IacProvider
	policyEngine policy.Engine
	notifiers    []notifications.Notifier
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath, dirPath, configFile, policyPath string) (e *Executor, err error) {
	e = &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		policyPath: policyPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
		configFile: configFile,
	}

	// initialize executor
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

	// create new notifiers
	e.notifiers, err = notifications.NewNotifiers(e.configFile)
	if err != nil {
		zap.S().Errorf("failed to create notifier(s). error: '%s'", err)
		return err
	}

	// create a new policy engine based on IaC type
	e.policyEngine, err = opa.NewEngine(e.policyPath)
	if err != nil {
		zap.S().Errorf("failed to create policy engine. error: '%s'", err)
		return err
	}

	zap.S().Debug("initialized executor")
	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() (results policy.EngineOutput, err error) {

	// create results output from Iac
	var normalized output.AllResourceConfigs
	if e.filePath != "" {
		normalized, err = e.iacProvider.LoadIacFile(e.filePath)
	} else {
		normalized, err = e.iacProvider.LoadIacDir(e.dirPath)
	}
	if err != nil {
		return results, err
	}

	// evaluate policies
	results, err = e.policyEngine.Evaluate(policy.EngineInput{InputData: &normalized})
	if err != nil {
		return results, err
	}

	// send notifications, if configured
	if err = e.SendNotifications(results); err != nil {
		return results, err
	}

	// successful
	return results, nil
}
