package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	v1 "k8s.io/api/admission/v1"
)

func TestUWebhooks(t *testing.T) {
	testFilePath := "./k8s_testdata/testconfig.json"
	testApiKey := "Test-API-KEY"
	testEnvApiKey := "Test-API-KEY"
	testConfigFile := ""

	table := []struct {
		name               string
		contentRequestPath string
		apiKey             string
		envApiKey          string
		wantStatus         int
		configFile         string
		warnings           bool
		allowed            bool
		statusCode         int32
		statusMessage      bool
	}{
		{
			name:               "missing api key",
			contentRequestPath: testFilePath,
			apiKey:             "",
			envApiKey:          testEnvApiKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "missing K8S_WEBHOOK_API_KEY",
			contentRequestPath: testFilePath,
			apiKey:             testApiKey,
			envApiKey:          "",
			wantStatus:         http.StatusInternalServerError,
			configFile:         testConfigFile,
		},
		{
			name:               "invalid api key",
			contentRequestPath: testFilePath,
			apiKey:             testApiKey,
			envApiKey:          "Invalid API KEY",
			wantStatus:         http.StatusUnauthorized,
			configFile:         testConfigFile,
		},
		{
			name:               "invalid api key",
			contentRequestPath: testFilePath,
			apiKey:             testApiKey,
			envApiKey:          "Invalid API KEY",
			wantStatus:         http.StatusUnauthorized,
			configFile:         testConfigFile,
		},
		{
			name:               "invalid request json content",
			contentRequestPath: "./k8s_testdata/invalid.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "empty request json content",
			contentRequestPath: "./k8s_testdata/empty.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			wantStatus:         http.StatusBadRequest,
			configFile:         testConfigFile,
		},
		{
			name:               "request with empty object",
			contentRequestPath: "./k8s_testdata/empty_object.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			wantStatus:         http.StatusOK,
			configFile:         testConfigFile,
			warnings:           false,
			allowed:            true,
		},
		{
			name:               "safe request object",
			contentRequestPath: testFilePath,
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			wantStatus:         http.StatusOK,
			configFile:         testConfigFile,
			warnings:           false,
			allowed:            true,
		},
		{
			name:               "risky request object without config",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         testConfigFile,
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with config that make it safe",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         "./k8s_testdata/config-specific-rule.toml",
			warnings:           false,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with config that just removes some of the violations",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         "./k8s_testdata/config-medium-severity.toml",
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied severity",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         "./k8s_testdata/config-deny-high.toml",
			warnings:           false,
			allowed:            false,
			statusCode:         403,
			statusMessage:      true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied categories",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         "./k8s_testdata/config-deny-category.toml",
			warnings:           false,
			allowed:            false,
			statusCode:         403,
			statusMessage:      true,
			wantStatus:         http.StatusOK,
		},
		{
			name:               "risky request object with denied categories that does not exist",
			contentRequestPath: "./k8s_testdata/risky_testconfig.json",
			apiKey:             testApiKey,
			envApiKey:          testEnvApiKey,
			configFile:         "./k8s_testdata/config-deny-non-existing-category.toml",
			warnings:           true,
			allowed:            true,
			wantStatus:         http.StatusOK,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("K8S_WEBHOOK_API_KEY", tt.envApiKey)

			// test file to upload
			path := tt.contentRequestPath
			jsonFile, err := os.Open(path)
			if err != nil {
				t.Error(err)
				return
			}
			defer jsonFile.Close()
			logger := WebhookScanLogger{
				test: true,
			}
			defer logger.clearDbFilePath()

			byteValue, _ := ioutil.ReadAll(jsonFile)

			var admissionRequest v1.AdmissionReview
			json.Unmarshal(byteValue, &admissionRequest)

			var url string
			if len(tt.apiKey) > 0 {
				url = fmt.Sprintf("/v1/k8s/webhooks/%v/scan", tt.apiKey)
			} else {
				url = fmt.Sprintf("/v1/k8s/webhooks/scan")
			}

			req := httptest.NewRequest("POST", url, bytes.NewReader(byteValue))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{
				"apiKey": tt.apiKey,
			})
			res := httptest.NewRecorder()
			// new api handler
			h := &APIHandler{test: true, configFile: tt.configFile}
			h.validateK8SWebhook(res, req)

			if res.Code != tt.wantStatus {
				t.Errorf("incorrect status code, got: '%v', want: '%v', error: '%v'", res.Code, tt.wantStatus, res.Body)
			}

			var response v1.AdmissionReview
			_ = json.Unmarshal(res.Body.Bytes(), &response)

			if res.Code == http.StatusOK {
				if tt.warnings && response.Response.Warnings == nil {
					t.Errorf("Expected warnings but received None")
				}

				if tt.allowed != response.Response.Allowed {
					t.Errorf("Mismach in allowed. Got: %v, expected: %v", response.Response.Allowed, tt.allowed)
				}

				if tt.statusCode != 0 && tt.statusCode != response.Response.Result.Code {
					t.Errorf("Mismach Statud code Got: %v, expected: %v", response.Response.Result.Code, tt.statusCode)
				}

				if tt.warnings || tt.statusMessage {
					var logPath string
					if tt.warnings {
						logPath = response.Response.Warnings[0]
					} else if tt.statusMessage {
						logPath = response.Response.Result.Message
					}

					expectedLogPath := fmt.Sprintf("https://%v/k8s/webhooks/logs/705ab4f5-6393-11e8-b7cc-42010a800002", req.Host)

					if logPath != expectedLogPath {
						t.Errorf("Mismach Log path. Got: %v, expected: %v", logPath, expectedLogPath)
					}
				}
			}
		})
	}
}
