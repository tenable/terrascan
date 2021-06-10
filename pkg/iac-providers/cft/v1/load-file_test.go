/*
    Copyright (C) 2021 Accurics, Inc.

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

package cftv1

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {
	testDataDir := "testdata"
	testFile, _ := filepath.Abs(path.Join(testDataDir, "testfile"))
	invalidFile, _ := filepath.Abs(path.Join(testDataDir, "deploy.yaml"))
	validFile, _ := filepath.Abs(path.Join(testDataDir, "templates", "s3", "deploy.template"))

	testErrString1 := fmt.Sprintf("unsupported extension for file %s", testFile)
	testErrString2 := "unable to read file nonexistent.txt"
	testErrString3 := "invalid YAML template: yaml: line 27: did not find expected alphabetic or numeric character"

	table := []struct {
		wantErr  error
		want     output.AllResourceConfigs
		cftv1    CFTV1
		name     string
		filePath string
		typeOnly bool
	}{
		{
			wantErr:  fmt.Errorf(testErrString1),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "invalid extension",
			filePath: testFile,
			typeOnly: false,
		}, {
			wantErr:  fmt.Errorf(testErrString2),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "nonexistent file",
			filePath: "nonexistent.txt",
			typeOnly: false,
		}, {
			wantErr:  fmt.Errorf(testErrString3),
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "invalid file",
			filePath: invalidFile,
			typeOnly: false,
		}, {
			wantErr:  nil,
			want:     output.AllResourceConfigs{},
			cftv1:    CFTV1{},
			name:     "invalid file",
			filePath: validFile,
			typeOnly: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.cftv1.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%+v', wantErr: '%+v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
