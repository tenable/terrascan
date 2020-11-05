package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/gorilla/mux"
)

var (
	someURL  = "some-url"
	someType = "some-type"
)

func TestScanRemoteRepo(t *testing.T) {

	var (
		d        = downloader.NewDownloader()
		noOutput interface{}
		// noErr    error = nil
	)

	table := []struct {
		name       string
		iacType    string
		iacVersion string
		cloudType  []string
		s          *scanRemoteRepoReq
		wantOutput interface{}
		wantErr    error
	}{
		{
			name: "remote url empty",
			s: &scanRemoteRepoReq{
				RemoteURL:  "",
				RemoteType: someType,
				d:          d,
			},
			wantOutput: noOutput,
			wantErr:    downloader.ErrEmptyURLDest,
		},
		{
			name: "remote type empty",
			s: &scanRemoteRepoReq{
				RemoteURL:  someURL,
				RemoteType: "",
				d:          d,
			},
			wantOutput: noOutput,
			wantErr:    downloader.ErrEmptyURLDest,
		},
		{
			name: "remote type and url empty",
			s: &scanRemoteRepoReq{
				RemoteURL:  "",
				RemoteType: "",
				d:          d,
			},
			wantOutput: noOutput,
			wantErr:    downloader.ErrEmptyURLType,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotErr := tt.s.ScanRemoteRepo(tt.iacType, tt.iacVersion, tt.cloudType, []string{})
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("output got: '%v', want: '%v'", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestScanRemoteRepoHandler(t *testing.T) {

	table := []struct {
		name       string
		iacType    string
		iacVersion string
		cloudType  string
		remoteURL  string
		remoteType string
		wantStatus int
	}{
		{
			name:       "empty url and type",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			remoteURL:  "",
			remoteType: "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty type",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			remoteURL:  someURL,
			remoteType: "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url and type",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			remoteURL:  someURL,
			remoteType: someType,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "valid url and type",
			iacType:    "terraform",
			iacVersion: "v12",
			cloudType:  "aws",
			remoteURL:  "https://github.com/kanchwala-yusuf/Damn-Vulnerable-Terraform-Project.git",
			remoteType: "git",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			// http request of the type "/v1/{iacType}/{iacVersion}/{cloudType}/remote/dir/scan"

			// request url
			url := fmt.Sprintf("/v1/%s/%s/%s/remote/dir/scan", tt.iacType, tt.iacVersion, tt.cloudType)

			// request body
			s := scanRemoteRepoReq{
				RemoteURL:  tt.remoteURL,
				RemoteType: tt.remoteType,
			}
			reqBody, _ := json.Marshal(s)

			// http request
			req := httptest.NewRequest("POST", url, bytes.NewBuffer(reqBody))

			// set headers
			req.Header.Set("Content-Type", "application/json")

			// set URL params
			req = mux.SetURLVars(req, map[string]string{
				"iac":        tt.iacType,
				"iacVersion": tt.iacVersion,
				"cloud":      tt.cloudType,
			})
			res := httptest.NewRecorder()

			// new api handler
			h := &APIHandler{test: true}
			h.scanRemoteRepo(res, req)

			if res.Code != tt.wantStatus {
				t.Errorf("incorrect status code, got: '%v', want: '%v', error: '%v'", res.Code, http.StatusOK, res.Body)
			}
		})
	}
}
