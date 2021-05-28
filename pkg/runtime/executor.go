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
	"sort"

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
	nonRecursive  bool
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion string, cloudType []string, filePath, dirPath string, policyPath, scanRules, skipRules, categories []string, severity string, nonRecursive bool) (e *Executor, err error) {
	e = &Executor{
		filePath:     filePath,
		dirPath:      dirPath,
		policyPath:   policyPath,
		cloudType:    cloudType,
		iacType:      iacType,
		iacVersion:   iacVersion,
		iacProviders: make([]iacProvider.IacProvider, 0),
		nonRecursive: nonRecursive,
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
			// skip tfplan because it doesn't support directory scanning
			if ip == "tfplan" {
				continue
			}

			// initialize iac providers with default versions
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
	var resourceConfig output.AllResourceConfigs

	// when dir path has value, only then it will 'all iac' scan
	// when file path has value, we will go with the only iac provider in the list
	if e.dirPath != "" {
		// get all resource configs in the directory
		resourceConfig, merr = e.getResourceConfigs()
	} else {
		// create results output from Iac provider
		// iac providers will contain one element
		resourceConfig, err = e.iacProviders[0].LoadIacFile(e.filePath)
		if err != nil {
			return results, err
		}
	}

	// for the iac providers that don't implement sub folder scanning
	// return the error to the caller
	if !implementsSubFolderScan(e.iacType, e.nonRecursive) {
		if err := merr.ErrorOrNil(); err != nil {
			return results, err
		}
	}

	// update results with resource config
	results.ResourceConfig = resourceConfig

	if err := e.findViolations(&results); err != nil {
		return results, err
	}

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

	// we want to display the dir scan errors with all the iac providers
	// that support sub folder scanning, which includes 'all' iac scan
	if err := merr.ErrorOrNil(); err != nil {
		// sort multi errors
		sort.Sort(merr)
		results.Violations.ViolationStore.AddLoadDirErrors(merr.WrappedErrors())
	}

	// successful
	return results, nil
}

// getResourceConfigs is a helper method to get all resource configs
func (e *Executor) getResourceConfigs() (output.AllResourceConfigs, *multierror.Error) {
	var merr *multierror.Error
	resourceConfig := make(output.AllResourceConfigs)

	// channel for directory scan response
	scanRespChan := make(chan dirScanResp)

	// create results output from Iac provider[s]
	for _, iacP := range e.iacProviders {
		go func(ip iacProvider.IacProvider) {
			rc, err := ip.LoadIacDir(e.dirPath, e.nonRecursive)
			scanRespChan <- dirScanResp{err, rc}
		}(iacP)
	}

	for i := 0; i < len(e.iacProviders); i++ {
		sr := <-scanRespChan
		merr = multierror.Append(merr, sr.err)
		// deduplication of resources
		if len(resourceConfig) > 0 {
			for key, r := range sr.rc {
				resourceConfig.UpdateResourceConfigs(key, r)
			}
		} else {
			for key := range sr.rc {
				resourceConfig[key] = append(resourceConfig[key], sr.rc[key]...)
			}
		}
	}

	return resourceConfig, merr
}

// findViolations is a helper method to find all violations in the resource config
func (e *Executor) findViolations(results *Output) error {
	// evaluate policies
	results.Violations = policy.EngineOutput{}
	violations := results.Violations.AsViolationStore()

	// channel for engine evaluation result
	evalResultChan := make(chan engineEvalResult)

	for _, engine := range e.policyEngines {
		go func(eng policy.Engine) {
			output, err := eng.Evaluate(policy.EngineInput{InputData: &results.ResourceConfig})
			evalResultChan <- engineEvalResult{err, output}
		}(engine)
	}

	for i := 0; i < len(e.policyEngines); i++ {
		evalR := <-evalResultChan
		if evalR.err != nil {
			return evalR.err
		}
		violations = violations.Add(evalR.output.AsViolationStore())
	}

	results.Violations = policy.EngineOutputFromViolationStore(&violations)
	return nil
}

// implementsSubFolderScan checks if given iac type supports sub folder scanning
func implementsSubFolderScan(iacType string, nonRecursive bool) bool {
	// iac providers that support sub folder scanning
	// this needs be updated when other iac providers implement
	// sub folder scanning
	if nonRecursive && iacType == "terraform" {
		return false
	}

	iacWithSubFolderScan := []string{"all", "cft", "k8s", "helm", "terraform"}
	for _, v := range iacWithSubFolderScan {
		if v == iacType {
			return true
		}
	}
	return false
}
