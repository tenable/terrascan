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

package admissionwebhook

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/k8s/dblogs"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/runtime"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeK8s "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// ValidatingWebhook handles the incoming validating admission webhook from
// the kubernetes API server and decides whether the admission request from
// the kubernetes client should be allowed or not
type ValidatingWebhook struct {
	requestBody              []byte
	dblogger                 *dblogs.WebhookScanLogger
	notificationWebhookURL   string
	notificationWebhookToken string
	repoURL                  string
	repoRef                  string
}

// NewValidatingWebhook returns a new, empty ValidatingWebhook struct
func NewValidatingWebhook(body []byte, notificationWebhookURL, notificationWebhookToken, repoURL, repoRef string) AdmissionWebhook {
	return ValidatingWebhook{
		dblogger:                 dblogs.NewWebhookScanLogger(),
		requestBody:              body,
		notificationWebhookURL:   notificationWebhookURL,
		notificationWebhookToken: notificationWebhookToken,
		repoURL:                  repoURL,
		repoRef:                  repoRef,
	}
}

var (
	// ErrAPIKeyMissing indicates that API key is missing in webhook request
	ErrAPIKeyMissing = fmt.Errorf("apiKey is missing in validating admission webhook url")

	// ErrAPIKeyEnvNotSet indicates K8S_WEBHOOK_API_KEY is not set in terrascan server env
	ErrAPIKeyEnvNotSet = fmt.Errorf("variable K8S_WEBHOOK_API_KEY not set in terrascan server environment")

	// ErrUnauthorized means user is not authorized to make this call
	ErrUnauthorized = fmt.Errorf("invalid API key in validating admission webhook url")

	// ErrEmptyAdmissionReview empty admission review request
	ErrEmptyAdmissionReview = fmt.Errorf("empty admission review request")

	// test policies path
	testPoliciesPath = filepath.Join("../", "policies", "opa", "rego", "k8s")
)

// Authorize checks if the incoming webhooks have valid apiKey
func (w ValidatingWebhook) Authorize(apiKey string) error {

	// check if key exists in API request
	if len(apiKey) < 1 {
		zap.S().Error(ErrAPIKeyMissing)
		return ErrAPIKeyMissing
	}

	// API key not set in terrascan env
	saveAPIKey := os.Getenv("K8S_WEBHOOK_API_KEY")
	if len(saveAPIKey) < 1 {
		zap.S().Error(ErrAPIKeyEnvNotSet)
		return ErrAPIKeyEnvNotSet
	}

	// invalid api key
	if apiKey != saveAPIKey {
		zap.S().Error(ErrUnauthorized)
		return ErrUnauthorized
	}

	return nil
}

// DecodeAdmissionReviewRequest reads the incoming admission request body,
// decodes it and returns an AdmissionReviewRequest struct
func (w ValidatingWebhook) DecodeAdmissionReviewRequest(requestBody []byte) (admissionv1.AdmissionReview, error) {

	var (
		scheme                   = runtimeK8s.NewScheme()
		codecs                   = serializer.NewCodecFactory(scheme)
		deserializer             = codecs.UniversalDeserializer()
		requestedAdmissionReview admissionv1.AdmissionReview
	)
	admissionv1.AddToScheme(scheme)

	// decode incoming admission request
	_, _, err := deserializer.Decode(requestBody, nil, &requestedAdmissionReview)
	if err != nil {
		errMsg := "failed to decode validating admission webhook request body"
		zap.S().Error(errMsg, zap.Error(err))
		return requestedAdmissionReview, fmt.Errorf("%s, error: %w", errMsg, err)
	}

	return requestedAdmissionReview, nil
}

