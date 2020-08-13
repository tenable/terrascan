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
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	errEmptyTFConfigDir = fmt.Errorf("directory has no terraform files")
	errLoadConfigDir    = fmt.Errorf("failed to load terraform allResourcesConfig dir")
	errBuildTFConfigDir = fmt.Errorf("failed to build terraform allResourcesConfig")
)

// LoadIacDir starts traversing from the given rootDir and traverses through
// all the descendant modules present to create an output list of all the
// resources present in rootDir and descendant modules
func (*TfV12) LoadIacDir(absRootDir string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// create a new config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	// check if the directory has any tf config files (.tf or .tf.json)
	if !parser.IsConfigDir(absRootDir) {
		zap.S().Errorf("directory '%s' has no terraform config files", absRootDir)
		return allResourcesConfig, errEmptyTFConfigDir
	}

	// load root config directory
	rootMod, diags := parser.LoadConfigDir(absRootDir)
	if diags.HasErrors() {
		zap.S().Errorf("failed to load terraform config dir '%s'. error:\n%+v\n", absRootDir, diags)
		return allResourcesConfig, errLoadConfigDir
	}

	// using the BuildConfig and ModuleWalkerFunc to traverse through all
	// descendant modules from the root module and create a unified
	// configuration of type *configs.Config
	// Note: currently, only Local paths are supported for Module Sources
	versionI := 0
	unified, diags := hclConfigs.BuildConfig(rootMod, hclConfigs.ModuleWalkerFunc(
		func(req *hclConfigs.ModuleRequest) (*hclConfigs.Module, *version.Version, hcl.Diagnostics) {

			// Note: currently only local paths are supported for Module Sources

			// determine the absolute path from root module to the sub module
			// using *configs.ModuleRequest.Path field
			var (
				pathArr      = strings.Split(req.Path.String(), ".")
				pathToModule = absRootDir
			)
			for _, subPath := range pathArr {
				pathToModule = filepath.Join(pathToModule, subPath)
			}

			// load sub module directory
			subMod, diags := parser.LoadConfigDir(pathToModule)
			version, _ := version.NewVersion(fmt.Sprintf("1.0.%d", versionI))
			versionI++
			return subMod, version, diags
		},
	))
	if diags.HasErrors() {
		zap.S().Errorf("failed to build unified config. errors:\n%+v\n", diags)
		return allResourcesConfig, errBuildTFConfigDir
	}

	/*
		The "unified" config created from BuildConfig in the previous step
		represents a tree structure with rootDir module being at its root and
		all the sub modules being its children, and these children can have
		more children and so on...

		Now, using BFS we traverse through all the submodules using the classic
		approach of using a queue data structure
	*/

	// queue of for BFS, add root module config to it
	configsQ := []*hclConfigs.Config{unified.Root}

	// initialize normalized output
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	// using BFS traverse through all modules in the unified config tree
	zap.S().Debug("traversing through all modules in config tree")
	for len(configsQ) > 0 {

		// pop first element from the queue
		current := configsQ[0]
		configsQ = configsQ[1:]

		// traverse through all current's resources
		for _, managedResource := range current.Module.ManagedResources {

			// create output.ResourceConfig from hclConfigs.Resource
			resourceConfig, err := CreateResourceConfig(managedResource)
			if err != nil {
				return allResourcesConfig, fmt.Errorf("failed to create ResourceConfig")
			}

			// trimFilePath
			resourceConfig.Source = trimFilePath(resourceConfig.Source, absRootDir)

			// append to normalized output
			if _, present := allResourcesConfig[resourceConfig.Type]; !present {
				allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
			} else {
				allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
			}
		}

		// add all current's children to the queue
		for _, childModule := range current.Children {
			configsQ = append(configsQ, childModule)
		}
	}

	// successful
	return allResourcesConfig, nil
}

// trimFilePath returns relative file path wrt to the base path
func trimFilePath(fullPath, basePath string) string {
	basePath = strings.TrimSuffix(basePath, "/")
	splits := bytes.Split([]byte(fullPath), []byte(basePath))
	return strings.TrimPrefix(string(splits[1]), "/")
}
