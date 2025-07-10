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

package tfv12

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/iac-providers/terraform/commons"
	commons_test "github.com/tenable/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/tenable/terrascan/pkg/utils"
)

const errStrLoadingConfigDir = `diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s:2,3-21: Duplicate required providers configuration; A module may have only one required providers configuration. The required providers were previously configured at %s:2,3-21.
`

const errStrLoadingConfigDir1 = `diagnostic errors while loading terraform config dir '%s'. error from terraform:
%s:8,12-22: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:42,15-35: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:39,14-24: Invalid reference from destroy provisioner; Destroy-time provisioners and their connection configurations may only reference attributes of the related resource, via 'self', 'count.index', or 'each.key'.

References to other resources during the destroy phase can cause dependency cycles and interact poorly with create_before_destroy.
%s:34,1-29: Duplicate resource "null_resource" configuration; A null_resource resource named "b" was already declared at testdata/destroy-provisioners/main.tf:23,1-29. Resource names must be unique per type in each module.
`

const errStrFiledBuildingUnifiedConfig = `failed to build unified config. errors:
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module cloudfront
<nil>: Invalid module config directory; Module directory '%s' has no terraform config files for module m1
`

func TestLoadIacDir(t *testing.T) {
	var nilMultiErr *multierror.Error = nil

	destroyProvisionersDir := filepath.Join(testDataDir, "destroy-provisioners")
	destroyProvisionersMainFile := filepath.Join(destroyProvisionersDir, "main.tf")

	testErrorString1 := commons.GenerateTerraformLoadError(testDataDir, emptyTfFilePath, emptyTfFilePath)

	multipleProvidersDir := filepath.Join(testDataDir, "multiple-required-providers")

	testErrorString2 := fmt.Errorf(errStrLoadingConfigDir, multipleProvidersDir, filepath.Join(multipleProvidersDir, "b.tf"), filepath.Join(multipleProvidersDir, "a.tf"))

	errStringInvalidModuleConfigs := commons.GenerateInvalidModuleConfigError(testDataDir)

	errStringDestroyProvisioners := fmt.Errorf(errStrLoadingConfigDir1, destroyProvisionersDir, destroyProvisionersMainFile, destroyProvisionersMainFile, destroyProvisionersMainFile, destroyProvisionersMainFile)

	errStringModuleSourceInvalid := fmt.Errorf(errStrFiledBuildingUnifiedConfig, filepath.Join(testDataDir, "invalid-module-source"), filepath.Join(testDataDir, "invalid-module-source"))

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
		options map[string]interface{}
	}{
		{
			name:    "invalid dirPath",
			dirPath: testDirPath1,
			tfv12:   TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath1)),
		},
		{
			name:    "invalid dirPath recursive",
			dirPath: testDirPath1,
			tfv12:   TfV12{},
			wantErr: multierror.Append(pathErr),
		},
		{
			name:    "empty config",
			dirPath: testDirPath2,
			tfv12:   TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf(invalidDirErrStringTemplate, testDirPath2)),
		},
		{
			name:    "empty config recursive",
			dirPath: testDirPath2,
			tfv12:   TfV12{},
			wantErr: nilMultiErr,
		},
		{
			name:    "incorrect module structure",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv12:   TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(fmt.Errorf("failed to build terraform allResourcesConfig")),
		},
		{
			name:    "incorrect module structure recursive",
			dirPath: filepath.Join(testDataDir, "invalid-moduleconfigs"),
			tfv12:   TfV12{},
			// same error is loaded two times because, both root module and a child module will generated same error
			wantErr: multierror.Append(errStringInvalidModuleConfigs, errStringInvalidModuleConfigs),
		},
		{
			name:    "load invalid config dir",
			dirPath: testDataDir,
			tfv12:   TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(testErrorString1),
		},
		{
			name:    "load invalid config dir recursive",
			dirPath: testDataDir,
			tfv12:   TfV12{},
			wantErr: multierror.Append(testErrorString1,
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules")),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules", "m4", "modules")),
				errStringDestroyProvisioners,
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "invalid-module-source")),
				errStringModuleSourceInvalid,
				errStringInvalidModuleConfigs,
				errStringInvalidModuleConfigs,
				testErrorString2,
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "relative-moduleconfigs")),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "tfjson")),
			),
		},
		{
			name:    "load multiple provider config dir",
			dirPath: multipleProvidersDir,
			tfv12:   TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: multierror.Append(testErrorString2),
		},
		{
			name:    "load multiple provider config dir recursive",
			dirPath: multipleProvidersDir,
			tfv12:   TfV12{},
			wantErr: multierror.Append(testErrorString2),
		},
		{
			name:    "invalid module source directory",
			dirPath: filepath.Join(testDataDir, "invalid-module-source", "invalid_source"),
			tfv12:   TfV12{},
			wantErr: multierror.Append(errStringModuleSourceInvalid),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacDir(tt.dirPath, tt.options)
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

	nestedModuleErr1 := fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules"))
	nestedModuleErr2 := fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "deep-modules", "modules", "m4", "modules"))

	table2 := []struct {
		name        string
		tfConfigDir string
		tfJSONFile  string
		tfv12       TfV12
		options     map[string]interface{}
		wantErr     error
	}{
		{
			name:        "config1",
			tfConfigDir: filepath.Join(testDataDir, "tfconfigs"),
			tfJSONFile:  filepath.Join(tfJSONDir, "fullconfig.json"),
			tfv12:       TfV12{},
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
			tfv12:      TfV12{},
			wantErr:    nilMultiErr,
		},
		{
			name:        "module directory",
			tfConfigDir: filepath.Join(testDataDir, "moduleconfigs"),
			tfJSONFile:  filepath.Join(tfJSONDir, "moduleconfigs.json"),
			tfv12:       TfV12{},
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
			tfv12:      TfV12{},
			wantErr:    nilMultiErr,
		},
		{
			name:        "nested module directory",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules.json"),
			tfv12:       TfV12{},
			options: map[string]interface{}{
				"nonRecursive": true,
			},
			wantErr: nilMultiErr,
		},
		{
			name:        "variables of list type",
			tfConfigDir: filepath.Join(testDataDir, "list-type-vars-test"),
			tfJSONFile:  filepath.Join(tfJSONDir, "list-vars-test.json"),
			tfv12:       TfV12{},
			wantErr:     nilMultiErr,
		},
		{
			name:        "nested module directory recursive",
			tfConfigDir: filepath.Join(testDataDir, "deep-modules"),
			tfJSONFile:  filepath.Join(tfJSONDir, "deep-modules-recursive.json"),
			tfv12:       TfV12{},
			wantErr:     multierror.Append(nestedModuleErr1, nestedModuleErr2),
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacDir(tt.tfConfigDir, tt.options)
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
