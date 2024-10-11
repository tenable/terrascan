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
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

var multiStageDockerfileConfig = output.AllResourceConfigs{
	"docker_COPY": []output.ResourceConfig{
		{
			ID:          "docker_COPY.e7031c0312fd3a6a88ac071b8ab508d0",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        7,
			Type:        "docker_COPY",
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
			Config: []string{
				"FROM",
				"RUN",
				"FROM",
				"COPY",
				"ENTRYPOINT"},
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
	"docker_ENTRYPOINT": []output.ResourceConfig{
		{
			ID:          "docker_ENTRYPOINT.21a71774cdc858e44224d9d490498d49",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        8,
			Type:        "docker_ENTRYPOINT",
			Config:      "/go/bin/main",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
	"docker_FROM": []output.ResourceConfig{
		{
			ID:              "docker_FROM.273cb3c947150fa2365b39346e207035",
			Name:            "dockerfile-withmultiple-stages",
			ModuleName:      "",
			Source:          "dockerfile-withmultiple-stages",
			PlanRoot:        "",
			Line:            2,
			Type:            "docker_FROM",
			Config:          "golang:alpine AS builder",
			SkipRules:       []output.SkipRule(nil),
			MaxSeverity:     "",
			MinSeverity:     "",
			ContainerImages: []output.ContainerDetails{},
		},
		{
			ID:              "docker_FROM.aaa14f2bb7549c35cdb047282de7e26b",
			Name:            "dockerfile-withmultiple-stages",
			ModuleName:      "",
			Source:          "dockerfile-withmultiple-stages",
			PlanRoot:        "",
			Line:            6,
			Type:            "docker_FROM",
			Config:          "alpine:3.14",
			SkipRules:       []output.SkipRule(nil),
			MaxSeverity:     "",
			MinSeverity:     "",
			ContainerImages: []output.ContainerDetails{},
		},
	},
	"docker_RUN": []output.ResourceConfig{
		{
			ID:          "docker_RUN.556176d13b816800a50fb2998cac92ec",
			Name:        "dockerfile-withmultiple-stages",
			ModuleName:  "",
			Source:      "dockerfile-withmultiple-stages",
			PlanRoot:    "",
			Line:        3,
			Type:        "docker_RUN",
			Config:      "go build main.go",
			SkipRules:   []output.SkipRule(nil),
			MaxSeverity: "",
			MinSeverity: "",
		},
	},
}

func TestLoadIacFile(t *testing.T) {
	RegisterFailHandler(Fail)
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
				gotdata, _ := json.MarshalIndent(got, "", "  ")
				wantData, _ := json.MarshalIndent(tt.want, "", "  ")
				fmt.Println(string(gotdata))
				fmt.Println(string(wantData))
				Expect(string(gotdata)).To(MatchJSON(wantData))

			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
