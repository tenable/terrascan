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

package initialize

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const TESTDIR = "test_data"

func TestGetCommercialPolicy(t *testing.T) {
	table := []struct {
		name       string
		configFile string
		wantErr    error
	}{
		{
			name:       "invalid policy json",
			configFile: filepath.Join(TESTDIR, "invalid_policies.json"),
			wantErr:    errors.New("failed to unmarshal policies into structure"),
		},
		{
			name:       "invalid ruleArgument data type",
			configFile: filepath.Join(TESTDIR, "invalid_ruleArg_type_policies.json"),
			wantErr:    errors.New("incorrect rule argument type, must be a string"),
		},
		{
			name:       "invalid ruleArgument format",
			configFile: filepath.Join(TESTDIR, "invalid_ruleArg_policies.json"),
			wantErr:    errors.New("error occurred while unmarshaling rule arguments into map[string]interface{}"),
		},
		{
			name:       "valid policy json",
			configFile: filepath.Join(TESTDIR, "policies.json"),
			wantErr:    nil,
		},
	}

	tempDir, err := os.MkdirTemp("", "test_policies")
	if err != nil {
		t.Errorf("unable to create temporary dir for testing")
	}
	defer os.RemoveAll(tempDir)

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			policies, err := os.ReadFile(tt.configFile)
			if err != nil {
				t.Errorf("unable to read test file")
			}

			err = convertEnvironmentPolicies(policies, tempDir)
			if err == nil {
				validPoliciesTest(t, tempDir)
				return
			}

			if !strings.HasPrefix(err.Error(), tt.wantErr.Error()) {
				t.Errorf("convertEnvironmentPolicies() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func validPoliciesTest(t *testing.T, tempDir string) {
	expectedCspMap := map[string]string{
		"aws":   "aws_cloudwatch_log_group",
		"azure": "azurerm_security_center_setting",
		"gcp":   "google_compute_firewall",
	}

	providers, err := os.ReadDir(tempDir)
	if err != nil {
		t.Errorf("unable to read temp dir: '%v'", err)
	}
	if len(providers) != 3 {
		t.Errorf("expected length 3; got: '%d'", len(providers))
	}

	for _, provider := range providers {
		cspdir := filepath.Join(tempDir, provider.Name())
		rscs, err := os.ReadDir(cspdir)
		if err != nil {
			t.Errorf("unable to read csp dir: '%v'", err)
		}
		if len(rscs) != 1 {
			t.Errorf("expected length 1; got '%d'", len(rscs))
		}

		val, ok := expectedCspMap[provider.Name()]
		if !ok {
			t.Errorf("unable to find expected csp")

		}
		if rscs[0].Name() != val {
			t.Errorf("unable to find expected resource type")
		}

		rscdir := filepath.Join(cspdir, rscs[0].Name())
		files, err := os.ReadDir(rscdir)
		if err != nil {
			t.Errorf("unable to read resource type dir: '%v'", err)
		}

		err = verifyFiles(files, rscdir)
		if err != nil {
			t.Error(err)
		}
	}
}

func verifyFiles(files []fs.DirEntry, rscdir string) error {
	expectedMetaMap := map[string]bool{
		"AC_AWS_0452.json":   true,
		"AC_AZURE_0330.json": true,
		"AC_GCP_0095.json":   true,
	}

	expectedRegoMap := map[string]bool{
		"awsNoRetentionPolicyCloudwatchLogGroup.rego":   true,
		"ensureProperSettings_25082021.rego":            true,
		"networkPortExposedToInternetGCP_25082021.rego": true,
	}

	if len(files) != 2 {
		return fmt.Errorf("expected length 2; got: '%d'", len(files))
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			_, ok := expectedMetaMap[file.Name()]
			if !ok {
				return fmt.Errorf("expected metadata file not found")
			}
		}

		if strings.HasSuffix(file.Name(), ".rego") {
			_, ok := expectedRegoMap[file.Name()]
			if !ok {
				return fmt.Errorf("expected rego file not found")
			}
		}

		testpath := filepath.Join(TESTDIR, file.Name())
		convpath := filepath.Join(rscdir, file.Name())

		testbytes, err := os.ReadFile(testpath)
		if err != nil {
			return fmt.Errorf("unable to read test file: '%s', err: '%w'", testpath, err)
		}

		databytes, err := os.ReadFile(convpath)
		if err != nil {
			return fmt.Errorf("unable to read converted file: '%s', err: '%w'", convpath, err)
		}

		if !bytes.Equal(testbytes, databytes) {
			return fmt.Errorf("test file and converted file do not match")
		}
	}

	return nil
}
