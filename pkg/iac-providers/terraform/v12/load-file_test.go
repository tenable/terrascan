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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

var testDataDir = "testdata"
var emptyTfFilePath = filepath.Join(testDataDir, "empty.tf")
var nonEmptyTfFilePath = filepath.Join(testDataDir, "destroy-provisioners", "main.tf")
var tfJSONDir = filepath.Join(testDataDir, "tfjson")

func TestLoadIacFile(t *testing.T) {

	testErrorString1 := `error occured while loading config file 'not-there'. error:
<nil>: Failed to read file; The file "not-there" could not be read.
`
	testErrorString2 := fmt.Sprintf(`failed to load config file '%s'. error:
%s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`, emptyTfFilePath, emptyTfFilePath, emptyTfFilePath)

	testErrorString3 := fmt.Sprintf(`failed to load config file '%s'. error:
%s:8,12-22: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:42,15-35: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:39,14-24: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
`, nonEmptyTfFilePath, nonEmptyTfFilePath, nonEmptyTfFilePath, nonEmptyTfFilePath)

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
			wantErr:  fmt.Errorf(testErrorString1),
		},
		{
			name:     "empty config",
			filePath: filepath.Join(testDataDir, "testfile"),
			tfv12:    TfV12{},
		},
		{
			name:     "invalid config",
			filePath: emptyTfFilePath,
			tfv12:    TfV12{},
			wantErr:  fmt.Errorf(testErrorString2),
		},
		{
			name:     "destroy-provisioners",
			filePath: nonEmptyTfFilePath,
			tfv12:    TfV12{},
			wantErr:  fmt.Errorf(testErrorString3),
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
			tfConfigFile: filepath.Join(testDataDir, "tfconfigs", "config1.tf"),
			tfJSONFile:   filepath.Join(tfJSONDir, "config1.json"),
			tfv12:        TfV12{},
			wantErr:      nil,
		},
		{
			name:         "dummyconfig",
			tfConfigFile: filepath.Join(testDataDir, "dummyconfig", "dummyconfig.tf"),
			tfJSONFile:   filepath.Join(tfJSONDir, "dummyconfig.json"),
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

			if utils.IsWindowsPlatform() {
				gotBytes = utils.ReplaceCarriageReturnBytes(gotBytes)
				wantBytes = utils.ReplaceWinNewLineBytes(wantBytes)
			}

			if !reflect.DeepEqual(bytes.TrimSpace(gotBytes), bytes.TrimSpace(wantBytes)) {
				t.Errorf("unexpected error; got '%v', want: '%v'", string(gotBytes), string(wantBytes))
			}
		})
	}
}
