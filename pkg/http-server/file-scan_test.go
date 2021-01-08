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
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpload(t *testing.T) {
	testFilePath := "./testdata/testconfig.tf"
	testIacType := "terraform"
	testIacVersion := "v12"
	testCloudType := "aws"
	testParamName := "file"

	table := []struct {
		name       string
		path       string
		param      string
		iacType    string
		iacVersion string
		cloudType  string
		scanRules  []string
		skipRules  []string
		wantStatus int
	}{
		{
			name:       "valid file scan",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid file scan default iac type",
			path:       testFilePath,
			param:      testParamName,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid file scan default iac version",
			path:       testFilePath,
			param:      testParamName,
			iacType:    testIacType,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid iacType",
			path:       testFilePath,
			param:      testParamName,
			iacType:    "notthere",
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid file param",
			path:       testFilePath,
			param:      "someparam",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid file config",
			path:       "./testdata/invalid.tf",
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "empty file config",
			path:       "./testdata/empty.tf",
			param:      testParamName,
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			wantStatus: http.StatusOK,
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
				t.Errorf("incorrect status code, got: '%v', want: '%v', error: '%v'", res.Code, http.StatusOK, res.Body)
			}
		})
	}
}
