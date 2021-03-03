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
	"io/ioutil"
	"os"

	"github.com/accurics/terrascan/pkg/config"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	basePath       = config.GetPolicyRepoPath()
	basePolicyPath = config.GetPolicyBasePath()
	repoURL        = config.GetPolicyRepoURL()
	branch         = config.GetPolicyBranch()
)

// Run initializes terrascan if not done already
func Run(isScanCmd bool) error {
	zap.S().Debug("initializing terrascan")

	// check if policy paths exist
	if path, err := os.Stat(basePolicyPath); err == nil && path.IsDir() {
		if isScanCmd {
			return nil
		}
	}

	// download policies
	if err := DownloadPolicies(); err != nil {
		return err
	}

	zap.S().Debug("intialized successfully")
	return nil
}

// DownloadPolicies clones the policies to a local folder
func DownloadPolicies() error {
	zap.S().Debug("downloading policies")

	tempPath, err := ioutil.TempDir("", "terrascan-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory. error: '%v'", err)
	}

	defer os.RemoveAll(tempPath)

	zap.S().Debugf("cloning terrascan repo at %s", tempPath)

	// clone the repo
	r, err := git.PlainClone(tempPath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return fmt.Errorf("failed to download policies. error: '%v'", err)
	}

	// create working tree
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to create working tree. error: '%v'", err)
	}

	// fetch references
	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch references from git repo. error: '%v'", err)
	}

	// checkout policies branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout git branch '%v'. error: '%v'", branch, err)
	}

	// cleaning the existing cached policies at basePath
	if err = os.RemoveAll(basePath); err != nil {
		return fmt.Errorf("failed to clean up the directory '%s'. error: '%v'", basePath, err)
	}

	// move the freshly cloned repo from tempPath to basePath
	if err = os.Rename(tempPath, basePath); err != nil {
		return fmt.Errorf("failed to install policies to '%s'. error: '%v'", basePath, err)
	}

	return nil
}
