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
	"strings"

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
func LoadIacDir(absRootDir string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// map to hold the download paths of remote modules
	remoteModPaths := make(map[string]string)

	// create a new config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	// check if the directory has any tf config files (.tf or .tf.json)
	if !parser.IsConfigDir(absRootDir) {
		errMessage := fmt.Sprintf("directory '%s' has no terraform config files", absRootDir)
		zap.S().Debug(errMessage)
		return allResourcesConfig, fmt.Errorf(errMessage)
	}

	// load root config directory
	rootMod, diags := parser.LoadConfigDir(absRootDir)
	if diags.HasErrors() {
		errMessage := fmt.Sprintf("failed to load terraform config dir '%s'. error:\n%+v\n", absRootDir, diags)
		zap.S().Debug(errMessage)
		return allResourcesConfig, fmt.Errorf(errMessage)
	}

	// create a new downloader to install remote modules
	r := downloader.NewRemoteDownloader()
	defer r.CleanUp()

	// using the BuildConfig and ModuleWalkerFunc to traverse through all
	// descendant modules from the root module and create a unified
	// configuration of type *configs.Config
	// Note: currently, only Local paths are supported for Module Sources
	versionI := 0
	unified, diags := hclConfigs.BuildConfig(rootMod, hclConfigs.ModuleWalkerFunc(
		func(req *hclConfigs.ModuleRequest) (*hclConfigs.Module, *version.Version, hcl.Diagnostics) {

			// figure out path sub module directory, if it's remote then download it locally
			var pathToModule string
			if downloader.IsLocalSourceAddr(req.SourceAddr) {

				pathToModule = processLocalSource(req, remoteModPaths, absRootDir)
				zap.S().Debugf("processing local module %q", pathToModule)
			} else if downloader.IsRegistrySourceAddr(req.SourceAddr) {
				// temp dir to download the remote repo
				tempDir := generateTempDir()

				pathToModule, err = processTerraformRegistrySource(req, remoteModPaths, tempDir, r)
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
	if diags.HasErrors() {
		zap.S().Errorf("failed to build unified config. errors:\n%+v\n", diags)
		return allResourcesConfig, ErrBuildTFConfigDir
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

	// initialize normalized output
	allResourcesConfig = make(map[string][]output.ResourceConfig)

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
			resourceConfig.Source, err = filepath.Rel(absRootDir, resourceConfig.Source)
			if err != nil {
				return allResourcesConfig, fmt.Errorf("failed to get resource: %s", err)
			}

			// append to normalized output
			if _, present := allResourcesConfig[resourceConfig.Type]; !present {
				allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
			} else {
				allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
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

func processLocalSource(req *hclConfigs.ModuleRequest, remoteModPaths map[string]string, absRootDir string) string {
	// determine the absolute path from root module to the sub module
	// since we start at the end of the path, we need to assemble
	// the parts in reverse order

	pathToModule := req.SourceAddr
	for p := req.Parent; p != nil; p = p.Parent {
		pathToModule = filepath.Join(p.SourceAddr, pathToModule)
	}

	// check if pathToModule consists of a remote module downloaded
	// if yes, then update the module path, with the remote module
	// download path
	keyFound := false
	for remoteSourceAddr, downloadPath := range remoteModPaths {
		if strings.Contains(pathToModule, remoteSourceAddr) {
			pathToModule = strings.Replace(pathToModule, remoteSourceAddr, "", 1)
			pathToModule = filepath.Join(downloadPath, pathToModule)
			keyFound = true
			break
		}
	}
	if !keyFound {
		pathToModule = filepath.Join(absRootDir, pathToModule)
	}
	return pathToModule
}

func processTerraformRegistrySource(req *hclConfigs.ModuleRequest, remoteModPaths map[string]string, tempDir string, m downloader.ModuleDownloader) (string, error) {
	// regsrc.ParseModuleSource func returns a terraform registry module source
	// error check is not required as the source address is already validated
	module, _ := regsrc.ParseModuleSource(req.SourceAddr)

	pathToModule, err := m.DownloadRemoteModule(req.VersionConstraint, tempDir, module)
	if err != nil {
		return pathToModule, err
	}

	// add the entry of remote module's source address to the map
	remoteModPaths[req.SourceAddr] = pathToModule
	return pathToModule, nil
}
