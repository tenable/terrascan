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
	"github.com/accurics/terrascan/pkg/results"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"time"
)

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
	LogUrl    string
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
	// Return an HTML page including all the logs history

	params := mux.Vars(r)

	apiKey := params["apiKey"]
	if !g.validateAuthorization(apiKey, w) {
		return
	}

	logger := WebhookScanLogger{
		test: g.test,
	}

	// The templates are saved in the docker in this location
	t, _ := template.ParseFiles("/go/terrascan/index.html")

	logs, err := logger.fetchLogs()
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
			LogUrl:    g.getLogPath(r.Host, apiKey, log.UID),
			Reasoning: g.getLogReasoning(log),
			Request:   g.getLogRequest(log),
		})
	}

	t.Execute(w, logsData)
}

func (g *APIHandler) getLogByUID(w http.ResponseWriter, r *http.Request) {
	// Return an HTML page including the selected log

	params := mux.Vars(r)

	if !g.validateAuthorization(params["apiKey"], w) {
		return
	}

	var uid = params["uid"]
	if len(uid) < 1 {
		apiErrorResponse(w, "Log UID is missing", http.StatusBadRequest)
		return
	}

	logger := WebhookScanLogger{
		test: g.test,
	}

	log, err := logger.fetchLogById(uid)
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

func (g *APIHandler) getLogPath(host string, apiKey string, logUID string) string {
	// Use this as the link to show the a specific log
	return fmt.Sprintf("https://%v/k8s/webhooks/%v/logs/%v", host, apiKey, logUID)
}

func (g *APIHandler) getLogStatus(log webhookScanLog) string {
	// Calculate a log status:
	// 1. !Allowed -> Rejected
	// 2. Allowed -> if there are violations -> Allowed with Warnings. Otherwise -> Allowed
	if !log.Allowed {
		return "Rejected"
	}

	var violationStore results.ViolationStore
	err := json.Unmarshal([]byte(log.ViolationsSummary), &violationStore)
	if err != nil {
		zap.S().Errorf("Failed to ..")
	}

	if len(violationStore.Violations) > 0 {
		return "Allowed with warnings"
	}

	return "Allowed"
}

func (g *APIHandler) getLogReasoning(log webhookScanLog) string {
	// Reasoning:
	// - In case the request is denied (rejected), show the violations that cause the denial.
	// - Otherwise, if there are violations, show the full violations list was found
	// - Otherwise, reasoning is empty

	var violations []*results.Violation
	if !log.Allowed {
		err := json.Unmarshal([]byte(log.DeniableViolations), &violations)
		if err != nil {
			zap.S().Errorf("Failed to deserialize deniable violations summary. Error: %v", err.Error())
			return ""
		}
	} else {
		var violationStore results.ViolationStore
		err := json.Unmarshal([]byte(log.ViolationsSummary), &violationStore)
		if err != nil {
			zap.S().Errorf("Failed to deserialize violations summary. Error: %v", err.Error())
			return ""
		}

		violations = violationStore.Violations
	}

	var result []webhookDisplayedViolation

	if len(violations) < 1 {
		return ""
	} else {
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

	return ""
}

func (g *APIHandler) getLogRequest(log webhookScanLog) string {
	var review webhookDisplayedReview

	err := json.Unmarshal([]byte(log.Request), &review)

	if err != nil {
		zap.S().Errorf("Failed to deserialize request. Error: %v", err.Error())
		return "{}"
	}

	result, err := json.Marshal(review.Request)
	if err != nil {
		zap.S().Errorf("Failed to serialize request. Error: %v", err.Error())
		return "{}"
	}

	return string(result)
}
