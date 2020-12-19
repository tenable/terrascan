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

func TestRun(t *testing.T) {
	table := []struct {
		name        string
		configFile  string
		format      string
		scanCommand *ScanCommand
		stdOut      string
		want        string
		wantErr     error
	}{
		{
			name: "normal terraform run",
			scanCommand: &ScanCommand{
				policyType: []string{"terraform"},
				iacDirPath: "testdata/run-test",
			},
		},
		{
			name: "normal k8s run",
			scanCommand: &ScanCommand{
				policyType: []string{"k8s"},
				iacDirPath: "testdata/run-test",
			},
		},
		{
			name: "config-only flag terraform",
			scanCommand: &ScanCommand{
				policyType: []string{"terraform"},
				iacDirPath: "testdata/run-test/config-only.tf",
				configOnly: true,
			},
		},
		{
			name: "config-only flag k8s",
			scanCommand: &ScanCommand{
				policyType: []string{"k8s"},
				iacDirPath: "testdata/run-test/config-only.yaml",
				configOnly: true,
			},
		},
		{
			name: "config-only flag true with human readable format",
			scanCommand: &ScanCommand{
				policyType: []string{"terraform"},
				iacDirPath: "testdata/run-test/config-only.tf",
				configOnly: true,
			},
			format: "human",
		},
		{
			name: "config-only flag false with human readable format",
			scanCommand: &ScanCommand{
				policyType: []string{"k8s"},
				iacDirPath: "testdata/run-test/config-only.yaml",
			},
			format: "human",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			tt.scanCommand.Run()
		})
	}
}

func TestScanCommand_downloadRemoteRepository(t *testing.T) {
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
			name: "invalid input parameters",
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
			s := ScanCommand{
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

func TestScanCommand_writeResults(t *testing.T) {
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
			s := ScanCommand{
				configOnly: tt.fields.ConfigOnly,
				outputType: tt.fields.OutputType,
			}
			if err := s.writeResults(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ScanOptions.writeResults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanCommand_validate(t *testing.T) {
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
			s := ScanCommand{
				configOnly: tt.fields.configOnly,
				outputType: tt.fields.outputType,
			}
			if err := s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("ScanCommand.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScanCommand_initColor(t *testing.T) {
	type fields struct {
		useColors string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "auto",
			fields: fields{
				useColors: "auto",
			},
		},
		{
			name: "true",
			fields: fields{
				useColors: "true",
			},
			want: true,
		},
		{
			name: "1",
			fields: fields{
				useColors: "1",
			},
			want: true,
		},
		{
			name: "false",
			fields: fields{
				useColors: "false",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanCommand{
				useColors: tt.fields.useColors,
			}
			s.initColor()
			if s.useColors != "auto" {
				if s.UseColors != tt.want {
					t.Errorf("ScanCommand.initColor() incorrect value for UseColors, got: %v, want %v", s.useColors, tt.want)
				}
			}
		})
	}
}

func TestScanCommand_Init(t *testing.T) {
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
			name: "test for init fail",
			fields: fields{
				useColors:  "auto",
				outputType: "human",
				configOnly: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScanCommand{
				configOnly: tt.fields.configOnly,
				outputType: tt.fields.outputType,
				useColors:  tt.fields.useColors,
			}
			if err := s.Init(); (err != nil) != tt.wantErr {
				t.Errorf("ScanCommand.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
