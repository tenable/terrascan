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

package tfv15

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	commons_test "github.com/tenable/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/tenable/terrascan/pkg/utils"
)

var testInvalidAttribdir = "terraform-with-attrib-errors/module/firewall-test"
var unsupportedArgument = `Unsupported argument; An argument named "replace_triggered_by" is not expected here`

func TestLoadIacDir(t *testing.T) {
	var nilMultiErr *multierror.Error = nil

	testErrorMessage := fmt.Sprintf(`diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s:1,21-2,1: Invalid block definition; A block definition must have block content delimited by "{" and "}", starting on the same line as the block header.
%s:1,1-5: Unsupported block type; Blocks of type "some" are not expected here.
`, testDataDir, emptyTfFilePath, emptyTfFilePath)

	errStringInvalidModuleConfigs := fmt.Sprintf(`failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`, filepath.Join(testDataDir, "invalid-moduleconfigs", "cloudfront", "sub-cloudfront"))

	errStringDependsOnDir := fmt.Sprintf(`failed to build unified config. errors:
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
<nil>: Failed to read module directory; Module directory %s does not exist or cannot be read.
`, filepath.Join(testDataDir, "depends_on", "live", "log"), filepath.Join(testDataDir, "depends_on", "live", "security"))

	errStringModuleSourceInvalid := fmt.Sprintf(`failed to build unified config. errors:
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module cloudfront
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module m1
`, filepath.Join(testDataDir, "invalid-module-source"), filepath.Join(testDataDir, "invalid-module-source"))

	errStringUnifiedInvalidattrib := fmt.Sprintf(`failed to build unified config. errors:
%s/firewall-module.tf:7,5-25: %s.
%s/firewall-module.tf:25,5-25: %s.
`, filepath.Join(testDataDir, testInvalidAttribdir), unsupportedArgument, filepath.Join(testDataDir, testInvalidAttribdir), unsupportedArgument)

	errDiagnostincMessageAttrib := fmt.Sprintf(`diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s/firewall-module.tf:7,5-25: %s.
%s/firewall-module.tf:25,5-25: %s.
`, filepath.Join(testDataDir, testInvalidAttribdir), filepath.Join(testDataDir, testInvalidAttribdir), unsupportedArgument, filepath.Join(testDataDir, testInvalidAttribdir), unsupportedArgument)

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
		tfv15   TfV15
		want    output.AllResourceConfigs
		wantErr error
		options map[string]interface{}
	}{
		{
			name:    "invalid dirPath",
			dirPath: testDirPath1,
			tfv15:   TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath1)),
		},
		{
			name:    "invalid dirPath recursive",
			dirPath: testDirPath1,
			tfv15:   TfV15{},
			wantErr: multierror.Append(pathErr),
		},
		{
			name:    "empty config",
			dirPath: testDirPath2,
			tfv15:   TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath2)),
		},
		{
			name:    "empty config recursive",
			dirPath: testDirPath2,
			tfv15:   TfV15{},
			wantErr: nilMultiErr,
		},
		{
			name:    "incorrect module structure",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv15:   TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf("failed to build terraform allResourcesConfig")),
		},
		{
			name:    "incorrect module structure recursive",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv15:   TfV15{},
			// same error is loaded two times because, both root module and a child module will generated same error
			wantErr: multierror.Append(fmt.Errorf(errStringInvalidModuleConfigs), fmt.Errorf(errStringInvalidModuleConfigs)),
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			tfv15:   TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf(testErrorMessage)),
		},
		{
			name:    "load invalid config dir recursive",
			dirPath: testDataDir,
			tfv15:   TfV15{},
			wantErr: multierror.Append(fmt.Errorf(testErrorMessage),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules")),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules", "m4", "modules")),
				fmt.Errorf(errStringDependsOnDir),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "invalid-module-source")),
				fmt.Errorf(errStringModuleSourceInvalid),
				fmt.Errorf(errStringInvalidModuleConfigs),
				fmt.Errorf(errStringInvalidModuleConfigs),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "relative-moduleconfigs")),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "terraform-with-attrib-errors")),           // 9 good
				fmt.Errorf(errStringUnifiedInvalidattrib),                                                                     // 10 good
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "terraform-with-attrib-errors", "module")), // 11 good
				fmt.Errorf(errDiagnostincMessageAttrib),                                                                       // 12
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "tfjson")),                                 // 13 good
			),
		},
		{
			name:    "invalid module source directory",
			dirPath: filepath.Join(testDataDir, "invalid-module-source", "invalid_source"),
			tfv15:   TfV15{},
			wantErr: multierror.Append(fmt.Errorf(errStringModuleSourceInvalid)),
		},
		{
			name:    "provider block with only alias",
			dirPath: filepath.Join(testDataDir, "provider-with-only-alias"),
			tfv15:   TfV15{},
			wantErr: nilMultiErr,
		},
		{
			name:    "required provider with configuration aliases",
			dirPath: filepath.Join(testDataDir, "required-providers", "configuration-alias"),
			tfv15:   TfV15{},
			wantErr: nilMultiErr,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv15.LoadIacDir(tt.dirPath, tt.options)
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

	tfJSONDir := filepath.Join(testDataDir, "tfjson")
	nestedModuleErr1 := fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules"))
	nestedModuleErr2 := fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules", "m4", "modules"))

	table2 := []struct {
		name        string
		tfConfigDir string
		tfJSONFile  string
		tfv15       TfV15
		wantErr     error
		options     map[string]interface{}
	}{
		{
			name:        "config1",
			tfConfigDir: filepath.Join(testDataDir, "tfconfigs"),
			tfJSONFile:  filepath.Join(tfJSONDir, "fullconfig.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "config1 recursive",
			tfConfigDir: filepath.Join(testDataDir, "tfconfigs"),
			// no change in the output expected as the config dir doesn't contain subfolder
			tfJSONFile: filepath.Join(tfJSONDir, "fullconfig.json"),
			tfv15:      TfV15{},
			wantErr:    nilMultiErr,
		},
		{
			name:        "module directory",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			tfJSONFile:  filepath.Join(tfJSONDir, "moduleconfigs.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "module directory recursive",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			// no change in the output expected as the config dir doesn't contain subfolder
			tfJSONFile: filepath.Join(tfJSONDir, "moduleconfigs.json"),
			tfv15:      TfV15{},
			wantErr:    nilMultiErr,
		},
		{
			name:        "nested module directory",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "nested module directory recursive",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules-recursive.json"),
			tfv15:       TfV15{},
			wantErr:     multierror.Append(nestedModuleErr1, nestedModuleErr2),
		},
		{
			name:        "complex variables",
			tfConfigDir: filepath.Join(testDataDir, "complex-variables"),
			tfJSONFile:  filepath.Join(tfJSONDir, "complex-variables.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "recursive loop while resolving variables",
			tfConfigDir: filepath.Join(testDataDir, "recursive-loop-variables"),
			tfJSONFile:  filepath.Join(tfJSONDir, "recursive-loop-variables.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "recursive loop while resolving locals",
			tfConfigDir: filepath.Join(testDataDir, "recursive-loop-locals"),
			tfJSONFile:  filepath.Join(tfJSONDir, "recursive-loop-locals.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "Invalid attribute error should be ignored and should return resources",
			tfConfigDir: filepath.Join(testDataDir, "terraform-with-attrib-errors/firewall"),
			tfJSONFile:  filepath.Join(tfJSONDir, "firewall-with-attrib-error.json"),
			tfv15:       TfV15{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf("failed to build terraform allResourcesConfig")),
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv15.LoadIacDir(tt.tfConfigDir, tt.options)
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
			contents, _ := os.ReadFile(tt.tfJSONFile)
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
