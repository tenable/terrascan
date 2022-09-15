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

package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

var (
	// ErrInitFailed error
	ErrInitFailed = fmt.Errorf("failed to initialize OPA policy engine")
)

// NewEngine returns a new OPA policy engine
func NewEngine() (*Engine, error) {

	// opa engine struct
	engine := &Engine{}

	// successful
	return engine, nil
}

// LoadRegoMetadata Loads rego metadata from a given file
func (e *Engine) LoadRegoMetadata(metaFilename string) (*policy.RegoMetadata, error) {
	// Load metadata file if it exists
	metadata, err := os.ReadFile(metaFilename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			zap.S().Error("failed to load rego metadata", zap.String("file", metaFilename), zap.Error(err))
		}
		return nil, err
	}

	// Read metadata into struct
	regoMetadata := policy.RegoMetadata{}
	if err = json.Unmarshal(metadata, &regoMetadata); err != nil {
		zap.S().Error("failed to unmarshal rego metadata", zap.String("file", metaFilename), zap.Error(err))
		return nil, err
	}
	return &regoMetadata, err
}

// loadRawRegoFilesIntoMap imports raw rego files into a map
func (e *Engine) loadRawRegoFilesIntoMap(currentDir string, regoDataList []*policy.RegoData, regoFileMap *map[string][]byte) error {
	for i := range regoDataList {
		regoPath := filepath.Join(currentDir, regoDataList[i].Metadata.File)
		rawRegoData, err := os.ReadFile(regoPath)
		if err != nil {
			zap.S().Error("failed to load rego file", zap.String("file", regoPath), zap.Error(err))
			continue
		}

		// Load the raw rego into the map
		_, ok := (*regoFileMap)[regoPath]
		if ok {
			// Already loaded this file, so continue
			zap.S().Debug("skipping already loaded rego file", zap.String("file", regoPath))
			continue
		}

		(*regoFileMap)[regoPath] = rawRegoData
	}
	return nil
}

