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

package armv1

import (
	"fmt"
	"github.com/infracloudio/mapper"
	"strings"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

// LoadIacFile loads the specified ARM template file.
// Note that a single ARM template json file may contain multiple resource definitions.
func (a *ARMV1) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	var iacDocuments []*utils.IacDocument

	fileExt := a.getFileType(absFilePath)
	switch fileExt {
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

	for _, doc := range iacDocuments {
		var config *output.ResourceConfig = &output.ResourceConfig{}

		// TODO: make required calls to mapper API
		m := mapper.NewMapper("arm")
		err := m.Validate(doc.Data)

		if err != nil {
			zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
			continue
		}

		config.Line = doc.StartLine
		config.Source = absFilePath

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}
	return allResourcesConfig, nil
}

func (*ARMV1) getFileType(file string) string {
	if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}
