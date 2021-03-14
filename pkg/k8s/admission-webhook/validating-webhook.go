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

package admissionwebhook

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	runtimeK8s "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// ValidatingWebhook handles the incoming validating admission webhook from
// the kubernetes API server and decides whether the admission request from
// the kubernetes client should be allowed or not
type ValidatingWebhook struct {
	configFile string
}

// NewValidatingWebhook returns a new, empty ValidatingWebhook struct
func NewValidatingWebhook(configFile string) AdmissionWebhook {
	return ValidatingWebhook{configFile: configFile}
}

var (
	ErrAPIKeyMissing   = fmt.Errorf("apiKey is missing in validating admission webhook url")
	ErrAPIKeyEnvNotSet = fmt.Errorf("variable K8S_WEBHOOK_API_KEY not set in terrascan server environment")
	ErrUnAuthorized    = fmt.Errorf("invalid API key in validating admission webhook url")
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
		zap.S().Error(ErrUnAuthorized)
		return ErrUnAuthorized
	}

	return nil
}

// DecodeAdmissionReviewRequest reads the incoming admission request body,
// decodes it and returns an AdmissionReviewRequest struct
func (w ValidatingWebhook) DecodeAdmissionReviewRequest(payload []byte) (admissionv1.AdmissionReview, error) {

	var (
		scheme                   = runtimeK8s.NewScheme()
		codecs                   = serializer.NewCodecFactory(scheme)
		deserializer             = codecs.UniversalDeserializer()
		requestedAdmissionReview admissionv1.AdmissionReview
	)
	admissionv1.AddToScheme(scheme)

	// decode incoming admission request
	_, _, err := deserializer.Decode(payload, nil, &requestedAdmissionReview)
	if err != nil {
		errMsg := "failed to decode validating admission webhook payload"
		zap.S().Error(errMsg, zap.Error(err))
		return requestedAdmissionReview, fmt.Errorf("%s, error: %w", errMsg, err)
	}

	return requestedAdmissionReview, nil
}

// ProcessWebhook processes the incoming AdmissionReview and creates
// a response
func (w ValidatingWebhook) ProcessWebhook(review admissionv1.AdmissionReview) (output *runtime.Output, allowed bool, denyViolations []*results.Violation, err error) {

	// In case the object is nil => an operation of DELETE happened, just return 'allow' since there is nothing to check
	if len(review.Request.Object.Raw) < 1 {
		zap.S().Info("recieved empty validating admission review request", zap.Any("admission review object", review))
		return output, true, denyViolations, nil
	}

	// Save the object into a temp file for the policy engines
	tempFile, err := w.writeObjectToTempFile(review.Request.Object.Raw)
	defer os.Remove(tempFile.Name())
	if err != nil {
		msg := "failed to create temp file for validating admission review request"
		zap.S().Error(msg, zap.Error(err))
		return output, true, denyViolations, fmt.Errorf("%s; error: %w", msg, err)
	}

	// Run the policy engines
	output, err = w.executeEngines(*tempFile)
	if err != nil {
		msg := "failed to evaluate terrascan policies"
		zap.S().Errorf(msg, zap.Error(err))
		return output, allowed, denyViolations, fmt.Errorf("%s; error: %w", msg, err)
	}

	// Calculate if there are anydeny violations
	denyViolations, err = w.getDenyViolations(*output)
	allowed = len(denyViolations) < 1

	return output, allowed, denyViolations, nil
}

func (w ValidatingWebhook) executeEngines(tempFile os.File) (*runtime.Output, error) {
	var executor *runtime.Executor
	var err error

	executor, err = runtime.NewExecutor("k8s", "v1", []string{"k8s"},
		tempFile.Name(), "", w.configFile, []string{}, []string{}, []string{}, []string{}, "")

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

func (w ValidatingWebhook) getDenyViolations(output runtime.Output) ([]*results.Violation, error) {

	// Calcualte the deny violations according to the configuration specified in the config file
	configReader, err := config.NewTerrascanConfigReader(w.configFile)
	if err != nil {
		zap.S().Errorf("error loading config file: '%v'", err)
		return nil, err
	}

	denyViolations := w.getDeniedViolations(*output.Violations.ViolationStore, configReader.GetK8sDenyRules())

	return denyViolations, nil
}

func (w ValidatingWebhook) getDeniedViolations(violations results.ViolationStore, denyRules config.K8sDenyRules) []*results.Violation {
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

type webhookDenyRuleMatcher struct {
}

// This class should check if one of the violations found is relevant for the specified K8s deny rules
func (g *webhookDenyRuleMatcher) match(violation results.Violation, denyRules config.K8sDenyRules) bool {
	if &denyRules == nil {
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

func (w ValidatingWebhook) writeObjectToTempFile(objectBytes []byte) (*os.File, error) {
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
