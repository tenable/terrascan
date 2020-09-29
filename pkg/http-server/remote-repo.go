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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// scanRemoteRepoReq contains request body for remote repository scanning
type scanRemoteRepoReq struct {
	remoteType string `json:"remote_type"`
	remoteURL  string `json:"remote_url"`
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
		cloudType  = params["cloud"]
	)

	// read request body
	var s scanRemoteRepoReq
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("request body:\n%v\n", string(body))
	err := json.Unmarshal(body, &s)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	zap.S().Debugf("scanning remote repository request: %+v", s)

	// temp destination directory to download remote repo
	tempDir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
	defer os.RemoveAll(tempDir)

	// download remote repository
	d := downloader.NewDownloader()
	iacDirPath, err := d.DownloadWithType(s.remoteType, s.remoteURL, tempDir)
	if err != nil {
		errMsg := fmt.Sprintf("failed to download remote repo. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	// create a new runtime executor for scanning the remote repo
	var executor *runtime.Executor
	if g.test {
		// executor, err = runtime.NewExecutor(iacType, iacVersion, cloudType,
		//	tempFile.Name(), "", "", "./testdata/testpolicies")
	} else {
		executor, err = runtime.NewExecutor(iacType, iacVersion, cloudType,
			"", iacDirPath, "", "")
	}
	if err != nil {
		zap.S().Error(err)
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := executor.Execute()
	if err != nil {
		errMsg := fmt.Sprintf("failed to scan uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	j, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("failed to create JSON. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	// return that we have successfully uploaded our file!
	apiResponse(w, string(j), http.StatusOK)
}
