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
	"reflect"

	"go.uber.org/zap"

	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/policy"
	opa "github.com/accurics/terrascan/pkg/policy/opa"
	"github.com/hashicorp/go-multierror"
)

// Executor object
type Executor struct {
	filePath      string
	dirPath       string
	policyPath    []string
	cloudType     []string
	iacType       string
	iacVersion    string
	scanRules     []string
	skipRules     []string
	iacProviders  []iacProvider.IacProvider
	policyEngines []policy.Engine
	notifiers     []notifications.Notifier
	categories    []string
	severity      string
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion string, cloudType []string, filePath, dirPath string, policyPath, scanRules, skipRules, categories []string, severity string) (e *Executor, err error) {
	e = &Executor{
		filePath:     filePath,
		dirPath:      dirPath,
		policyPath:   policyPath,
		cloudType:    cloudType,
		iacType:      iacType,
		iacVersion:   iacVersion,
		iacProviders: make([]iacProvider.IacProvider, 0),
	}

	// read config file and update scan and skip rules
	if err := e.loadRuleSetFromConfig(); err != nil {
		zap.S().Error("error initialising scan and skip rules", zap.Error(err))
		return nil, err
	}

	if len(scanRules) > 0 {
		e.scanRules = scanRules
	}

	if len(skipRules) > 0 {
		e.skipRules = skipRules
	}

	if len(severity) > 0 {
		e.severity = severity
	}

	if len(categories) > 0 {
		e.categories = categories
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

	// create new IacProviders
	if e.iacType == "all" {
		for _, ip := range iacProvider.SupportedIacProviders() {
			if ip == "tfplan" {
				continue
			}
			defaultIacVersion := iacProvider.GetDefaultIacVersion(ip)
			iacP, err := iacProvider.NewIacProvider(ip, defaultIacVersion)
			if err != nil {
				zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", e.iacType, err)
				return err
			}
			e.iacProviders = append(e.iacProviders, iacP)
		}
	} else {
		iacP, err := iacProvider.NewIacProvider(e.iacType, e.iacVersion)
		if err != nil {
			zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", e.iacType, err)
			return err
		}
		e.iacProviders = append(e.iacProviders, iacP)
	}

	// create new notifiers
	e.notifiers, err = notifications.NewNotifiers()
	if err != nil {
		zap.S().Debug("failed to create notifier(s).", zap.Error(err))
		// do not return an error if a key is not present in the config file
		if err != notifications.ErrNotificationNotPresent {
			zap.S().Error("failed to create notifier(s).", zap.Error(err))
			return err
		}
	}

	// create a new policy engine based on IaC type
	zap.S().Debugf("using policy path %v", e.policyPath)
	for _, policyPath := range e.policyPath {
		engine, err := opa.NewEngine()
		if err != nil {
			zap.S().Errorf("failed to create policy engine. error: '%s'", err)
			return err
		}

		// initialize the engine
		if err := engine.Init(policyPath, e.scanRules, e.skipRules, e.categories, e.severity); err != nil {
			zap.S().Errorf("%s", err)
			return err
		}
		e.policyEngines = append(e.policyEngines, engine)
	}

	zap.S().Debug("initialized executor")
	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() (results Output, err error) {

	var merr *multierror.Error
	resourceConfig := make(output.AllResourceConfigs)
	// create results output from Iac provider[s]
	for _, iacP := range e.iacProviders {
		var err error
		var rc output.AllResourceConfigs
		if e.filePath != "" {
			rc, err = iacP.LoadIacFile(e.filePath)
		} else {
			rc, err = iacP.LoadIacDir(e.dirPath)
		}
		if err != nil {
			merr = multierror.Append(merr, err)
		}

		// deduplication
		if len(resourceConfig) > 0 {
			for key, r := range rc {
				for _, v := range r {
					if !isConfigPresent(resourceConfig[key], v) {
						resourceConfig[key] = append(resourceConfig[key], v)
					}
				}
			}
		} else {
			for key := range rc {
				resourceConfig[key] = append(resourceConfig[key], rc[key]...)
			}
		}
	}

	if len(e.iacProviders) == 1 && !(e.iacType == "k8s" || e.iacType == "helm") {
		if err := merr.ErrorOrNil(); err != nil {
			return results, err
		}
	}

	results.ResourceConfig = resourceConfig
	// evaluate policies
	results.Violations = policy.EngineOutput{}
	violations := results.Violations.AsViolationStore()
	for _, engine := range e.policyEngines {
		output, err := engine.Evaluate(policy.EngineInput{InputData: &results.ResourceConfig})
		if err != nil {
			return results, err
		}
		violations = violations.Add(output.AsViolationStore())
	}

	results.Violations = policy.EngineOutputFromViolationStore(&violations)

	resourcePath := e.filePath
	if resourcePath == "" {
		resourcePath = e.dirPath
	}

	// add other summary details after policies are evaluated
	results.Violations.ViolationStore.AddSummary(e.iacType, resourcePath)

	// send notifications, if configured
	if err = e.SendNotifications(results); err != nil {
		return results, err
	}

	if e.iacType == "all" {
		if err := merr.ErrorOrNil(); err != nil {
			results.Violations.ViolationStore.AddLoadDirErrors(merr.WrappedErrors())
		}
	}

	// successful
	return results, nil
}

// isConfigPresent checks whether a resource is already present in the list of configs or not
// the equality of a resource is based on name, source and config of the resource
func isConfigPresent(resources []output.ResourceConfig, resourceConfig output.ResourceConfig) bool {
	for _, resource := range resources {
		if resource.Name == resourceConfig.Name && resource.Source == resourceConfig.Source {
			if reflect.DeepEqual(resource.Config, resourceConfig.Config) {
				return true
			}
		}
	}
	return false
}
