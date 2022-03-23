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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/config"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
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

	// download policies
	if err := DownloadPolicies(); err != nil {
		return err
	}

	zap.S().Debug("initialized successfully")
	return nil
}

// DownloadPolicies clones the policies to a local folder
func DownloadPolicies() error {
	accessToken := config.GetPolicyAccessToken()
	policyBasePath := config.GetPolicyBasePath()

	zap.S().Debug("downloading policies")
	zap.S().Debugf("base directory path : %s", policyBasePath)

	err := os.RemoveAll(policyBasePath)
	if err != nil {
		return fmt.Errorf("unable to delete base folder. error: '%v'", err)
	}

	if accessToken == "" {
		return downloadTerrascanPolicies(policyBasePath)
	}

	return downloadCommercialPolicies(policyBasePath, accessToken)
}

func downloadCommercialPolicies(policyBasePath, accessToken string) error {
	err := ensureDir(policyBasePath)
	if err != nil {
		return err
	}

	const apiPath = "/v1/api/rule?default=true"
	environment := config.GetPolicyEnvironment()

	zap.S().Debugf("policy environment : %s", environment)
	zap.S().Debugf("downloading commercial policies in %s", policyBasePath)

	var client http.Client

	req, err := http.NewRequest(http.MethodGet, environment+apiPath, nil)
	if err != nil {
		return fmt.Errorf("error constructing request object. error: '%v'", err)
	}

	var cookie http.Cookie
	cookie.Name = "x-siac-session"
	cookie.Value = accessToken
	req.AddCookie(&cookie)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading commercial policies. error: '%v', response code: '%d'", err, res.StatusCode)
	}

	policies, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading api call response for commercial policies. error: '%v'", err)
	}

	return convertCommercialPolicies(policies, policyBasePath)
}

func convertCommercialPolicies(policies []byte, policyBasePath string) error {
	var rules []commercialPolicyMetadata

	err := json.Unmarshal(policies, &rules)
	if err != nil {
		return fmt.Errorf("failed to unmarshal policies into structure. error: '%v'", err)
	}

	var policy commercialPolicy
	for _, rule := range rules {
		policy.regoTemplate = rule.RuleTemplate
		policy.metadataFileName = rule.RuleReferenceId + ".json"
		policy.resourceType = rule.ResourceType

		policy.policyMetadata.Name = rule.RuleName
		policy.policyMetadata.File = rule.RuleTemplateID + ".rego"
		policy.policyMetadata.ResourceType = rule.ResourceType
		policy.policyMetadata.Severity = rule.Severity
		policy.policyMetadata.Description = rule.RuleDisplayName
		policy.policyMetadata.ReferenceID = rule.RuleReferenceId
		policy.policyMetadata.ID = rule.RuleReferenceId
		policy.policyMetadata.Category = rule.Category
		policy.policyMetadata.Version = rule.Version

		templateArgs, ok := rule.RuleArgument.(map[string]interface{})
		if !ok && templateArgs != nil {
			return fmt.Errorf("incorrect rule argument type, must be a map[string]interface{}")
		}
		policy.policyMetadata.TemplateArgs = templateArgs

		err = saveCommercialPolicies(policy, policyBasePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveCommercialPolicies(policy commercialPolicy, policyBasePath string) error {
	const tabSpace = "    "
	csp := strings.ToLower(strings.Split(policy.policyMetadata.ReferenceID, "_")[1])

	cspDir := filepath.Join(policyBasePath, csp)
	err := ensureDir(cspDir)
	if err != nil {
		return err
	}

	resourceDir := filepath.Join(cspDir, policy.resourceType)
	err = ensureDir(resourceDir)
	if err != nil {
		return err
	}

	metadata, err := json.MarshalIndent(policy.policyMetadata, "", tabSpace)
	if err != nil {
		return fmt.Errorf("error could not marshal json object into byte array: '%v'", err)
	}
	metaDataPath := filepath.Join(resourceDir, policy.metadataFileName)
	err = ioutil.WriteFile(metaDataPath, metadata, FilePermissionBits)
	if err != nil {
		return fmt.Errorf("error could not write rule metadata file on disk: '%v'", err)
	}

	regoPath := filepath.Join(resourceDir, policy.policyMetadata.File)
	err = ioutil.WriteFile(regoPath, []byte(policy.regoTemplate), FilePermissionBits)
	if err != nil {
		return fmt.Errorf("error could not write rego code file on disk: '%v'", err)
	}

	return nil
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, FilePermissionBits)
		if err != nil {
			return fmt.Errorf("error unable to create requested directory: '%v'", err)
		}
	}
	return nil
}

func downloadTerrascanPolicies(policyBasePath string) error {
	if !connected(terrascanReadmeURL) {
		return errNoConnection
	}

	repoURL := config.GetPolicyRepoURL()
	branch := config.GetPolicyBranch()

	zap.S().Debugf("policy directory path : %s", repoURL)
	zap.S().Debugf("policy repo url : %s", repoURL)
	zap.S().Debugf("policy repo git branch : %s", branch)
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
