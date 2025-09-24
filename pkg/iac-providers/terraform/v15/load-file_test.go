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

package tfv15

import (
	"encoding/json"
	"fmt"
	"github.com/tenable/terrascan/pkg/iac-providers/terraform/commons"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/tenable/terrascan/pkg/utils"
)

var testDataDir = "testdata"
var emptyTfFilePath = filepath.Join(testDataDir, "empty.tf")

func TestLoadIacFile(t *testing.T) {

	testErrorString1 := fmt.Errorf(commons.ErrMsgFailedLoadingConfigFile)
	testErrorString2 := fmt.Errorf(commons.ErrMsgFailedLoadingIACFile, emptyTfFilePath, emptyTfFilePath, emptyTfFilePath)

	table := []struct {
		name     string
		filePath string
		options  map[string]interface{}
		tfv15    TfV15
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfv15:    TfV15{},
			wantErr:  testErrorString1,
		},
		{
			name:     "empty config",
			filePath: filepath.Join(testDataDir, "testfile"),
			tfv15:    TfV15{},
		},
		{
			name:     "invalid config",
			filePath: filepath.Join(testDataDir, "empty.tf"),
			tfv15:    TfV15{},
			wantErr:  testErrorString2,
		},
		{
			name:     "depends_on",
			filePath: filepath.Join(testDataDir, "depends_on", "main.tf"),
			tfv15:    TfV15{},
		},
		{
			name:     "count",
			filePath: filepath.Join(testDataDir, "count", "main.tf"),
			tfv15:    TfV15{},
		},
		{
			name:     "for_each",
			filePath: filepath.Join(testDataDir, "for_each", "main.tf"),
			tfv15:    TfV15{},
		},
		{
			name:     "required_providers",
			filePath: filepath.Join(testDataDir, "required-providers", "main.tf"),
			tfv15:    TfV15{},
		},
		{
			name:     "required_providers with configuration alias",
			filePath: filepath.Join(testDataDir, "required-providers", "configuration-alias", "main.tf"),
			tfv15:    TfV15{},
		},
		{
			name:     "provider with only alias",
			filePath: filepath.Join(testDataDir, "provider-with-only-alias", "main.tf"),
			tfv15:    TfV15{},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv15.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name         string
		tfConfigFile string
		tfJSONFile   string
		options      map[string]interface{}
		tfv15        TfV15
		wantErr      error
	}{
		{
			name:         "config1",
			tfConfigFile: filepath.Join(testDataDir, "tfconfigs", "config1.tf"),
			tfJSONFile:   filepath.Join(testDataDir, "tfjson", "config1.json"),
			tfv15:        TfV15{},
			wantErr:      nil,
		},
		{
			name:         "dummyconfig",
			tfConfigFile: filepath.Join(testDataDir, "dummyconfig", "dummyconfig.tf"),
			tfJSONFile:   filepath.Join(testDataDir, "tfjson", "dummyconfig.json"),
			tfv15:        TfV15{},
			wantErr:      nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv15.LoadIacFile(tt.tfConfigFile, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			var want output.AllResourceConfigs

			wantBytes, _ := os.ReadFile(tt.tfJSONFile)
			if utils.IsWindowsPlatform() {
				wantBytes = utils.ReplaceWinNewLineBytes(wantBytes)
			}

			err := json.Unmarshal(wantBytes, &want)
			if err != nil {
				t.Errorf("unexpected error unmarshalling want: %v", err)
			}

			match, err := test.IdenticalAllResourceConfigs(got, want)
			if err != nil {
				t.Errorf("unexpected error checking result: %v", err)
			}
			if !match {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(want, "", "  ")
				t.Errorf("got '%v', want: '%v'", string(g), string(w))
			}
		})
	}
}
