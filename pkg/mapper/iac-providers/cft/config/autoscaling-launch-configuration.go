/*
    Copyright (C) 2022 Accurics, Inc.

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
	"encoding/base64"

	"github.com/awslabs/goformation/v4/cloudformation/autoscaling"
)

// EbsBlockDeviceBlock hold config for EbsBlockDevice
type EbsBlockDeviceBlock struct {
	DeviceName string `json:"device_name"`
	Encrypted  bool   `json:"encrypted"`
}

// MetadataOptionsBlock hold config for MetadataOptions
type MetadataOptionsBlock struct {
	HttpEndpoint string `json:"http_endpoint"`
	HttpTokens   string `json:"http_tokens"`
}

// AutoScalingLaunchConfigurationConfig hold config for AutoScalingLaunchConfiguration
type AutoScalingLaunchConfigurationConfig struct {
	Config
	EnableMonitoring bool                  `json:"enable_monitoring"`
	UserDataBase64   string                `json:"user_data_base64"`
	MetadataOptions  MetadataOptionsBlock  `json:"metadata_options"`
	EbsBlockDevice   []EbsBlockDeviceBlock `json:"ebs_block_device"`
}

// GetAutoScalingLaunchConfigurationConfig returns config for AutoScalingLaunchConfiguration
func GetAutoScalingLaunchConfigurationConfig(l *autoscaling.LaunchConfiguration) []AWSResourceConfig {
	userDataBase64 := base64.StdEncoding.EncodeToString([]byte(l.UserData))

	ebsBlockDevice := make([]EbsBlockDeviceBlock, len(l.BlockDeviceMappings))

	for i := range l.BlockDeviceMappings {
		if l.BlockDeviceMappings[i].Ebs != nil {
			ebsBlockDevice[i].Encrypted = l.BlockDeviceMappings[i].Ebs.Encrypted
		}
		ebsBlockDevice[i].DeviceName = l.BlockDeviceMappings[i].DeviceName
	}

	var metadataOptions MetadataOptionsBlock
	metadataOptions.HttpEndpoint = l.MetadataOptions.HttpEndpoint
	metadataOptions.HttpTokens = l.MetadataOptions.HttpTokens

	cf := AutoScalingLaunchConfigurationConfig{
		Config: Config{
			Name: l.LaunchConfigurationName,
		},
		EnableMonitoring: l.InstanceMonitoring,
		UserDataBase64:   userDataBase64,
		MetadataOptions:  metadataOptions,
		EbsBlockDevice:   ebsBlockDevice,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: l.AWSCloudFormationMetadata,
	}}
}
