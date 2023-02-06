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
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/runtime"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// scanFile accepts uploaded file and runs scan on it
func (g *APIHandler) scanFile(w http.ResponseWriter, r *http.Request) {
	zap.S().Debug("handle: file scan request")

	// get url params
	params := mux.Vars(r)
	var (
		iacType             = params["iac"]
		iacVersion          = params["iacVersion"]
		cloudType           = strings.Split(params["cloud"], ",")
		scanRules           = []string{}
		skipRules           = []string{}
		configOnly          = false
		showPassed          = false
		findVulnerabilities = false
		categories          = []string{}
		configWithError     = false
	)

	// parse multipart form, 10 << 20 specifies maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		errMsg := fmt.Sprintf("failed to retreive uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// fileExtension will include the period. (eg ".yaml")
	fileExtension := path.Ext(handler.Filename)

	zap.S().Debugf("uploaded file: %+v", handler.Filename)
	zap.S().Debugf("uploaded file extension: %+v", fileExtension)
	zap.S().Debugf("file size: %+v", handler.Size)
	zap.S().Debugf("MIME header: %+v", handler.Header)

	// Create a temporary file within temp directory
	tempFileTemplate := fmt.Sprintf("terrascan-*%s", fileExtension)
	tempFile, err := os.CreateTemp("", tempFileTemplate)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create temp file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	zap.S().Debugf("create temp config file at '%s'", tempFile.Name())

	// read all of the contents of uploaded file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// read scan and skip rules from the form data
	// scan and skip rules are comma separated rule id's in the request body
	scanRulesValue := r.FormValue("scan_rules")
	skipRulesValue := r.FormValue("skip_rules")
	notificationWebhookURL := r.FormValue("webhook_url")
	notificationWebhookToken := r.FormValue("webhook_token")

	// categories is the list categories of violations that the user want to get informed about: low, medium or high
	categoriesValue := r.FormValue("categories")

	// severity is the minimum severity level of violations that the user want to get informed about: low, medium or high
	severity := r.FormValue("severity")

	findVulnerabilitiesValue := r.FormValue("find_vulnerabilities")
	if findVulnerabilitiesValue != "" {
		findVulnerabilities, err = strconv.ParseBool(findVulnerabilitiesValue)
		if err != nil {
			errMsg := fmt.Sprintf("error while reading 'find_vulnerabilities' value. error: '%v'", err)
			zap.S().Error(errMsg)
			apiErrorResponse(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	// read config_only from the form data
	configOnlyValue := r.FormValue("config_only")
	if configOnlyValue != "" {
		configOnly, err = strconv.ParseBool(configOnlyValue)
		if err != nil {
			errMsg := fmt.Sprintf("error while reading 'config_only' value. error: '%v'", err)
			zap.S().Error(errMsg)
			apiErrorResponse(w, errMsg, http.StatusBadRequest)
			return
		}
	}
	// read config_with_error from the form data
	configWithErrorValue := r.FormValue("config_with_error")
	if configWithErrorValue != "" {
		configWithError, err = strconv.ParseBool(configWithErrorValue)
		if err != nil {
			errMsg := fmt.Sprintf("error while reading 'config_with_error' value. error: '%v'", err)
			zap.S().Error(errMsg)
			apiErrorResponse(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	// read show_passed from the form data
	showPassedValue := r.FormValue("show_passed")
	if showPassedValue != "" {
		showPassed, err = strconv.ParseBool(showPassedValue)
		if err != nil {
			errMsg := fmt.Sprintf("error while reading 'show_passed' value. error: '%v'", err)
			zap.S().Error(errMsg)
			apiErrorResponse(w, errMsg, http.StatusBadRequest)
			return
		}
	}

	if scanRulesValue != "" {
		scanRules = strings.Split(scanRulesValue, ",")
	}

	if skipRulesValue != "" {
		skipRules = strings.Split(skipRulesValue, ",")
	}

	if categoriesValue != "" {
		categories = strings.Split(categoriesValue, ",")
	}

	if severity != "" {
		severity = utils.EnsureUpperCaseTrimmed(severity)
	}

	// create a new runtime executor for scanning the uploaded file
	var executor *runtime.Executor
	if g.test {
		executor, err = runtime.NewExecutor(iacType, iacVersion, cloudType,
			tempFile.Name(), "", []string{"./testdata/testpolicies"}, scanRules, skipRules, categories, severity, false, false, false, notificationWebhookURL, notificationWebhookToken, "", "", []string{})
	} else {
		executor, err = runtime.NewExecutor(iacType, iacVersion, cloudType,
			tempFile.Name(), "", getPolicyPathFromConfig(), scanRules, skipRules, categories, severity, false, false, findVulnerabilities, notificationWebhookURL, notificationWebhookToken, "", "", []string{})
	}
	if err != nil {
		zap.S().Error(err)
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	normalized, err := executor.Execute(configOnly, configWithError)
	if err != nil {
		errMsg := fmt.Sprintf("failed to scan uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	var output interface{}

	// if config-with-error return config as well as dir errors,for config only, return resource config else return violations
	if configWithError {
		output = normalized
	} else if configOnly {
		output = normalized.ResourceConfig
	} else {
		if !showPassed {
			normalized.Violations.ViolationStore.PassedRules = nil
		}
		output = normalized.Violations
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

// getPolicyPathFromConfig returns the policy path from config
func getPolicyPathFromConfig() []string {
	return []string{config.GetPolicyRepoPath()}
}
