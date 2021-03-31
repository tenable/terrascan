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
	"os"
	"path/filepath"
	"reflect"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	// ErrBuildTFConfigDir error
	ErrBuildTFConfigDir = fmt.Errorf("failed to build terraform allResourcesConfig")
)

// ModuleConfig contains the *hclConfigs.Config for every module in the
// unified config tree along with *hclConfig.ModuleCall made by the parent
// module. The ParentModuleCall helps in resolving references for variables
// initilaized in the parent ModuleCall
type ModuleConfig struct {
	Config           *hclConfigs.Config
	ParentModuleCall *hclConfigs.ModuleCall
}

// LoadIacDir starts traversing from the given rootDir and traverses through
// all the descendant modules present to create an output list of all the
// resources present in rootDir and descendant modules
func LoadIacDir(absRootDir string, nonRecursiveScan bool) (allResourcesConfig output.AllResourceConfigs, err error) {

	// create a new config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	// create a new downloader to install remote modules
	r := downloader.NewRemoteDownloader()
	defer r.CleanUp()

	if nonRecursiveScan {
		return loadDirNonRecursive(absRootDir, parser, r)
	}

	// Walk the file path and find all directories
	dirList, err := utils.FindAllDirectories(absRootDir)
	if err != nil {
		return nil, err
	}
	dirList = utils.FilterHiddenDirectories(dirList, absRootDir)

	return loadDirRecursive(dirList, absRootDir, parser, r)
}

