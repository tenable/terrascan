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

package k8sv1

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var testDataDir = "testdata"
var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name     string
		filePath string
		options  map[string]interface{}
		k8sV1    K8sV1
		typeOnly bool
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			// file is skipped if no kind is specified or bad
			name:     "empty config file",
			filePath: filepath.Join(fileTestDataDir, "empty-file.yaml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			name:     "yaml with multiple documents",
			filePath: filepath.Join(testDataDir, "yaml-with-multiple-documents", "test_pod.yaml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			name:     "pod with the yml extension",
			filePath: filepath.Join(testDataDir, "yaml-extension2", "test_pod.yml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			// file is skipped if no kind is specified or bad
			name:     "yaml with no kind",
			filePath: filepath.Join(fileTestDataDir, "test_no_kind.yml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			// file is skipped if no kind is specified or bad
			name:     "yaml with bad kind",
			filePath: filepath.Join(fileTestDataDir, "test_bad_kind.yml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
		{
			name:     "file with skip rules in annotations",
			filePath: filepath.Join(fileTestDataDir, "test_pod_skip_rules.yaml"),
			k8sV1:    K8sV1{},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.k8sV1.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}

}

func Test_getSourceRelativePath(t *testing.T) {
	dir1, dir2 := "Dir1", "Dir2"
	sourcePath1 := filepath.Join(dir1, dir2, "filename.yaml")

	type args struct {
		absRootDir string
		sourcePath string
	}
	tests := []struct {
		name            string
		expectedRelPath string
		args            args
	}{
		{
			name: "empty root directory",
			args: args{
				sourcePath: sourcePath1,
			},
			expectedRelPath: "filename.yaml",
		},
		{
			name: "root directory not empty",
			args: args{
				absRootDir: dir1,
				sourcePath: sourcePath1,
			},
			expectedRelPath: filepath.Join(dir2, "filename.yaml"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &K8sV1{absRootDir: tt.args.absRootDir}
			gotRelPath := k.getSourceRelativePath(tt.args.sourcePath)
			if gotRelPath != tt.expectedRelPath {
				t.Errorf("Test_getSourceRelativePath() = unexpected relative path; want relPath: %s, got relPath: %s", tt.expectedRelPath, gotRelPath)
			}
		})
	}
}
