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

	version "github.com/hashicorp/go-version"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/regsrc"
	"go.uber.org/zap"
)

// DownloadRemoteModule will download remote modules from public and private terraform registries
// this function takes similar approach taken by terraform init for downloading terraform registry modules
func (r *RemoteModuleInstaller) DownloadRemoteModule(requiredVersion hclConfigs.VersionConstraint, destPath string, module *regsrc.Module) (string, error) {
	// Terraform doesn't allow the hostname to contain Punycode
	// module.SvcHost returns an error for such case
	_, err := module.SvcHost()
	if err != nil {
		zap.S().Errorf("hostname for the module %s is invalid", module.String())
		return "", err
	}

	// get terraform registry client.
	// terraform registry client provides methods for querying the terraform module registry
	regClient := registry.NewClient(nil, nil)

	// get all the available module versions from the terraform registry
	moduleVersions, err := regClient.ModuleVersions(module)
	if err != nil {
		if registry.IsModuleNotFound(err) {
			zap.S().Errorf("module: %s, not be found at registry: %s", module.String(), module.Host().Display())
		} else {
			zap.S().Errorf("error while fetching available modules for module: %s, at registry: %s", module.String(), module.Host().Display())
		}
		return "", err
	}

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
		return "", fmt.Errorf("no versions for module: %s, found at registry: %s", module.String(), module.Host().Display())
	}

	if requiredVersion.Required != nil && latestMatch == nil {
		return "", fmt.Errorf("no versions matching: %s, for module: %s, found at registry: %s, latest version found: %s", requiredVersion.Required.String(), module.String(), module.Host().Display(), latestVersion.String())
	}

	versionToDownload = latestVersion
	if latestMatch != nil {
		versionToDownload = latestMatch
	}

	// get the source location for the matched version
	sourceLocation, err := regClient.ModuleLocation(module, versionToDownload.String())
	if err != nil {
		zap.S().Errorf("error while getting the source location for module: %s, at registry: %s", module.String(), module.Host().Display())
		return "", err
	}

	downloadLocation, err := r.DownloadModule(sourceLocation, destPath)
	if err != nil {
		zap.S().Errorf("error while downloading module: %s, with source location: %s", module.String(), sourceLocation)
		return "", nil
	}

	if module.RawSubmodule != "" {
		// Append the user's requested subdirectory to any subdirectory that
		// was implied by any of the nested layers we expanded within go-getter.
		downloadLocation = filepath.Join(downloadLocation, module.RawSubmodule)
	}

	return downloadLocation, nil
}

func getGreaterVersion(latestVersion *version.Version, currentVersion *version.Version) *version.Version {
	if latestVersion == nil || currentVersion.GreaterThan(latestVersion) {
		latestVersion = currentVersion
	}
	return latestVersion
}
