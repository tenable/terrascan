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

package cli

import (
	"testing"
)

func TestRun(t *testing.T) {
	table := []struct {
		name        string
		iacType     string
		iacVersion  string
		cloudType   []string
		iacFilePath string
		iacDirPath  string
		configFile  string
		configOnly  bool
		stdOut      string
		want        string
		wantErr     error
	}{
		{
			name:       "normal terraform run",
			cloudType:  []string{"terraform"},
			iacDirPath: "testdata/run-test",
		},
		{
			name:       "normal k8s run",
			cloudType:  []string{"k8s"},
			iacDirPath: "testdata/run-test",
		},
		{
			name:        "config-only flag terraform",
			cloudType:   []string{"terraform"},
			iacFilePath: "testdata/run-test/config-only.tf",
			configOnly:  true,
		},
		{
			name:        "config-only flag k8s",
			cloudType:   []string{"k8s"},
			iacFilePath: "testdata/run-test/config-only.yaml",
			configOnly:  true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.iacType, tt.iacVersion, tt.cloudType, tt.iacFilePath, tt.iacDirPath, tt.configFile, []string{}, "", "", "", tt.configOnly, false)
		})
	}
}
