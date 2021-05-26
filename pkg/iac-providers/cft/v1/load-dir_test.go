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

package cftv1

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacDir(t *testing.T) {
	table := []struct {
		name    string
		dirPath string
		cftv1   CFTV1
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "empty config",
			dirPath: "./testdata/testfile",
			cftv1:   CFTV1{},
			wantErr: fmt.Errorf("no directories found for path ./testdata/testfile"),
		},
		{
			name:    "load invalid config dir",
			dirPath: "./testdata",
			cftv1:   CFTV1{},
			wantErr: nil,
		},
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			cftv1:   CFTV1{},
			wantErr: &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "not-there"},
		},
		{
			name:    "key-vault",
			dirPath: "./testdata/key-vault",
			cftv1:   CFTV1{},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.cftv1.LoadIacDir(tt.dirPath, false)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
