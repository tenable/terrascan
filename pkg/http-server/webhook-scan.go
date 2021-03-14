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

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeK8s "k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/runtime/serializer"

	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	v1 "k8s.io/api/admission/v1"
)

// validateK8SWebhook handles the incoming validating admission webhook from kubernetes API server
func (g *APIHandler) validateK8SWebhook(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()

	params := mux.Vars(r)
	apiKey := params["apiKey"]

	// Validate if authorized (API key is specified and matched the server one (saved in an environment variable)
	if !g.validateAuthorization(apiKey, w) {
		return
	}

	// Read the request into byte array
	bytesRequestAdmissionReview, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Failed to read admission review: '%v'", err)
		apiErrorResponse(w, msg, http.StatusBadRequest)
		return
	}

	zap.S().Debugf("scanning configuration webhook request: %+v", string(bytesRequestAdmissionReview))

	// Unmarshal the byte array into a v1.AdmissionReview object
	requestedAdmissionReview, err := g.deserializeAdmissionReviewRequest(bytesRequestAdmissionReview)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In case the object is nil => an operation of DELETE happened, just return 'allow' since there is nothing to check
	if len(requestedAdmissionReview.Request.Object.Raw) < 1 {
		g.sendResponseAdmissionReview(w, *requestedAdmissionReview, true, nil, "")
		return
	}

	// Save the object into a temp file for the policy engines
	tempFile, err := g.writeObjectToTempFile(requestedAdmissionReview.Request.Object.Raw)
	defer os.Remove(tempFile.Name())
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Run the policy engines
	output, err := g.executeEngines(*tempFile)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate if there are anydeny violations
	denyViolations, err := g.getDenyViolations(*output)
	allowed := len(denyViolations) < 1
	logPath := g.getLogPath(r.Host, apiKey, string(requestedAdmissionReview.Request.UID))

	// Log the request in the DB
	err = g.logWebhook(*output, string(requestedAdmissionReview.Request.UID), bytesRequestAdmissionReview, denyViolations, currentTime, allowed)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the correct response according to the result
	g.sendResponseAdmissionReview(w, *requestedAdmissionReview, allowed, output, logPath)
}

func (g *APIHandler) validateAuthorization(apiKey string, w http.ResponseWriter) bool {
	if len(apiKey) < 1 {
		msg := "apiKey is missing"
		zap.S().Error(msg)
		apiErrorResponse(w, msg, http.StatusBadRequest)
		return false
	}

	savedApiKey := os.Getenv("K8S_WEBHOOK_API_KEY")
	if len(savedApiKey) < 1 {
		msg := "K8S_WEBHOOK_API_KEY environment variable MUST be declared"
		zap.S().Error(msg)
		apiErrorResponse(w, msg, http.StatusInternalServerError)
		return false
	}

	if apiKey != savedApiKey {
		msg := "Invalid apiKey"
		zap.S().Error(msg)
		apiErrorResponse(w, msg, http.StatusUnauthorized)
		return false
	}

	return true
}

func (g *APIHandler) getDeniedViolations(violations results.ViolationStore, denyRules config.K8sDenyRules) []*results.Violation {
	// Check whether one of the violations matches the deny violations configuration

	var denyViolations []*results.Violation

	denyRuleMatcher := webhookDenyRuleMatcher{}

	for _, violation := range violations.Violations {
		if denyRuleMatcher.match(*violation, denyRules) {
			denyViolations = append(denyViolations, violation)
		}
	}

	return denyViolations
}

func (g *APIHandler) writeObjectToTempFile(objectBytes []byte) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", "terrascan-*.json")
	if err != nil {
		zap.S().Errorf("failed to create temp file: '%v'", err)
		return nil, err
	}

	zap.S().Debugf("created temp config file at '%s'", tempFile.Name())

	_, err = tempFile.Write(objectBytes)
	if err != nil {
		zap.S().Errorf("failed to write object to temp file: '%v'", err)
		return nil, err
	}

	return tempFile, nil
}

