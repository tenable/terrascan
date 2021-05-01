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
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	commons_test "github.com/accurics/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

func TestLoadIacDir(t *testing.T) {
	var nilMultiErr *multierror.Error = nil

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

	pathErr := &os.PathError{Op: "lstat", Path: "not-there", Err: syscall.ENOENT}
	if utils.IsWindowsPlatform() {
		pathErr = &os.PathError{Op: "CreateFile", Path: "not-there", Err: syscall.ENOENT}
	}

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
			tfv12:   TfV12{true},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath1)),
		},
		{
			name:    "invalid dirPath recursive",
			dirPath: testDirPath1,
			tfv12:   TfV12{false},
			wantErr: multierror.Append(pathErr),
		},
		{
			name:    "empty config",
			dirPath: testDirPath2,
			tfv12:   TfV12{true},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath2)),
		},
		{
			name:    "empty config recursive",
			dirPath: testDirPath2,
			tfv12:   TfV12{false},
			wantErr: nilMultiErr,
		},
		{
			name:    "incorrect module structure",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv12:   TfV12{true},
			wantErr: multierror.Append(fmt.Errorf("failed to build terraform allResourcesConfig")),
		},
		{
			name:    "incorrect module structure recursive",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv12:   TfV12{false},
			wantErr: nilMultiErr,
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			tfv12:   TfV12{true},
			wantErr: multierror.Append(fmt.Errorf(testErrorString1)),
		},
		{
			name:    "load invalid config dir recursive",
			dirPath: testDataDir,
			tfv12:   TfV12{false},
			wantErr: nilMultiErr,
		},
		{
			name:    "load invalid config dir",
			dirPath: multipleProvidersDir,
			tfv12:   TfV12{true},
			wantErr: multierror.Append(fmt.Errorf(testErrorString2)),
		},
		{
			name:    "load invalid config dir recursive",
			dirPath: multipleProvidersDir,
			tfv12:   TfV12{false},
			wantErr: nilMultiErr,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacDir(tt.dirPath)
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.wantErr == nilMultiErr {
				if err := me.ErrorOrNil(); err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else if me.Error() != tt.wantErr.Error() {
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
			tfJSONFile:  filepath.Join(tfJSONDir, "fullconfig.json"),
			tfv12:       TfV12{true},
			wantErr:     nilMultiErr,
		},
		{
			name:        "config1 recursive",
			tfConfigDir: filepath.Join(testDataDir, "tfconfigs"),
			// no change in the output expected as the config dir doesn't contain subfolder
			tfJSONFile: filepath.Join(tfJSONDir, "fullconfig.json"),
			tfv12:      TfV12{false},
			wantErr:    nilMultiErr,
		},
		{
			name:        "module directory",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			tfJSONFile:  filepath.Join(tfJSONDir, "moduleconfigs.json"),
			tfv12:       TfV12{true},
			wantErr:     nilMultiErr,
		},
		{
			name:        "module directory recursive",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			// no change in the output expected as the config dir doesn't contain subfolder
			tfJSONFile: filepath.Join(tfJSONDir, "moduleconfigs.json"),
			tfv12:      TfV12{false},
			wantErr:    nilMultiErr,
		},
		{
			name:        "nested module directory",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules.json"),
			tfv12:       TfV12{true},
			wantErr:     nilMultiErr,
		},
		{
			name:        "nested module directory recursive",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules-recursive.json"),
			tfv12:       TfV12{false},
			wantErr:     nilMultiErr,
		},
		{
			name:        "variables of list type",
			tfConfigDir: filepath.Join(testDataDir, "list-type-vars-test"),
			tfJSONFile:  filepath.Join(tfJSONDir, "list-vars-test.json"),
			tfv12:       TfV12{true},
			wantErr:     nilMultiErr,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacDir(tt.tfConfigDir)
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.wantErr == nilMultiErr {
				if err := me.ErrorOrNil(); err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else if me.Error() != tt.wantErr.Error() {
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
