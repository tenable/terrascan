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

package tfv12

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	commons_test "github.com/accurics/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/accurics/terrascan/pkg/utils"
)

func TestLoadIacDir(t *testing.T) {
	testErrorString1 := fmt.Sprintf(`failed to load terraform config dir '%s'. error from terraform:
%s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`, testDataDir, emptyTfFilePath, emptyTfFilePath)

	multipleProvidersDir := filepath.Join(testDataDir, "multiple-required-providers")

	testErrorString2 := fmt.Sprintf(`failed to load terraform config dir '%s'. error from terraform:
%s:2,3-21: Duplicate required providers configuration; A module may have only one required providers configuration. The required providers were previously configured at %s:2,3-21.
`, multipleProvidersDir, filepath.Join(multipleProvidersDir, "b.tf"), filepath.Join(multipleProvidersDir, "a.tf"))

	testDirPath1 := "not-there"

	testDirPath2 := filepath.Join(testDataDir, "testfile")

	invalidDirErrStringTemplate := "directory '%s' has no terraform config files"

	table := []struct {
		name    string
		dirPath string
		tfv12   TfV12
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "invalid dirPath",
			dirPath: testDirPath1,
			tfv12:   TfV12{},
			wantErr: fmt.Errorf(invalidDirErrStringTemplate, testDirPath1),
		},
		{
			name:    "empty config",
			dirPath: testDirPath2,
			tfv12:   TfV12{},
			wantErr: fmt.Errorf(invalidDirErrStringTemplate, testDirPath2),
		},
		{
			name:    "incorrect module structure",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv12:   TfV12{},
			wantErr: fmt.Errorf("failed to build terraform allResourcesConfig"),
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			tfv12:   TfV12{},
			wantErr: fmt.Errorf(testErrorString1),
		},
		{
			name:    "load invalid config dir",
			dirPath: multipleProvidersDir,
			tfv12:   TfV12{},
			wantErr: fmt.Errorf(testErrorString2),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name        string
		tfConfigDir string
		tfJSONFile  string
		tfv12       TfV12
		wantErr     error
	}{
		{
			name:        "config1",
			tfConfigDir: filepath.Join(testDataDir, "tfconfigs"),
			tfJSONFile:  filepath.Join(tfJsonDir, "fullconfig.json"),
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "module directory",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			tfJSONFile:  filepath.Join(tfJsonDir, "moduleconfigs.json"),
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "nested module directory",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJsonDir, "deep-modules.json"),
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "variables of list type",
			tfConfigDir: filepath.Join(testDataDir, "list-type-vars-test"),
			tfJSONFile:  filepath.Join(tfJsonDir, "list-vars-test.json"),
			tfv12:       TfV12{},
			wantErr:     nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacDir(tt.tfConfigDir)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			var want output.AllResourceConfigs

			// Read the expected value and unmarshal into want
			contents, _ := ioutil.ReadFile(tt.tfJSONFile)
			if utils.IsWindowsPlatform() {
				contents = utils.ReplaceWinNewLineBytes(contents)
			}

			err := json.Unmarshal(contents, &want)
			if err != nil {
				t.Errorf("unexpected error unmarshalling want: %v", err)
			}

			match, err := commons_test.IdenticalAllResourceConfigs(got, want)
			if err != nil {
				t.Errorf("unexpected error checking result: %v", err)
			}
			if !match {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(want, "", "  ")
				t.Errorf("got '%v', want: '%v'", string(g), string(w))
			}
		})
	}
}
