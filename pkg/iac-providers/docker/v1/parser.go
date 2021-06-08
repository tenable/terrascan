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
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/command"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"go.uber.org/zap"
)

// DockerConfig holds configuration of dockerfile
type DockerConfig struct {
	Args        []string `json:"args"`
	Cmd         []string `json:"cmd"`
	From        []string `json:"from"`
	Labels      []string `json:"labels"`
	Run         []string `json:"run"`
	Expose      []string `json:"expose"`
	Env         []string `json:"env"`
	Add         []string `json:"add"`
	Copy        []string `json:"copy"`
	Entrypoint  []string `json:"entrypoint"`
	Volume      []string `json:"volume"`
	User        []string `json:"user"`
	WorkDir     []string `json:"work_dir"`
	Onbuild     []string `json:"onBuild"`
	Maintainer  []string `json:"maintainer"`
	HealthCheck []string `json:"healthCheck"`
	Shell       []string `json:"shell"`
	StopSignal  []string `json:"stopSignal"`
}

const (
	stringJoinCharacter = " "
)

func (dc *DockerV1) ValidateInstruction(node *parser.Node) error {
	_, err := instructions.ParseInstruction(node)
	return err
}

// Parse parses the given dockerfile and gives docker config.
func (dc *DockerV1) Parse(filepath string) (DockerConfig, string, error) {
	dockerConfig := DockerConfig{}
	data, err := ioutil.ReadFile(filepath)
	comments := ""
	if err != nil {
		zap.S().Error("error loading docker file", filepath, zap.Error(err))
		return DockerConfig{}, "", err
	}
	r := bytes.NewReader(data)
	res, err := parser.Parse(r)
	if err != nil {
		zap.S().Errorf("error while parsing iac file", filepath, zap.Error(err))
		return DockerConfig{}, "", err
	}

	for _, child := range res.AST.Children {
		values := []string{}
		err = dc.ValidateInstruction(child)
		if err != nil {
			return DockerConfig{}, "", err
		}

		for _, comment := range child.PrevComment {
			comments = comments + "#" + comment + "\n"
		}

		for i := child.Next; i != nil; i = i.Next {
			values = append(values, i.Value)
		}
		value := strings.Join(values, stringJoinCharacter)
		switch child.Value {
		case command.Arg:
			dockerConfig.Args = append(dockerConfig.Args, value)
		case command.Cmd:
			dockerConfig.Cmd = append(dockerConfig.Cmd, value)
		case command.From:
			dockerConfig.From = append(dockerConfig.From, value)
		case command.Label:
			dockerConfig.Labels = append(dockerConfig.Labels, value)
		case command.Run:
			dockerConfig.Run = append(dockerConfig.Run, value)
		case command.Expose:
			dockerConfig.Expose = append(dockerConfig.Expose, value)
		case command.Env:
			dockerConfig.Env = append(dockerConfig.Env, value)
		case command.Add:
			dockerConfig.Add = append(dockerConfig.Add, value)
		case command.Copy:
			dockerConfig.Copy = append(dockerConfig.Copy, value)
		case command.Entrypoint:
			dockerConfig.Entrypoint = append(dockerConfig.Entrypoint, value)
		case command.Volume:
			dockerConfig.Volume = append(dockerConfig.Volume, value)
		case command.User:
			dockerConfig.User = append(dockerConfig.User, value)
		case command.Workdir:
			dockerConfig.WorkDir = append(dockerConfig.WorkDir, value)
		case command.Onbuild:
			dockerConfig.Onbuild = append(dockerConfig.Onbuild, value)
		case command.Healthcheck:
			dockerConfig.HealthCheck = append(dockerConfig.HealthCheck, value)
		case command.Maintainer:
			dockerConfig.Maintainer = append(dockerConfig.Maintainer, value)
		case command.Shell:
			dockerConfig.Shell = append(dockerConfig.Shell, value)
		case command.StopSignal:
			dockerConfig.StopSignal = append(dockerConfig.StopSignal, value)
		default:
			zap.S().Errorf("Unknow command %s", child.Value, nil)
		}
	}
	return dockerConfig, comments, nil
}
