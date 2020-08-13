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

package tfv12

import (
	"fmt"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	errLoadConfigFile = fmt.Errorf("failed to load config file")
)

// LoadIacFile parses the given terraform file from the given file path
func (*TfV12) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// new terraform config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	hclFile, diags := parser.LoadConfigFile(absFilePath)
	if diags != nil {
		zap.S().Errorf("failed to load config file '%s'. error:\n%v\n", absFilePath, diags)
		return allResourcesConfig, errLoadConfigFile
	}
	if hclFile == nil && diags.HasErrors() {
		zap.S().Errorf("error occured while loading config file. error:\n%v\n", diags)
		return allResourcesConfig, errLoadConfigFile
	}

	// initialize normalized output
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	// traverse through all current's resources
	for _, managedResource := range hclFile.ManagedResources {

		// create output.ResourceConfig from hclConfigs.Resource
		resourceConfig, err := CreateResourceConfig(managedResource)
		if err != nil {
			return allResourcesConfig, fmt.Errorf("failed to create ResourceConfig")
		}

		// extract file name from path
		resourceConfig.Source = getFileName(resourceConfig.Source)

		// append to normalized output
		if _, present := allResourcesConfig[resourceConfig.Type]; !present {
			allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
		} else {
			allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
		}
	}

	// successful
	return allResourcesConfig, nil
}

// getFileName return file name from the given file path
func getFileName(path string) string {
	_, file := filepath.Split(path)
	return file
}
