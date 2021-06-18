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

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"go.uber.org/zap"
)

// ResourceConfig holds information about individual docker instructions
type ResourceConfig struct {
	Cmd   string `json:"cmd"`
	Value string `json:"value"`
	Line  int    `json:"line"`
}

const (
	stringJoinCharacter = " "
	commentPrefix       = "#"
	newLine             = "\n"
)

// ValidateInstruction validates the dockerfile instructions
func (dc *DockerV1) ValidateInstruction(node *parser.Node) error {
	_, err := instructions.ParseInstruction(node)
	return err
}

// Parse parses the given dockerfile and gives docker config.
func (dc *DockerV1) Parse(filepath string) ([]ResourceConfig, string, error) {
	config := []ResourceConfig{}
	data, err := ioutil.ReadFile(filepath)
	comments := ""
	if err != nil {
		zap.S().Error("error loading docker file", filepath, zap.Error(err))
		return []ResourceConfig{}, "", err
	}
	r := bytes.NewReader(data)
	res, err := parser.Parse(r)
	if err != nil {
		zap.S().Errorf("error while parsing iac file", filepath, zap.Error(err))
		return []ResourceConfig{}, "", err
	}

	for _, child := range res.AST.Children {
		values := []string{}
		err = dc.ValidateInstruction(child)
		if err != nil {
			return []ResourceConfig{}, "", err
		}

		for _, comment := range child.PrevComment {
			comments = comments + commentPrefix + comment + newLine
		}

		for i := child.Next; i != nil; i = i.Next {
			values = append(values, i.Value)
		}
		value := strings.Join(values, stringJoinCharacter)
		tempConfig := ResourceConfig{
			Cmd:   child.Value,
			Value: value,
			Line:  child.StartLine,
		}
		config = append(config, tempConfig)
	}
	return config, comments, nil
}
