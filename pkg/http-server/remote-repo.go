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

package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// scanRemoteRepoReq contains request body for remote repository scanning
type scanRemoteRepoReq struct {
	RemoteType string `json:"remote_type"`
	RemoteURL  string `json:"remote_url"`
	ConfigOnly bool   `json:"config_only"`
	d          downloader.Downloader
}

// scanRemoteRepo downloads the remote Iac repository and scans it for
// violations
func (g *APIHandler) scanRemoteRepo(w http.ResponseWriter, r *http.Request) {

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
	if g.test {
		results, err = s.ScanRemoteRepo(iacType, iacVersion, cloudType, []string{"./testdata/testpolicies"})

	} else {
		results, err = s.ScanRemoteRepo(iacType, iacVersion, cloudType, []string{})
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
	apiResponse(w, string(j), http.StatusOK)
}

// ScanRemoteRepo is the actual method where a remote repo is downloaded and
// scanned for violations
func (s *scanRemoteRepoReq) ScanRemoteRepo(iacType, iacVersion string, cloudType []string, policyPath []string) (interface{}, error) {

	// return params
	var (
		output interface{}
		err    error
	)

	// temp destination directory to download remote repo
	tempDir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
	defer os.RemoveAll(tempDir)

	// download remote repository
	iacDirPath, err := s.d.DownloadWithType(s.RemoteType, s.RemoteURL, tempDir)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download remote repo. error: '%v'", err)
		zap.S().Error(errMsg)
		return output, err
	}

	// create a new runtime executor for scanning the remote repo
	executor, err := runtime.NewExecutor(iacType, iacVersion, cloudType,
		"", iacDirPath, "", policyPath)
	if err != nil {
		zap.S().Error(err)
		return output, err
	}

	// evaluate policies IaC for violations
	results, err := executor.Execute()
	if err != nil {
		errMsg := fmt.Sprintf("failed to scan uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		return output, err
	}

	// if config only, return only config else return only violations
	if s.ConfigOnly {
		output = results.ResourceConfig
	} else {
		output = results.Violations
	}

	// succesful
	return output, nil
}
