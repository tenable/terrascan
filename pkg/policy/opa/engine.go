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

package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/accurics/terrascan/pkg/utils"

	"github.com/open-policy-agent/opa/ast"

	"go.uber.org/zap"

	"github.com/open-policy-agent/opa/rego"
)

// LoadRegoMetadata Loads rego metadata from a given file
func (e *Engine) LoadRegoMetadata(metaFilename string) (*RegoMetadata, error) {
	// Load metadata file if it exists
	metadata, err := ioutil.ReadFile(metaFilename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			zap.S().Warn("failed to load rego metadata", zap.String("file", metaFilename))
		}
		return nil, err
	}

	// Read metadata into struct
	regoMetadata := RegoMetadata{}
	if err = json.Unmarshal(metadata, &regoMetadata); err != nil {
		zap.S().Warn("failed to unmarshal rego metadata", zap.String("file", metaFilename))
		return nil, err
	}
	return &regoMetadata, err
}

// loadRawRegoFilesIntoMap imports raw rego files into a map
func (e *Engine) loadRawRegoFilesIntoMap(currentDir string, regoDataList []*RegoData, regoFileMap *map[string][]byte) error {
	for i := range regoDataList {
		regoPath := filepath.Join(currentDir, regoDataList[i].Metadata.File)
		rawRegoData, err := ioutil.ReadFile(regoPath)
		if err != nil {
			zap.S().Warn("failed to load rego file", zap.String("file", regoPath))
			continue
		}

		// Load the raw rego into the map
		_, ok := (*regoFileMap)[regoPath]
		if ok {
			// Already loaded this file, so continue
			continue
		}

		(*regoFileMap)[regoPath] = rawRegoData
	}
	return nil
}

// LoadRegoFiles Loads all related rego files from the given policy path into memory
func (e *Engine) LoadRegoFiles(policyPath string) error {
	// Walk the file path and find all directories
	dirList, err := utils.FindAllDirectories(policyPath)
	if err != nil {
		return err
	}

	if len(dirList) == 0 {
		return fmt.Errorf("no directories found for path %s", policyPath)
	}

	e.RegoFileMap = make(map[string][]byte)
	e.RegoDataMap = make(map[string]*RegoData)

	// Load rego data files from each dir
	// First, we read the metadata file, which contains info about the associated rego rule. The .rego file data is
	// stored in a map in its raw format.
	sort.Strings(dirList)
	for i := range dirList {
		// Find all files in the current dir
		var fileInfo []os.FileInfo
		fileInfo, err = ioutil.ReadDir(dirList[i])
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				zap.S().Error("error while searching for files", zap.String("dir", dirList[i]))
			}
			continue
		}

		// Load the rego metadata first (*.json)
		metadataFiles := utils.FilterFileInfoBySuffix(&fileInfo, RegoMetadataFileSuffix)
		if metadataFiles == nil {
			return fmt.Errorf("no metadata files were found")
		}

		var regoDataList []*RegoData
		for j := range *metadataFiles {
			filePath := filepath.Join(dirList[i], (*metadataFiles)[j])

			var regoMetadata *RegoMetadata
			regoMetadata, err = e.LoadRegoMetadata(filePath)
			if err != nil {
				continue
			}

			regoData := RegoData{
				Metadata: *regoMetadata,
			}

			regoDataList = append(regoDataList, &regoData)
			e.stats.metadataFileCount++
		}

		// Read in raw rego data from associated rego files
		if err = e.loadRawRegoFilesIntoMap(dirList[i], regoDataList, &e.RegoFileMap); err != nil {
			continue
		}

		for j := range regoDataList {
			e.stats.metadataCount++
			// Apply templates if available
			var templateData bytes.Buffer
			t := template.New("opa")
			_, err = t.Parse(string(e.RegoFileMap[filepath.Join(dirList[i], regoDataList[j].Metadata.RuleTemplate+".rego")]))
			if err != nil {
				zap.S().Warn("unable to parse template", zap.String("template", regoDataList[j].Metadata.RuleTemplate))
				continue
			}
			if err = t.Execute(&templateData, regoDataList[j].Metadata.RuleTemplateArgs); err != nil {
				zap.S().Warn("unable to execute template", zap.String("template", regoDataList[j].Metadata.RuleTemplate))
				continue
			}

			regoDataList[j].RawRego = templateData.Bytes()
			e.RegoDataMap[regoDataList[j].Metadata.RuleName] = regoDataList[j]
		}
	}

	e.stats.ruleCount = len(e.RegoDataMap)
	e.stats.regoFileCount = len(e.RegoFileMap)
	zap.S().Infof("loaded %d Rego rules from %d rego files (%d metadata files).", e.stats.ruleCount, e.stats.regoFileCount, e.stats.metadataFileCount)

	return err
}

// CompileRegoFiles Compiles rego files for faster evaluation
func (e *Engine) CompileRegoFiles() error {
	for k := range e.RegoDataMap {
		compiler, err := ast.CompileModules(map[string]string{
			e.RegoDataMap[k].Metadata.RuleName: string(e.RegoDataMap[k].RawRego),
		})
		if err != nil {
			return err
		}

		r := rego.New(
			rego.Query(RuleQueryBase+"."+e.RegoDataMap[k].Metadata.RuleName),
			rego.Compiler(compiler),
		)

		// Create a prepared query that can be evaluated.
		query, err := r.PrepareForEval(e.Context)
		if err != nil {
			return err
		}

		e.RegoDataMap[k].PreparedQuery = &query
	}

	return nil
}

// Initialize Initializes the Opa engine
// Handles loading all rules, filtering, compiling, and preparing for evaluation
func (e *Engine) Initialize(policyPath string) error {
	e.Context = context.Background()

	if err := e.LoadRegoFiles(policyPath); err != nil {
		return err
	}

	err := e.CompileRegoFiles()
	if err != nil {
		return err
	}

	return nil
}

// Configure Configures the OPA engine
func (e *Engine) Configure() error {
	return nil
}

// GetResults Fetches results from OPA engine policy evaluation
func (e *Engine) GetResults() error {
	return nil
}

// Release Performs any tasks required to free resources
func (e *Engine) Release() error {
	return nil
}

// Evaluate Executes compiled OPA queries against the input JSON data
func (e *Engine) Evaluate(inputData *interface{}) error {

	sortedKeys := make([]string, len(e.RegoDataMap))
	x := 0
	for k := range e.RegoDataMap {
		sortedKeys[x] = k
		x++
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		// Execute the prepared query.
		rs, err := e.RegoDataMap[k].PreparedQuery.Eval(e.Context, rego.EvalInput(inputData))
		//		rs, err := r.Eval(o.Context)
		if err != nil {
			zap.S().Warn("failed to run prepared query", zap.String("rule", "'"+k+"'"), zap.Any("input", inputData))
			continue
		}

		if len(rs) > 0 {
			results := rs[0].Expressions[0].Value.([]interface{})
			if len(results) > 0 {
				r := e.RegoDataMap[k].Metadata
				fmt.Printf("\nResource(s): %v\n[%s] [%s] %s\n    %s\n", results, r.Severity, r.RuleReferenceID, r.RuleName, r.Description)
				continue
			}
		}

		// Store results
	}

	_, err := json.MarshalIndent(inputData, "", "  ")
	if err != nil {
		return err
	}
	//fmt.Printf("InputData:\n%v\n", string(b))

	return nil
}
