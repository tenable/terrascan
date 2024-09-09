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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestParse(t *testing.T) {
	RegisterFailHandler(Fail)
	tests := []struct {
		name     string
		filePath string
		dockerv1 DockerV1
		want     []DockerConfig
		wantErr  error
	}{
		{
			name:     "valid docker file",
			filePath: filepath.Join(fileTestDataDir, "dockerfile-testparse-function"),
			dockerv1: DockerV1{},
			wantErr:  nil,
			want: []DockerConfig{
				{Cmd: "FROM", Value: "runatlantis/atlantis:v0.16.1", Line: 1},
				{Cmd: "MAINTAINER", Value: "tenable", Line: 2},
				{Cmd: "LABEL", Value: "key \"value\"", Line: 3},
				{Cmd: "WORKDIR", Value: "test", Line: 4},
				{Cmd: "ENV", Value: "DEFAULT_TERRASCAN_VERSION 1.5.1", Line: 5},
				{Cmd: "ENV", Value: "PLANFILE tfplan", Line: 6},
				{Cmd: "ADD", Value: "setup.sh terrascan.sh launch-atlantis.sh entrypoint.sh /usr/local/bin/", Line: 7},
				{Cmd: "RUN", Value: "mkdir -p /etc/atlantis/ &&     chmod +x /usr/local/bin/*.sh &&     /usr/local/bin/setup.sh", Line: 8},
				{Cmd: "Copy", Value: "terrascan-workflow.yaml /etc/atlantis/workflow.yaml", Line: 11},
				{Cmd: "USER", Value: "atlantis", Line: 13},
				{Cmd: "ARG", Value: "name=defaultValue", Line: 14},
				{Cmd: "RUN", Value: "terrascan init", Line: 15},
				{Cmd: "VOLUME", Value: "/temp", Line: 16},
				{Cmd: "HEALTHCHECK", Value: "--interval=30s --timeout=30s --start-period=5s --retries=3 CMD executable", Line: 17},
				{Cmd: "ENTRYPOINT", Value: "/bin/bash entrypoint.sh", Line: 18},
				{Cmd: "SHELL", Value: "cd", Line: 19},
				{Cmd: "ONBUILD", Value: "", Line: 20},
				{Cmd: "expose", Value: "9090", Line: 21},
				{Cmd: "STOPSIGNAL", Value: "1", Line: 22},
				{Cmd: "CMD", Value: "server", Line: 23}},
		},
		{
			name:     "invalid  docker file path",
			filePath: filepath.Join(fileTestDataDir, "dockerfile-testparse-function1"),
			dockerv1: DockerV1{},
			want:     []DockerConfig{},
			wantErr:  fmt.Errorf("open %s: no such file or directory", filepath.Join(fileTestDataDir, "dockerfile-testparse-function1")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _, err := tt.dockerv1.Parse(tt.filePath)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", err, tt.wantErr)
					return
				}
			} else if err.Error() != tt.wantErr.Error() {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DockerV1.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
