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

package cftv1

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/infracloudio/mapper"
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
			continue
		}
		for t, resource := range arc {
			config = &resource[0]
			config.Type = t
			config.Source = absFilePath
			allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
			break
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
	}
	return UnknownExtension
}
