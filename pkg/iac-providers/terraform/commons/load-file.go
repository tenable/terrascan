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

package commons

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

// LoadIacFile parses the given terraform file from the given file path
func LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// new terraform config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	hclFile, diags := parser.LoadConfigFile(absFilePath)
	if hclFile == nil {
		errMessage := fmt.Sprintf("error occured while loading config file '%s'. error:\n%v\n", absFilePath, getErrorMessagesFromDiagnostics(diags))
		zap.S().Debug(errMessage)
		return allResourcesConfig, fmt.Errorf(errMessage)
	}
	if diags != nil {
		errMessage := fmt.Sprintf("failed to load config file '%s'. error:\n%v\n", absFilePath, getErrorMessagesFromDiagnostics(diags))
		zap.S().Debug(errMessage)
		return allResourcesConfig, fmt.Errorf(errMessage)
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

		// set module name
		// module name for the file scan will always be root
		resourceConfig.ModuleName = "root"

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

// getErrorMessagesFromDiagnostics should be called when diags.HasErrors is true
func getErrorMessagesFromDiagnostics(diags hcl.Diagnostics) string {
	var errMsgs []string
	for _, v := range diags.Errs() {
		errMsgs = append(errMsgs, v.Error())
	}
	return strings.Join(errMsgs, "\n")
}
