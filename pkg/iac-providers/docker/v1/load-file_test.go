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

package dockerv1

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

var multiStageDockerfileConfig = output.AllResourceConfigs{
	"docker_copy": []output.ResourceConfig{
		{
			ID:          "docker_copy.e7031c0312fd3a6a88ac071b8ab508d0",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        7,
			Type:        "docker_copy",
			Config:      "--from=builder /go/bin/terrascan /go/bin/terrascan",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
	"docker_dockerfile": []output.ResourceConfig{
		{
			ID:         "docker_dockerfile.37b4650bc01b676c9d26b3596fc1e6bd",
			Name:       "dockerfile-withmultiple-stages",
			ModuleName: "",
			Source:     "dockerfile-withmultiple-stages",
			PlanRoot:   "",
			Line:       1,
			Type:       "docker_dockerfile",
			Config: []string{"from",
				"run",
				"from",
				"copy",
				"entrypoint"},
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
	"docker_entrypoint": []output.ResourceConfig{
		{
			ID:          "docker_entrypoint.21a71774cdc858e44224d9d490498d49",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        8,
			Type:        "docker_entrypoint",
			Config:      "/go/bin/main",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
	"docker_from": []output.ResourceConfig{
		{
			ID:          "docker_from.273cb3c947150fa2365b39346e207035",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        2,
			Type:        "docker_from",
			Config:      "golang:alpine AS builder",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
			ContainerImages: []output.ContainerDetails{
				{
					Image: "golang:alpine",
				},
			},
		},
		{
			ID:          "docker_from.aaa14f2bb7549c35cdb047282de7e26b",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        6,
			Type:        "docker_from",
			Config:      "alpine:3.14",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
			ContainerImages: []output.ContainerDetails{
				{
					Image: "alpine:3.14",
				},
			},
		},
	},
	"docker_run": []output.ResourceConfig{
		{
			ID:          "docker_run.556176d13b816800a50fb2998cac92ec",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        3,
			Type:        "docker_run",
			Config:      "go build main.go",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
}

func TestLoadIacFile(t *testing.T) {

	tests := []struct {
		name        string
		absFilePath string
		options     map[string]interface{}
		dockerV1    DockerV1
		want        output.AllResourceConfigs
		wantErr     error
		typeOnly    bool
	}{
		{
			name:        "empty config file",
			absFilePath: filepath.Join(fileTestDataDir, "Dockerfile"),
			dockerV1:    DockerV1{},
			wantErr:     fmt.Errorf("error while parsing dockerfile %s, error: file with no instructions", filepath.Join(fileTestDataDir, "Dockerfile")),
		},
		{
			name:        "valid docker file",
			absFilePath: filepath.Join(fileTestDataDir, "valid-Dockerfile"),
			dockerV1:    DockerV1{},
			wantErr:     nil,
		},
		{
			name:        "docker file with multiple stages",
			absFilePath: filepath.Join(fileTestDataDir, "dockerfile-withmultiple-stages"),
			dockerV1:    DockerV1{},
			wantErr:     nil,
			want:        multiStageDockerfileConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, gotErr := tt.dockerV1.LoadIacFile(tt.absFilePath, tt.options)
			if tt.want != nil {
				if got == nil || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("unexpected result; got: '%#v', want: '%v'", got, tt.want)
				}
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
