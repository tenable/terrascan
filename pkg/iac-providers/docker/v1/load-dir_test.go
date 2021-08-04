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

package dockerv1

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

var testDataDir = "testdata"

func TestLoadIacDir(t *testing.T) {
	invalidDirErr := &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "not-there"}
	if utils.IsWindowsPlatform() {
		invalidDirErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: "not-there"}
	}
	errString := fmt.Sprintf("error while parsing file %s", filepath.Join(testDataDir, "valid-directory-with-invalid-file", "Dockerfile"))
	tests := []struct {
		name     string
		dirPath  string
		dockerV1 DockerV1
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "empty config",
			dirPath:  filepath.Join(testDataDir, "testfile"),
			dockerV1: DockerV1{},
			wantErr:  multierror.Append(fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile"))),
		},
		{
			name:     "invalid dirPath",
			dirPath:  "not-there",
			dockerV1: DockerV1{},
			wantErr:  multierror.Append(invalidDirErr),
		},
		{
			name:     "valid dirPath",
			dirPath:  filepath.Join(testDataDir, "valid-directory"),
			dockerV1: DockerV1{},
			wantErr:  nil,
		},
		{
			name:     "valid dirPath with invalid  file",
			dirPath:  filepath.Join(testDataDir, "valid-directory-with-invalid-file"),
			dockerV1: DockerV1{},
			wantErr:  multierror.Append(errors.New(errString)),
		},
		{
			name:     "valid dirPath having dockerfile with in-file instrumentation",
			dirPath:  filepath.Join(testDataDir, "valid-directory-with-in-file-instrumentation"),
			dockerV1: DockerV1{},
			want: output.AllResourceConfigs{
				"docker_cmd": []output.ResourceConfig{
					{ID: "docker_cmd.55ceacedc5f1c0df6951723a7401a74e",
						Name:       "Dockerfile",
						ModuleName: "",
						Source:     "Dockerfile",
						PlanRoot:   "", Line: 5,
						Type:   "docker_cmd",
						Config: "server",
						SkipRules: []output.SkipRule{{Rule: "AWS.S3Bucket.DS.High.1041",
							Comment: "This rule does not belong to dockerfile will add correct once dockerfile policy added."}},
						MaxSeverity: "None",
						MinSeverity: "High"}},
				"docker_dockerfile": []output.ResourceConfig{{ID: "docker_dockerfile.96052d48e5364a05995aaec1e5d53f2d", Name: "Dockerfile", ModuleName: "", Source: "Dockerfile", PlanRoot: "", Line: 1, Type: "docker_dockerfile", Config: []string{"from", "cmd"}, SkipRules: []output.SkipRule{{Rule: "AWS.S3Bucket.DS.High.1041", Comment: "This rule does not belong to dockerfile will add correct once dockerfile policy added."}}, MaxSeverity: "None", MinSeverity: "High"}},
				"docker_from":       []output.ResourceConfig{{ID: "docker_from.68be487d8ad02b4e09b46d29c8dbef3b", Name: "Dockerfile", ModuleName: "", Source: "Dockerfile", PlanRoot: "", Line: 1, Type: "docker_from", Config: "runatlantis/atlantis:v0.16.1", SkipRules: []output.SkipRule{{Rule: "AWS.S3Bucket.DS.High.1041", Comment: "This rule does not belong to dockerfile will add correct once dockerfile policy added."}}, MaxSeverity: "None", MinSeverity: "High"}}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := tt.dockerV1.LoadIacDir(tt.dirPath, false)
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.want != nil {
				if got == nil || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unexpected result; got: '%#v', want: '%v'", got, tt.want)
				}
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
