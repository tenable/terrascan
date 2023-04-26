package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/downloader"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/runtime"
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
			gotOutput, _, gotErr := tt.s.ScanRemoteRepo(tt.iacType, tt.iacVersion, tt.cloudType, []string{})
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
	validRepo := "https://github.com/kanchwala-yusuf/Damn-Vulnerable-Terraform-Project.git"
	testIacType := "terraform"
	testIacVersion := "v14"
	testCloudType := "aws"

	table := []struct {
		name            string
		iacType         string
		iacVersion      string
		cloudType       string
		remoteURL       string
		remoteType      string
		scanRules       []string
		skipRules       []string
		showPassed      bool
		configOnly      bool
		configWithError bool
		nonRecursive    bool
		wantStatus      int
	}{
		{
			name:       "empty url and type",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  "",
			remoteType: "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty type",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  someURL,
			remoteType: "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url and type",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  someURL,
			remoteType: someType,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "valid url and type",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  validRepo,
			remoteType: "git",
			wantStatus: http.StatusOK,
		},
		{
			name:         "iac type terraform with non-recursive scan",
			iacType:      testIacType,
			iacVersion:   testIacVersion,
			cloudType:    testCloudType,
			remoteURL:    validRepo,
			remoteType:   "git",
			nonRecursive: true,
			wantStatus:   http.StatusOK,
		},
		{
			name:       "valid url and type with scan and skip rules",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  validRepo,
			remoteType: "git",
			scanRules: []string{"AWS.CloudFront.EncryptionandKeyManagement.High.0407", "AWS.CloudFront.EncryptionandKeyManagement.High.0408",
				"AWS.CloudFront.Logging.Medium.0567", "AWS.CloudFront.Network Security.Low.0568"},
			skipRules:  []string{"AWS.CloudFront.Network Security.Low.0568"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "test show passed rules and config only",
			iacType:    testIacType,
			iacVersion: testIacVersion,
			cloudType:  testCloudType,
			remoteURL:  validRepo,
			remoteType: "git",
			showPassed: true,
			configOnly: true,
			wantStatus: http.StatusOK,
		},
		{
			name:            "test show config with error",
			iacType:         testIacType,
			iacVersion:      testIacVersion,
			cloudType:       testCloudType,
			remoteURL:       validRepo,
			remoteType:      "git",
			showPassed:      false,
			configWithError: true,
			wantStatus:      http.StatusOK,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			// http request of the type "/v1/{iacType}/{iacVersion}/{cloudType}/remote/dir/scan"

			// request url
			url := fmt.Sprintf("/v1/%s/%s/%s/remote/dir/scan", tt.iacType, tt.iacVersion, tt.cloudType)

			// request body
			s := scanRemoteRepoReq{
				RemoteURL:       tt.remoteURL,
				RemoteType:      tt.remoteType,
				ScanRules:       tt.scanRules,
				SkipRules:       tt.skipRules,
				ShowPassed:      tt.showPassed,
				ConfigOnly:      tt.configOnly,
				ConfigWithError: tt.configWithError,
				NonRecursive:    tt.nonRecursive,
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

func TestHasK8sAdmissionDeniedViolations(t *testing.T) {
	k8sTestData := "k8s_testdata"
	configFileWithCategoryDenied := filepath.Join(k8sTestData, "config-deny-category.toml")

	type args struct {
		o runtime.Output
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		configFile string
	}{
		{
			name: "result with no violations",
			args: args{
				o: runtime.Output{
					Violations: policy.EngineOutput{
						ViolationStore: &results.ViolationStore{},
					},
				},
			},
			want: false,
		},
		{
			name: "result contains denied violations",
			args: args{
				o: runtime.Output{
					Violations: policy.EngineOutput{
						ViolationStore: &results.ViolationStore{
							Violations: []*results.Violation{
								{
									Category: "Identity and Access Management",
								},
							},
						},
					},
				},
			},
			want:       true,
			configFile: configFileWithCategoryDenied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.LoadGlobalConfig(tt.configFile); err != nil {
				t.Errorf("error while loading the config file '%s'", tt.configFile)
			}
			if got := hasK8sAdmissionDeniedViolations(tt.args.o); got != tt.want {
				t.Errorf("hasK8sAdmissionDeniedViolations() = %v, want %v", got, tt.want)
			}
		})
	}
}
