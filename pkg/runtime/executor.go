/*
    Copyright (C) 2022 Tenable, Inc.

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
	"strings"

	"github.com/tenable/terrascan/pkg/notifications/webhook"
	"github.com/tenable/terrascan/pkg/policy/opa"
	"github.com/tenable/terrascan/pkg/vulnerability"

	"go.uber.org/zap"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/filters"
	iacProvider "github.com/tenable/terrascan/pkg/iac-providers"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/notifications"
	"github.com/tenable/terrascan/pkg/policy"
	res "github.com/tenable/terrascan/pkg/results"
)

const (
	useTerraformCache = "useTerraformCache"
	nonRecursive      = "nonRecursive"
	valuesFiles       = "valuesFiles"
)

// Executor object
type Executor struct {
	filePath                 string
	dirPath                  string
	policyPath               []string
	iacType                  string
	iacVersion               string
	scanRules                []string
	skipRules                []string
	iacProviders             []iacProvider.IacProvider
	policyEngines            []policy.Engine
	notifiers                []notifications.Notifier
	categories               []string
	policyTypes              []string
	severity                 string
	nonRecursive             bool
	useTerraformCache        bool
	findVulnerabilities      bool
	vulnerabilityEngine      vulnerability.Engine
	notificationWebhookURL   string
	notificationWebhookToken string
	repoURL                  string
	repoRef                  string
	valuesFiles              []string
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion string, policyTypes []string, filePath, dirPath string, policyPath, scanRules, skipRules, categories []string, severity string, nonRecursive, useTerraformCache, findVulnerabilities bool, notificationWebhookURL, notificationWebhookToken, repoURL, repoRef string, valuesFiles []string) (e *Executor, err error) {
	e = &Executor{
		filePath:                 filePath,
		dirPath:                  dirPath,
		policyPath:               policyPath,
		policyTypes:              policyTypes,
		iacType:                  iacType,
		iacVersion:               iacVersion,
		iacProviders:             make([]iacProvider.IacProvider, 0),
		nonRecursive:             nonRecursive,
		useTerraformCache:        useTerraformCache,
		findVulnerabilities:      findVulnerabilities,
		notificationWebhookURL:   notificationWebhookURL,
		notificationWebhookToken: notificationWebhookToken,
		repoURL:                  repoURL,
		repoRef:                  repoRef,
		valuesFiles:              valuesFiles,
	}

	// assigning vulnerabilityEngine
	vulnerabilityEngine, err := vulnerability.NewVulEngine()
	if err == nil {
		e.vulnerabilityEngine = vulnerabilityEngine
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
	if e.notificationWebhookURL != "" {
		e.notifiers = []notifications.Notifier{
			&webhook.Webhook{
				URL:   e.notificationWebhookURL,
				Token: e.notificationWebhookToken,
			},
		}
	} else {
		e.notifiers, err = notifications.NewNotifiers()
		if err != nil {
			zap.S().Debug("failed to create notifier(s).", zap.Error(err))
			// do not return an error if a key is not present in the config file
			if err != notifications.ErrNotificationNotPresent {
				zap.S().Error("failed to create notifier(s).", zap.Error(err))
				return err
			}
		}
	}

	zap.S().Debug("initialized executor")
	return nil
}

func (e *Executor) initPolicyEngines() (err error) {
	// create a new policy engine based on IaC type
	zap.S().Debugf("using policy path %v", e.policyPath)
	for _, policyPath := range e.policyPath {
		engine, err := opa.NewEngine()
		if err != nil {
			zap.S().Errorf("failed to create policy engine. error: '%s'", err)
			return err
		}

		// create a new RegoMetadata pre load filter
		preloadFilter := filters.NewRegoMetadataPreLoadFilter(e.scanRules, e.skipRules, e.categories, e.policyTypes, e.severity)

		// initialize the engine
		if err := engine.Init(policyPath, preloadFilter); err != nil {
			zap.S().Errorf("failed to initialize policy engine for path %s, error: %s", policyPath, err)
			zap.S().Error("perform 'terrascan init' command and then try running the scan command again")
			return err
		}
		e.policyEngines = append(e.policyEngines, engine)
	}
	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute(configOnly, configWithError bool) (results Output, err error) {

	var merr *multierror.Error
	var iacTypes []string
	var resourceConfig output.AllResourceConfigs
	options := e.buildOptions()
	// when dir path has value, only then it will 'all iac' scan
	// when file path has value, we will go with the only iac provider in the list
	if e.dirPath != "" {
		// get all resource configs in the directory
		resourceConfig, iacTypes, merr = e.getResourceConfigs()
	} else {
		// create results output from Iac provider
		// iac providers will contain one element
		resourceConfig, err = e.iacProviders[0].LoadIacFile(e.filePath, options)
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

	if e.findVulnerabilities {
		results.ResourceConfig = e.fetchVulnerabilities(&results, options)
	}

	if configWithError {
		results.Violations.ViolationStore = res.NewViolationStore()
		if e.iacType == "all" {
			results.Violations.ViolationStore.Summary.IacType = strings.Join(iacTypes, ",")
		}
		if err := merr.ErrorOrNil(); err != nil {
			sort.Sort(merr)
			results.Violations.ViolationStore.AddLoadDirErrors(merr.WrappedErrors())
		}
		return results, nil
	}

	if configOnly {
		return results, nil
	}

	if err := e.initPolicyEngines(); err != nil {
		return results, err
	}

	if err := e.findViolations(&results); err != nil {
		return results, err
	}

	if e.findVulnerabilities {
		e.reportVulnerabilities(&results, options)
	}
	resourcePath := e.filePath
	if resourcePath == "" {
		resourcePath = e.dirPath
	}

	// add other summary details after policies are evaluated
	results.Violations.ViolationStore.AddSummary(e.iacType, resourcePath)

	// we want to display the dir scan errors with all the iac providers
	// that support sub folder scanning, which includes 'all' iac scan
	if err := merr.ErrorOrNil(); err != nil {
		// sort multi errors
		sort.Sort(merr)
		results.Violations.ViolationStore.AddLoadDirErrors(merr.WrappedErrors())
	}

	if e.iacType == "all" {
		results.Violations.ViolationStore.Summary.IacType = strings.Join(iacTypes, ",")
	}

	// send notifications, if configured
	if e.repoURL != "" {
		results.Violations.Summary.ResourcePath = e.repoURL
		results.Violations.Summary.Branch = e.repoRef
	}
	e.SendNotifications(results)

	// successful
	return results, nil
}

// getResourceConfigs is a helper method to get all resource configs
func (e *Executor) getResourceConfigs() (output.AllResourceConfigs, []string, *multierror.Error) {
	var merr *multierror.Error
	var iacTypes []string
	resourceConfig := make(output.AllResourceConfigs)

	// channel for directory scan response
	scanRespChan := make(chan dirScanResp)

	options := e.buildOptions()
	// create results output from Iac provider[s]
	for _, iacP := range e.iacProviders {
		go func(ip iacProvider.IacProvider) {
			rc, err := ip.LoadIacDir(e.dirPath, options)
			scanRespChan <- dirScanResp{err, rc, ip.Name()}
		}(iacP)
	}

	for i := 0; i < len(e.iacProviders); i++ {
		sr := <-scanRespChan
		merr = multierror.Append(merr, sr.err)
		if len(sr.rc) > 0 {
			iacTypes = append(iacTypes, sr.iacType)
		}
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

	return resourceConfig, iacTypes, merr
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
			output, err := eng.Evaluate(policy.EngineInput{InputData: &results.ResourceConfig}, &filters.RegoDataFilter{})
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

// fetchVulnerabilities adds vulnerability findings in resource config object
func (e *Executor) fetchVulnerabilities(results *Output, options map[string]interface{}) output.AllResourceConfigs {

	return e.vulnerabilityEngine.FetchVulnerabilities(results.ResourceConfig, options)

}

// reportVulnerabilities adds the found vulnerability data to scan summary output
func (e *Executor) reportVulnerabilities(results *Output, options map[string]interface{}) {
	violations := results.Violations.AsViolationStore()
	vulnerabilityData := e.vulnerabilityEngine.ReportVulnerability(vulnerability.EngineInput{InputData: &results.ResourceConfig}, options)
	if vulnerabilityData.ViolationStore != nil {
		violations.Vulnerabilities = vulnerabilityData.Vulnerabilities
		vulnCount := len(vulnerabilityData.Vulnerabilities)
		violations.Summary.Vulnerabilities = &vulnCount
	}
	results.Violations = policy.EngineOutputFromViolationStore(&violations)
}

// implementsSubFolderScan checks if given iac type supports sub folder scanning
func implementsSubFolderScan(iacType string, nonRecursive bool) bool {
	// iac providers that support sub folder scanning
	// this needs be updated when other iac providers implement
	// sub folder scanning
	if nonRecursive && iacType == "terraform" {
		return false
	}

	iacWithSubFolderScan := []string{"all", "arm", "cft", "k8s", "helm", "terraform", "docker"}
	for _, v := range iacWithSubFolderScan {
		if v == iacType {
			return true
		}
	}
	return false
}

// buildOptions builds map of scan options from executors config
func (e *Executor) buildOptions() map[string]interface{} {
	options := make(map[string]interface{})
	options[useTerraformCache] = e.useTerraformCache
	options[nonRecursive] = e.nonRecursive
	options[valuesFiles] = e.valuesFiles
	return options
}
