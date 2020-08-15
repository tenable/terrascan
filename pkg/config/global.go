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
)

const (
	policyRepoURL = "https://github.com/accurics/terrascan.git"
	policyBranch  = "master"
)

var (
	policyRepoPath = os.Getenv("HOME") + "/.terrascan"
	policyBasePath = policyRepoPath + "/pkg/policies/opa/rego"
)

func init() {
	Global.Policy = PolicyConfig{
		BasePath: policyBasePath,
		RepoPath: policyRepoPath,
		RepoURL:  policyRepoURL,
		Branch:   policyBranch,
	}
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
