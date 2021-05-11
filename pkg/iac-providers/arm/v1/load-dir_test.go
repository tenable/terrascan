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

package armv1

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/hashicorp/go-multierror"
)

func TestLoadIacDir(t *testing.T) {
	table := []struct {
		name    string
		dirPath string
		armv1   ARMV1
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "empty config",
			dirPath: filepath.Join(testDataDir, "testfile"),
			armv1:   ARMV1{},
			wantErr: multierror.Append(fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile"))),
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			armv1:   ARMV1{},
			wantErr: nil,
		},
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			armv1:   ARMV1{},
			wantErr: &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "not-there"},
		},
		{
			name:    "key-vault",
			dirPath: "./testdata/key-vault",
			armv1:   ARMV1{},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.armv1.LoadIacDir(tt.dirPath, false)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
