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

package runtime

import (
	"reflect"
	"testing"
)

func TestValidateInputs(t *testing.T) {

	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr: nil,
		},
		{
			name: "valid dirPath",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr: nil,
		},
		{
			name: "empty iac path",
			executor: Executor{
				filePath: "",
				dirPath:  "",
			},
			wantErr: errEmptyIacPath,
		},
		{
			name: "filepath does not exist",
			executor: Executor{
				filePath: "./testdata/notthere",
			},
			wantErr: errFileNotExists,
		},
		{
			name: "directory does not exist",
			executor: Executor{
				dirPath: "./testdata/notthere",
			},
			wantErr: errDirNotExists,
		},
		{
			name: "invalid iac type",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  []string{"aws"},
				iacType:    "notthere",
				iacVersion: "v12",
			},
			wantErr: errIacNotSupported,
		},
		{
			name: "invalid iac version",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  []string{"aws"},
				iacType:    "terraform",
				iacVersion: "notthere",
			},
			wantErr: errIacNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.ValidateInputs()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error, gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
