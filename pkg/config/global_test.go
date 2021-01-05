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
	"testing"
)

func TestLoadGlobalConfig(t *testing.T) {
	testConfigFile := "./testdata/terrascan-config.toml"

	type args struct {
		configFile string
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		policyBasePath string
		policyRepoPath string
		repoURL        string
		branchName     string
	}{
		{
			// no error expected
			name: "global config file not specified",
			args: args{
				configFile: "",
			},
			policyBasePath: policyBasePath,
			policyRepoPath: policyRepoPath,
			repoURL:        policyRepoURL,
			branchName:     policyBranch,
		},
		{
			name: "global config file specified but doesn't exist",
			args: args{
				configFile: "test.toml",
			},
			wantErr:        true,
			policyBasePath: policyBasePath,
			policyRepoPath: policyRepoPath,
			repoURL:        policyRepoURL,
			branchName:     policyBranch,
		},
		{
			name: "valid global config file specified",
			args: args{
				configFile: testConfigFile,
			},
			policyBasePath: "custom-path",
			policyRepoPath: "rego-subdir",
			repoURL:        "https://repository/url",
			branchName:     "branch-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadGlobalConfig(tt.args.configFile); (err != nil) != tt.wantErr {
				t.Errorf("LoadGlobalConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if GetPolicyBasePath() != tt.policyBasePath || GetPolicyRepoPath() != tt.policyRepoPath || GetPolicyRepoURL() != tt.repoURL || GetPolicyBranch() != tt.branchName {
				t.Errorf("LoadGlobalConfig() error = got BasePath: %v, RepoPath: %v, RepoURL: %v, BranchName: %v, want BasePath: %v, RepoPath: %v, RepoURL: %v, BranchName: %v", GetPolicyBasePath(), GetPolicyRepoPath(), GetPolicyRepoURL(), GetPolicyBranch(), tt.policyBasePath, tt.policyRepoPath, tt.repoURL, tt.branchName)
			}
		})
	}
}
