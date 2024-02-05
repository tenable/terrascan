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

package initialize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tenable/terrascan/pkg/config"
	"go.uber.org/zap"
)

var (
	errNoConnection = fmt.Errorf("could not connect to github.com")
)

const terrascanReadmeURL string = "https://raw.githubusercontent.com/tenable/terrascan/master/README.md"
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
		return fmt.Errorf("unable to delete base folder. error: '%w'", err)
	}

	if accessToken == "" {
		return downloadDefaultPolicies(policyBasePath)
	}

	return downloadEnvironmentPolicies(policyBasePath, accessToken)
}

func downloadEnvironmentPolicies(policyBasePath, accessToken string) error {
	err := ensureDir(policyBasePath)
	if err != nil {
		return err
	}

	policyRepoPath := config.GetPolicyRepoPath()
	err = os.MkdirAll(policyRepoPath, filePermissionBits)
	if err != nil {
		return fmt.Errorf("unable to prepare directories representing policyRepoPath. err: '%w', policyRepoPath: '%s'", err, policyRepoPath)
	}

	const apiPath = "/v1/api/app/rules?default=true"
	environment := config.GetPolicyEnvironment()

	zap.S().Debugf("policy environment : %s", environment)
	zap.S().Debugf("downloading environment policies in %s", policyBasePath)

	var client http.Client

	req, err := http.NewRequest(http.MethodGet, environment+apiPath, nil)
	if err != nil {
		return fmt.Errorf("error constructing request object. error: '%w'", err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading environment policies. error: '%w'", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error downloading environment policies, response status code: '%d'", res.StatusCode)
	}

	policies, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading api call response for environment policies. error: '%w'", err)
	}

	err = convertEnvironmentPolicies(policies, policyRepoPath)
	if err != nil {
		return err
	}

	// creating empty docker folder to satisfy folder structure dep
	dockerPath := filepath.Join(policyRepoPath, "docker")
	err = os.Mkdir(dockerPath, filePermissionBits)
	if err != nil {
		return fmt.Errorf("unable to create empty docker dir. error: '%w'", err)
	}

	return nil
}

func convertEnvironmentPolicies(policies []byte, policyRepoPath string) error {
	var ruleMetadataList []environmentPolicyMetadata

	err := json.Unmarshal(policies, &ruleMetadataList)
	if err != nil {
		return fmt.Errorf("failed to unmarshal policies into structure. error: '%w'", err)
	}

	for _, ruleMetadata := range ruleMetadataList {
		policy, err := newPolicy(ruleMetadata)
		if err != nil {
			return err
		}

		err = saveEnvironmentPolicies(policy, policyRepoPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveEnvironmentPolicies(policy environmentPolicy, policyRepoPath string) error {
	policy.policyMetadata.PolicyType = policy.getType()
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
	encoder.SetIndent("", "    ")
	err = encoder.Encode(policy.policyMetadata)
	if err != nil {
		return fmt.Errorf("could not marshal json object into byte array. error: '%w'", err)
	}

	metadata := buffer.Bytes()
	metadata = bytes.TrimRight(metadata, "\n")
	metaDataPath := filepath.Join(resourceDir, policy.metadataFileName)
	err = os.WriteFile(metaDataPath, metadata, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write rule metadata file on disk. error: '%w'", err)
	}

	regoPath := filepath.Join(resourceDir, policy.policyMetadata.File)

	if _, err := os.Stat(regoPath); os.IsExist(err) {
		zap.S().Debug("rego code file %s exists, skipping", regoPath)
		return nil
	}

	err = os.WriteFile(regoPath, []byte(policy.regoTemplate), filePermissionBits)
	if err != nil {
		return fmt.Errorf("could not write rego code file on disk. error: '%w'", err)
	}

	return nil
}

func ensureDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, filePermissionBits)
		if err != nil {
			return fmt.Errorf("unable to create requested directory error: '%w' for path: '%s'", err, path)
		}
	}
	return nil
}

func downloadDefaultPolicies(policyBasePath string) error {
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
		return fmt.Errorf("failed to download policies. error: '%w'", err)
	}

	// create working tree
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to create working tree. error: '%w'", err)
	}

	// fetch references
	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch references from git repo. error: '%w'", err)
	}

	// checkout policies branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout git branch '%s'. error: '%w'", branch, err)
	}

	return nil
}

func connected(url string) bool {
	_, err := http.Get(url)
	return err == nil
}
