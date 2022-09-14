/*
    Copyright (C) 2022 Tenable, Inc.

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
	"os"
	"path/filepath"
	"testing"
)

func TestAuthorize(t *testing.T) {

	table := []struct {
		name         string
		apiKey       string
		envK8sAPIKey string
		wantErr      error
	}{
		{
			name:         "empty api key",
			apiKey:       "",
			envK8sAPIKey: "valid",
			wantErr:      ErrAPIKeyMissing,
		},
		{
			name:         "K8S_WEBHOOK_API_KEY env not set",
			apiKey:       "valid",
			envK8sAPIKey: "",
			wantErr:      ErrAPIKeyEnvNotSet,
		},
		{
			name:         "invalid api key",
			apiKey:       "invalid",
			envK8sAPIKey: "valid",
			wantErr:      ErrUnauthorized,
		},
		{
			name:         "valid api key",
			apiKey:       "valid",
			envK8sAPIKey: "valid",
			wantErr:      nil,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {

			// validating webhook object
			var w ValidatingWebhook

			// set K8S_WEBHOOK_API_KEY env if not empty
			if test.envK8sAPIKey != "" {
				os.Setenv("K8S_WEBHOOK_API_KEY", test.envK8sAPIKey)
			}
			defer os.Unsetenv("K8S_WEBHOOK_API_KEY")

			err := w.Authorize(test.apiKey)
			if err != test.wantErr {
				t.Errorf("unexpected error; got: '%v', want: '%v'", err, test.wantErr)
			}
		})
	}
}

func TestDecodeAdmissionReviewRequest(t *testing.T) {

	table := []struct {
		name        string
		requestFile string
		wantErr     bool
	}{
		{
			name:        "empty review request",
			requestFile: filepath.Join("testdata", "empty.json"),
			wantErr:     false,
		},
		{
			name:        "invalid review request",
			requestFile: filepath.Join("testdata", "invalid.json"),
			wantErr:     true,
		},
		{
			name:        "valid review request",
			requestFile: filepath.Join("testdata", "valid.json"),
			wantErr:     false,
		},
	}

	for _, test := range table {

		// read test request from file
		requestBody, err := os.ReadFile(test.requestFile)
		if err != nil {
			t.Errorf("failed to read test data, error: '%v'", err)
		}

		var w ValidatingWebhook
		_, err = w.DecodeAdmissionReviewRequest(requestBody)
		if (err == nil) == test.wantErr {
			t.Errorf("unexpected error '%v'", err)
		}
	}
}
