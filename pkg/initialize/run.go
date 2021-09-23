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
	"net/http"
	"os"

	"github.com/accurics/terrascan/pkg/config"
	"go.uber.org/zap"
	git "gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	errNoConnection = fmt.Errorf("could not connect to github.com")
)

const terrascanReadmeURL string = "https://raw.githubusercontent.com/accurics/terrascan/master/README.md"

// Run initializes terrascan if not done already
func Run(isNonInitCmd bool) error {
	// check if policy paths exist
	if path, err := os.Stat(config.GetPolicyRepoPath()); err == nil && path.IsDir() {
		if isNonInitCmd {
			return nil
		}
	}

	zap.S().Debug("initializing terrascan")

	if !connected(terrascanReadmeURL) {
		return errNoConnection
	}

	// download policies
	if err := DownloadPolicies(); err != nil {
		return err
	}

	zap.S().Debug("initialized successfully")
	return nil
}

// DownloadPolicies clones the policies to a local folder
func DownloadPolicies() error {

	policyBasePath := config.GetPolicyBasePath()
	repoURL := config.GetPolicyRepoURL()
	branch := config.GetPolicyBranch()

	zap.S().Debug("downloading policies")

	zap.S().Debugf("base directory path : %s", policyBasePath)
	zap.S().Debugf("policy directory path : %s", config.GetPolicyRepoPath())
	zap.S().Debugf("policy repo url : %s", repoURL)
	zap.S().Debugf("policy repo git branch : %s", branch)

	os.RemoveAll(policyBasePath)

	zap.S().Debugf("cloning terrascan repo at %s", policyBasePath)

	// clone the repo
	r, err := git.PlainClone(policyBasePath, false, &git.CloneOptions{
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

	return nil
}

func connected(url string) bool {
	_, err := http.Get(url)
	return err == nil
}
