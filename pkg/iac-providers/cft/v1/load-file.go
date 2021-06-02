/*
    Copyright (C) 2021 Accurics, Inc.

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

package cftv1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/mapper"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// LoadIacFile loads the specified CFT template file.
// Note that a single CFT template json file may contain multiple resource definitions.
func (a *CFTV1) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	var iacDocuments []*utils.IacDocument
	fileExt := a.getFileType(absFilePath)
	switch fileExt {
	case YAMLExtension:
		fallthrough
	case YAMLExtension2:
		iacDocuments, err = utils.LoadYAML(absFilePath)
	case JSONExtension:
		iacDocuments, err = utils.LoadJSON(absFilePath)
	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return allResourcesConfig, fmt.Errorf("unknown file extension for file %s", absFilePath)
	}
	if err != nil {
		zap.S().Debug("failed to load file", zap.String("file", absFilePath))
		return allResourcesConfig, err
	}
	allResourcesConfig = make(map[string][]output.ResourceConfig)
	for _, doc := range iacDocuments {

		// replacing yaml data as the default yaml.v3 removes
		// intrinsic tags for cloudformation templates
		// (!Ref, !Fn::<> etc are removed and resolved to a string
		// which disables parameter resolution by goformation)
		if fileExt != JSONExtension {
			templateData, err := ioutil.ReadFile(absFilePath)
			if err != nil {
				zap.S().Debug("unable to read template data", zap.Error(err), zap.String("file", absFilePath))
				return allResourcesConfig, err
			}
			doc.Data = templateData
		}

		var config *output.ResourceConfig
		m := mapper.NewMapper("cft")
		arc, err := m.Map(doc)
		if err != nil {
			zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
			return allResourcesConfig, err
		}
		for t, resources := range arc {
			for _, resource := range resources {
				config = &resource
				config.Type = t
				config.Source = a.getSourceRelativePath(absFilePath)
				allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
			}
		}
	}
	return allResourcesConfig, nil
}

func (*CFTV1) getFileType(file string) string {
	if strings.HasSuffix(file, YAMLExtension) {
		return YAMLExtension
	} else if strings.HasSuffix(file, YAMLExtension2) {
		return YAMLExtension2
	} else if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	} else if strings.HasSuffix(file, TXTExtension) || strings.HasSuffix(file, TemplateExtension) {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			zap.S().Debug("unable to read file", zap.Error(err), zap.String("file", file))
			return UnknownExtension
		}
		if isJSON(string(f)) {
			return JSONExtension
		}
		return YAMLExtension
	}
	return UnknownExtension
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// getSourceRelativePath fetches the relative path of file being loaded
func (a *CFTV1) getSourceRelativePath(sourceFile string) string {

	// rootDir should be empty when file scan was initiated by user
	if a.absRootDir == "" {
		return filepath.Base(sourceFile)
	}
	relPath, err := filepath.Rel(a.absRootDir, sourceFile)
	if err != nil {
		zap.S().Debug("error while getting the relative path for", zap.String("IAC file", sourceFile), zap.Error(err))
		return sourceFile
	}
	return relPath
}