// ProcessWebhook processes the incoming AdmissionReview and creates
// a response
func (w ValidatingWebhook) ProcessWebhook(review admissionv1.AdmissionReview, serverURL string) (*admissionv1.AdmissionReview, error) {

	var (
		output         runtime.Output
		denyViolations []results.Violation
		logMsg         string
		allowed        = false
		errMsg         = "%s; error: %w"
	)

	// when admission controller webhook runs in dashboard mode, request and response are logged
	// to database and we would display the url where the logs would be available to the user
	if config.GetK8sAdmissionControl().Dashboard {
		logMsg = w.dblogger.GetLogURL(serverURL, string(review.Request.UID))
	}

	// In case the object is nil => an operation of DELETE happened, just return 'allow' since there is nothing to check
	if len(review.Request.Object.Raw) < 1 {
		zap.S().Info(ErrEmptyAdmissionReview, zap.Any("admission review object", review))
		return w.createResponseAdmissionReview(review, true, output, logMsg), ErrEmptyAdmissionReview
	}

	// Save the object into a temp file for the policy engines
	tempFile, err := utils.CreateTempFile(review.Request.Object.Raw, "json")
	defer os.Remove(tempFile.Name())
	if err != nil {
		msg := "failed to create temp file for validating admission review request"
		zap.S().Error(msg, zap.Error(err))
		return w.createResponseAdmissionReview(review, allowed, output, logMsg), fmt.Errorf(errMsg, msg, err)
	}

	// Run the policy engines
	output, err = w.scanK8sFile(tempFile.Name())
	if err != nil {
		msg := "failed to evaluate terrascan policies"
		zap.S().Errorf(msg, zap.Error(err))
		return w.createResponseAdmissionReview(review, allowed, output, logMsg), fmt.Errorf(errMsg, msg, err)
	}

	// Calculate if there are any deny violations
	denyViolations, err = w.getDenyViolations(output)
	if err != nil {
		msg := "failed to figure out denied violations"
		zap.S().Errorf(msg, zap.Error(err))
		return w.createResponseAdmissionReview(review, allowed, output, logMsg), fmt.Errorf(errMsg, msg, err)
	}
	allowed = len(denyViolations) < 1

	// Log the request in the DB
	if config.GetK8sAdmissionControl().Dashboard {
		err = w.logWebhook(output, string(review.Request.UID), denyViolations, allowed)
		if err != nil {
			msg := "failed to log validating admission review request into database"
			zap.S().Error(msg, zap.Error(err))
		}
	}

	return w.createResponseAdmissionReview(review, allowed, output, logMsg), nil
}

func (w ValidatingWebhook) scanK8sFile(filePath string) (runtime.Output, error) {

	var (
		executor *runtime.Executor
		err      error
		result   runtime.Output
	)

	if flag.Lookup("test.v") != nil {
		executor, err = runtime.NewExecutor("k8s", "v1", []string{"k8s"},
			filePath, "", []string{testPoliciesPath}, []string{}, []string{}, []string{}, "", false, false, false, w.notificationWebhookURL, w.notificationWebhookToken, w.repoURL, w.repoRef, []string{})
	} else {
		executor, err = runtime.NewExecutor("k8s", "v1", []string{"k8s"},
			filePath, "", []string{}, []string{}, []string{}, []string{}, "", false, false, false, w.notificationWebhookURL, w.notificationWebhookToken, w.repoURL, w.repoRef, []string{})
	}
	if err != nil {
		zap.S().Errorf("failed to create runtime executer: '%v'", err)
		return result, err
	}

	result, err = executor.Execute(false, false)
	if err != nil {
		zap.S().Error("failed to scan resource object. error: '%v'", err)
		return result, err
	}

	return result, nil
}

func (w ValidatingWebhook) getDenyViolations(output runtime.Output) ([]results.Violation, error) {

	// Calcualte the deny violations according to the configuration specified in the config file
	denyViolations := w.getDeniedViolations(*output.Violations.ViolationStore, config.GetK8sAdmissionControl())

	return denyViolations, nil
}

func (w ValidatingWebhook) getDeniedViolations(violations results.ViolationStore, denyRules config.K8sAdmissionControl) []results.Violation {
	// Check whether one of the violations matches the deny violations configuration

	var denyViolations []results.Violation

	denyRuleMatcher := WebhookDenyRuleMatcher{}

	for _, violation := range violations.Violations {
		if denyRuleMatcher.Match(*violation, denyRules) {
			denyViolations = append(denyViolations, *violation)
		}
	}

	return denyViolations
}