// LoadRegoFiles Loads all related rego files from the given policy path into memory
func (e *Engine) LoadRegoFiles(policyPath string, filter policy.PreLoadFilter) error {
	// Walk the file path and find all directories
	dirList, err := utils.FindAllDirectories(policyPath)
	if err != nil {
		return err
	}

	if len(dirList) == 0 {
		return fmt.Errorf("no directories found for path %s", policyPath)
	}

	e.regoFileMap = make(map[string][]byte)
	e.regoDataMap = make(map[string]*policy.RegoData)

	// Load rego data files from each dir
	// First, we read the metadata file, which contains info about the associated rego rule. The .rego file data is
	// stored in a map in its raw format.
	sort.Strings(dirList)
	for i := range dirList {
		// Find all files in the current dir
		var dirEntries []os.DirEntry
		dirEntries, err = os.ReadDir(dirList[i])
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				zap.S().Debug("error while searching for files", zap.String("dir", dirList[i]))
			}
			continue
		}

		// Load the rego metadata first (*.json)
		metadataFiles := utils.FilterFileInfoBySuffix(&dirEntries, []string{RegoMetadataFileSuffix})
		if len(metadataFiles) == 0 {
			zap.S().Debug("no metadata files were found", zap.String("dir", dirList[i]))
			continue
		}

		var regoDataList []*policy.RegoData
		for j := range metadataFiles {
			filePath := filepath.Join(dirList[i], *metadataFiles[j])

			var regoMetadata *policy.RegoMetadata
			regoMetadata, err = e.LoadRegoMetadata(filePath)
			if err != nil {
				zap.S().Error("error loading rego metadata", zap.String("file", filePath), zap.Error(err))
				continue
			}

			// check if the rego metadata is allowed
			// this check is for scan rules, categories, policy types, and severity
			if !filter.IsAllowed(regoMetadata) {
				continue
			}

			// check if the rego metadata should be filtered
			// this check is for skip rules
			if filter.IsFiltered(regoMetadata) {
				continue
			}

			// Perform some sanity checks
			if strings.Contains(regoMetadata.Name, ".") {
				zap.S().Error("error loading rego metadata: rule name must not contain a dot character", zap.String("name", regoMetadata.Name), zap.String("file", filePath))
				continue
			}

			// Check for default template variable values specified
			if _, ok := regoMetadata.TemplateArgs["name"]; !ok {
				if regoMetadata.TemplateArgs == nil {
					regoMetadata.TemplateArgs = make(map[string]interface{})
				}
				// Add reserved template variable values
				regoMetadata.TemplateArgs["name"] = regoMetadata.Name
			}

			regoData := policy.RegoData{
				Metadata: *regoMetadata,
			}

			regoDataList = append(regoDataList, &regoData)
			e.stats.metadataFileCount++
		}

		// Read in raw rego data from associated rego files
		if err = e.loadRawRegoFilesIntoMap(dirList[i], regoDataList, &e.regoFileMap); err != nil {
			zap.S().Error("error loading raw rego data", zap.String("dir", dirList[i]), zap.Error(err))
			continue
		}

		for j := range regoDataList {
			e.stats.metadataCount++

			// Check if the template file exists
			templateFile := filepath.Join(dirList[i], regoDataList[j].Metadata.File)

			// Apply templates if available
			var templateData bytes.Buffer
			t := template.New("opa")
			_, err = t.Parse(string(e.regoFileMap[templateFile]))
			if err != nil {
				zap.S().Error("unable to parse template", zap.String("template", regoDataList[j].Metadata.File), zap.Error(err))
				continue
			}
			if err = t.Execute(&templateData, regoDataList[j].Metadata.TemplateArgs); err != nil {
				zap.S().Error("unable to execute template", zap.String("template", regoDataList[j].Metadata.File), zap.Error(err))
				continue
			}

			regoDataList[j].RawRego = templateData.Bytes()
			if regoDataList[j].RawRego == nil {
				zap.S().Debug("raw rego data was null after applying template", zap.String("template", regoDataList[j].Metadata.File))
				continue
			}
			e.regoDataMap[regoDataList[j].Metadata.ReferenceID] = regoDataList[j]
		}
	}

	e.stats.ruleCount = len(e.regoDataMap)
	e.stats.regoFileCount = len(e.regoFileMap)
	zap.S().Debugf("loaded %d Rego rules from %d rego files (%d metadata files).", e.stats.ruleCount, e.stats.regoFileCount, e.stats.metadataFileCount)

	return err
}

// CompileRegoFiles Compiles rego files for faster evaluation
func (e *Engine) CompileRegoFiles() error {
	for k := range e.regoDataMap {
		compiler, err := ast.CompileModules(map[string]string{
			e.regoDataMap[k].Metadata.Name: string(e.regoDataMap[k].RawRego),
		})
		if err != nil {
			zap.S().Error("error compiling rego files", zap.String("rule", e.regoDataMap[k].Metadata.Name),
				zap.String("raw rego", string(e.regoDataMap[k].RawRego)), zap.Error(err))
			return err
		}

		r := rego.New(
			rego.Query(RuleQueryBase+"."+e.regoDataMap[k].Metadata.Name),
			rego.Compiler(compiler),
		)

		// Create a prepared query that can be evaluated.
		query, err := r.PrepareForEval(e.context)
		if err != nil {
			zap.S().Error("error creating prepared query", zap.String("rule", e.regoDataMap[k].Metadata.Name),
				zap.String("raw rego", string(e.regoDataMap[k].RawRego)), zap.Error(err))
			return err
		}

		e.regoDataMap[k].PreparedQuery = &query
	}

	return nil
}

