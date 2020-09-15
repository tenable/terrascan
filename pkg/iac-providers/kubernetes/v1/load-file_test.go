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

package k8sv1

import (
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name     string
		filePath string
		k8sV1    K8sV1
		typeOnly bool
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			// file is skipped if no kind is specified or bad
			name:     "empty config file",
			filePath: "./testdata/file-test-data/empty-file.yaml",
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			name:     "yaml with multiple documents",
			filePath: "./testdata/yaml-with-multiple-documents/test_pod.yaml",
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			name:     "pod with the yml extension",
			filePath: "./testdata/yaml-extension2/test_pod.yml",
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			// file is skipped if no kind is specified or bad
			name:     "yaml with no kind",
			filePath: "./testdata/file-test-data/test_no_kind.yml",
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			// file is skipped if no kind is specified or bad
			name:     "yaml with bad kind",
			filePath: "./testdata/file-test-data/test_bad_kind.yml",
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.k8sV1.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}

}
