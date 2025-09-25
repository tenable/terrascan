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

package tfv12

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tenable/terrascan/pkg/iac-providers/terraform/commons"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

var testDataDir = "testdata"
var emptyTfFilePath = filepath.Join(testDataDir, "empty.tf")
var nonEmptyTfFilePath = filepath.Join(testDataDir, "destroy-provisioners", "main.tf")
var tfJSONDir = filepath.Join(testDataDir, "tfjson")

const errMsgFailedLoadingIACFile2 = `failed to load iac file '%s'. error: %s:8,12-22: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:42,15-35: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:39,14-24: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.`

func TestLoadIacFile(t *testing.T) {

	testErrorString1 := fmt.Errorf(commons.ErrMsgFailedLoadingConfigFile)
	testErrorString2 := fmt.Errorf(commons.ErrMsgFailedLoadingIACFile, emptyTfFilePath, emptyTfFilePath, emptyTfFilePath)
	testErrorString3 := fmt.Errorf(errMsgFailedLoadingIACFile2, nonEmptyTfFilePath, nonEmptyTfFilePath, nonEmptyTfFilePath, nonEmptyTfFilePath)

	table := []struct {
		name     string
		filePath string
		options  map[string]interface{}
		tfv12    TfV12
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfv12:    TfV12{},
			wantErr:  testErrorString1,
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
			wantErr:  testErrorString2,
		},
		{
			name:     "destroy-provisioners",
			filePath: nonEmptyTfFilePath,
			tfv12:    TfV12{},
			wantErr:  testErrorString3,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacFile(tt.filePath, tt.options)
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
			got, gotErr := tt.tfv12.LoadIacFile(tt.tfConfigFile, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			gotBytes, _ := json.MarshalIndent(got, "", "  ")
			gotBytes = append(gotBytes, []byte{'\n'}...)
			wantBytes, _ := os.ReadFile(tt.tfJSONFile)

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
