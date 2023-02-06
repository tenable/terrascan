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

package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/downloader"
	admissionwebhook "github.com/tenable/terrascan/pkg/k8s/admission-webhook"
	"github.com/tenable/terrascan/pkg/runtime"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// scanRemoteRepoReq contains request body for remote repository scanning
type scanRemoteRepoReq struct {
	RemoteType               string   `json:"remote_type"`
	RemoteURL                string   `json:"remote_url"`
	ConfigOnly               bool     `json:"config_only"`
	ScanRules                []string `json:"scan_rules"`
	SkipRules                []string `json:"skip_rules"`
	Categories               []string `json:"categories"`
	Severity                 string   `json:"severity"`
	ShowPassed               bool     `json:"show_passed"`
	NonRecursive             bool     `json:"non_recursive"`
	FindVulnerabilities      bool     `json:"find_vulnerabilities"`
	d                        downloader.Downloader
	NotificationWebhookURL   string `json:"webhook_url"`
	NotificationWebhookToken string `json:"webhook_token"`
	RepoRef                  string `json:"repo-ref"`
	ConfigWithError          bool   `json:"config_with_error"`
}

// scanRemoteRepo downloads the remote Iac repository and scans it for
// violations
func (g *APIHandler) scanRemoteRepo(w http.ResponseWriter, r *http.Request) {
	zap.S().Debug("handle: remote repository scan request")

	// get url params
	params := mux.Vars(r)
	var (
		// url params
		iacType    = params["iac"]
		iacVersion = params["iacVersion"]
		cloudType  = strings.Split(params["cloud"], ",")
	)

	// read request body
	var s scanRemoteRepoReq
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zap.S().Debugf("scanning remote repository request: %+v", s)

	// scan remote repo
	s.d = downloader.NewDownloader()
	var results interface{}
	var isAdmissionDenied bool
	if g.test {
		results, isAdmissionDenied, err = s.ScanRemoteRepo(iacType, iacVersion, cloudType, []string{"./testdata/testpolicies"})

	} else {
		results, isAdmissionDenied, err = s.ScanRemoteRepo(iacType, iacVersion, cloudType, getPolicyPathFromConfig())
	}
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert results into JSON
	j, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("failed to create JSON. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	// return with results
	// if result contain violations denied by admission controller return 403 status code
	if isAdmissionDenied {
		apiResponse(w, string(j), http.StatusForbidden)
		return
	}
	apiResponse(w, string(j), http.StatusOK)
}

// ScanRemoteRepo is the actual method where a remote repo is downloaded and
// scanned for violations
func (s *scanRemoteRepoReq) ScanRemoteRepo(iacType, iacVersion string, cloudType []string, policyPath []string) (interface{}, bool, error) {

	// return params
	var (
		output            interface{}
		err               error
		isAdmissionDenied bool
	)

	// temp destination directory to download remote repo
	tempDir := utils.GenerateTempDir()
	defer os.RemoveAll(tempDir)

	// download remote repository
	iacDirPath, err := s.d.DownloadWithType(s.RemoteType, s.RemoteURL, tempDir)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download remote repo. error: '%v'", err)
		zap.S().Error(errMsg)
		return output, isAdmissionDenied, err
	}

	// create a new runtime executor for scanning the remote repo
	executor, err := runtime.NewExecutor(iacType, iacVersion, cloudType,
		"", iacDirPath, policyPath, s.ScanRules, s.SkipRules, s.Categories, s.Severity, s.NonRecursive, false, s.FindVulnerabilities, s.NotificationWebhookURL, s.NotificationWebhookToken, s.RemoteURL, s.RepoRef, []string{})
	if err != nil {
		zap.S().Error(err)
		return output, isAdmissionDenied, err
	}

	// evaluate policies IaC for violations
	results, err := executor.Execute(s.ConfigOnly, s.ConfigWithError)
	if err != nil {
		errMsg := fmt.Sprintf("failed to scan uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		return output, isAdmissionDenied, err
	}

	// if config-with-error return config as well as dir errors,for config only, return resource config else return violations
	if s.ConfigWithError {
		output = results
	} else if s.ConfigOnly {
		output = results.ResourceConfig
	} else {
		// set remote url in case remote repo is scanned
		if s.RemoteURL != "" {
			results.Violations.Summary.ResourcePath = s.RemoteURL
		}

		if !s.ShowPassed {
			results.Violations.ViolationStore.PassedRules = nil
		}

		isAdmissionDenied = hasK8sAdmissionDeniedViolations(results)
		output = results.Violations
	}

	// succesful
	return output, isAdmissionDenied, nil
}

// hasK8sAdmissionDeniedViolations checks if violations have denied by k8s admission controller
func hasK8sAdmissionDeniedViolations(o runtime.Output) bool {
	denyRuleMatcher := admissionwebhook.WebhookDenyRuleMatcher{}
	for _, v := range o.Violations.ViolationStore.Violations {
		if denyRuleMatcher.Match(*v, config.GetK8sAdmissionControl()) {
			return true
		}
	}
	return false
}
