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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

const (
	dockerDirectory        string = "docker"
	resourceTypeDockerfile string = "dockerfile"
)

// LoadIacDir loads the docker file specified in given folder.
func (dc *DockerV1) LoadIacDir(absRootDir string, nonRecursive bool) (output.AllResourceConfigs, error) {
	allResourcesConfig := make(map[string][]output.ResourceConfig)
	fileMap, err := utils.FindFilesBySuffix(absRootDir, []string{DockerFileName})
	if err != nil {
		zap.S().Errorf("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, multierror.Append(dc.errIacLoadDirs, results.DirScanErr{IacType: "docker", Directory: absRootDir, ErrMessage: err.Error()})
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])
			data, comments, err := dc.Parse(file)
			if err != nil {
				errMsg := fmt.Sprintf("error while parsing file %s", file)
				zap.S().Errorf("error while searching for iac files", zap.String("root dir", absRootDir), errMsg)
				dc.errIacLoadDirs = multierror.Append(dc.errIacLoadDirs, results.DirScanErr{IacType: "docker", Directory: absRootDir, ErrMessage: errMsg})
				continue
			}
			minSeverity, maxSeverity := utils.GetMinMaxSeverity(comments)
			sourcePath := file
			sourcePath, err = filepath.Rel(absRootDir, file)
			if err != nil {
				zap.S().Debug("error while getting the relative path for", zap.String("IAC file", file), zap.Error(err))
			}
			skipRules := utils.GetSkipRules(comments)

			dockerCommand := []string{}
			for j := 0; j < len(data); j++ {
				dockerCommand = append(dockerCommand, data[j].Cmd)
				config := output.ResourceConfig{
					Name:        *files[i],
					Type:        data[j].Cmd,
					Line:        data[j].Line,
					ID:          data[j].Cmd + "." + GetresourceIdforDockerfile(file, data[j].Value),
					Source:      sourcePath,
					Config:      data[j].Value,
					SkipRules:   skipRules,
					MinSeverity: minSeverity,
					MaxSeverity: maxSeverity,
				}
				allResourcesConfig[data[j].Cmd] = append(allResourcesConfig[data[j].Cmd], config)

			}
			config := output.ResourceConfig{
				Name:        *files[i],
				Type:        resourceTypeDockerfile,
				Line:        1,
				ID:          dockerDirectory + "." + GetresourceIdforDockerfile(file, ""),
				Source:      sourcePath,
				Config:      dockerCommand,
				SkipRules:   skipRules,
				MinSeverity: minSeverity,
				MaxSeverity: maxSeverity,
			}
			allResourcesConfig[dockerDirectory] = append(allResourcesConfig[dockerDirectory], config)
		}
	}

	return allResourcesConfig, dc.errIacLoadDirs

}

// GetresourceIdforDockerfile Generates hash of the string to be used as the reference id for docker file
func GetresourceIdforDockerfile(filepath string, value string) (referenceID string) {
	hasher := md5.New()
	hasher.Write([]byte(filepath + value))
	referenceID = strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))
	return
}
