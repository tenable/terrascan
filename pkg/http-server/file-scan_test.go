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
	"testing"

	"github.com/gorilla/mux"
)

func TestUpload(t *testing.T) {

	table := []struct {
		name       string
		path       string
		param      string
		iacType    string
		iacVersion string
		cloudType  string
		wantStatus int
	}{
		{
			name:       "valid file scan",
			path:       "./testdata/testconfig.tf",
			param:      "file",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid file scan default iac type",
			path:       "./testdata/testconfig.tf",
			param:      "file",
			cloudType:  "aws",
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid file scan default iac version",
			path:       "./testdata/testconfig.tf",
			param:      "file",
			iacType:    "terraform",
			cloudType:  "aws",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid iacType",
			path:       "./testdata/testconfig.tf",
			param:      "file",
			iacType:    "notthere",
			iacVersion: "v12",
			cloudType:  "aws",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid file param",
			path:       "./testdata/testconfig.tf",
			param:      "someparam",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid file config",
			path:       "./testdata/invalid.tf",
			param:      "file",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "empty file config",
			path:       "./testdata/empty.tf",
			param:      "file",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			wantStatus: http.StatusOK,
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
