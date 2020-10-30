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
	"fmt"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
	"os"
)

const (
	policyRepoURL    = "https://github.com/accurics/terrascan.git"
	policyBranch     = "master"
	configEnvvarName = "TERRASCAN_CONFIG"
	policyConfigKey  = "policy"
)

var (
	policyRepoPath       = os.Getenv("HOME") + "/.terrascan"
	policyBasePath       = policyRepoPath + "/pkg/policies/opa/rego"
	errTomlKeyNotPresent = fmt.Errorf("%s key not present in toml config", policyConfigKey)
)

func init() {
	// If the user specifies a config file in TERRASCAN_CONFIG,
	// overwrite the defaults with the values from that file.
	// Retain the defaults for members not specified in the file.
	LoadGlobalConfig(os.Getenv(configEnvvarName))
}

// LoadGlobalConfig loads policy configuration from specified configFile
// into var Global.Policy.  Members of Global.Policy that are not specified
// in configFile will get default values
func LoadGlobalConfig(configFile string) {
	// Start with the defaults
	Global.Policy = PolicyConfig{
		BasePath: policyBasePath,
		RepoPath: policyRepoPath,
		RepoURL:  policyRepoURL,
		Branch:   policyBranch,
	}

	if len(configFile) > 0 {
		p, err := loadConfigFile(configFile)
		if err != nil {
			zap.S().Error(err)
			return
		}
		if len(p.Policy.BasePath) > 0 {
			Global.Policy.BasePath = p.Policy.BasePath
		}
		if len(p.Policy.RepoPath) > 0 {
			Global.Policy.RepoPath = p.Policy.RepoPath
		}
		if len(p.Policy.RepoURL) > 0 {
			Global.Policy.RepoURL = p.Policy.RepoURL
		}
		if len(p.Policy.Branch) > 0 {
			Global.Policy.Branch = p.Policy.Branch
		}
	}
}

func loadConfigFile(configFile string) (GlobalConfig, error) {
	p := GlobalConfig{}

	config, err := LoadConfig(configFile)
	if err != nil {
		return p, ErrNotPresent
	}

	keyConfig := config.Get(policyConfigKey)
	if keyConfig == nil {
		return p, errTomlKeyNotPresent
	}

	keyTomlConfig := keyConfig.(*toml.Tree)

	// We want to treat missing keys as empty strings
	str := func(x interface{}) string {
		if x == nil {
			return ""
		}
		return x.(string)
	}

	// path = path where repo will be checked out
	p.Policy.BasePath = str(keyTomlConfig.Get("path"))

	// repo_url = git url to policy repository
	p.Policy.RepoURL = str(keyTomlConfig.Get("repo_url"))

	// rego_subdir = subdir of <path> where rego files are located
	p.Policy.RepoPath = str(keyTomlConfig.Get("rego_subdir"))

	// branch = git branch where policies are stored
	p.Policy.Branch = str(keyTomlConfig.Get("branch"))

	return p, nil
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
