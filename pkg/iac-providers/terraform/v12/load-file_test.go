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

package tfv12

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name     string
		filePath string
		tfv12    TfV12
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfv12:    TfV12{},
			wantErr:  errLoadConfigFile,
		},
		{
			name:     "empty config",
			filePath: "./testdata/testfile",
			tfv12:    TfV12{},
			wantErr:  nil,
		},
		{
			name:     "invalid config",
			filePath: "./testdata/empty.tf",
			tfv12:    TfV12{},
			wantErr:  errLoadConfigFile,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name         string
		tfConfigFile string
		tfJSONFile   string
		tfv12        TfV12
		wantErr      error
	}{
		{
			name:         "config1",
			tfConfigFile: "./testdata/tfconfigs/config1.tf",
			tfJSONFile:   "./testdata/tfjson/config1.json",
			tfv12:        TfV12{},
			wantErr:      nil,
		},
		{
			name:         "dummyconfig",
			tfConfigFile: "./testdata/dummyconfig/dummyconfig.tf",
			tfJSONFile:   "./testdata/tfjson/dummyconfig.json",
			tfv12:        TfV12{},
			wantErr:      nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacFile(tt.tfConfigFile)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			gotBytes, _ := json.MarshalIndent(got, "", "  ")
			gotBytes = append(gotBytes, []byte{'\n'}...)
			wantBytes, _ := ioutil.ReadFile(tt.tfJSONFile)
			if !reflect.DeepEqual(bytes.TrimSpace(gotBytes), bytes.TrimSpace(wantBytes)) {
				t.Errorf("unexpected error; got '%v', want: '%v'", string(gotBytes), string(wantBytes))
			}
		})
	}
}