func (g *APIHandler) executeEngines(tempFile os.File) (*runtime.Output, error) {
	var executor *runtime.Executor
	var err error
	if g.test {
		executor, err = runtime.NewExecutor("k8s", "v1", []string{"k8s"},
			tempFile.Name(), "", g.configFile, []string{"./k8s_testdata/testpolicies"}, []string{}, []string{}, []string{}, "")
	} else {
		executor, err = runtime.NewExecutor("k8s", "v1", []string{"k8s"},
			tempFile.Name(), "", g.configFile, []string{}, []string{}, []string{}, []string{}, "")
	}

	if err != nil {
		zap.S().Errorf("failed to create runtime executer: '%v'", err)
		return nil, err
	}

	result, err := executor.Execute()
	if err != nil {
		zap.S().Error("failed to scan resource object. error: '%v'", err)
		return nil, err
	}

	return &result, nil
}

func (g *APIHandler) getDenyViolations(output runtime.Output) ([]*results.Violation, error) {
	// Calcualte the deny violations according to the configuration specified in the config file
	configReader, err := config.NewTerrascanConfigReader(g.configFile)
	if err != nil {
		zap.S().Errorf("error loading config file: '%v'", err)
		return nil, err
	}

	denyViolations := g.getDeniedViolations(*output.Violations.ViolationStore, configReader.GetK8sDenyRules())

	return denyViolations, nil
}

func (g *APIHandler) sendResponseAdmissionReview(w http.ResponseWriter,
	requestedAdmissionReview v1.AdmissionReview,
	allowed bool,
	output *runtime.Output,
	logPath string) {
	responseAdmissionReview := &v1.AdmissionReview{}
	responseAdmissionReview.SetGroupVersionKind(requestedAdmissionReview.GroupVersionKind())

	responseAdmissionReview.Response = &v1.AdmissionResponse{
		UID:     requestedAdmissionReview.Request.UID,
		Allowed: allowed,
	}

	if output != nil {
		// Means we ran the engines and we have results
		if allowed {
			if len(output.Violations.ViolationStore.Violations) > 0 {
				// In case there are no denial violations, just return the log URL as a warning
				responseAdmissionReview.Response.Warnings = []string{logPath}
			}
		} else {
			// In case the request was denied, return 403 and the log URL as an error message
			responseAdmissionReview.Response.Result = &metav1.Status{Message: logPath, Code: 403}
		}
	}

	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		msg := fmt.Sprintf("failed to serialize admission review response: %v", err)
		zap.S().Error(msg)
		apiErrorResponse(w, msg, http.StatusInternalServerError)
	}

	zap.S().Debugf("Response result: %+v", string(respBytes))

	apiResponse(w, string(respBytes), http.StatusOK)
}

func (g *APIHandler) logWebhook(output runtime.Output,
	uid string,
	bytesAdmissionReview []byte,
	denyViolations []*results.Violation,
	currentTime time.Time,
	allowed bool) error {
	var deniedViolationsEncoded string

	if len(denyViolations) < 1 {
		deniedViolationsEncoded = ""
	} else {
		d, _ := json.Marshal(denyViolations)
		deniedViolationsEncoded = string(d)
	}

	encodedViolationsSummary, _ := json.Marshal(output.Violations.ViolationStore)

	logger := WebhookScanLogger{
		test: g.test,
	}

	err := logger.log(webhookScanLog{
		UID:                uid,
		Request:            string(bytesAdmissionReview),
		Allowed:            allowed,
		DeniableViolations: deniedViolationsEncoded,
		ViolationsSummary:  string(encodedViolationsSummary),
		CreatedAt:          currentTime,
	})
	if err != nil {
		zap.S().Error("error logging scan result: '%v'", err)
		return err
	}

	return nil
}

func (g *APIHandler) deserializeAdmissionReviewRequest(bytesAdmissionReview []byte) (*v1.AdmissionReview, error) {
	var scheme = runtimeK8s.NewScheme()
	v1.AddToScheme(scheme)

	var codecs = serializer.NewCodecFactory(scheme)
	deserializer := codecs.UniversalDeserializer()

	obj, _, err := deserializer.Decode(bytesAdmissionReview, nil, nil)
	if err != nil {
		zap.S().Errorf("Request could not be decoded: %v", err)
		return nil, err
	}

	requestedAdmissionReview, ok := obj.(*v1.AdmissionReview)
	if !ok {
		zap.S().Errorf("Failed to deserialize request body to v1.AdmissionReview. Obj: %v", obj)
		return nil, err
	}

	return requestedAdmissionReview, nil
}
