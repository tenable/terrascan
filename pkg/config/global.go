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

package config

import (
	"os"

	"go.uber.org/zap"
)

const (
	policyRepoURL    = "https://github.com/accurics/terrascan.git"
	policyBranch     = "master"
	configEnvvarName = "TERRASCAN_CONFIG"
	policyConfigKey  = "policy"
)

var (
	policyRepoPath = os.Getenv("HOME") + "/.terrascan"
	policyBasePath = policyRepoPath + "/pkg/policies/opa/rego"
)

func init() {
	// If the user specifies a config file in TERRASCAN_CONFIG,
	// overwrite the defaults with the values from that file.
	// Retain the defaults for members not specified in the file.
	if err := LoadGlobalConfig(os.Getenv(configEnvvarName)); err != nil {
		zap.S().Error("error while loading global config", zap.Error(err))
	}
}

// LoadGlobalConfig loads policy configuration from specified configFile
// into var Global.Policy.  Members of Global.Policy that are not specified
// in configFile will get default values
func LoadGlobalConfig(configFile string) error {
	// Start with the defaults
	Global.Policy = Policy{
		BasePath: policyBasePath,
		RepoPath: policyRepoPath,
		RepoURL:  policyRepoURL,
		Branch:   policyBranch,
	}

	if configFile == "" {
		zap.S().Debug("global config env variable is not specified")
		return nil
	}

	configReader, err := NewTerrascanConfigReader(configFile)
	if err != nil {
		return err
	}

	if len(configReader.GetPolicyConfig().BasePath) > 0 {
		Global.Policy.BasePath = configReader.GetPolicyConfig().BasePath
	}
	if len(configReader.GetPolicyConfig().RepoPath) > 0 {
		Global.Policy.RepoPath = configReader.GetPolicyConfig().RepoPath
	}
	if len(configReader.GetPolicyConfig().RepoURL) > 0 {
		Global.Policy.RepoURL = configReader.GetPolicyConfig().RepoURL
	}
	if len(configReader.GetPolicyConfig().Branch) > 0 {
		Global.Policy.Branch = configReader.GetPolicyConfig().Branch
	}
	return nil
}

// GetPolicyBasePath returns policy base path as set in global config
func GetPolicyBasePath() string {
	return Global.Policy.BasePath
}

// GetPolicyRepoPath return path to the policies repo locally downloaded
func GetPolicyRepoPath() string {
	return Global.Policy.RepoPath
}

// GetPolicyRepoURL returns policy repo url
func GetPolicyRepoURL() string {
	return Global.Policy.RepoURL
}

// GetPolicyBranch returns policy repo url
func GetPolicyBranch() string {
	return Global.Policy.Branch
}
