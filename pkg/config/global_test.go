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

package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/utils"
)

func TestLoadGlobalConfig(t *testing.T) {
	testConfigFile := filepath.Join(testDataDir, "terrascan-config-all-fields.toml")
	absDefaultBasePolicyPath, absDefaultPolicyRepoPath, _ := utils.GetAbsPolicyConfigPaths(defaultBasePolicyPath, defaultPolicyRepoPath)
	absCustomPath, absRegoSubdirPath, _ := utils.GetAbsPolicyConfigPaths("custom-path", "rego-subdir")

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
		scanRules      []string
		skipRules      []string
		severity       string
		categories     []string
		notifications  map[string]Notifier
		k8sAdmControl  K8sAdmissionControl
	}{
		{
			// no error expected
			name: "global config file not specified",
			args: args{
				configFile: "",
			},
			policyBasePath: absDefaultBasePolicyPath,
			policyRepoPath: absDefaultPolicyRepoPath,
			repoURL:        defaultPolicyRepoURL,
			branchName:     defaultPolicyBranch,
		},
		{
			name: "global config file specified but doesn't exist",
			args: args{
				configFile: "test.toml",
			},
			wantErr:        true,
			policyBasePath: defaultBasePolicyPath,
			policyRepoPath: defaultPolicyRepoPath,
			repoURL:        defaultPolicyRepoURL,
			branchName:     defaultPolicyBranch,
		},
		{
			name: "valid global config file specified",
			args: args{
				configFile: testConfigFile,
			},
			policyBasePath: absCustomPath,
			policyRepoPath: absRegoSubdirPath,
			repoURL:        "https://repository/url",
			branchName:     "branch-name",
			scanRules:      testRules.ScanRules,
			skipRules:      testRules.SkipRules,
			severity:       highSeverity.Level,
			categories:     testCategoryList.List,
			notifications: map[string]Notifier{
				"webhook1": testNotifier,
			},
			k8sAdmControl: testK8sAdmControl,
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

			if !utils.IsSliceEqual(GetScanRules(), tt.scanRules) || !utils.IsSliceEqual(GetSkipRules(), tt.skipRules) || !utils.IsSliceEqual(GetCategoryList(), tt.categories) || GetSeverityLevel() != tt.severity {
				t.Errorf("LoadGlobalConfig() error = got scan rules: %v, skip rules: %v, categories: %v, severity: %v, want scan rules: %v, skip rules: %v, categories: %v, severity: %v", GetScanRules(), GetSkipRules(), GetCategoryList(), GetSeverityLevel(), tt.scanRules, tt.skipRules, tt.categories, tt.severity)
			}

			if !reflect.DeepEqual(GetNotifications(), tt.notifications) || !reflect.DeepEqual(GetK8sAdmissionControl(), tt.k8sAdmControl) {
				t.Errorf("LoadGlobalConfig() error = got notifications: %v, k8s admission control: %v, want notifications: %v, k8s admission control: %v", GetNotifications(), GetK8sAdmissionControl(), tt.notifications, tt.k8sAdmControl)
			}
		})
	}
}
