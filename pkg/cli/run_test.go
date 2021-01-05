/*
    Copyright (C) 2020 Accurics, Inc.

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

package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/runtime"
	"github.com/accurics/terrascan/pkg/utils"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// to download the policies for Run test
	// downloads the policies at $HOME/.terrascan
	initial(nil, nil)
}

func shutdown() {
	// remove the downloaded policies
	os.RemoveAll(os.Getenv("HOME") + "/.terrascan")
}

func TestRun(t *testing.T) {
	testDirPath := "testdata/run-test"
	kustomizeTestDirPath := testDirPath + "/kustomize-test"
	testTerraformFilePath := testDirPath + "/config-only.tf"
	ruleSlice := []string{"AWS.ECR.DataSecurity.High.0579", "AWS.SecurityGroup.NetworkPortsSecurity.Low.0561"}

	table := []struct {
		name        string
		configFile  string
		scanOptions *ScanOptions
		stdOut      string
		want        string
		wantErr     bool
	}{
		{
			name: "normal terraform run",
			scanOptions: &ScanOptions{
				// policy type terraform is not supported, error expected
				policyType: []string{"terraform"},
				iacDirPath: testDirPath,
			},
			wantErr: true,
		},
		{
			name: "normal terraform run with successful output",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: testDirPath,
				outputType: "json",
			},
		},
		{
			name: "normal k8s run",
			scanOptions: &ScanOptions{
				policyType: []string{"k8s"},
				// kustomization.y(a)ml file not present under the dir path, error expected
				iacDirPath: testDirPath,
			},
			wantErr: true,
		},
		{
			name: "normal k8s run with successful output",
			scanOptions: &ScanOptions{
				policyType: []string{"k8s"},
				iacDirPath: kustomizeTestDirPath,
				outputType: "human",
			},
		},
		{
			name: "config-only flag terraform",
			scanOptions: &ScanOptions{
				policyType:  []string{"all"},
				iacFilePath: testTerraformFilePath,
				configOnly:  true,
				outputType:  "yaml",
			},
		},
		{
			name: "config-only flag k8s",
			scanOptions: &ScanOptions{
				policyType: []string{"k8s"},
				iacDirPath: kustomizeTestDirPath,
				configOnly: true,
				outputType: "json",
			},
		},
		{
			// xml doesn't support config-only, error expected
			// modify the test results when xml supports config-only
			name: "config-only flag true with xml output format",
			scanOptions: &ScanOptions{
				policyType:  []string{"all"},
				iacFilePath: testTerraformFilePath,
				configOnly:  true,
				outputType:  "xml",
			},
			wantErr: true,
		},
		{
			name: "fail to download remote repository",
			scanOptions: &ScanOptions{
				policyType:  []string{"all"},
				iacFilePath: testTerraformFilePath,
				remoteURL:   "test",
				remoteType:  "test",
			},
			wantErr: true,
		},
		{
			name: "incorrect config file",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: testTerraformFilePath,
				outputType: "json",
				configFile: "invalidFile",
			},
			wantErr: true,
		},
		{
			name: "run with skip rules",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: testDirPath,
				outputType: "json",
				skipRules:  ruleSlice,
			},
		},
		{
			name: "run with scan rules",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: testDirPath,
				outputType: "yaml",
				scanRules:  ruleSlice,
			},
		},
		{
			name: "config file with rules",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: testDirPath,
				outputType: "yaml",
				configFile: "testdata/configFile.toml",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.scanOptions.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScanOptionsDownloadRemoteRepository(t *testing.T) {
	testTempdir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
	defer os.RemoveAll(testTempdir)

	type fields struct {
		RemoteType string
		RemoteURL  string
	}
	tests := []struct {
		name    string
		fields  fields
		tempDir string
		want    string
		wantErr bool
	}{
		{
			name: "blank input parameters",
			fields: fields{
				RemoteType: "",
				RemoteURL:  "",
			},
			tempDir: "",
		},
		{
			name: "invalid input parameters",
			fields: fields{
				RemoteType: "test",
				RemoteURL:  "test",
			},
			tempDir: "test",
			wantErr: true,
		},
		{
			name: "valid input parameters",
			fields: fields{
				RemoteType: "git",
				RemoteURL:  "github.com/accurics/terrascan",
			},
			tempDir: testTempdir,
			want:    testTempdir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ScanOptions{
				remoteType: tt.fields.RemoteType,
				remoteURL:  tt.fields.RemoteURL,
			}
			err := s.downloadRemoteRepository(tt.tempDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.downloadRemoteRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if s.iacDirPath != tt.want {
				t.Errorf("ScanOptions.downloadRemoteRepository() = %v, want %v", s.iacDirPath, tt.want)
			}
		})
	}
}

func TestScanOptionsWriteResults(t *testing.T) {
	testInput := runtime.Output{
		ResourceConfig: output.AllResourceConfigs{},
		Violations: policy.EngineOutput{
			ViolationStore: &results.ViolationStore{},
		},
	}

	type fields struct {
		ConfigOnly bool
		OutputType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    runtime.Output
		wantErr bool
	}{
		{
			name: "config only true",
			fields: fields{
				ConfigOnly: true,
				OutputType: "yaml",
			},
			args: testInput,
		},
		{
			name: "config only false",
			fields: fields{
				ConfigOnly: false,
				OutputType: "json",
			},
			args: testInput,
		},
		{
			// until we support config only flag for xml, this test case is for expected failure
			name: "config only true for xml",
			fields: fields{
				ConfigOnly: true,
				OutputType: "xml",
			},
			args:    testInput,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ScanOptions{
				configOnly: tt.fields.ConfigOnly,
				outputType: tt.fields.OutputType,
			}
			if err := s.writeResults(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.writeResults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanOptionsValidate(t *testing.T) {
	type fields struct {
		configOnly bool
		outputType string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate --config-only with human readable output",
			fields: fields{
				configOnly: true,
				outputType: "human",
			},
			wantErr: true,
		},
		{
			name: "validate --config-only with non human readable output",
			fields: fields{
				configOnly: true,
				outputType: "json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ScanOptions{
				configOnly: tt.fields.configOnly,
				outputType: tt.fields.outputType,
			}
			if err := s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanOptionsInitColor(t *testing.T) {
	type fields struct {
		useColors string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test for auto as input",
			fields: fields{
				useColors: "auto",
			},
		},
		{
			name: "test for true as input",
			fields: fields{
				useColors: "true",
			},
			want: true,
		},
		{
			name: "test for 1 as input",
			fields: fields{
				useColors: "1",
			},
			want: true,
		},
		{
			name: "test for false as input",
			fields: fields{
				useColors: "false",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanOptions{
				useColors: tt.fields.useColors,
			}
			s.initColor()
			if s.useColors != "auto" {
				if s.UseColors != tt.want {
					t.Errorf("ScanOptions.initColor() incorrect value for UseColors, got: %v, want %v", s.useColors, tt.want)
				}
			}
		})
	}
}

func TestScanOptionsInit(t *testing.T) {
	type fields struct {
		configOnly bool
		outputType string
		useColors  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test for init fail",
			fields: fields{
				useColors:  "auto",
				outputType: "human",
				configOnly: true,
			},
			wantErr: true,
		},
		{
			name: "test for init success",
			fields: fields{
				useColors:  "auto",
				outputType: "human",
				configOnly: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanOptions{
				configOnly: tt.fields.configOnly,
				outputType: tt.fields.outputType,
				useColors:  tt.fields.useColors,
			}
			if err := s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanOptionsScan(t *testing.T) {
	type fields struct {
		policyType []string
		iacDirPath string
		configOnly bool
		outputType string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "failure in init",
			fields: fields{
				configOnly: true,
				outputType: "human",
			},
			wantErr: true,
		},
		{
			name: "failure in run",
			fields: fields{
				policyType: []string{"terraform"},
				iacDirPath: "testdata/run-test",
			},
			wantErr: true,
		},
		{
			name: "successful scan",
			fields: fields{
				policyType: []string{"all"},
				iacDirPath: "testdata/run-test",
				outputType: "json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanOptions{
				policyType: tt.fields.policyType,
				iacDirPath: tt.fields.iacDirPath,
				configOnly: tt.fields.configOnly,
				outputType: tt.fields.outputType,
			}
			if err := s.Scan(); (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
