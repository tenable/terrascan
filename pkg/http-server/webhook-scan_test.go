package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/k8s/dblogs"
	v1 "k8s.io/api/admission/v1"
)

func TestUWebhooks(t *testing.T) {
	k8sTestData := "k8s_testdata"
	testFilePath := filepath.Join(k8sTestData, "testconfig.json")
	testAPIKey := "Test-API-KEY"
	testEnvAPIKey := "Test-API-KEY"
	testConfigFile := ""

	table := []struct {
		name               string
		contentRequestPath string
		apiKey             string
		envAPIKey          string
		wantStatus         int
		configFile         string
		warnings           bool
		allowed            bool
		statusCode         int32
		statusMessage      bool
		dashboardCheck     bool
	}{
		{
			name:               "missing api key",
			contentRequestPath: testFilePath,
			apiKey:             "",
			envAPIKey:          testEnvAPIKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "missing K8S_WEBHOOK_API_KEY",
			contentRequestPath: testFilePath,
			apiKey:             testAPIKey,
			envAPIKey:          "",
			wantStatus:         http.StatusInternalServerError,
			configFile:         testConfigFile,
		},
		{
			name:               "invalid api key",
			contentRequestPath: testFilePath,
			apiKey:             testAPIKey,
			envAPIKey:          "Invalid API KEY",
			wantStatus:         http.StatusUnauthorized,
			configFile:         testConfigFile,
		},
		{
			name:               "invalid request json content",
			contentRequestPath: filepath.Join(k8sTestData, "invalid.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "empty request json content",
			contentRequestPath: filepath.Join(k8sTestData, "empty.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "request with empty object",
			contentRequestPath: filepath.Join(k8sTestData, "empty_object.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			wantStatus:         http.StatusOK,
			configFile:         testConfigFile,
			warnings:           false,
			allowed:            true,
		},
		{
			name:               "safe request object",
			contentRequestPath: testFilePath,
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			wantStatus:         http.StatusOK,
			configFile:         testConfigFile,
			warnings:           false,
			allowed:            true,
		},
		{
			name:               "risky request object without config",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         testConfigFile,
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with config that make it safe",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-specific-rule.toml"),
			warnings:           false,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with config that just removes some of the violations",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-medium-severity.toml"),
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied severity",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-deny-high.toml"),
			warnings:           false,
			allowed:            false,
			statusCode:         403,
			statusMessage:      true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied categories",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-deny-category.toml"),
			warnings:           false,
			allowed:            false,
			statusCode:         403,
			statusMessage:      true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied categories that does not exist",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-deny-non-existing-category.toml"),
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with dashboard true",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-dashboard-true.toml"),
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
			dashboardCheck:     true,
		},
		{
			name:               "risky request object with denied categories and dashboard true",
			contentRequestPath: filepath.Join(k8sTestData, "risky_testconfig.json"),
			apiKey:             testAPIKey,
			envAPIKey:          testEnvAPIKey,
			configFile:         filepath.Join(k8sTestData, "config-with-dashboard-deny-categories.toml"),
			warnings:           false,
			allowed:            false,
			statusCode:         403,
			statusMessage:      true,
			wantStatus:         http.StatusOK,
			dashboardCheck:     true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := config.LoadGlobalConfig(tt.configFile)
			if err != nil {
				t.Errorf("error while loading the config file '%s', error: %v", tt.configFile, err)
			}

			os.Setenv("K8S_WEBHOOK_API_KEY", tt.envAPIKey)

			// test file to upload
			path := tt.contentRequestPath
			jsonFile, err := os.Open(path)
			if err != nil {
				t.Error(err)
				return
			}
			defer jsonFile.Close()

			logger := dblogs.WebhookScanLogger{
				Test: true,
			}
			defer logger.ClearDbFilePath()

			byteValue, _ := io.ReadAll(jsonFile)

			var admissionRequest v1.AdmissionReview
			json.Unmarshal(byteValue, &admissionRequest)

			var url string
			if len(tt.apiKey) > 0 {
				url = fmt.Sprintf("/v1/k8s/webhooks/%v/scan", tt.apiKey)
			} else {
				url = "/v1/k8s/webhooks/scan"
			}

			req := httptest.NewRequest("POST", url, bytes.NewReader(byteValue))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{
				"apiKey": tt.apiKey,
			})
			res := httptest.NewRecorder()
			// new api handler
			h := &APIHandler{test: true}
			h.validateK8SWebhook(res, req)

			if res.Code != tt.wantStatus {
				t.Errorf("incorrect status code, got: '%v', want: '%v', error: '%v'", res.Code, tt.wantStatus, res.Body)
			}

			var response v1.AdmissionReview
			_ = json.Unmarshal(res.Body.Bytes(), &response)

			if res.Code == http.StatusOK {
				if tt.warnings && response.Response.Warnings == nil {
					t.Errorf("expected warnings but received None")
				}

				if tt.allowed != response.Response.Allowed {
					t.Errorf("mismatch in allowed. Got: %v, expected: %v", response.Response.Allowed, tt.allowed)
				}

				if tt.statusCode != 0 && tt.statusCode != response.Response.Result.Code {
					t.Errorf("mismatch Status code Got: %v, expected: %v", response.Response.Result.Code, tt.statusCode)
				}

				if tt.dashboardCheck {
					if tt.warnings || tt.statusMessage {
						var logPath string
						if tt.warnings {
							logPath = strings.TrimSpace(response.Response.Warnings[0])
						} else if tt.statusMessage {
							logPath = strings.TrimSpace(response.Response.Result.Message)
						}

						subLogPath := fmt.Sprintf("https://%v/k8s/webhooks/logs/705ab4f5-6393-11e8-b7cc-42010a800002", req.Host)
						expectedLogPath := fmt.Sprintf("For more details please visit %q", subLogPath)

						if logPath != expectedLogPath {
							t.Errorf("mismatch Log path. Got: %v, expected: %v", logPath, expectedLogPath)
						}
					}
				}
			}
		})
	}
}
