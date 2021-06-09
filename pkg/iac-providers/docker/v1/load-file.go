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
	"path/filepath"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

// LoadIacFile loads the docker file specified
func (dc *DockerV1) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)
	data, comments, err := dc.Parse(absFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("error while parsing file %s, error: %v", absFilePath, err)
		zap.S().Errorf("error while parsing file %s", absFilePath, err)
		return allResourcesConfig, errors.New(errMsg)
	}
	minSeverity, maxSeverity := utils.GetMinMaxSeverity(comments)
	config := output.ResourceConfig{
		Name:        filepath.Base(absFilePath),
		Type:        resourceTypeDockerfile,
		Line:        1,
		ID:          dockerDirectory + "." + GetresourceIdforDockerfile(absFilePath),
		Source:      filepath.Base(absFilePath),
		Config:      data,
		SkipRules:   utils.GetSkipRules(comments),
		MinSeverity: minSeverity,
		MaxSeverity: maxSeverity,
	}
	allResourcesConfig[dockerDirectory] = append(allResourcesConfig[dockerDirectory], config)
	return allResourcesConfig, nil

}
