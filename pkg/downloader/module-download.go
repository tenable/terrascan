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

package downloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/hashicorp/terraform-svchost/auth"
	"github.com/hashicorp/terraform-svchost/disco"
	"github.com/hashicorp/terraform/command/cliconfig"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/hashicorp/terraform/registry/response"
	"go.uber.org/zap"
)

// newRemoteModuleInstaller returns a RemoteModuleInstaller initialized with a
// new cache, downloader and terraform registry client
func newRemoteModuleInstaller() *remoteModuleInstaller {
	return &remoteModuleInstaller{
		cache:                   make(map[string]string),
		downloader:              NewDownloader(),
		terraformRegistryClient: newClientRegistry(),
	}
}

// newTerraformRegistryClient returns a client to query terraform registries
func newTerraformRegistryClient() terraformRegistryClient {
	// get terraform registry client.
	// terraform registry client provides methods for querying the terraform module registry
	// if terraformrc file is found, attempt to load credentials from it
	homedir := utils.GetHomeDir()
	if len(homedir) == 0 {
		return registry.NewClient(nil, nil)
	}
	// TODO: Need to add support for Windows %APPDIR%/terraform.rc
	terraformrc := filepath.Join(homedir, ".terraformrc")

	if _, err := os.Stat(terraformrc); os.IsNotExist(err) {
		return registry.NewClient(nil, nil)
	}
	zap.S().Debugf("Found terraform rc file at %s, attempting to parse", terraformrc)

	return NewAuthenticatedRegistryClient(terraformrc)
}

// NewAuthenticatedRegistryClient parses the contents of a terraformrc file and builds an authenticated
// registry client using the credentials found in the rcfile
func NewAuthenticatedRegistryClient(rcFile string) terraformRegistryClient {
	b, err := ioutil.ReadFile(rcFile)
	if err != nil {
		zap.S().Infof("Error reading %s: %s", rcFile, err)
		return registry.NewClient(nil, nil)
	}

	services, err := buildDiscoServices(b)
	if err != nil {
		zap.S().Infof("Error building terraform credentials %s: %s", rcFile, err)
		return registry.NewClient(nil, nil)
	}
	return registry.NewClient(services, nil)
}

// buildDiscoServices builds terraform services struct for later authentication to a terraform registry
func buildDiscoServices(b []byte) (*disco.Disco, error) {
	obj, err := hcl.Parse(string(b))
	if err != nil {
		return nil, err
	}

	result := &cliconfig.Config{}
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	hostCredentialMap := convertCredentialMapToHostMap(result.Credentials)
	if hostCredentialMap == nil {
		return nil, fmt.Errorf("error converting credential map to host map")
	}
	credSource := auth.StaticCredentialsSource(hostCredentialMap)
	services := disco.NewWithCredentialsSource(credSource)
	return services, nil
}

// convertCredentialMapToHostMap converts map key to pass to terraform auth builder
func convertCredentialMapToHostMap(credentialMap map[string]map[string]interface{}) map[svchost.Hostname]map[string]interface{} {
	hostMap := make(map[svchost.Hostname]map[string]interface{})
	if credentialMap == nil {
		return nil
	}

	for k := range credentialMap {
		hostMap[svchost.Hostname(k)] = credentialMap[k]
	}
	return hostMap
}

// DownloadModule retrieves the package referenced in the given address
// into the installation path and then returns the full path to any subdir
// indicated in the address.
func (r *remoteModuleInstaller) DownloadModule(addr, destPath string) (string, error) {

	// split url and subdir
	URLWithType, subDir, err := r.downloader.GetURLSubDir(addr, destPath)
	if err != nil {
		return "", err
	}

	// check if the module has already been downloaded
	if prevDir, exists := r.cache[URLWithType]; exists {
		zap.S().Debugf("module %q already installed at %q", URLWithType, prevDir)
		destPath = prevDir
	} else {
		destPath, err := r.downloader.Download(URLWithType, destPath)
		if err != nil {
			zap.S().Debugf("failed to download remote module. error: '%v'", err)
			return "", err
		}
		// Remember where we installed this so we might reuse this directory
		// on subsequent calls to avoid re-downloading.
		r.cache[URLWithType] = destPath
	}

	// Our subDir string can contain wildcards until this point, so that
	// e.g. a subDir of * can expand to one top-level directory in a .tar.gz
	// archive. Now that we've expanded the archive successfully we must
	// resolve that into a concrete path.
	var finalDir string
	if subDir != "" {
		finalDir, err = r.downloader.SubDirGlob(destPath, subDir)
		if err != nil {
			return "", err
		}
		zap.S().Debugf("expanded %q to %q", subDir, finalDir)
	} else {
		finalDir = destPath
	}

	// If we got this far then we have apparently succeeded in downloading
	// the requested object!
	return filepath.Clean(finalDir), nil
}

