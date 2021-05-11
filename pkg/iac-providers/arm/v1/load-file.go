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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/mapper/core"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// LoadIacFile loads the specified ARM template file.
// Note that a single ARM template json file may contain multiple resource definitions.
func (a *ARMV1) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
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

	m := arm.Mapper()
	for _, doc := range iacDocuments {
		template, err := a.extractTemplate(doc)
		if err != nil {
			zap.S().Debug("unable to parse template", zap.Error(err), zap.String("file", absFilePath))
			continue
		}

		// set template parameters with default values if not found
		if a.templateParameters == nil {
			a.templateParameters = make(map[string]interface{})
		}
		for key, param := range template.Parameters {
			if _, ok := a.templateParameters[key]; !ok {
				a.templateParameters[key] = param.DefaultValue
			}
		}

		for _, r := range template.Resources {
			config := a.getConfig(doc, absFilePath, m, r, template.Variables)
			if config == nil {
				continue
			}

			cf, ok := config.Config.(map[string]interface{})
			if !ok {
				zap.S().Debug("unable to parse config.Config data",
					zap.String("resource", r.Type), zap.String("file", absFilePath),
				)
				continue
			}

			for _, nr := range r.Resources {
				if !strings.HasPrefix(nr.Type, "Microsoft.") {
					nr.Type = r.Type + "/" + nr.Type
				}

				config := a.getConfig(doc, absFilePath, m, nr, template.Variables)
				if config != nil {
					cf[config.Type] = config.Config
				}
			}
			allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
		}
	}
	return allResourcesConfig, nil
}

func (*ARMV1) getFileType(file string) string {
	if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}

func (ARMV1) extractTemplate(doc *utils.IacDocument) (*types.Template, error) {
	const errUnsupportedDoc = "unsupported document type"

	if doc.Type == utils.JSONDoc {
		var t types.Template
		err := json.Unmarshal(doc.Data, &t)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, errors.New(errUnsupportedDoc)
}

// getSourceRelativePath fetches the relative path of file being loaded
func (a *ARMV1) getSourceRelativePath(sourceFile string) string {

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

func (a *ARMV1) getConfig(doc *utils.IacDocument, path string, m core.Mapper, r types.Resource,
	vars map[string]interface{}) *output.ResourceConfig {
	if _, ok := types.ResourceTypes[r.Type]; !ok {
		return nil
	}

	config := &output.ResourceConfig{
		Line:   doc.StartLine,
		Source: a.getSourceRelativePath(path),
	}
	err := m.Map(r, config, vars, a.templateParameters)
	if err != nil {
		zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", path))
		return nil
	}
	return config
}
