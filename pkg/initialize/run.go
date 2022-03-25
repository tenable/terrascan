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
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
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
const filePermissionBits fs.FileMode = 0755

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

	policyRepoPath := config.GetPolicyRepoPath()
	err = os.MkdirAll(policyRepoPath, filePermissionBits)
	if err != nil {
		return err
	}

	const apiPath = "/v1/api/app/rules?default=true"
	environment := config.GetPolicyEnvironment()

	zap.S().Debugf("policy environment : %s", environment)
	zap.S().Debugf("downloading commercial policies in %s", policyBasePath)

	var client http.Client

	req, err := http.NewRequest(http.MethodGet, environment+apiPath, nil)
	if err != nil {
		return fmt.Errorf("error constructing request object. error: '%v'", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading commercial policies. error: '%v', response code: '%d'", err, res.StatusCode)
	}

	if res == nil {
		return fmt.Errorf("error could not download policies, please check your network connection")
	}

	policies, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading api call response for commercial policies. error: '%v'", err)
	}

	return convertCommercialPolicies(policies, policyRepoPath)
}

func convertCommercialPolicies(policies []byte, policyRepoPath string) error {
	var ruleMetadataList []commercialPolicyMetadata

	err := json.Unmarshal(policies, &ruleMetadataList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal policies into structure. error: '%v'", err)
	}

	for _, ruleMetadata := range ruleMetadataList {
		policy, err := getCommercialPolicy(ruleMetadata)
		if err != nil {
			return err
		}

		err = saveCommercialPolicies(policy, policyRepoPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCommercialPolicy(ruleMetadata commercialPolicyMetadata) (commercialPolicy, error) {
	var policy commercialPolicy
	var templateArgs map[string]interface{}

	policy.regoTemplate = "package accurics\n\n" + ruleMetadata.RuleTemplate
	policy.metadataFileName = ruleMetadata.RuleReferenceID + ".json"
	policy.resourceType = ruleMetadata.ResourceType

	policy.policyMetadata.Name = ruleMetadata.RuleName
	policy.policyMetadata.File = ruleMetadata.RegoName + ".rego"
	policy.policyMetadata.ResourceType = ruleMetadata.ResourceType
	policy.policyMetadata.Severity = ruleMetadata.Severity
	policy.policyMetadata.Description = ruleMetadata.RuleDisplayName
	policy.policyMetadata.ReferenceID = ruleMetadata.RuleReferenceID
	policy.policyMetadata.ID = ruleMetadata.RuleReferenceID
	policy.policyMetadata.Category = ruleMetadata.Category
	policy.policyMetadata.Version = ruleMetadata.Version

	templateString, ok := ruleMetadata.RuleArgument.(string)
	if !ok && templateString != "" {
		return policy, fmt.Errorf("incorrect rule argument type, must be a []byte")
	}
	err := json.Unmarshal([]byte(templateString), &templateArgs)
	if err != nil {
		return policy, fmt.Errorf("error could not unmarshal rule arguments into map[string]interface{}: '%v'", err)
	}
	policy.policyMetadata.TemplateArgs = templateArgs

	return policy, nil
}

func saveCommercialPolicies(policy commercialPolicy, policyRepoPath string) error {
	const tabSpace = "    "

	policy.policyMetadata.PolicyType = getCSP(policy.resourceType)
	cspDir := filepath.Join(policyRepoPath, policy.policyMetadata.PolicyType)
	err := ensureDir(cspDir)
	if err != nil {
		return err
	}

	resourceDir := filepath.Join(cspDir, policy.resourceType)
	err = ensureDir(resourceDir)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", tabSpace)
	err = encoder.Encode(policy.policyMetadata)
	if err != nil {
		return fmt.Errorf("error could not marshal json object into byte array: '%v'", err)
	}

	metadata := buffer.Bytes()
	metadata = bytes.TrimRight(metadata, "\n")
	metaDataPath := filepath.Join(resourceDir, policy.metadataFileName)
	err = ioutil.WriteFile(metaDataPath, metadata, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error could not write rule metadata file on disk: '%v'", err)
	}

	regoPath := filepath.Join(resourceDir, policy.policyMetadata.File)
	err = ioutil.WriteFile(regoPath, []byte(policy.regoTemplate), filePermissionBits)
	if err != nil {
		return fmt.Errorf("error could not write rego code file on disk: '%v'", err)
	}

	return nil
}

func getCSP(resourceType string) string {
	csp := strings.ToLower(resourceType)

	if strings.HasPrefix(csp, "azure") {
		return "azure"
	}

	if strings.HasPrefix(csp, "google") {
		return "gcp"
	}

	if strings.HasPrefix(csp, "kubernetes") {
		return "k8s"
	}

	return strings.Split(csp, "_")[0]
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, filePermissionBits)
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
