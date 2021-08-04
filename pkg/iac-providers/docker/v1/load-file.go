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
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

const (
	docker                 string = "docker"
	resourceTypeDockerfile string = "dockerfile"
	underScoreSeparator    string = "_"

	// IDConnectorString is string connector used in id creation
	IDConnectorString        string = "."
	dockerResourceTypePrefix        = docker + underScoreSeparator
)

// LoadIacFile loads the docker file specified and create ResourceConfig for each dockerfile
func (dc *DockerV1) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	data, comments, err := dc.Parse(absFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("error while parsing dockerfile %s, error: %v", absFilePath, err)
		zap.S().Errorf("error while parsing dockerfile %s", absFilePath, err)
		return allResourcesConfig, errors.New(errMsg)
	}

	minSeverity, maxSeverity := utils.GetMinMaxSeverity(comments)

	skipRules := utils.GetSkipRules(comments)

	// create an array of all the instructions present in the docker file
	dockerCommand := []string{}

	// create config for each instruction of  dockerfile
	for i := 0; i < len(data); i++ {
		dockerCommand = append(dockerCommand, data[i].Cmd)

		config := output.ResourceConfig{
			Name:        filepath.Base(absFilePath),
			Type:        dockerResourceTypePrefix + data[i].Cmd,
			Line:        data[i].Line,
			ID:          dockerResourceTypePrefix + data[i].Cmd + IDConnectorString + GetresourceIdforDockerfile(absFilePath, data[i].Value, data[i].Line),
			Source:      dc.getSourceRelativePath(absFilePath),
			Config:      data[i].Value,
			SkipRules:   skipRules,
			MinSeverity: minSeverity,
			MaxSeverity: maxSeverity,
		}
		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], config)

	}

	// Creates config for entire dockerfile which has array of instructions against the Config field.
	// Created to use against policies which checks for availablility of command/instruction in dockerfile
	// if command is not present line no also doesnot have any importance thats why set to 1.
	config := output.ResourceConfig{
		Name:        filepath.Base(absFilePath),
		Type:        dockerResourceTypePrefix + resourceTypeDockerfile,
		Line:        1,
		ID:          dockerResourceTypePrefix + resourceTypeDockerfile + IDConnectorString + GetresourceIdforDockerfile(absFilePath, "", 1),
		Source:      dc.getSourceRelativePath(absFilePath),
		Config:      dockerCommand,
		SkipRules:   skipRules,
		MinSeverity: minSeverity,
		MaxSeverity: maxSeverity,
	}

	allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], config)
	return allResourcesConfig, nil

}

// getSourceRelativePath fetches the relative path of file being loaded
func (dc *DockerV1) getSourceRelativePath(sourceFile string) string {

	// rootDir should be empty when file scan was initiated by user
	if dc.absRootDir == "" {
		return filepath.Base(sourceFile)
	}
	relPath, err := filepath.Rel(dc.absRootDir, sourceFile)
	if err != nil {
		zap.S().Debug("error while getting the relative path for", zap.String("IAC file", sourceFile), zap.Error(err))
		return sourceFile
	}
	return relPath
}

// GetresourceIdforDockerfile Generates hash of the string to be used as the reference id for docker file
// added line no in creating hash because dockerfile may have same command multiple times with same value
func GetresourceIdforDockerfile(filepath string, value string, lineNumber int) (referenceID string) {
	hasher := md5.New()
	hasher.Write([]byte(filepath + value + strconv.Itoa(lineNumber)))
	referenceID = strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	return
}
