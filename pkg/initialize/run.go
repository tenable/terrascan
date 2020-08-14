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

package initialize

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	repoURL = "https://github.com/accurics/terrascan.git"
	branch  = "terrascan-v1.0"
)

var (
	basePath       = os.Getenv("HOME") + "/.terrascan"
	basePolicyPath = basePath + "/pkg/policies/opa/rego"

	// AWSDefaultPolicyPath default policy path for aws
	AWSDefaultPolicyPath = basePolicyPath + "/aws"

	// AzureDefaultPolicyPath  default policy path for azure
	AzureDefaultPolicyPath = basePolicyPath + "/azure"
)

var policyPathMap = make(map[string]string)

func init() {
	policyPathMap["aws"] = AWSDefaultPolicyPath
	policyPathMap["azure"] = AzureDefaultPolicyPath
}

// IsCloudSupported checks if cloud provider is present in policyPathMap
func IsCloudSupported(cloud string) bool {
	_, supported := policyPathMap[cloud]
	return supported
}

// GetPolicyPath returns the policy path for a cloud provider
func GetPolicyPath(cloud string) string {
	return policyPathMap[cloud]
}

// Run initializes terrascan if not done already
func Run() error {

	// check if policy paths exist
	if path, err := os.Stat(basePolicyPath); err == nil && path.IsDir() {
		return nil
	}

	// download policies
	os.RemoveAll(basePath)
	if err := DownloadPolicies(); err != nil {
		return err
	}

	return nil
}

// DownloadPolicies clones the policies to a local folder
func DownloadPolicies() error {

	// clone the repo
	r, err := git.PlainClone(basePath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		zap.S().Errorf("failed to download policies. error: '%v'", err)
		return err
	}

	// create working tree
	w, err := r.Worktree()
	if err != nil {
		zap.S().Errorf("failed to create working tree. error: '%v'", err)
		return err
	}

	// fetch references
	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		zap.S().Errorf("failed to fetch references from repo. error: '%v'", err)
		return err
	}

	// checkout policies branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	if err != nil {
		zap.S().Errorf("failed to checkout branch '%v'. error: '%v'", branch, err)
		return err
	}

	return nil
}
