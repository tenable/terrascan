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
	"time"

	admissionWebhook "github.com/accurics/terrascan/pkg/k8s/admission-webhook"
	admissionwebhook "github.com/accurics/terrascan/pkg/k8s/admission-webhook"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// validateK8SWebhook handles the incoming validating admission webhook from kubernetes API server
func (g *APIHandler) validateK8SWebhook(w http.ResponseWriter, r *http.Request) {

	var (
		currentTime       = time.Now()
		params            = mux.Vars(r)
		apiKey            = params["apiKey"]
		validatingWebhook = admissionWebhook.NewValidatingWebhook(g.configFile)
	)

	// Validate if authorized (API key is specified and matched the server one (saved in an environment variable)
	if err := validatingWebhook.Authorize(apiKey); err != nil {
		switch err {
		case admissionWebhook.ErrAPIKeyMissing:
			apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		case admissionwebhook.ErrAPIKeyEnvNotSet:
			apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		case admissionWebhook.ErrUnAuthorized:
			apiErrorResponse(w, err.Error(), http.StatusUnauthorized)
		}
		return
	}

	// Read the request into byte array
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read validating admission webhook body, error: '%v'", err)
		apiErrorResponse(w, msg, http.StatusBadRequest)
		return
	}

	zap.S().Debugf("scanning configuration webhook request: %+v", string(payload))

	// decode incoming admission review request
	requestedAdmissionReview, err := validatingWebhook.DecodeAdmissionReviewRequest(payload)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// process the admission review request
	output, allowed, denyViolations, err := validatingWebhook.ProcessWebhook(requestedAdmissionReview)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logPath := g.getLogPath(r.Host, string(requestedAdmissionReview.Request.UID))

	// Log the request in the DB
	err = g.logWebhook(*output, string(requestedAdmissionReview.Request.UID), payload, denyViolations, currentTime, allowed)
	if err != nil {
		apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the correct response according to the result
	g.sendResponseAdmissionReview(w, requestedAdmissionReview, allowed, output, logPath)
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
