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

package dockerv1

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

func TestLoadIacFile(t *testing.T) {

	tests := []struct {
		name        string
		absFilePath string
		dockerV1    DockerV1
		want        output.AllResourceConfigs
		wantErr     error
		typeOnly    bool
	}{
		{
			name:        "empty config file",
			absFilePath: filepath.Join(fileTestDataDir, "Dockerfile"),
			dockerV1:    DockerV1{},
			wantErr:     fmt.Errorf("error while parsing file %s, error: file with no instructions", filepath.Join(fileTestDataDir, "Dockerfile")),
		},
		{
			name:        "valid docker file",
			absFilePath: filepath.Join(fileTestDataDir, "valid-Dockerfile"),
			dockerV1:    DockerV1{},
			wantErr:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, gotErr := tt.dockerV1.LoadIacFile(tt.absFilePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
