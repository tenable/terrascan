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
	"strconv"
	"strings"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// scanFile accepts uploaded file and runs scan on it
func (g *APIHandler) scanFile(w http.ResponseWriter, r *http.Request) {

	// get url params
	params := mux.Vars(r)
	var (
		iacType    = params["iac"]
		iacVersion = params["iacVersion"]
		cloudType  = strings.Split(params["cloud"], ",")
		scanRules  = []string{}
		skipRules  = []string{}
		configOnly = false
		showPassed = false
		categories = []string{}
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

	zap.S().Debugf("uploaded file: %+v", handler.Filename)
	zap.S().Debugf("file size: %+v", handler.Size)
	zap.S().Debugf("MIME header: %+v", handler.Header)

	// Create a temporary file within temp directory
	tempFile, err := ioutil.TempFile("", "terrascan-*.tf")
	if err != nil {
		errMsg := fmt.Sprintf("failed to create temp file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())
	zap.S().Debugf("create temp config file at '%s'", tempFile.Name())

	// read all of the contents of uploaded file
	fileBytes, err := ioutil.ReadAll(file)
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

	// categories is the list categories of violations that the user want to get informed about: low, medium or high
	categoriesValue := r.FormValue("categories")

	// severity is the minimum severity level of violations that the user want to get informed about: low, medium or high
	severity := r.FormValue("severity")

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
			tempFile.Name(), "", "", []string{"./testdata/testpolicies"}, scanRules, skipRules, categories, severity)
	} else {
		executor, err = runtime.NewExecutor(iacType, iacVersion, cloudType,
			tempFile.Name(), "", "", getPolicyPathFromEnv(), scanRules, skipRules, categories, severity)
	}
	if err != nil {
		zap.S().Error(err)
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	normalized, err := executor.Execute()
	if err != nil {
		errMsg := fmt.Sprintf("failed to scan uploaded file. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	var output interface{}

	if !showPassed {
		normalized.Violations.ViolationStore.PassedRules = nil
	}

	// if config only, return resource config else return violations
	if configOnly {
		output = normalized.ResourceConfig
	} else {
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

// getPolicyPathFromEnv reads the TERRASCAN_CONFIG env variable (if present) and returns the policy path
func getPolicyPathFromEnv() []string {
	policyPath := []string{}

	// read policy path from TERRASCAN_CONFIG env variable
	terrascanConfigFile := os.Getenv("TERRASCAN_CONFIG")
	if terrascanConfigFile != "" {
		terrascanConfigReader, err := config.NewTerrascanConfigReader(terrascanConfigFile)
		if err != nil {
			zap.S().Errorf("error while reading config file, %s; err %v", terrascanConfigFile, err)
		} else {
			policyPath = append(policyPath, terrascanConfigReader.GetPolicyConfig().RepoPath)
		}
	}
	return policyPath
}