// DownloadRemoteModule will download remote modules from public and private terraform registries
// this function takes similar approach taken by terraform init for downloading terraform registry modules
func (r *remoteModuleInstaller) DownloadRemoteModule(requiredVersion hclConfigs.VersionConstraint, destPath string, module *regsrc.Module) (string, error) {
	// Terraform doesn't allow the hostname to contain Punycode
	// for more information on Punycode refer https://en.wikipedia.org/wiki/Punycode
	// module.SvcHost returns an error for such case
	_, err := module.SvcHost()
	if err != nil {
		zap.S().Errorf("hostname for the module %s is invalid", module.String())
		return "", err
	}

	// get all the available module versions from the terraform registry
	moduleVersions, err := r.terraformRegistryClient.ModuleVersions(module)
	if err != nil {
		if registry.IsModuleNotFound(err) {
			zap.S().Errorf("module: %s, not be found at registry: %s", module.String(), module.Host().Display())
		} else {
			zap.S().Errorf("error while fetching available modules for module: %s, at registry: %s. Error: %s", module.String(), module.Host().Display(), err.Error())
		}
		return "", err
	}

	// get the version to download
	versionToDownload, err := getVersionToDownload(moduleVersions, requiredVersion, module)
	if err != nil {
		zap.S().Error("error while fetching the version to download,", zap.Error(err))
		return "", err
	}

	// get the source location for the matched version
	sourceLocation, err := r.terraformRegistryClient.ModuleLocation(module, versionToDownload.String())
	if err != nil {
		zap.S().Errorf("error while getting the source location for module: %s, at registry: %s", module.String(), module.Host().Display())
		return "", err
	}

	downloadLocation, err := r.DownloadModule(sourceLocation, destPath)
	if err != nil {
		zap.S().Errorf("error while downloading module: %s, with source location: %s", module.String(), sourceLocation)
		return "", err
	}

	if module.RawSubmodule != "" {
		// Append the user's requested subdirectory
		downloadLocation = filepath.Join(downloadLocation, module.RawSubmodule)
	}

	return downloadLocation, nil
}

// CleanUp cleans up all the locally downloaded modules
func (r *remoteModuleInstaller) CleanUp() {
	for url, path := range r.cache {
		zap.S().Debugf("deleting %q installed at %q", url, path)
		os.RemoveAll(path)
	}
}

// helper func to compare and update the version
func getGreaterVersion(latestVersion *version.Version, currentVersion *version.Version) *version.Version {
	if latestVersion == nil || currentVersion.GreaterThan(latestVersion) {
		latestVersion = currentVersion
	}
	return latestVersion
}

// helper func to get the module version to download
func getVersionToDownload(moduleVersions *response.ModuleVersions, requiredVersion hclConfigs.VersionConstraint, module *regsrc.Module) (*version.Version, error) {
	// terraform init command pulls all the available versions of a module,
	// and downloads the latest non pre-release (unless a pre-release version was
	// specified in tf file) version, if a version constraint is not provided in the tf file.
	// we are following what terraform does
	allModules := moduleVersions.Modules[0]

	var latestMatch *version.Version
	var latestVersion *version.Version
	var versionToDownload *version.Version
	for _, moduleVersion := range allModules.Versions {
		currentVersion, err := version.NewVersion(moduleVersion.Version)
		if err != nil {
			// if error is received for a version, then skip the current version
			zap.S().Errorf("invalid version: %s, for module: %s, at registry: %s", moduleVersion.Version, module.String(), module.Host().Display())
			continue
		}

		if requiredVersion.Required == nil {
			// skip the pre release version
			if currentVersion.Prerelease() != "" {
				continue
			}

			// update the latest version
			latestVersion = getGreaterVersion(latestVersion, currentVersion)
		} else {
			// skip the pre-release version, unless specified in the tf file
			if currentVersion.Prerelease() != "" && requiredVersion.Required.String() != currentVersion.String() {
				continue
			}

			// update the latest version
			latestVersion = getGreaterVersion(latestVersion, currentVersion)

			// update latest match
			if requiredVersion.Required.Check(currentVersion) {
				latestMatch = getGreaterVersion(latestMatch, currentVersion)
			}
		}
	}

	if latestVersion == nil {
		return nil, fmt.Errorf("no versions for module: %s, found at registry: %s", module.String(), module.Host().Display())
	}

	if requiredVersion.Required != nil && latestMatch == nil {
		return nil, fmt.Errorf("no versions matching: %s, for module: %s, found at registry: %s, latest version found: %s", requiredVersion.Required.String(), module.String(), module.Host().Display(), latestVersion.String())
	}

	versionToDownload = latestVersion
	if latestMatch != nil {
		versionToDownload = latestMatch
	}

	return versionToDownload, nil
}

// GetModuleInstalledCache - returns the cache of dowloaded modules mapping
func (r *remoteModuleInstaller) GetModuleInstalledCache() map[string]string {
	return r.cache
}
