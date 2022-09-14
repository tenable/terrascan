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

package cli

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/runtime"
	"github.com/tenable/terrascan/pkg/utils"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	// set default config values before policy download
	config.LoadGlobalConfig("")

	// to download the policies for Run test
	// downloads the policies at $HOME/.terrascan
	initial(nil, nil, false)
}

func shutdown() {
	// remove the downloaded policies
	os.RemoveAll(config.GetPolicyBasePath())
	// cleanup the loaded config values
}

var testDataDir = "testdata"
var runTestDir = filepath.Join(testDataDir, "run-test")

func TestRun(t *testing.T) {
	// disable terraform logs when TF_LOG env variable is not set
	if os.Getenv("TF_LOG") == "" {
		log.SetOutput(io.Discard)
	}

	kustomizeTestDirPath := filepath.Join(runTestDir, "kustomize-test")
	testTerraformFilePath := filepath.Join(runTestDir, "config-only.tf")
	testRemoteModuleFilePath := filepath.Join(runTestDir, "remote-modules.tf")
	testTFJSONFilePath := filepath.Join(runTestDir, "tf-plan.json")

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
				iacDirPath: runTestDir,
			},
			wantErr: true,
		},
		{
			name: "normal terraform run with successful output",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: runTestDir,
				outputType: "json",
			},
		},
		{
			name: "terraform run with --non-recursive flag",
			scanOptions: &ScanOptions{
				iacType:      "terraform",
				policyType:   []string{"all"},
				iacDirPath:   testDataDir,
				outputType:   "json",
				nonRecursive: true,
			},
			wantErr: true,
		},
		{
			name: "normal k8s run",
			scanOptions: &ScanOptions{
				policyType: []string{"k8s"},
				// kustomization.y(a)ml file not present under the dir path, error expected
				iacDirPath: runTestDir,
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
			name: "normal k8s run with successful output for junit-xml with passed tests",
			scanOptions: &ScanOptions{
				policyType:      []string{"k8s"},
				iacDirPath:      kustomizeTestDirPath,
				outputType:      "junit-xml",
				showPassedRules: true,
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
			// test for https://github.com/tenable/terrascan/issues/718
			// a valid tfplan file is supplied, error is not expected
			name: "iac type is tfplan and -f option used to specify the tfplan.json",
			scanOptions: &ScanOptions{
				policyType:  []string{"all"},
				iacType:     "tfplan",
				iacFilePath: testTFJSONFilePath,
				outputType:  "yaml",
			},
			wantErr: false,
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
				iacDirPath: runTestDir,
				outputType: "json",
				skipRules:  ruleSlice,
			},
		},
		{
			name: "run with scan rules",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: runTestDir,
				outputType: "yaml",
				scanRules:  ruleSlice,
			},
		},
		{
			name: "config file with rules",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				iacDirPath: runTestDir,
				outputType: "yaml",
				configFile: filepath.Join(testDataDir, "configFile.toml"),
			},
		},
		{
			name: "scan file with remote module",
			scanOptions: &ScanOptions{
				policyType:  []string{"all"},
				iacFilePath: testRemoteModuleFilePath,
				outputType:  "human",
				configFile:  filepath.Join(testDataDir, "configFile.toml"),
			},
		},
		{
			name: "invalid remote type",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				remoteType: "test",
				remoteURL:  "test",
				outputType: "human",
			},
			wantErr: true,
		},
		{
			name: "valid remote type with invalid remote url",
			scanOptions: &ScanOptions{
				policyType: []string{"all"},
				remoteType: "terraform-registry",
				remoteURL:  "terraform-aws-modules/eks",
				outputType: "human",
			},
			wantErr: true,
		},
		{
			name: "config-with-error flag terraform",
			scanOptions: &ScanOptions{
				policyType:      []string{"all"},
				iacFilePath:     testTerraformFilePath,
				configWithError: true,
				outputType:      "json",
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			config.LoadGlobalConfig(tt.scanOptions.configFile)

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
				RemoteURL:  "github.com/tenable/terrascan",
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
		configOnly      bool
		configWithError bool
		outputType      string
		useColors       string
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
		{
			name: "init fail for --config-with-error with human readable output",
			fields: fields{
				useColors:       "auto",
				outputType:      "human",
				configWithError: true,
			},
			wantErr: true,
		},
		{
			name: "init success for --config-with-error with yaml readable output",
			fields: fields{
				useColors:       "auto",
				outputType:      "yaml",
				configWithError: true,
			},
			wantErr: false,
		},
		{
			name: "init success for --config-with-error with json readable output",
			fields: fields{
				useColors:       "auto",
				outputType:      "json",
				configWithError: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanOptions{
				configOnly:      tt.fields.configOnly,
				configWithError: tt.fields.configWithError,
				outputType:      tt.fields.outputType,
				useColors:       tt.fields.useColors,
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
				iacDirPath: runTestDir,
			},
			wantErr: true,
		},
		{
			name: "successful scan",
			fields: fields{
				policyType: []string{"all"},
				iacDirPath: runTestDir,
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

func Test_getExitCode(t *testing.T) {
	testDirScanErrors := []results.DirScanErr{
		{
			IacType:    "all",
			Directory:  "test",
			ErrMessage: "error occurred",
		},
	}

	testScanSummary := results.ScanSummary{
		ViolatedPolicies: 1,
	}

	scanOutputWithDirErrorsOnly := runtime.Output{
		Violations: policy.EngineOutput{
			ViolationStore: &results.ViolationStore{
				DirScanErrors: testDirScanErrors,
			},
		},
	}

	scanOutputWithDirErrorsAndViolatedPolicies := runtime.Output{
		Violations: policy.EngineOutput{
			ViolationStore: &results.ViolationStore{
				DirScanErrors: testDirScanErrors,
				Summary:       testScanSummary,
			},
		},
	}

	scanOutputWithViolatedPoliciesOnly := runtime.Output{
		Violations: policy.EngineOutput{
			ViolationStore: &results.ViolationStore{
				Summary: testScanSummary,
			},
		},
	}

	type args struct {
		o runtime.Output
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "has directory scan errors without violated policies",
			args: args{
				o: scanOutputWithDirErrorsOnly,
			},
			want: 4,
		},
		{
			name: "has directory scan errors with violated policies",
			args: args{
				o: scanOutputWithDirErrorsAndViolatedPolicies,
			},
			want: 5,
		},
		{
			name: "has violated policies but no directory scan errors",
			args: args{
				o: scanOutputWithViolatedPoliciesOnly,
			},
			want: 3,
		},
		{
			name: "neither violations nor directory scan errors",
			args: args{
				o: runtime.Output{
					Violations: policy.EngineOutput{
						ViolationStore: &results.ViolationStore{},
					},
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExitCode(tt.args.o); got != tt.want {
				t.Errorf("getExitCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
