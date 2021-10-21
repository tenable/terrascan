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
	"html/template"
	"net/http"
	"time"

	"github.com/accurics/terrascan/pkg/config"
	admissionWebhook "github.com/accurics/terrascan/pkg/k8s/admission-webhook"
	"github.com/accurics/terrascan/pkg/k8s/dblogs"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// ErrDashboardDisabled would be the error returned back when log endpoint
// is hit while the dashboard mode is disabled
var ErrDashboardDisabled = fmt.Errorf("set 'dashboard=true' in terrascan config file to enable database logs")

type webhookDisplayedViolation struct {
	RuleName    string `json:"rule_name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type webhookDisplayedReview struct {
	Request webhookDisplayedRequest `json:"request"`
}

type webhookDisplayedRequest struct {
	Operation string                 `json:"operation"`
	Object    map[string]interface{} `json:"object"`
}

type webhookDisplayedIndexScanLog struct {
	CreatedAt time.Time
	LogURL    string
	Status    string
	Request   string
	Reasoning string
}

type webhookDisplayedShowLog struct {
	CreatedAt          time.Time
	UID                string
	Status             string
	Request            string
	Violations         string
	DeniableViolations string
}

func (g *APIHandler) getLogs(w http.ResponseWriter, r *http.Request) {
	zap.S().Debug("handle: validating webhook's get logs request")

	if !config.GetK8sAdmissionControl().Dashboard {
		apiErrorResponse(w, ErrDashboardDisabled.Error(), http.StatusBadRequest)
		return
	}

	var (
		params = mux.Vars(r)
		apiKey = params["apiKey"]
	)

	// Validate if authorized (API key is specified and matched the server one (saved in an environment variable)
	validatingWebhook := admissionWebhook.NewValidatingWebhook([]byte(""), "", "", "", "")
	if err := validatingWebhook.Authorize(apiKey); err != nil {
		switch err {
		case admissionWebhook.ErrAPIKeyMissing:
			apiErrorResponse(w, err.Error(), http.StatusBadRequest)
		case admissionWebhook.ErrUnauthorized:
			apiErrorResponse(w, err.Error(), http.StatusUnauthorized)
		default:
			apiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return an HTML page including all the logs history
	logger := dblogs.NewWebhookScanLogger()

	// The templates are saved in the docker in this location
	t, err := template.ParseFiles("/go/terrascan/index.html")
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse index.html file; error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	logs, err := logger.FetchLogs()
	if err != nil {
		errMsg := fmt.Sprintf("error reading logs from DB: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	var logsData []webhookDisplayedIndexScanLog
	for _, log := range logs {
		logsData = append(logsData, webhookDisplayedIndexScanLog{
			CreatedAt: log.CreatedAt,
			Status:    g.getLogStatus(log),
			LogURL:    logger.GetLogURL(r.Host, log.UID),
			Reasoning: g.getLogReasoning(log),
			Request:   g.getLogRequest(log),
		})
	}

	if err := t.Execute(w, logsData); err != nil {
		errMsg := fmt.Sprintf("failed to execute html template; error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (g *APIHandler) getLogByUID(w http.ResponseWriter, r *http.Request) {
	zap.S().Debug("handle: validating webhook's get log by uid request")

	if !config.GetK8sAdmissionControl().Dashboard {
		apiErrorResponse(w, ErrDashboardDisabled.Error(), http.StatusBadRequest)
		return
	}

	// Return an HTML page including the selected log
	var (
		params = mux.Vars(r)
		uid    = params["uid"]
		logger = dblogs.NewWebhookScanLogger()
	)

	if len(uid) < 1 {
		apiErrorResponse(w, "Log UID is missing", http.StatusBadRequest)
		return
	}

	log, err := logger.FetchLogByID(uid)
	if err != nil {
		errMsg := fmt.Sprintf("error reading logs from DB: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	if len(log.UID) < 1 {
		apiErrorResponse(w, "Log is not found", http.StatusNotFound)
		return
	}

	displayedScanLog := webhookDisplayedShowLog{
		UID:                log.UID,
		CreatedAt:          log.CreatedAt,
		Status:             g.getLogStatus(*log),
		Request:            log.Request,
		Violations:         log.ViolationsSummary,
		DeniableViolations: log.DeniableViolations,
	}

	t, _ := template.ParseFiles("/go/terrascan/show.html")

	t.Execute(w, displayedScanLog)
}

func (g *APIHandler) getLogStatus(log dblogs.WebhookScanLog) string {
	// Calculate a log status:
	// 1. !Allowed -> Rejected
	// 2. Allowed -> if there are violations -> Allowed with Warnings. Otherwise -> Allowed
	if !log.Allowed {
		return "Rejected"
	}

	var violationStore results.ViolationStore
	err := json.Unmarshal([]byte(log.ViolationsSummary), &violationStore)
	if err != nil {
		zap.S().Errorf("failed to decode violation results", zap.Error(err))
	}

	if len(violationStore.Violations) > 0 {
		return "Allowed with warnings"
	}

	return "Allowed"
}

func (g *APIHandler) getLogReasoning(log dblogs.WebhookScanLog) string {
	// Reasoning:
	// - In case the request is denied (rejected), show the violations that cause the denial.
	// - Otherwise, if there are violations, show the full violations list was found
	// - Otherwise, reasoning is empty

	var violations []*results.Violation
	if !log.Allowed {
		err := json.Unmarshal([]byte(log.DeniableViolations), &violations)
		if err != nil {
			zap.S().Errorf("failed to deserialize deniable violations summary. Error: %v", err.Error())
			return ""
		}
	} else {
		var violationStore results.ViolationStore
		err := json.Unmarshal([]byte(log.ViolationsSummary), &violationStore)
		if err != nil {
			zap.S().Errorf("failed to deserialize violations summary. Error: %v", err.Error())
			return ""
		}

		violations = violationStore.Violations
	}

	var result []webhookDisplayedViolation

	if len(violations) < 1 {
		return ""
	}
	for _, v := range violations {
		result = append(result, webhookDisplayedViolation{
			Category:    v.Category,
			Description: v.Description,
			RuleName:    v.RuleName,
			Severity:    v.Severity,
		})
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		zap.S().Errorf("failed to serialize violations: '%v'", err)
		return ""
	}

	return string(encoded)
}

func (g *APIHandler) getLogRequest(log dblogs.WebhookScanLog) string {
	var review webhookDisplayedReview

	if !config.GetK8sAdmissionControl().SaveRequests {
		return "{}"
	}

	err := json.Unmarshal([]byte(log.Request), &review)

	if err != nil {
		zap.S().Errorf("failed to deserialize request. Error: %v", err.Error())
		return "{}"
	}

	result, err := json.Marshal(review.Request)
	if err != nil {
		zap.S().Errorf("failed to serialize request. Error: %v", err.Error())
		return "{}"
	}

	return string(result)
}
