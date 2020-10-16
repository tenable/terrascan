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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log" // we log from init(), so can't rely on zap to be available
	"os"
)

const (
	policyRepoURL    = "https://github.com/accurics/terrascan.git"
	policyBranch     = "master"
	configEnvvarName = "TERRASCAN_CONFIG"
)

var (
	policyRepoPath = os.Getenv("HOME") + "/.terrascan"
	policyBasePath = policyRepoPath + "/pkg/policies/opa/rego"
)

func init() {
	loadGlobalConfig()
}

func loadGlobalConfig() {
	// Start with the defaults
	Global.Policy = PolicyConfig{
		BasePath: policyBasePath,
		RepoPath: policyRepoPath,
		RepoURL:  policyRepoURL,
		Branch:   policyBranch,
	}

	// If the user specifies a config file in TERRASCAN_CONFIG,
	// overwrite the defaults with the values from that file.
	// Retain the defaults for members not specified in the file.
	configFile := os.Getenv(configEnvvarName)

	if len(configFile) > 0 {
		p, err := loadConfigFile(configFile)
		if err != nil {
			log.Println(err)
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
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return p, fmt.Errorf("unable to read config file: %v", err)
	}

	if err = json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("unable to unmarshal config file %s: %v", configFile, err)
	}
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
