/*
    Copyright (C) 2022 Accurics, Inc.

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
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const TESTDIR = "test_data"

func TestGetCommercialPolicy(t *testing.T) {
	policiesPath := filepath.Join(TESTDIR, "policies.json")
	policies, err := ioutil.ReadFile(policiesPath)
	if err != nil {
		t.Errorf("unable to read test file")
	}

	tempDir, err := os.MkdirTemp("", "test_policies")
	if err != nil {
		t.Errorf("unable to create temporary dir for testing")
	}
	defer os.RemoveAll(tempDir)

	expectedCspMap := map[string]string{
		"aws":   "aws_cloudwatch_log_group",
		"azure": "azurerm_security_center_setting",
		"gcp":   "google_compute_firewall",
	}

	err = convertEnvironmentPolicies(policies, tempDir)
	if err != nil {
		t.Errorf("unable to convert and save policies: '%v'", err)
	}

	csps, err := os.ReadDir(tempDir)
	if err != nil {
		t.Errorf("unable to read temp dir: '%v'", err)
	}
	if len(csps) != 3 {
		t.Errorf("expected length 3; got: '%d'", len(csps))
	}

	for _, csp := range csps {
		cspdir := filepath.Join(tempDir, csp.Name())
		rscs, err := os.ReadDir(cspdir)
		if err != nil {
			t.Errorf("unable to read csp dir: '%v'", err)
			break
		}
		if len(rscs) != 1 {
			t.Errorf("expected length 1; got '%d'", len(rscs))
			break
		}

		val, ok := expectedCspMap[csp.Name()]
		if !ok {
			t.Errorf("unable to find expected csp")
			break
		}
		if rscs[0].Name() != val {
			t.Errorf("unable to find expected resource type")
			break
		}

		rscdir := filepath.Join(cspdir, rscs[0].Name())
		files, err := os.ReadDir(rscdir)
		if err != nil {
			t.Errorf("unable to read resource type dir: '%v'", err)
		}

		err = verifyFiles(files, rscdir)
		if err != nil {
			t.Error(err)
			break
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

		testbytes, err := ioutil.ReadFile(testpath)
		if err != nil {
			return fmt.Errorf("unable to read test file: '%s', err: '%w'", testpath, err)
		}

		databytes, err := ioutil.ReadFile(convpath)
		if err != nil {
			return fmt.Errorf("unable to read converted file: '%s', err: '%w'", convpath, err)
		}

		if bytes.Compare(testbytes, databytes) != 0 {
			return fmt.Errorf("test file and converted file do not match")
		}
	}

	return nil
}
