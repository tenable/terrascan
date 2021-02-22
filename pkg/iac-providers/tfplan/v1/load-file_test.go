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

package tfplan

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

	table := []struct {
		name     string
		filePath string
		tfplan   TFPlan
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "invalid filepath",
			filePath: "not-there",
			tfplan:   TFPlan{},
			wantErr:  fmt.Errorf("failed to read tfplan JSON file. error: 'open not-there: no such file or directory'"),
		},
		{
			name:     "invalid json",
			filePath: "./testdata/invalid-json.json",
			tfplan:   TFPlan{},
			wantErr:  fmt.Errorf("failed to process tfplan JSON. error: 'failed to decode input JSON. error: 'invalid character 'I' looking for beginning of value''"),
		},
		{
			name:     "valid tfplan json",
			filePath: "./testdata/valid-tfplan.json",
			tfplan:   TFPlan{},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfplan.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	t.Run("validate tfplan iac output", func(t *testing.T) {
		var (
			tfplan             = TFPlan{}
			tfplanFile         = "./testdata/valid-tfplan.json"
			tfplanOutput       = "./testdata/valid-tfplan-resource-config.json"
			wantErr      error = nil
		)

		got, err := tfplan.LoadIacFile(tfplanFile)
		if !reflect.DeepEqual(err, wantErr) {
			t.Errorf("error want: '%v', got: '%v'", wantErr, err)
		}

		gotBytes, _ := json.MarshalIndent(got, "", "  ")
		gotBytes = append(gotBytes, []byte{'\n'}...)
		wantBytes, _ := ioutil.ReadFile(tfplanOutput)
		if !reflect.DeepEqual(bytes.TrimSpace(gotBytes), bytes.TrimSpace(wantBytes)) {
			t.Errorf("unexpected error; got '%v', want: '%v'", string(gotBytes), string(wantBytes))
		}
	})
}
