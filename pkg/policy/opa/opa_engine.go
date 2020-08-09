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

type Violation struct {
	Name        string
	Description string
	LineNumber  int
	Category    string
	Data        interface{}
	RuleData    interface{}
}

type ResultData struct {
	EngineType string
	Provider   string
	Violations []*Violation
}

type RegoMetadata struct {
	RuleName         string                 `json:"ruleName"`
	File             string                 `json:"file"`
	RuleTemplate     string                 `json:"ruleTemplate"`
	RuleTemplateArgs map[string]interface{} `json:"ruleTemplateArgs"`
	Severity         string                 `json:"severity"`
	Description      string                 `json:"description"`
	RuleReferenceID  string                 `json:"ruleReferenceId"`
	Category         string                 `json:"category"`
	Version          int                    `json:"version"`
}

type RegoData struct {
	Metadata      RegoMetadata
	RawRego       []byte
	PreparedQuery *rego.PreparedEvalQuery
}

type EngineStats struct {
	ruleCount         int
	regoFileCount     int
	metadataFileCount int
	metadataCount     int
}

type OpaEngine struct {
	Context     context.Context
	RegoFileMap map[string][]byte
	RegoDataMap map[string]*RegoData
	stats       EngineStats
}

func (o *OpaEngine) LoadRegoMetadata(metaFilename string) (*RegoMetadata, error) {
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

func (o *OpaEngine) loadRawRegoFilesIntoMap(currentDir string, regoDataList []*RegoData, regoFileMap *map[string][]byte) error {
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

func (o *OpaEngine) LoadRegoFiles(policyPath string) error {
	// Walk the file path and find all directories
	dirList, err := utils.FindAllDirectories(policyPath)
	if err != nil {
		return err
	}

	if len(dirList) == 0 {
		return fmt.Errorf("no directories found for path %s", policyPath)
	}

	o.RegoFileMap = make(map[string][]byte)
	o.RegoDataMap = make(map[string]*RegoData)

	// Load rego data files from each dir
	// First, we read the metadata file, which contains info about the associated rego rule. The .rego file data is
	// stored in a map in its raw format.
	sort.Strings(dirList)
	for i := range dirList {
		// Find all files in the current dir
		fileInfo, err := ioutil.ReadDir(dirList[i])
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
			regoMetadata, err = o.LoadRegoMetadata(filePath)
			if err != nil {
				continue
			}

			regoData := RegoData{
				Metadata: *regoMetadata,
			}

			regoDataList = append(regoDataList, &regoData)
			o.stats.metadataFileCount++
		}

		// Read in raw rego data from associated rego files
		if err = o.loadRawRegoFilesIntoMap(dirList[i], regoDataList, &o.RegoFileMap); err != nil {
			continue
		}

		for j := range regoDataList {
			o.stats.metadataCount++
			// Apply templates if available
			var templateData bytes.Buffer
			t := template.New("opa")
			t.Parse(string(o.RegoFileMap[filepath.Join(dirList[i], regoDataList[j].Metadata.RuleTemplate+".rego")]))
			t.Execute(&templateData, regoDataList[j].Metadata.RuleTemplateArgs)

			regoDataList[j].RawRego = templateData.Bytes()
			o.RegoDataMap[regoDataList[j].Metadata.RuleName] = regoDataList[j]
		}
	}

	o.stats.ruleCount = len(o.RegoDataMap)
	zap.S().Infof("Loaded %d Rego rules from %d rego files (%d metadata files).", o.stats.ruleCount, o.stats.regoFileCount, o.stats.metadataCount)

	return err
}

func (o *OpaEngine) CompileRegoFiles() error {
	for k := range o.RegoDataMap {
		compiler, err := ast.CompileModules(map[string]string{
			o.RegoDataMap[k].Metadata.RuleName: string(o.RegoDataMap[k].RawRego),
		})

		r := rego.New(
			rego.Query(RuleQueryBase+"."+o.RegoDataMap[k].Metadata.RuleName),
			rego.Compiler(compiler),
		)

		// Create a prepared query that can be evaluated.
		query, err := r.PrepareForEval(o.Context)
		if err != nil {
			return err
		}

		o.RegoDataMap[k].PreparedQuery = &query
	}

	return nil
}

// Initialize Initializes the Opa engine
// Handles loading all rules, filtering, compiling, and preparing for evaluation
func (o *OpaEngine) Initialize(policyPath string) error {
	o.Context = context.Background()

	if err := o.LoadRegoFiles(policyPath); err != nil {
		return err
	}

	err := o.CompileRegoFiles()
	if err != nil {
		return err
	}

	return nil
}

func (o *OpaEngine) Configure() error {
	return nil
}

func (o *OpaEngine) GetResults() error {
	return nil
}

func (o *OpaEngine) Release() error {
	return nil
}

func (o *OpaEngine) Evaluate(inputData *interface{}) error {

	sortedKeys := make([]string, len(o.RegoDataMap))
	x := 0
	for k := range o.RegoDataMap {
		sortedKeys[x] = k
		x++
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		// Execute the prepared query.
		rs, err := o.RegoDataMap[k].PreparedQuery.Eval(o.Context, rego.EvalInput(inputData))
		//		rs, err := r.Eval(o.Context)
		if err != nil {
			zap.S().Warn("failed to run prepared query", zap.String("rule", "'"+k+"'"), zap.Any("input", inputData))
			continue
		}

		if len(rs) > 0 {
			results := rs[0].Expressions[0].Value.([]interface{})
			if len(results) > 0 {
				r := o.RegoDataMap[k].Metadata
				fmt.Printf("\nResource(s): %v\n[%s] [%s] %s\n    %s\n", results, r.Severity, r.RuleReferenceID, r.RuleName, r.Description)
				continue
			}
			//			fmt.Printf("   [%s] %v\n", k, results)
		} else {
			//			fmt.Printf("No Result [%s] \n", k)
		}

		// Store results
	}

	b, _ := json.MarshalIndent(inputData, "", "  ")
	//fmt.Printf("InputData:\n%v\n", string(b))

	return nil
}
