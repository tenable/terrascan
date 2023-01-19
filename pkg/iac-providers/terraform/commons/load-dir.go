/*
    Copyright (C) 2022 Tenable, Inc.

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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	version "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/spf13/afero"
	"github.com/tenable/terrascan/pkg/downloader"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

var (
	// ErrBuildTFConfigDir error
	ErrBuildTFConfigDir = fmt.Errorf("failed to build terraform allResourcesConfig")
)

const (
	terraformModuleInstallDir             = ".terraform/modules"
	terraformInstalledModulelMetaFileName = "modules.json"
)

// TerraformInstalledModuleMetaData metadata about the module downloaded and present in terraform cache.
type TerraformInstalledModuleMetaData struct {
	Key        string `json:"Key"`
	SourceAddr string `json:"Source"`
	VersionStr string `json:"Version,omitempty"`
	Dir        string `json:"Dir"`
}

// TerraformModuleManifest holds details of all modules downloaded by terraform
type TerraformModuleManifest struct {
	Modules []TerraformInstalledModuleMetaData `json:"Modules"`
}

// ModuleConfig contains the *hclConfigs.Config for every module in the
// unified config tree along with *hclConfig.ModuleCall made by the parent
// module. The ParentModuleCall helps in resolving references for variables
// initilaized in the parent ModuleCall
type ModuleConfig struct {
	Config           *hclConfigs.Config
	ParentModuleCall *hclConfigs.ModuleCall
	Name             string
}

// TerraformDirectoryLoader implements terraform directory loading
type TerraformDirectoryLoader struct {
	absRootDir               string
	nonRecursive             bool
	useTerraformCache        bool
	remoteDownloader         downloader.ModuleDownloader
	parser                   *hclConfigs.Parser
	errIacLoadDirs           *multierror.Error
	terraformInitModuleCache map[string]TerraformModuleManifest
	terraformVersion         string
}

// NewTerraformDirectoryLoader creates a new terraformDirectoryLoader
func NewTerraformDirectoryLoader(rootDirectory, terraformVersion string, options map[string]interface{}) TerraformDirectoryLoader {
	terraformDirectoryLoader := TerraformDirectoryLoader{
		absRootDir:               rootDirectory,
		remoteDownloader:         downloader.NewRemoteDownloader(),
		parser:                   hclConfigs.NewParser(afero.NewOsFs()),
		terraformInitModuleCache: make(map[string]TerraformModuleManifest),
		terraformVersion:         terraformVersion,
	}
	for key, val := range options {
		// keeping switch case in case more flags are added
		switch key {
		case "useTerraformCache":
			terraformDirectoryLoader.useTerraformCache = val.(bool)
		case "nonRecursive":
			terraformDirectoryLoader.nonRecursive = val.(bool)
		}
	}
	return terraformDirectoryLoader
}

// LoadIacDir starts traversing from the given rootDir and traverses through
// all the descendant modules present to create an output list of all the
// resources present in rootDir and descendant modules
func (t TerraformDirectoryLoader) LoadIacDir() (allResourcesConfig output.AllResourceConfigs, err error) {

	defer t.remoteDownloader.CleanUp()

	if t.nonRecursive || t.useTerraformCache {
		return t.loadDirNonRecursive()
	}

	// Walk the file path and find all directories
	dirList, err := utils.FindAllDirectories(t.absRootDir)
	if err != nil {
		return nil, multierror.Append(t.errIacLoadDirs, err)
	}
	dirList = utils.FilterHiddenDirectories(dirList, t.absRootDir)

	return t.loadDirRecursive(dirList)
}

func (t TerraformDirectoryLoader) loadDirRecursive(dirList []string) (output.AllResourceConfigs, error) {

	// initialize normalized output
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	for _, dir := range dirList {
		// check if the directory has any tf config files (.tf or .tf.json)
		if !t.parser.IsConfigDir(dir) {
			// log a debug message and continue with other directories
			errMessage := fmt.Sprintf("directory '%s' has no terraform config files", dir)
			zap.S().Debug(errMessage)
			t.addError(errMessage, dir)
			continue
		}

		// load current config directory
		rootMod, diags := t.parser.LoadConfigDir(dir)
		if rootMod == nil && diags.HasErrors() {
			// log a debug message and continue with other directories
			errMessage := fmt.Sprintf("failed to load terraform config dir '%s'. error from terraform:\n%+v\n", dir, getErrorMessagesFromDiagnostics(diags))
			zap.S().Debug(errMessage)
			t.addError(errMessage, dir)
			continue
		}

		// rootMod can be considered for static analysis
		if rootMod != nil && diags.HasErrors() {
			// log a debug message and continue with analysis of the root mod.
			errMessage := fmt.Sprintf("diagnostic errors while loading terraform config dir '%s'. error from terraform:\n%+v\n", dir, getErrorMessagesFromDiagnostics(diags))
			zap.S().Debug(errMessage)
			t.addError(errMessage, dir)
		}

		// getting provider version for the root module
		providerVersion := GetModuleProviderVersion(rootMod)

		// get unified config for the current directory
		unified, diags := t.buildUnifiedConfig(rootMod, dir)
		// Get the downloader chache
		remoteURLMapping := t.remoteDownloader.GetDownloaderCache()

		if diags.HasErrors() {
			// log a warn message in this case because there are errors in
			// loading the config dir, and continue with other directories
			errMessage := fmt.Sprintf("failed to build unified config. errors:\n%+v\n", getErrorMessagesFromDiagnostics(diags))
			zap.S().Warnf(errMessage)
			t.addError(errMessage, dir)
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
		root := &ModuleConfig{Config: unified.Root, Name: "root"}
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
					t.addError(err.Error(), dir)
					continue
				}

				resourceConfig.TerraformVersion = t.terraformVersion
				resourceConfig.ProviderVersion = providerVersion

				// if root module do not have provider contraints fetch the latest compatible version
				if resourceConfig.ProviderVersion == "" {
					resourceConfig.ProviderVersion = LatestProviderVersion(managedResource.Provider, t.terraformVersion)
				}
				// set module name
				resourceConfig.ModuleName = current.Name

				// resolve references
				resourceConfig.Config = r.ResolveRefs(resourceConfig.Config.(jsonObj))

				var isRemoteModule bool
				// source file path
				resourceConfig.Source, isRemoteModule, err = GetConfigSource(remoteURLMapping, resourceConfig, t.absRootDir)
				if err != nil {
					t.addError(err.Error(), dir)
					continue
				}
				if isRemoteModule {
					resourceConfig.IsRemoteModule = &isRemoteModule
				}

				// tf plan directory relative path
				planRoot, err := filepath.Rel(t.absRootDir, dir)
				if err != nil {
					t.addError(err.Error(), dir)
					continue
				}
				if t.absRootDir == dir {
					planRoot = fmt.Sprintf(".%s", string(os.PathSeparator))
				}
				resourceConfig.PlanRoot = planRoot

				// append to normalized output
				if _, present := allResourcesConfig[resourceConfig.Type]; !present {
					allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
				} else {
					resources := allResourcesConfig[resourceConfig.Type]
					if !output.IsConfigPresent(resources, resourceConfig) {
						allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
					}
				}
			}

			// add all current's children to the queue
			configsQ = append(configsQ, current.getChildConfigs()...)
		}
	}

	// successful
	return allResourcesConfig, t.errIacLoadDirs
}

// loadDirNonRecursive has duplicate code
// this function will be removed when we deprecate non recursive scan
func (t TerraformDirectoryLoader) loadDirNonRecursive() (output.AllResourceConfigs, error) {

	// initialize normalized output
	allResourcesConfig := make(map[string][]output.ResourceConfig)

	// check if the directory has any tf config files (.tf or .tf.json)
	if !t.parser.IsConfigDir(t.absRootDir) {
		// log a debug message and continue with other directories
		errMessage := fmt.Sprintf("directory '%s' has no terraform config files", t.absRootDir)
		zap.S().Debug(errMessage)
		return nil, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "terraform", Directory: t.absRootDir, ErrMessage: errMessage})
	}

	// load current config directory
	rootMod, diags := t.parser.LoadConfigDir(t.absRootDir)
	if rootMod == nil && diags.HasErrors() {
		// log a debug message and continue with other directories
		errMessage := fmt.Sprintf("failed to load terraform config dir '%s'. error from terraform:\n%+v\n", t.absRootDir, getErrorMessagesFromDiagnostics(diags))
		zap.S().Debug(errMessage)
		return nil, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "terraform", Directory: t.absRootDir, ErrMessage: errMessage})
	}

	// rootMod can be considered for static analysis
	if rootMod != nil && diags.HasErrors() {
		// log a debug message and continue with analysis of the root mod.
		errMessage := fmt.Sprintf("diagnostic errors while loading terraform config dir '%s'. error from terraform:\n%+v\n", t.absRootDir, getErrorMessagesFromDiagnostics(diags))
		zap.S().Debug(errMessage)
		t.addError(errMessage, t.absRootDir)
	}

	// getting provider version for the root module
	providerVersion := GetModuleProviderVersion(rootMod)

	// get unified config for the current directory
	unified, diags := t.buildUnifiedConfig(rootMod, t.absRootDir)

	// Get the downloader chache
	remoteURLMapping := t.remoteDownloader.GetDownloaderCache()

	if diags.HasErrors() {
		// log a warn message in this case because there are errors in
		// loading the config dir, and continue with other directories
		errMessage := fmt.Sprintf("failed to build unified config. errors:\n%+v\n", getErrorMessagesFromDiagnostics(diags))
		zap.S().Warnf(errMessage)
		t.addError(ErrBuildTFConfigDir.Error(), t.absRootDir)
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
	root := &ModuleConfig{Config: unified.Root, Name: "root"}
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
				return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "terraform", Directory: t.absRootDir, ErrMessage: "failed to create ResourceConfig"})
			}

			// set module name
			resourceConfig.ModuleName = current.Name

			// resolve references
			resourceConfig.Config = r.ResolveRefs(resourceConfig.Config.(jsonObj))
			var isRemoteModule bool
			// source file path
			resourceConfig.Source, isRemoteModule, err = GetConfigSource(remoteURLMapping, resourceConfig, t.absRootDir)
			if err != nil {
				errMessage := fmt.Sprintf("failed to get resource's filepath: %v", err)
				return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "terraform", Directory: t.absRootDir, ErrMessage: errMessage})
			}

			resourceConfig.TerraformVersion = t.terraformVersion
			resourceConfig.ProviderVersion = providerVersion

			// if root module do not have provider contraints fetch the latest compatible version
			if resourceConfig.ProviderVersion == "" {
				resourceConfig.ProviderVersion = LatestProviderVersion(managedResource.Provider, t.terraformVersion)
			}

			if isRemoteModule {
				resourceConfig.IsRemoteModule = &isRemoteModule
			}

			// add tf plan directory relative path
			resourceConfig.PlanRoot = fmt.Sprintf(".%s", string(os.PathSeparator))

			// append to normalized output
			if _, present := allResourcesConfig[resourceConfig.Type]; !present {
				allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
			} else {
				resources := allResourcesConfig[resourceConfig.Type]
				if !output.IsConfigPresent(resources, resourceConfig) {
					allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
				}
			}
		}

		// add all current's children to the queue
		configsQ = append(configsQ, current.getChildConfigs()...)
	}

	// successful
	return allResourcesConfig, t.errIacLoadDirs
}

// buildUnifiedConfig builds a unified config from *hclConfigs.Module object specified for a dir
func (t TerraformDirectoryLoader) buildUnifiedConfig(rootMod *hclConfigs.Module, dir string) (*hclConfigs.Config, hcl.Diagnostics) {
	// using the BuildConfig and ModuleWalkerFunc to traverse through all
	// descendant modules from the root module and create a unified
	// configuration of type *configs.Config
	versionI := 0
	return hclConfigs.BuildConfig(rootMod, hclConfigs.ModuleWalkerFunc(
		func(req *hclConfigs.ModuleRequest) (*hclConfigs.Module, *version.Version, hcl.Diagnostics) {

			// figure out path sub module directory, if it's remote then download it locally
			var (
				pathToModule   string
				err            error
				moduleDirDiags hcl.Diagnostics
			)
			if downloader.IsLocalSourceAddr(req.SourceAddr) {

				pathToModule = t.processLocalSource(req)
				zap.S().Debugf("processing local module %q", pathToModule)
			} else if t.useTerraformCache {
				// check if module is present in terraform cache
				if _, dest := t.GetRemoteModuleIfPresentInTerraformSrc(req); dest != "" {
					pathToModule = dest
				}
			}
			if pathToModule == "" {
				if downloader.IsRegistrySourceAddr(req.SourceAddr) {
					// temp dir to download the remote repo
					tempDir := utils.GenerateTempDir()

					pathToModule, err = t.processTerraformRegistrySource(req, tempDir)
					if err != nil {
						zap.S().Errorf("failed to download remote module %q. error: '%v'", req.SourceAddr, err)
					}
				} else {
					// temp dir to download the remote repo
					tempDir := utils.GenerateTempDir()

					// Download remote module
					pathToModule, err = t.remoteDownloader.DownloadModule(req.SourceAddr, tempDir)
					if err != nil {
						zap.S().Errorf("failed to download remote module %q. error: '%v'", req.SourceAddr, err)
					}
				}
			}

			// verify whether the module source directory has any .tf config files
			if utils.IsDirExists(pathToModule) && !t.parser.IsConfigDir(pathToModule) {
				moduleDirDiags = append(moduleDirDiags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid module config directory",
					Detail:   fmt.Sprintf("Module directory '%s' has no terraform config files for module %s", pathToModule, req.Name),
				})
				return nil, nil, moduleDirDiags
			}
			// load sub module directory
			subMod, diags := t.parser.LoadConfigDir(pathToModule)
			version, _ := version.NewVersion(fmt.Sprintf("1.0.%d", versionI))
			versionI++
			return subMod, version, diags
		},
	))
}

// addError adds error load dir errors
func (t *TerraformDirectoryLoader) addError(errMessage, directory string) {
	t.errIacLoadDirs = multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "terraform", Directory: directory, ErrMessage: errMessage})
}

func (t TerraformDirectoryLoader) processLocalSource(req *hclConfigs.ModuleRequest) string {
	// determine the absolute path from root module to the sub module

	// while building the unified config, recursive calls are made for all the paths resolved,
	// the source address in a tf file is relative while this func is called, hence, we get the
	// path of caller dir, and join the source address of current module request to get the path to module

	// get the caller dir path
	callDirPath := filepath.Dir(req.CallRange.Filename)

	// join source address to the caller dir
	return filepath.Join(callDirPath, req.SourceAddr)
}

func (t TerraformDirectoryLoader) processTerraformRegistrySource(req *hclConfigs.ModuleRequest, tempDir string) (string, error) {
	// regsrc.ParseModuleSource func returns a terraform registry module source
	// error check is not required as the source address is already validated
	module, _ := regsrc.ParseModuleSource(req.SourceAddr)

	pathToModule, err := t.remoteDownloader.DownloadRemoteModule(req.VersionConstraint, tempDir, module)
	if err != nil {
		return pathToModule, err
	}

	return pathToModule, nil
}

// getChildConfigs will get all child configs in a ModuleConfig
func (m *ModuleConfig) getChildConfigs() []*ModuleConfig {
	allConfigs := make([]*ModuleConfig, 0)
	for childName, childModule := range m.Config.Children {
		childModuleConfig := &ModuleConfig{
			Config:           childModule,
			ParentModuleCall: m.Config.Module.ModuleCalls[childName],
			Name:             childName,
		}
		allConfigs = append(allConfigs, childModuleConfig)
	}
	return allConfigs
}

// GetRemoteLocation checks wether the source belongs to the remote module present in downloader cache.
// cache has key = remoteURL and value = tempDir
func GetRemoteLocation(cache map[string]string, resourcePath string) (remoteURL, tmpDir string) {
	dir := filepath.Dir(resourcePath)
	for k, v := range cache {
		// this condition will never arise added for safe check
		if len(v) > 0 {
			// check dir length is greater than tempDir and file belongs to the same tempDir
			if len(dir) >= len(v) && v == dir[:len(v)] {
				return k, v
			}
		}
	}
	return
}

// GetConfigSource - get the source path for the resource
func GetConfigSource(remoteURLMapping map[string]string, resourceConfig output.ResourceConfig, absRootDir string) (string, bool, error) {
	var (
		source   string
		err      error
		rel      string
		isRemote bool
	)

	// Get source path if remote module used
	remoteURL, tempDir := GetRemoteLocation(remoteURLMapping, resourceConfig.Source)
	if remoteURL != "" {
		rel, err = filepath.Rel(tempDir, resourceConfig.Source)
		if err != nil {
			errMessage := fmt.Sprintf("failed to get remote resource's %s filepath: %v", resourceConfig.Name, err)
			return source, false, errors.New(errMessage)
		}
		isRemote = true

		source = filepath.Join(url.PathEscape(remoteURL), rel)
		source, err = url.PathUnescape(source)
		if err != nil {
			errMessage := fmt.Sprintf("failed to get remote resource's %s filepath: %v", resourceConfig.Name, err)
			return source, false, errors.New(errMessage)
		}
	} else {
		// source file path
		source, err = filepath.Rel(absRootDir, resourceConfig.Source)
		if err != nil {
			return source, false, err
		}
	}
	return source, isRemote, nil
}

// GetRemoteModuleIfPresentInTerraformSrc - Gets the remote module if present in terraform init cache
func (t *TerraformDirectoryLoader) GetRemoteModuleIfPresentInTerraformSrc(req *hclConfigs.ModuleRequest) (src string, destpath string) {
	terraformInitRegs := filepath.Join(t.absRootDir, terraformModuleInstallDir)
	modules := TerraformModuleManifest{}
	var ok bool
	if modules, ok = t.terraformInitModuleCache[terraformInitRegs]; !ok {
		if utils.IsDirExists(terraformInitRegs) {
			_, err := os.Stat(filepath.Join(terraformInitRegs, terraformInstalledModulelMetaFileName))
			if err != nil {
				if os.IsNotExist(err) {
					zap.S().Debug("found no terraform module metadata file in dir %s", terraformInitRegs)
					return
				}
				zap.S().Error("error reading terraform module metadata file", err)
				return
			}
			data, err := os.ReadFile(filepath.Join(terraformInitRegs, terraformInstalledModulelMetaFileName))
			if err == nil {
				err := json.Unmarshal(data, &modules)
				if err != nil {
					zap.S().Error("error unmarshalling terraform module metadata", err)
					return
				}
			}
		}
		// if the module metadata file was read first time add that to cache against the found working directory
		t.terraformInitModuleCache[terraformInitRegs] = modules
	}
	for _, m := range modules.Modules {
		if strings.EqualFold(m.SourceAddr, req.SourceAddr) {
			// if the module source is not registry then version check is not required
			if !downloader.IsRegistrySourceAddr(req.SourceAddr) {
				return req.SourceAddr, filepath.Join(t.absRootDir, m.Dir)
			} else if versionSatisfied(m.VersionStr, req.VersionConstraint) {
				return req.SourceAddr, filepath.Join(t.absRootDir, m.Dir)
			}
		}
	}
	zap.S().Debug("found no version matching for module: %s in terraform module cache %s", req.Name, filepath.Join(terraformInitRegs, "modules.json"))
	return
}

// versionSatisfied - check version in terraform init cache satisfies the required version constraints
func versionSatisfied(foundversion string, requiredVersion hclConfigs.VersionConstraint) bool {
	currentVersion, err := version.NewVersion(foundversion)
	if err != nil {
		return false
	}

	if requiredVersion.Required == nil && foundversion != "" {
		return true
	}

	if requiredVersion.Required.Check(currentVersion) {
		return true
	}

	return false
}
