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

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name    string
		dirPath string
		tfv12   TfV12
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			tfv12:   TfV12{},
			wantErr: errDirEmptyTFConfig,
		},
		{
			name:    "empty config",
			dirPath: "./testdata/testfile",
			tfv12:   TfV12{},
			wantErr: errDirEmptyTFConfig,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name        string
		tfConfigDir string
		tfJSONFile  string
		tfv12       TfV12
		wantErr     error
	}{
		{
			name:        "config1",
			tfConfigDir: "./testdata/tfconfigs",
			tfJSONFile:  "testdata/tfjson/fullconfig.json",
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "module directory",
			tfConfigDir: "./testdata/moduleconfigs",
			tfJSONFile:  "./testdata/tfjson/moduleconfigs.json",
			tfv12:       TfV12{},
			wantErr:     nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacDir(tt.tfConfigDir)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			gotBytes, _ := json.MarshalIndent(got, "", "  ")
			gotBytes = append(gotBytes, []byte{'\n'}...)
			wantBytes, _ := ioutil.ReadFile(tt.tfJSONFile)

			if !bytes.Equal(bytes.TrimSpace(gotBytes), bytes.TrimSpace(wantBytes)) {
				t.Errorf("got '%v', want: '%v'", string(gotBytes), string(wantBytes))
			}
		})
	}
}
