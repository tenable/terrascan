package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpload(t *testing.T) {
	testFilePath := filepath.Join("testdata", "testconfig.tf")
	testIacType := "terraform"
	testIacVersion := "v14"
	testCloudType := "aws"
	testParamName := "file"
	testCategories := []string{"COMPLIANCE VALIDATION", "DATA PROTECTION"}

	table := []struct {
		name                       string
		path                       string
		param                      string
		iacType                    string
		iacVersion                 string
		cloudType                  string
		scanRules                  []string
		skipRules                  []string
		severity                   string
		configOnly                 bool
		invalidConfigOnly          bool
		showPassed                 bool
		invalidShowPassed          bool
		wantStatus                 int
		categories                 []string
		invalidFindVulnerabilities bool
		findVulnerabilities        bool
		notificationWebhookURL     string
		notificationWebhookToken   string
		configWithError            bool
		invalidConfigWithError     bool
	}{
		{
			name:       "valid file scan",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan default iac type",
			path:       testFilePath,
			param:      testParamName,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan default iac version",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan default iac version, with invalid severity level input",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			severity:   "HGIH",
			cloudType:  testCloudType,
			wantStatus: http.StatusBadRequest,
			categories: testCategories,
		},
		{
			name:       "invalid iacType",
			path:       testFilePath,
			param:      testParamName,
			iacType:    "notthere",
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusBadRequest,
			categories: testCategories,
		},
		{
			name:       "invalid file param",
			path:       testFilePath,
			param:      "someparam",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid file config",
			path:       filepath.Join("testdata", "invalid.tf"),
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusInternalServerError,
			categories: testCategories,
		},
		{
			name:       "empty file config",
			path:       filepath.Join("testdata", "empty.tf"),
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan with scan and skip rules",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			scanRules: []string{"AWS.CloudFront.EncryptionandKeyManagement.High.0407", "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
				"AWS.CloudFront.Logging.Medium.0567", "AWS.CloudFront.Network Security.Low.0568"},
			skipRules: []string{"AWS.CloudFront.Network Security.Low.0568"},
		},
		{
			name:       "valid file scan default iac version",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			severity:   "low ",
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan default iac version  with MEDIUM severity",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			severity:   " MEDIUM ",
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan default iac version with high severity",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			severity:   "high",
			wantStatus: http.StatusOK,
			categories: testCategories,
		},
		{
			name:       "valid file scan with scan and skip rules with low severity",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			severity:   "low ",
			scanRules: []string{"AWS.CloudFront.EncryptionandKeyManagement.High.0407", "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
				"AWS.CloudFront.Logging.Medium.0567", "AWS.CloudFront.Network Security.Low.0568"},
			skipRules:  []string{"AWS.CloudFront.Network Security.Low.0568"},
			categories: testCategories,
		},
		{
			name:       "valid file scan with scan and skip rules with medium severity",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			severity:   " medium",
			scanRules: []string{"AWS.CloudFront.EncryptionandKeyManagement.High.0407", "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
				"AWS.CloudFront.Logging.Medium.0567", "AWS.CloudFront.Network Security.Low.0568"},
			skipRules:  []string{"AWS.CloudFront.Network Security.Low.0568"},
			categories: testCategories,
		},
		{
			name:       "valid file scan with scan and skip rules with HIGH severity",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			severity:   "HIGH",
			scanRules: []string{"AWS.CloudFront.EncryptionandKeyManagement.High.0407", "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
				"AWS.CloudFront.Logging.Medium.0567", "AWS.CloudFront.Network Security.Low.0568"},
			skipRules: []string{"AWS.CloudFront.Network Security.Low.0568"},
		},
		{
			name:       "test for config only",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
			configOnly: true,
		},
		{
			name:              "test for invalid value config only",
			path:              testFilePath,
			param:             testParamName,
			iacType:           testIacType,
			cloudType:         testCloudType,
			wantStatus:        http.StatusBadRequest,
			invalidConfigOnly: true,
		},
		{
			name:       "test for show passed attribute",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			showPassed: true,
			wantStatus: http.StatusOK,
		},
		{
			name:              "test for invalid show_passed value",
			path:              testFilePath,
			param:             testParamName,
			iacType:           testIacType,
			cloudType:         testCloudType,
			invalidShowPassed: true,
			wantStatus:        http.StatusBadRequest,
		},
		{
			name:       "scan valid kubernetes yaml",
			path:       filepath.Join("..", "iac-providers", "kubernetes", "v1", "testdata", "yaml-extension2", "test_pod.yml"),
			param:      testParamName,
			iacType:    "k8s",
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
		},
		{
			name:       "scan valid tfplan json",
			path:       filepath.Join("..", "iac-providers", "tfplan", "v1", "testdata", "valid-tfplan.json"),
			param:      testParamName,
			iacType:    "tfplan",
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
		},
		{
			name:                       "test for invalid value find vulnerability",
			path:                       testFilePath,
			param:                      testParamName,
			iacType:                    testIacType,
			cloudType:                  testCloudType,
			wantStatus:                 http.StatusBadRequest,
			invalidFindVulnerabilities: true,
		},
		{
			name:                     "valid file scan with notification webhook",
			path:                     testFilePath,
			param:                    testParamName,
			iacType:                  testIacType,
			iacVersion:               testIacVersion,
			cloudType:                testCloudType,
			wantStatus:               http.StatusOK,
			notificationWebhookURL:   "https://httpbin.org/post",
			notificationWebhookToken: "token",
		},
		{
			name:                   "test for config with error invalid",
			path:                   testFilePath,
			param:                  testParamName,
			iacType:                testIacType,
			cloudType:              testCloudType,
			wantStatus:             http.StatusBadRequest,
			invalidConfigWithError: true,
		},
		{
			name:                   "test for config with error",
			path:                   testFilePath,
			param:                  testParamName,
			iacType:                testIacType,
			cloudType:              testCloudType,
			wantStatus:             http.StatusOK,
			invalidConfigWithError: false,
			configWithError:        true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			// test file to upload
			path := tt.path
			file, err := os.Open(path)
			if err != nil {
				t.Error(err)
			}
			defer file.Close()
			// use buffer to store response body
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile(tt.param, filepath.Base(path))
			if err != nil {
				writer.Close()
				t.Error(err)
			}
			io.Copy(part, file)

			if len(tt.scanRules) > 0 {
				if err = writer.WriteField("scan_rules", strings.Join(tt.scanRules, ",")); err != nil {
					writer.Close()
					t.Error(err)
				}
			}
			if len(tt.skipRules) > 0 {
				if err = writer.WriteField("skip_rules", strings.Join(tt.scanRules, ",")); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if len(tt.severity) > 0 {
				if err = writer.WriteField("severity", tt.severity); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if len(tt.categories) > 0 {
				if err = writer.WriteField("categories", strings.Join(tt.categories, ",")); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if !tt.invalidConfigOnly {
				if err = writer.WriteField("config_only", strconv.FormatBool(tt.configOnly)); err != nil {
					writer.Close()
					t.Error(err)
				}
			} else {
				if err = writer.WriteField("config_only", "invalid"); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if !tt.invalidShowPassed {
				if err = writer.WriteField("show_passed", strconv.FormatBool(tt.showPassed)); err != nil {
					writer.Close()
					t.Error(err)
				}
			} else {
				if err = writer.WriteField("show_passed", "invalid"); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if len(tt.categories) > 0 {
				if err = writer.WriteField("categories", strings.Join(tt.categories, ",")); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if !tt.invalidFindVulnerabilities {
				if err = writer.WriteField("find_vulnerabilities", strconv.FormatBool(tt.findVulnerabilities)); err != nil {
					writer.Close()
					t.Error(err)
				}
			} else {
				if err = writer.WriteField("find_vulnerabilities", "invalid"); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if tt.notificationWebhookURL != "" {
				if err = writer.WriteField("webhook_url", tt.notificationWebhookURL); err != nil {
					writer.Close()
					t.Error(err)
				}
			}
			if tt.notificationWebhookToken != "" {
				if err = writer.WriteField("webhook_token", tt.notificationWebhookToken); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			if !tt.invalidConfigWithError {
				if err = writer.WriteField("config_with_error", strconv.FormatBool(tt.configWithError)); err != nil {
					writer.Close()
					t.Error(err)
				}
			} else {
				if err = writer.WriteField("config_with_error", "invalid"); err != nil {
					writer.Close()
					t.Error(err)
				}
			}

			writer.Close()

			// http request of the type "/v1/{iacType}/{iacVersion}/{cloudType}/file/scan"
			url := fmt.Sprintf("/v1/%s/%s/%s/local/file/scan", tt.iacType, tt.iacVersion, tt.cloudType)
			req := httptest.NewRequest("POST", url, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req = mux.SetURLVars(req, map[string]string{
				"iac":        tt.iacType,
				"iacVersion": tt.iacVersion,
				"cloud":      tt.cloudType,
			})
			res := httptest.NewRecorder()
			// new api handler
			h := &APIHandler{test: true}
			h.scanFile(res, req)

			if res.Code != tt.wantStatus {
				t.Errorf("incorrect status code, got: '%v', want: '%v', error: '%v'", res.Code, tt.wantStatus, res.Body)
			}
		})
	}
}