// Init initializes the Opa engine
// Handles loading all rules, filtering, compiling, and preparing for evaluation
func (e *Engine) Init(policyPath string, filter policy.PreLoadFilter) error {
	e.context = context.Background()

	if err := e.LoadRegoFiles(policyPath, filter); err != nil {
		zap.S().Error("error loading rego files", zap.String("policy path", policyPath), zap.Error(err))
		return ErrInitFailed
	}

	err := e.CompileRegoFiles()
	if err != nil {
		zap.S().Error("error compiling rego files", zap.String("policy path", policyPath), zap.Error(err))
		return ErrInitFailed
	}

	// initialize ViolationStore
	e.results.ViolationStore = results.NewViolationStore()

	return nil
}

// Configure Configures the OPA engine
func (e *Engine) Configure() error {
	return nil
}

// GetResults Fetches results from OPA engine policy evaluation
func (e *Engine) GetResults() policy.EngineOutput {
	return e.results
}

// Release Performs any tasks required to free resources
func (e *Engine) Release() error {
	return nil
}

// reportViolation Add a violation for a given resource
func (e *Engine) reportViolation(regoData *policy.RegoData, resource *output.ResourceConfig, isSkipped bool, skipComment string) {
	violation := results.Violation{
		RuleName:     regoData.Metadata.Name,
		Description:  regoData.Metadata.Description,
		RuleID:       regoData.Metadata.ID,
		Severity:     regoData.Metadata.Severity,
		Category:     regoData.Metadata.Category,
		RuleFile:     regoData.Metadata.File,
		RuleData:     regoData.RawRego,
		ResourceName: resource.Name,
		ResourceType: resource.Type,
		ResourceData: resource.Config,
		ModuleName:   resource.ModuleName,
		File:         resource.Source,
		PlanRoot:     resource.PlanRoot,
		LineNumber:   resource.Line,
	}

	if !strings.EqualFold(resource.MaxSeverity, "none") {
		// if both values are set then min severity will be applicable
		if resource.MinSeverity != "" {
			if utils.MinSeverityApplicable(regoData.Metadata.Severity, resource.MinSeverity) {
				violation.Severity = strings.ToUpper(resource.MinSeverity)
			}
		} else if utils.MaxSeverityApplicable(regoData.Metadata.Severity, resource.MaxSeverity) {
			violation.Severity = strings.ToUpper(resource.MaxSeverity)
		}
	}
	if !isSkipped && !strings.EqualFold(resource.MaxSeverity, "none") {
		severity := violation.Severity
		if strings.ToLower(severity) == "high" {
			e.results.ViolationStore.Summary.HighCount++
		} else if strings.ToLower(severity) == "medium" {
			e.results.ViolationStore.Summary.MediumCount++
		} else if strings.ToLower(severity) == "low" {
			e.results.ViolationStore.Summary.LowCount++
		} else {
			zap.S().Warn("invalid severity found in rule definition",
				zap.String("rule id", violation.RuleID), zap.String("severity", severity))
		}
		e.results.ViolationStore.AddResult(&violation, false)
		e.results.ViolationStore.Summary.ViolatedPolicies++
	} else {
		violation.Comment = skipComment
		e.results.ViolationStore.AddResult(&violation, true)
	}
}

// reportPassed Adds a passed rule which wasn't violated by all the resources
func (e *Engine) reportPassed(regoData *policy.RegoData) {
	passedRule := results.PassedRule{
		RuleName:    regoData.Metadata.Name,
		Description: regoData.Metadata.Description,
		RuleID:      regoData.Metadata.ID,
		Severity:    regoData.Metadata.Severity,
		Category:    regoData.Metadata.Category,
	}

	e.results.ViolationStore.AddPassedRule(&passedRule)
}

