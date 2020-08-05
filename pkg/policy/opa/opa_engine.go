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

package policy

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
	"strings"
	"text/template"

	"github.com/open-policy-agent/opa/ast"

	"go.uber.org/zap"

	"github.com/open-policy-agent/opa/rego"
)

type AccuricsRegoMetadata struct {
	Name             string                 `json:"ruleName"`
	DisplayName      string                 `json:"ruleDisplayName"`
	Category         string                 `json:"category"`
	ImpactedRes      []string               `json:"impactedRes"`
	PolicyRelevance  string                 `json:"policyRelevance"`
	Remediation      string                 `json:"remediation"`
	Row              int                    `json:"row"`
	Rule             string                 `json:"rule"`
	RuleTemplate     string                 `json:"ruleTemplate"`
	RuleTemplateArgs map[string]interface{} `json:"ruleArgument"`
	RuleReferenceID  string                 `json:"ruleReferenceId"`
	Severity         string                 `json:"severity"`
	Vulnerability    string                 `json:"vulnerability"`
}

type RegoData struct {
	Name             string                 `json:"ruleName"`
	DisplayName      string                 `json:"ruleDisplayName"`
	Category         string                 `json:"category"`
	Remediation      string                 `json:"remediation"`
	Rule             string                 `json:"rule"`
	RuleTemplate     string                 `json:"ruleTemplate"`
	RuleTemplateArgs map[string]interface{} `json:"ruleArgument"`
	RuleReferenceID  string                 `json:"ruleReferenceId"`
	Severity         string                 `json:"severity"`
	Vulnerability    string                 `json:"vulnerability"`
	RawRego          *[]byte
	PreparedQuery    *rego.PreparedEvalQuery
}

type ResultData struct {
}

type OpaEngine struct {
	Context     context.Context
	RegoFileMap map[string][]byte
	RegoDataMap map[string]*RegoData
}

func filterFileListBySuffix(allFileList *[]string, filter string) *[]string {
	fileList := make([]string, 0)

	for i := range *allFileList {
		if strings.HasSuffix((*allFileList)[i], filter) {
			fileList = append(fileList, (*allFileList)[i])
		}
	}
	return &fileList
}

func (o *OpaEngine) LoadRegoFiles(policyPath string) error {
	ruleCount := 0
	regoFileCount := 0
	metadataCount := 0

	// Walk the file path and find all directories
	dirList := make([]string, 0)
	err := filepath.Walk(policyPath, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo != nil && fileInfo.IsDir() {
			dirList = append(dirList, filePath)
		}
		return err
	})

	if len(dirList) == 0 {
		return fmt.Errorf("no directories found for path %s", policyPath)
	}

	o.RegoFileMap = make(map[string][]byte)
	o.RegoDataMap = make(map[string]*RegoData)

	// Load rego data files from each dir
	sort.Strings(dirList)
	for i := range dirList {
		metaFilename := filepath.Join(dirList[i], RegoMetadataFile)
		var metadata []byte
		metadata, err = ioutil.ReadFile(metaFilename)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				zap.S().Warn("failed to load rego metadata", zap.String("file", metaFilename))
			}
			continue
		}

		// Read metadata into struct
		regoMetadata := make([]*RegoData, 0)
		if err = json.Unmarshal(metadata, &regoMetadata); err != nil {
			zap.S().Warn("failed to unmarshal rego metadata", zap.String("file", metaFilename))
			continue
		}

		metadataCount++

		// Find all .rego files within the directory
		fileInfo, err := ioutil.ReadDir(dirList[i])
		if err != nil {
			zap.S().Error("error while finding rego files", zap.String("dir", dirList[i]))
			continue
		}

		files := make([]string, 0)
		for j := range fileInfo {
			files = append(files, fileInfo[j].Name())
		}

		// Load rego data for all rego files
		regoFileList := filterFileListBySuffix(&files, RegoFileSuffix)
		regoFileCount += len(*regoFileList)
		for j := range *regoFileList {
			regoFilename := (*regoFileList)[j]
			regoFullPath := filepath.Join(dirList[i], regoFilename)
			var rawRegoData []byte
			rawRegoData, err = ioutil.ReadFile(regoFullPath)
			if err != nil {
				zap.S().Warn("failed to load rego file", zap.String("file", regoFilename))
				continue
			}

			_, ok := o.RegoFileMap[regoFullPath]
			if ok {
				// Already loaded this file, so continue
				continue
			}

			// Set raw rego data
			o.RegoFileMap[regoFullPath] = rawRegoData
		}

		for j := range regoMetadata {
			//key := filepath.Join(dirList[i], regoMetadata[j].Rule)
			//regoData := o.RegoFileMap[key]
			metadataCount++
			// Apply templates if available
			var buf bytes.Buffer
			t := template.New("opa")
			t.Parse(string(o.RegoFileMap[filepath.Join(dirList[i], regoMetadata[j].RuleTemplate+".rego")]))
			t.Execute(&buf, regoMetadata[j].RuleTemplateArgs)

			templateData := buf.Bytes()
			regoMetadata[j].RawRego = &templateData
			o.RegoDataMap[regoMetadata[j].Name] = regoMetadata[j]
		}
	}

	ruleCount = len(o.RegoDataMap)
	zap.S().Infof("Loaded %d Rego rules from %d rego files (%d metadata files).", ruleCount, regoFileCount, metadataCount)

	return err
}

func (o *OpaEngine) CompileRegoFiles() error {
	for k := range o.RegoDataMap {
		compiler, err := ast.CompileModules(map[string]string{
			o.RegoDataMap[k].Rule: string(*(o.RegoDataMap[k].RawRego)),
		})

		r := rego.New(
			rego.Query(RuleQueryBase+"."+o.RegoDataMap[k].Name),
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
				r := o.RegoDataMap[k]
				fmt.Printf("\n[%s] [%s] %s\n    %s\n", r.Severity, r.RuleReferenceID, r.DisplayName, r.Vulnerability)
			}
			//			fmt.Printf("   [%s] %v\n", k, results)
		} else {
			//			fmt.Printf("No Result [%s] \n", k)
		}
		// Store results
	}

	return nil
}
