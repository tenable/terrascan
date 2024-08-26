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

package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EcsTaskDefinitionConfig holds config for aws_ecs_task_definition
type EcsTaskDefinitionConfig struct {
	Config
	ContainerDefinitions string         `json:"container_definitions"`
	NetworkMode          string         `json:"network_mode"`
	Volumes              []VolumeConfig `json:"volume"`
}

// VolumeConfig holds config for volume attribute of aws_ecs_task_definition
type VolumeConfig struct {
	EfsVolumeConfiguration EfsVolumeConfig `json:"efs_volume_configuration"`
}

// EfsVolumeConfig holds config for efs_volume_configuration attribute of volume
type EfsVolumeConfig struct {
	TransitEncryption string `json:"transit_encryption"`
}

// ContainerDefinitionConfig holds config for container_definitions
type ContainerDefinitionConfig struct {
	Environment []EnvironmentConfig `json:"environment"`
}

// EnvironmentConfig holds config for environment attribute for container_definitions
type EnvironmentConfig struct {
	Name string `json:"name"`
}

// GetEcsTaskDefinitionConfig returns config for aws_ecs_service and aws_ecs_task_definition
// aws_ecs_task_definition
func GetEcsTaskDefinitionConfig(t *ecs.TaskDefinition) []AWSResourceConfig {
	cf := EcsTaskDefinitionConfig{
		Config: Config{
			Tags: t.Tags,
		},
		NetworkMode: functions.GetVal(t.NetworkMode),
	}

	if t.ContainerDefinitions != nil {
		// add container_definitions as a json string with mapped values
		cDefs := make([]ContainerDefinitionConfig, 0)
		for _, cDef := range t.ContainerDefinitions {
			// add environment kn pairs
			if cDef.Environment != nil {
				env := make([]EnvironmentConfig, 0)
				for _, kvPair := range cDef.Environment {
					env = append(env, EnvironmentConfig{
						Name: functions.GetVal(kvPair.Name),
					})
				}
				cDefs = append(cDefs, ContainerDefinitionConfig{
					Environment: env,
				})
			}
		}
		definitions, err := json.Marshal(cDefs)
		if err == nil {
			cf.ContainerDefinitions = string(definitions)
		}
	}

	if t.Volumes != nil {
		volumes := make([]VolumeConfig, 0)
		for _, volume := range t.Volumes {
			if volume.EFSVolumeConfiguration != nil {
				volumes = append(volumes, VolumeConfig{
					EfsVolumeConfiguration: EfsVolumeConfig{
						TransitEncryption: functions.GetVal(volume.EFSVolumeConfiguration.TransitEncryption),
					},
				})
			}
		}
		cf.Volumes = volumes
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: t.AWSCloudFormationMetadata,
	}}
}
