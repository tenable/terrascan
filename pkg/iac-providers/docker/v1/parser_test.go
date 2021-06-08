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
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	tests := []struct {
		name     string
		filePath string
		dockerv1 DockerV1
		want     DockerConfig
		wantErr  error
	}{
		{
			name:     "valid docker file",
			filePath: filepath.Join(fileTestDataDir, "dockerfile-testparse-function"),
			dockerv1: DockerV1{},
			wantErr:  nil,
			want:     DockerConfig{Args: []string{"name=defaultValue"}, Cmd: []string{"server"}, From: []string{"runatlantis/atlantis:v0.16.1"}, Labels: []string{"key \"value\""}, Run: []string{"mkdir -p /etc/atlantis/ &&     chmod +x /usr/local/bin/*.sh &&     /usr/local/bin/setup.sh", "terrascan init"}, Expose: []string{"9090"}, Env: []string{"DEFAULT_TERRASCAN_VERSION 1.5.1", "PLANFILE tfplan"}, Add: []string{"setup.sh terrascan.sh launch-atlantis.sh entrypoint.sh /usr/local/bin/"}, Copy: []string{"terrascan-workflow.yaml /etc/atlantis/workflow.yaml"}, Entrypoint: []string{"/bin/bash entrypoint.sh"}, Volume: []string{"/temp"}, User: []string{"atlantis"}, WorkDir: []string{"test"}, Onbuild: []string{""}, Maintainer: []string{"accurics"}, HealthCheck: []string{"CMD executable"}, Shell: []string{"cd"}, StopSignal: []string{"1"}},
		},
		{
			name:     "invalid  docker file path",
			filePath: filepath.Join(fileTestDataDir, "dockerfile-testparse-function1"),
			dockerv1: DockerV1{},
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