func (w ValidatingWebhook) logWebhook(output runtime.Output,
	uid string,
	denyViolations []results.Violation,
	allowed bool) error {

	var (
		currentTime             = time.Now()
		deniedViolationsEncoded string
		requestBody             string
	)

	// encode denied violations into a string
	if len(denyViolations) < 1 {
		deniedViolationsEncoded = ""
	} else {
		d, _ := json.Marshal(denyViolations)
		deniedViolationsEncoded = string(d)
	}

	encodedViolationsSummary, _ := json.Marshal(output.Violations.ViolationStore)

	if config.GetK8sAdmissionControl().SaveRequests {
		requestBody = string(w.requestBody)
	}

	// insert the webhook log into db
	err := w.dblogger.Log(dblogs.WebhookScanLog{
		UID:                uid,
		Request:            requestBody,
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

// createAdmissionResponse creates a admission review response which is sent
// to calling kubernetes API server
func (w ValidatingWebhook) createResponseAdmissionReview(
	requestedAdmissionReview admissionv1.AdmissionReview,
	allowed bool,
	output runtime.Output,
	logMsg string) *admissionv1.AdmissionReview {

	// for dashboard mode, we display user a log endpoint where user can access the
	// admission request(provided requests are logged) and violation details
	errMsgs := []string{fmt.Sprintf("For more details please visit %q", logMsg)}

	// create an admission review request to be sent as response
	responseAdmissionReview := &admissionv1.AdmissionReview{}
	responseAdmissionReview.SetGroupVersionKind(requestedAdmissionReview.GroupVersionKind())

	// populate admission response
	responseAdmissionReview.Response = &admissionv1.AdmissionResponse{
		UID:     requestedAdmissionReview.Request.UID,
		Allowed: allowed,
	}

	if output.Violations.ViolationStore != nil {
		if !config.GetK8sAdmissionControl().Dashboard {
			errMsgs = w.buildErrors(output.Violations.ViolationStore.Violations)
		}

		// Means we ran the engines and we have results
		if allowed {
			if len(output.Violations.ViolationStore.Violations) > 0 {
				// In case there are no denial violations, just return the log URL as a warning
				responseAdmissionReview.Response.Warnings = errMsgs
			}
		} else {
			// In case the request was denied, return 403 and the log URL as an error message
			responseAdmissionReview.Response.Result = &metav1.Status{Message: "\n" + strings.Join(errMsgs, "\n"), Code: 403}
		}
	}

	return responseAdmissionReview
}

// WebhookDenyRuleMatcher helps in matching violated rules with k8s denied admission control rules
type WebhookDenyRuleMatcher struct {
}

// Match should check if one of the violations found is relevant for the specified K8s deny rules
func (g *WebhookDenyRuleMatcher) Match(violation results.Violation, denyRules config.K8sAdmissionControl) bool {

	if denyRules.DeniedSeverity == "" && len(denyRules.Categories) == 0 {
		return false
	}

	// Currently we support:
	// 1. A minimum severity level
	// 2. A category list
	// In case one of the conditions is met, we return true. (We perform an OR between the rules)
	if len(denyRules.DeniedSeverity) > 0 && utils.CheckSeverity(violation.Severity, denyRules.DeniedSeverity) {
		return true
	}

	if denyRules.Categories != nil {
		for _, category := range denyRules.Categories {
			if category == violation.Category {
				return true
			}
		}
	}

	return false
}

// buildErrors build a list of error messages from all the violations
func (w ValidatingWebhook) buildErrors(violations []*results.Violation) []string {
	errMsgs := make([]string, 0)
	if len(violations) > 0 {
		for _, v := range violations {
			// make the 'file' and 'line number' field blank, because it is a temporary file and would confuse the user
			v.File = ""
			v.LineNumber = 0
			out, _ := json.Marshal(v)
			errMsgs = append(errMsgs, string(out))
		}
	}
	return errMsgs
}