func loadDirRecursive(dirList []string, absRootDir string, parser *hclConfigs.Parser, r downloader.ModuleDownloader) (output.AllResourceConfigs, error) {

	// initialize normalized output
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	for _, dir := range dirList {
		// check if the directory has any tf config files (.tf or .tf.json)
		if !parser.IsConfigDir(dir) {
			// log a debug message and continue with other directories
			errMessage := fmt.Sprintf("directory '%s' has no terraform config files", dir)
			zap.S().Debug(errMessage)
			continue
		}

		// load current config directory
		rootMod, diags := parser.LoadConfigDir(dir)
		if diags.HasErrors() {
			// log a debug message and continue with other directories
			errMessage := fmt.Sprintf("failed to load terraform config dir '%s'. error from terraform:\n%+v\n", dir, getErrorMessagesFromDiagnostics(diags))
			zap.S().Debug(errMessage)
			continue
		}

		// get unified config for the current directory
		unified, diags := buildUnifiedConfig(rootMod, parser, dir, r)

		if diags.HasErrors() {
			// log a warn message in this case because there are errors in
			// loading the config dir, and continue with other directories
			errMessage := fmt.Sprintf("failed to build unified config. errors:\n%+v\n", diags)
			zap.S().Warnf(errMessage)
			continue
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
		root := &ModuleConfig{Config: unified.Root}
		configsQ := []*ModuleConfig{root}

		// using BFS traverse through all modules in the unified config tree
		zap.S().Debug("traversing through all modules in config tree")
		for len(configsQ) > 0 {

			// pop first element from the queue
			current := configsQ[0]
			configsQ = configsQ[1:]

			// reference resolver
			r := NewRefResolver(current.Config, current.ParentModuleCall)

			// traverse through all current's resources
			for _, managedResource := range current.Config.Module.ManagedResources {

				// create output.ResourceConfig from hclConfigs.Resource
				resourceConfig, err := CreateResourceConfig(managedResource)
				if err != nil {
					continue
				}

				// resolve references
				resourceConfig.Config = r.ResolveRefs(resourceConfig.Config.(jsonObj))

				// source file path
				resourceConfig.Source, err = filepath.Rel(absRootDir, resourceConfig.Source)
				if err != nil {
					continue
				}

				// tf plan directory relative path
				planRoot, err := filepath.Rel(absRootDir, dir)
				if err != nil {
					continue
				}
				if absRootDir == dir {
					planRoot = fmt.Sprintf(".%s", string(os.PathSeparator))
				}
				resourceConfig.PlanRoot = planRoot

				// append to normalized output
				if _, present := allResourcesConfig[resourceConfig.Type]; !present {
					allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
				} else {
					resources := allResourcesConfig[resourceConfig.Type]
					if !isConfigPresent(resources, resourceConfig) {
						allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
					}
				}
			}

			// add all current's children to the queue
			for childName, childModule := range current.Config.Children {
				childModuleConfig := &ModuleConfig{
					Config:           childModule,
					ParentModuleCall: current.Config.Module.ModuleCalls[childName],
				}
				configsQ = append(configsQ, childModuleConfig)
			}
		}
	}

	// successful
	return allResourcesConfig, nil
}

// loadDirNonRecursive has duplicate code
// this function will be removed when we deprecate non recursive scan
func loadDirNonRecursive(dir string, parser *hclConfigs.Parser, r downloader.ModuleDownloader) (output.AllResourceConfigs, error) {

	// initialize normalized output
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	// check if the directory has any tf config files (.tf or .tf.json)
	if !parser.IsConfigDir(dir) {
		// log a debug message and continue with other directories
		errMessage := fmt.Sprintf("directory '%s' has no terraform config files", dir)
		zap.S().Debug(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	// load current config directory
	rootMod, diags := parser.LoadConfigDir(dir)
	if diags.HasErrors() {
		// log a debug message and continue with other directories
		errMessage := fmt.Sprintf("failed to load terraform config dir '%s'. error from terraform:\n%+v\n", dir, getErrorMessagesFromDiagnostics(diags))
		zap.S().Debug(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	// get unified config for the current directory
	unified, diags := buildUnifiedConfig(rootMod, parser, dir, r)

	if diags.HasErrors() {
		// log a warn message in this case because there are errors in
		// loading the config dir, and continue with other directories
		errMessage := fmt.Sprintf("failed to build unified config. errors:\n%+v\n", diags)
		zap.S().Warnf(errMessage)
		return nil, ErrBuildTFConfigDir
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
	root := &ModuleConfig{Config: unified.Root}
	configsQ := []*ModuleConfig{root}

	// using BFS traverse through all modules in the unified config tree
	zap.S().Debug("traversing through all modules in config tree")
	for len(configsQ) > 0 {

		// pop first element from the queue
		current := configsQ[0]
		configsQ = configsQ[1:]

		// reference resolver
		r := NewRefResolver(current.Config, current.ParentModuleCall)

		// traverse through all current's resources
		for _, managedResource := range current.Config.Module.ManagedResources {

			// create output.ResourceConfig from hclConfigs.Resource
			resourceConfig, err := CreateResourceConfig(managedResource)
			if err != nil {
				return allResourcesConfig, fmt.Errorf("failed to create ResourceConfig")
			}

			// resolve references
			resourceConfig.Config = r.ResolveRefs(resourceConfig.Config.(jsonObj))

			// source file path
			resourceConfig.Source, err = filepath.Rel(dir, resourceConfig.Source)
			if err != nil {
				return allResourcesConfig, fmt.Errorf("failed to get resource's filepath: %s", err)
			}

			// add tf plan directory relative path
			resourceConfig.PlanRoot = fmt.Sprintf(".%s", string(os.PathSeparator))

			// append to normalized output
			if _, present := allResourcesConfig[resourceConfig.Type]; !present {
				allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
			} else {
				resources := allResourcesConfig[resourceConfig.Type]
				if !isConfigPresent(resources, resourceConfig) {
					allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
				}
			}
		}

		// add all current's children to the queue
		for childName, childModule := range current.Config.Children {
			childModuleConfig := &ModuleConfig{
				Config:           childModule,
				ParentModuleCall: current.Config.Module.ModuleCalls[childName],
			}
			configsQ = append(configsQ, childModuleConfig)
		}
	}

	// successful
	return allResourcesConfig, nil
}

func generateTempDir() string {
	return filepath.Join(os.TempDir(), utils.GenRandomString(6))
}

func processLocalSource(req *hclConfigs.ModuleRequest) string {
	// determine the absolute path from root module to the sub module

	// while building the unified config, recursive calls are made for all the paths resolved,
	// the source address in a tf file is relative while this func is called, hence, we get the
	// path of caller dir, and join the source address of current module request to get the path to module

	// get the caller dir path
	callDirPath := filepath.Dir(req.CallRange.Filename)

	// join source address to the caller dir
	return filepath.Join(callDirPath, req.SourceAddr)
}

func processTerraformRegistrySource(req *hclConfigs.ModuleRequest, tempDir string, m downloader.ModuleDownloader) (string, error) {
	// regsrc.ParseModuleSource func returns a terraform registry module source
	// error check is not required as the source address is already validated
	module, _ := regsrc.ParseModuleSource(req.SourceAddr)

	pathToModule, err := m.DownloadRemoteModule(req.VersionConstraint, tempDir, module)
	if err != nil {
		return pathToModule, err
	}

	return pathToModule, nil
}

// buildUnifiedConfig builds a unified config from *hclConfigs.Module object specified for a dir
func buildUnifiedConfig(rootMod *hclConfigs.Module, parser *hclConfigs.Parser, dir string, r downloader.ModuleDownloader) (*hclConfigs.Config, hcl.Diagnostics) {
	// using the BuildConfig and ModuleWalkerFunc to traverse through all
	// descendant modules from the root module and create a unified
	// configuration of type *configs.Config
	versionI := 0
	return hclConfigs.BuildConfig(rootMod, hclConfigs.ModuleWalkerFunc(
		func(req *hclConfigs.ModuleRequest) (*hclConfigs.Module, *version.Version, hcl.Diagnostics) {

			// figure out path sub module directory, if it's remote then download it locally
			var pathToModule string
			var err error
			if downloader.IsLocalSourceAddr(req.SourceAddr) {

				pathToModule = processLocalSource(req)
				zap.S().Debugf("processing local module %q", pathToModule)
			} else if downloader.IsRegistrySourceAddr(req.SourceAddr) {
				// temp dir to download the remote repo
				tempDir := generateTempDir()

				pathToModule, err = processTerraformRegistrySource(req, tempDir, r)
				if err != nil {
					zap.S().Errorf("failed to download remote module %q. error: '%v'", req.SourceAddr, err)
				}
			} else {
				// temp dir to download the remote repo
				tempDir := generateTempDir()

				// Download remote module
				pathToModule, err = r.DownloadModule(req.SourceAddr, tempDir)
				if err != nil {
					zap.S().Errorf("failed to download remote module %q. error: '%v'", req.SourceAddr, err)
				}
			}

			// load sub module directory
			subMod, diags := parser.LoadConfigDir(pathToModule)
			version, _ := version.NewVersion(fmt.Sprintf("1.0.%d", versionI))
			versionI++
			return subMod, version, diags
		},
	))
}

// isConfigPresent checks whether a resource is already present in the list of configs or not
// the equality of a resource is based on name, source and config of the resource
func isConfigPresent(resources []output.ResourceConfig, resourceConfig output.ResourceConfig) bool {
	for _, resource := range resources {
		if resource.Name == resourceConfig.Name && resource.Source == resourceConfig.Source {
			if reflect.DeepEqual(resource.Config, resourceConfig.Config) {
				return true
			}
		}
	}
	return false
}
