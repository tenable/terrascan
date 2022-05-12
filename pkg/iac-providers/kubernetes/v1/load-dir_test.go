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
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
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
		options map[string]interface{}
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "empty config",
			dirPath: filepath.Join(testDataDir, "testfile"),
			k8sV1:   K8sV1{},
			wantErr: multierror.Append(fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile"))),
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
			wantErr: multierror.Append(invalidDirErr),
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
			_, gotErr := tt.k8sV1.LoadIacDir(tt.dirPath, tt.options)
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.wantErr == nil {
				if err := me.ErrorOrNil(); err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else if me.Error() != tt.wantErr.Error() {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

}
