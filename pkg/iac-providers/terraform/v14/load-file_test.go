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

package tfv14

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	testErrorString1 := `error occured while loading config file 'not-there'. error:
<nil>: Failed to read file; The file "not-there" could not be read.
`
	testErrorString2 := `failed to load config file './testdata/empty.tf'. error:
./testdata/empty.tf:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
./testdata/empty.tf:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`

	table := []struct {
		name     string
		filePath string
		tfv14    TfV14
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfv14:    TfV14{},
			wantErr:  fmt.Errorf(testErrorString1),
		},
		{
			name:     "empty config",
			filePath: "./testdata/testfile",
			tfv14:    TfV14{},
		},
		{
			name:     "invalid config",
			filePath: "./testdata/empty.tf",
			tfv14:    TfV14{},
			wantErr:  fmt.Errorf(testErrorString2),
		},
		{
			name:     "depends_on",
			filePath: "./testdata/depends_on/main.tf",
			tfv14:    TfV14{},
		},
		{
			name:     "count",
			filePath: "./testdata/count/main.tf",
			tfv14:    TfV14{},
		},
		{
			name:     "for_each",
			filePath: "./testdata/for_each/main.tf",
			tfv14:    TfV14{},
		},
		{
			name:     "required_providers",
			filePath: "./testdata/required-providers/main.tf",
			tfv14:    TfV14{},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv14.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name         string
		tfConfigFile string
		tfJSONFile   string
		tfv14        TfV14
		wantErr      error
	}{
		{
			name:         "config1",
			tfConfigFile: "./testdata/tfconfigs/config1.tf",
			tfJSONFile:   "./testdata/tfjson/config1.json",
			tfv14:        TfV14{},
			wantErr:      nil,
		},
		{
			name:         "dummyconfig",
			tfConfigFile: "./testdata/dummyconfig/dummyconfig.tf",
			tfJSONFile:   "./testdata/tfjson/dummyconfig.json",
			tfv14:        TfV14{},
			wantErr:      nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv14.LoadIacFile(tt.tfConfigFile)
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
