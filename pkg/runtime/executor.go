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
	"fmt"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/config"
	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/notifications"
	"github.com/accurics/terrascan/pkg/policy"
	opa "github.com/accurics/terrascan/pkg/policy/opa"
	"github.com/pelletier/go-toml"
)

const (
	rulesKey     = "rules"
	scanRulesKey = "scan-rules"
	skipRulesKey = "skip-rules"
)

var (
	errRuleNotString          = fmt.Errorf("each scan and skip rule must be string")
	errIncorrectValueForRules = fmt.Errorf("'scan-rules' and 'skip-rules' must be an array")
)

// Executor object
type Executor struct {
	filePath     string
	dirPath      string
	policyPath   []string
	cloudType    []string
	iacType      string
	iacVersion   string
	configFile   string
	scanRules    []string
	skipRules    []string
	iacProvider  iacProvider.IacProvider
	policyEngine []policy.Engine
	notifiers    []notifications.Notifier
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion string, cloudType []string, filePath, dirPath, configFile string, policyPath, scanRules, skipRules []string) (e *Executor, err error) {
	e = &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		policyPath: policyPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
		configFile: configFile,
		scanRules:  scanRules,
		skipRules:  skipRules,
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

	// read config file and update scan and skip rules
	if err := e.initScanAndSkipRules(); err != nil {
		if !(err == errIncorrectValueForRules || err == errRuleNotString) {
			return err
		}
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
		// do not return an error if a key is not present in the config file
		if err != notifications.ErrTomlKeyNotPresent {
			return err
		}
	}

	// create a new policy engine based on IaC type
	zap.S().Debugf("using policy path %v", e.policyPath)
	for _, policyPath := range e.policyPath {
		engine, err := opa.NewEngine(policyPath, e.scanRules, e.skipRules)
		if err != nil {
			zap.S().Errorf("failed to create policy engine. error: '%s'", err)
			return err
		}
		e.policyEngine = append(e.policyEngine, engine)
	}

	zap.S().Debug("initialized executor")
	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() (results Output, err error) {

	// create results output from Iac
	if e.filePath != "" {
		results.ResourceConfig, err = e.iacProvider.LoadIacFile(e.filePath)
	} else {
		results.ResourceConfig, err = e.iacProvider.LoadIacDir(e.dirPath)
	}
	if err != nil {
		return results, err
	}

	// evaluate policies
	results.Violations = policy.EngineOutput{}
	violations := results.Violations.AsViolationStore()
	for _, engine := range e.policyEngine {
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

	// successful
	return results, nil
}

// read the config file and update scan and skip rules
func (e *Executor) initScanAndSkipRules() error {
	if e.configFile != "" {
		configData, err := config.LoadConfig(e.configFile)
		if err != nil {
			zap.S().Error("error loading config file", zap.Error(err))
			return err
		}

		if configData.Has(rulesKey) {

			data := (configData.Get(rulesKey)).(*toml.Tree)

			// read scan rules in the toml tree
			if err := initRules(e, data, scanRulesKey); err != nil {
				return err
			}

			// read skip rules in the toml tree
			if err := initRules(e, data, skipRulesKey); err != nil {
				return err
			}
		}
	}
	return nil
}

func initRules(e *Executor, tree *toml.Tree, key string) error {
	rules, err := getRulesInTomlTree(tree, e.configFile, key)
	if err != nil {
		zap.S().Error("error reading config file", zap.Error(err))
		return err
	}
	if len(rules) > 0 {
		if key == scanRulesKey {
			e.scanRules = append(e.scanRules, rules...)
		} else {
			e.skipRules = append(e.skipRules, rules...)
		}
	} else {
		zap.S().Debugf("key '%s' not found in the config file: %s", key, e.configFile)
	}
	return nil
}

func getRulesInTomlTree(tree *toml.Tree, configFile, key string) ([]string, error) {
	ruleSlice := make([]string, 0)
	if tree.Has(key) {
		rules, ok := (tree.Get(key)).([]interface{})
		if !ok {
			zap.S().Errorf("key '%s' must be an array in the config file: %s", key, configFile)
			return nil, errIncorrectValueForRules
		}
		for _, rule := range rules {
			r, ok := rule.(string)
			if !ok {
				zap.S().Errorf("rules must be of type string for key: %s", key)
				return nil, errRuleNotString
			}
			ruleSlice = append(ruleSlice, r)
		}
	}
	return ruleSlice, nil
}
