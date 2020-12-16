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
		iacType     string
		iacVersion  string
		cloudType   []string
		format      string
		iacFilePath string
		iacDirPath  string
		configFile  string
		configOnly  bool
		verbose     bool
		stdOut      string
		want        string
		wantErr     error
	}{
		{
			name:       "normal terraform run",
			cloudType:  []string{"terraform"},
			iacDirPath: "testdata/run-test",
		},
		{
			name:       "normal k8s run",
			cloudType:  []string{"k8s"},
			iacDirPath: "testdata/run-test",
		},
		{
			name:        "config-only flag terraform",
			cloudType:   []string{"terraform"},
			iacFilePath: "testdata/run-test/config-only.tf",
			configOnly:  true,
		},
		{
			name:        "config-only flag k8s",
			cloudType:   []string{"k8s"},
			iacFilePath: "testdata/run-test/config-only.yaml",
			configOnly:  true,
		},
		{
			name:        "config-only flag true with human readable format",
			cloudType:   []string{"terraform"},
			iacFilePath: "testdata/run-test/config-only.tf",
			configOnly:  true,
			format:      "human",
		},
		{
			name:        "config-only flag false with human readable format",
			cloudType:   []string{"k8s"},
			iacFilePath: "testdata/run-test/config-only.yaml",
			format:      "human",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.iacType, tt.iacVersion, tt.cloudType, tt.iacFilePath, tt.iacDirPath, tt.configFile, []string{}, tt.format, "", "", tt.configOnly, false, tt.verbose)
		})
	}
}

func TestWriteResults(t *testing.T) {
	testInput := runtime.Output{
		ResourceConfig: output.AllResourceConfigs{},
		Violations: policy.EngineOutput{
			ViolationStore: &results.ViolationStore{},
		},
	}
	type args struct {
		results    runtime.Output
		useColors  bool
		verbose    bool
		configOnly bool
		format     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "config only true with human readable output format",
			args: args{
				results:    testInput,
				configOnly: true,
				format:     "human",
			},
			wantErr: true,
		},
		{
			name: "config only true with non human readable output format",
			args: args{
				results:    testInput,
				configOnly: true,
				format:     "json",
			},
			wantErr: false,
		},
		{
			name: "config only false",
			args: args{
				results:    testInput,
				configOnly: false,
				format:     "human",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeResults(tt.args.results, tt.args.useColors, tt.args.verbose, tt.args.configOnly, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("writeResults() error = gotErr: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownloadRemoteRepository(t *testing.T) {
	testTempdir := filepath.Join(os.TempDir(), utils.GenRandomString(6))

	type args struct {
		remoteType string
		remoteURL  string
		tempDir    string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "blank input paramters",
			args: args{
				remoteType: "",
				remoteURL:  "",
				tempDir:    "",
			},
		},
		{
			name: "invalid input parameters",
			args: args{
				remoteType: "test",
				remoteURL:  "test",
				tempDir:    "test",
			},
			wantErr: true,
		},
		{
			name: "valid inputs paramters",
			args: args{
				remoteType: "git",
				remoteURL:  "github.com/accurics/terrascan",
				tempDir:    testTempdir,
			},
			want: testTempdir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := downloadRemoteRepository(tt.args.remoteType, tt.args.remoteURL, tt.args.tempDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("downloadRemoteRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("downloadRemoteRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
