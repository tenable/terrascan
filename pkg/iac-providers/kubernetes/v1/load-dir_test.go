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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

func TestLoadIacDir(t *testing.T) {

	invalidDirErr := &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "not-there"}
	if utils.IsWindowsPlatform() {
		invalidDirErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: "not-there"}
	}

	table := []struct {
		name    string
		dirPath string
		k8sV1   K8sV1
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "empty config",
			dirPath: filepath.Join(testDataDir, "testfile"),
			k8sV1:   K8sV1{},
			wantErr: fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile")),
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			k8sV1:   K8sV1{},
			wantErr: nil,
		},
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			k8sV1:   K8sV1{},
			wantErr: invalidDirErr,
		},
		{
			name:    "yaml with multiple documents",
			dirPath: filepath.Join(testDataDir, "yaml-with-multiple-documents"),
			k8sV1:   K8sV1{},
			wantErr: nil,
		},
		{
			name:    "pod with the yml extension",
			dirPath: filepath.Join(testDataDir, "yaml-extension2"),
			k8sV1:   K8sV1{},
			wantErr: nil,
		},
		{
			name:    "yaml with no kind",
			dirPath: filepath.Join(testDataDir, "yaml-extension2"),
			k8sV1:   K8sV1{},
			wantErr: nil,
		},
		{
			name:    "pod with the json extension",
			dirPath: filepath.Join(testDataDir, "json-extension"),
			k8sV1:   K8sV1{},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.k8sV1.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

}

func TestMakeSourcePathRelative(t *testing.T) {
	dir1, dir2 := "Dir1", "Dir2"
	sourcePath1 := filepath.Join(dir1, dir2, "filename.yaml")
	sourcePath2 := filepath.Join(dir1, "someDir", "test.yaml")

	testResourceConfigs := []output.ResourceConfig{
		{
			Source: sourcePath1,
		},
		{
			Source: sourcePath2,
		},
	}

	type args struct {
		absRootDir      string
		resourceConfigs []output.ResourceConfig
	}
	tests := []struct {
		name                 string
		expectedSourceValues []string
		args                 args
	}{
		{
			name:                 "test to verify path becomes relative",
			expectedSourceValues: []string{filepath.Join(dir2, "filename.yaml"), filepath.Join("someDir", "test.yaml")},
			args: args{
				absRootDir:      dir1,
				resourceConfigs: testResourceConfigs,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			makeSourcePathRelative(tt.args.absRootDir, tt.args.resourceConfigs)
			updatedSourceValues := []string{tt.args.resourceConfigs[0].Source, tt.args.resourceConfigs[1].Source}
			if !utils.IsSliceEqual(tt.expectedSourceValues, updatedSourceValues) {
				t.Errorf("expected source values %v, got %v", tt.expectedSourceValues, updatedSourceValues)
			}
		})
	}
}
