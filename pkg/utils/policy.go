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

package utils

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetAbsPolicyConfigPaths transforms the provided policy base path and repo path into absolute paths
func GetAbsPolicyConfigPaths(policyBasePath, policyRepoPath string) (string, string, error) {
	absolutePolicyBasePath, err := GetAbsPath(policyBasePath)
	if err != nil {
		return policyBasePath, policyRepoPath, errors.Errorf("invalid policy path `%s`, error : `%v`", policyBasePath, err)
	}

	absolutePolicyRepoPath, err := GetAbsPath(policyRepoPath)
	if err != nil {
		return policyRepoPath, policyBasePath, errors.Errorf("invalid repository path `%s`, error : `%v`", policyRepoPath, err)
	}

	if strings.HasPrefix(absolutePolicyRepoPath, absolutePolicyBasePath) {
		return absolutePolicyBasePath, absolutePolicyRepoPath, nil
	}

	zap.S().Debugf("absolute rego_subdir path, `%s`, does not fall under base repo path's `%s` directory structure", absolutePolicyRepoPath, absolutePolicyBasePath)
	zap.S().Debugf("appending rego_subdir path: `%s` to the policy base path: `%s`. checking ...", policyRepoPath, policyBasePath)

	absolutePolicyRepoPath = filepath.Join(absolutePolicyBasePath, policyRepoPath)
	return absolutePolicyBasePath, absolutePolicyRepoPath, nil
}

// CheckPolicyType checks if supplied policy type matches desired policy types
func CheckPolicyType(rulePolicyType string, desiredPolicyTypes []string) bool {
	normDesiredPolicyTypes := make(map[string]bool, len(desiredPolicyTypes))
	normRulePolicyType := EnsureUpperCaseTrimmed(rulePolicyType)

	for _, desiredPolicyType := range desiredPolicyTypes {
		desiredPolicyType = EnsureUpperCaseTrimmed(desiredPolicyType)
		normDesiredPolicyTypes[desiredPolicyType] = true
	}

	if _, ok := normDesiredPolicyTypes["ALL"]; ok {
		return true
	}

	_, ok := normDesiredPolicyTypes[normRulePolicyType]
	return ok
}
