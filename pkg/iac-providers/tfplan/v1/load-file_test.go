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

package tfplan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

func TestLoadIacFile(t *testing.T) {

	invalidFilePathErr := os.PathError{Op: "open", Path: "not-there", Err: syscall.ENOENT}

	table := []struct {
		name     string
		filePath string
		options  map[string]interface{}
		tfplan   TFPlan
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfplan:   TFPlan{},
			wantErr:  fmt.Errorf("failed to read tfplan JSON file. error: '%s'", invalidFilePathErr.Error()),
		},
		{
			name:     "invalid json",
			filePath: filepath.Join("testdata", "invalid-json.json"),
			tfplan:   TFPlan{},
			wantErr:  fmt.Errorf("invalid terraform json file; error: 'failed to decode tfplan json. error: 'invalid character 'I' looking for beginning of value''"),
		},
		{
			name:     "invalid tfplan json",
			filePath: filepath.Join("testdata", "invalid-tfplan.json"),
			tfplan:   TFPlan{},
			wantErr:  fmt.Errorf("invalid terraform json file; error: 'terraform format version should be one of '0.1, 0.2''"),
		},
		{
			name:     "valid tfplan json",
			filePath: filepath.Join("testdata", "valid-tfplan.json"),
			tfplan:   TFPlan{},
			wantErr:  nil,
		},
		{
			name:     "valid tfplan v0.2 json",
			filePath: filepath.Join("testdata", "valid-tfplan-0.2.json"),
			tfplan:   TFPlan{},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfplan.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
	options := make(map[string]interface{})
	t.Run("validate tfplan iac output", func(t *testing.T) {
		var (
			tfplan             = TFPlan{}
			tfplanFile         = filepath.Join("testdata", "valid-tfplan.json")
			tfplanOutput       = filepath.Join("testdata", "valid-tfplan-resource-config.json")
			wantErr      error = nil
		)

		got, err := tfplan.LoadIacFile(tfplanFile, options)
		if !reflect.DeepEqual(err, wantErr) {
			t.Errorf("error want: '%v', got: '%v'", wantErr, err)
		}

		gotBytes, _ := json.MarshalIndent(got, "", "  ")
		gotBytes = append(gotBytes, []byte{'\n'}...)
		wantBytes, _ := os.ReadFile(tfplanOutput)
		if utils.IsWindowsPlatform() {
			wantBytes = utils.ReplaceWinNewLineBytes(wantBytes)
		}
		if !reflect.DeepEqual(bytes.TrimSpace(gotBytes), bytes.TrimSpace(wantBytes)) {
			t.Errorf("unexpected error; got '%v', want: '%v'", string(gotBytes), string(wantBytes))
		}
	})
}

func TestIsValidTFPlanJSON(t *testing.T) {

	tfplan := TFPlan{}

	table := []struct {
		name    string
		tfjson  []byte
		wantErr error
	}{
		{
			name:    "invalid json",
			tfjson:  []byte("I am invalid"),
			wantErr: fmt.Errorf("failed to decode tfplan json. error: 'invalid character 'I' looking for beginning of value'"),
		},
		{
			name:    "incorrect terraform format version",
			tfjson:  []byte(`{"format_version": "bad version"}`),
			wantErr: errIncorrectFormatVersion,
		},
		{
			name:    "empty terraform version",
			tfjson:  []byte(`{"format_version": "0.1"}`),
			wantErr: errEmptyTerraformVersion,
		},
		{
			name:    "valid tfplan json",
			tfjson:  []byte(`{"format_version": "0.1", "terraform_version": "0.12.3"}`),
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err := tfplan.isValidTFPlanJSON(tt.tfjson)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", err, tt.wantErr)
			}
		})
	}
}