// Evaluate Executes compiled OPA queries against the input JSON data
func (e *Engine) Evaluate(engineInput policy.EngineInput, filter policy.PreScanFilter) (policy.EngineOutput, error) {
	// Keep track of how long it takes to evaluate the policies
	start := time.Now()

	e.regoDataMap = filter.Filter(e.regoDataMap, engineInput)

	// update the rule count
	e.stats.ruleCount = len(e.regoDataMap)

	// Evaluate the policy against each resource type
	for k := range e.regoDataMap {
		// Execute the prepared query.
		rs, err := e.regoDataMap[k].PreparedQuery.Eval(e.context, rego.EvalInput(engineInput.InputData))
		if err != nil {
			// since the eval failed with the policy, we should decrement the total count by 1
			e.stats.ruleCount--
			zap.S().Debug("failed to run prepared query", zap.Error(err), zap.String("rule", "'"+k+"'"), zap.String("file", e.regoDataMap[k].Metadata.File))
			continue
		}

		if len(rs) == 0 || len(rs[0].Expressions) == 0 {
			zap.S().Debug("query executed but found no matches", zap.Error(err), zap.String("rule", "'"+k+"'"))
			continue
		}

		resourceViolations := rs[0].Expressions[0].Value.([]interface{})
		if len(resourceViolations) == 0 {
			zap.S().Debug("query executed but found no violations", zap.Error(err), zap.String("rule", "'"+k+"'"))
			// add the passed rule
			e.reportPassed(e.regoDataMap[k])
			continue
		}

		// Report a violation for each resource returned by the policy evaluation
		for i := range resourceViolations {
			var resourceID string

			// The return values come in two categories--either a map[string]interface{} type, where the "Id" key
			// contains the resource ID, or a string type which is the resource ID. This resource ID is where a
			// violation was found
			switch res := resourceViolations[i].(type) {
			case map[string]interface{}:
				_, ok := res["Id"]
				if !ok {
					zap.S().Warn("no Id key found in resource map", zap.Any("resource", res))
					continue
				}

				_, ok = res["Id"].(string)
				if !ok {
					zap.S().Warn("id key was invalid", zap.Any("resource", res))
					continue
				}
				resourceID = res["Id"].(string)
			case string:
				resourceID = res
			default:
				zap.S().Warn("resource ID format was invalid", zap.Any("resource", res))
				continue
			}

			// Locate the resource details within the input map
			resources, err := engineInput.InputData.FindAllResourcesByID(resourceID)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			if len(resources) == 0 {
				zap.S().Warn("resource was not found", zap.String("resource id", resourceID))
				continue
			}

			for _, resource := range resources {
				// add to skipped violations if rule is skipped for resource
				if len(resource.SkipRules) > 0 {
					found := false
					var skipComment string
					for _, rule := range resource.SkipRules {
						if strings.EqualFold(e.regoDataMap[k].Metadata.ID, rule.Rule) {
							found = true
							skipComment = rule.Comment
							break
						}
						if strings.EqualFold(k, rule.Rule) {
							found = true
							skipComment = rule.Comment
							zap.S().Warnf("Deprecation warning : Use 'id' (%s) instead of 'reference_id' (%s) to skip rules", e.regoDataMap[k].Metadata.ID, k)
							break
						}
					}
					if found {
						// report skipped
						e.reportViolation(e.regoDataMap[k], resource, true, skipComment)
						zap.S().Debugf("rule: %s skipped for resource: %s", k, resource.Name)
					} else {
						// Report the violation
						e.reportViolation(e.regoDataMap[k], resource, false, "")
					}
				} else {
					// Report the violation
					e.reportViolation(e.regoDataMap[k], resource, false, "")
				}
			}

			zap.S().Debug("violation found for rule with rego", zap.String("rego", string("\n")+string(e.regoDataMap[k].RawRego)+string("\n")))
		}
	}

	e.stats.runTime = time.Since(start)

	// add the rule count of the policy engine to result summary
	e.results.ViolationStore.Summary.TotalPolicies += e.stats.ruleCount

	// add the time taken to the result summary
	e.results.ViolationStore.Summary.TotalTime += int64(e.stats.runTime)
	return e.results, nil
}
