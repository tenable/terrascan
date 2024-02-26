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
	"github.com/awslabs/goformation/v7/cloudformation/autoscaling"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EbsBlockDeviceBlock holds config for EbsBlockDevice
type EbsBlockDeviceBlock struct {
	DeviceName          string `json:"device_name"`
	Encrypted           bool   `json:"encrypted"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

// MetadataOptionsBlock holds config for MetadataOptions
type MetadataOptionsBlock struct {
	HTTPEndpoint string `json:"http_endpoint"`
	HTTPTokens   string `json:"http_tokens"`
}

// AutoScalingLaunchConfigurationConfig holds config for AutoScalingLaunchConfiguration
type AutoScalingLaunchConfigurationConfig struct {
	Config
	EnableMonitoring bool                  `json:"enable_monitoring"`
	UserDataBase64   string                `json:"user_data_base64"`
	UserData         string                `json:"user_data"`
	MetadataOptions  MetadataOptionsBlock  `json:"metadata_options"`
	EbsBlockDevice   []EbsBlockDeviceBlock `json:"ebs_block_device"`
}

// GetAutoScalingLaunchConfigurationConfig returns config for AutoScalingLaunchConfiguration
// aws_launch_configuration
func GetAutoScalingLaunchConfigurationConfig(l *autoscaling.LaunchConfiguration) []AWSResourceConfig {
	var ebsBlockDevice []EbsBlockDeviceBlock
	if l.BlockDeviceMappings != nil {
		blockDeviceMappingLen := len(l.BlockDeviceMappings)
		ebsBlockDevice = make([]EbsBlockDeviceBlock, blockDeviceMappingLen)
		for i, blockDeviceMapping := range l.BlockDeviceMappings {
			if blockDeviceMapping.Ebs != nil {
				ebsBlockDevice[i].Encrypted = functions.GetVal(blockDeviceMapping.Ebs.Encrypted)
				ebsBlockDevice[i].DeleteOnTermination = functions.GetVal(blockDeviceMapping.Ebs.DeleteOnTermination)
			}
			ebsBlockDevice[i].DeviceName = blockDeviceMapping.DeviceName
		}
	}

	var metadataOptions MetadataOptionsBlock
	if l.MetadataOptions != nil {
		metadataOptions.HTTPEndpoint = functions.GetVal(l.MetadataOptions.HttpEndpoint)
		metadataOptions.HTTPTokens = functions.GetVal(l.MetadataOptions.HttpTokens)
	}

	cf := AutoScalingLaunchConfigurationConfig{
		Config: Config{
			Name: functions.GetVal(l.LaunchConfigurationName),
		},
		EnableMonitoring: functions.GetVal(l.InstanceMonitoring),
		MetadataOptions:  metadataOptions,
		EbsBlockDevice:   ebsBlockDevice,
	}

	if l.UserData != nil {
		cf.UserDataBase64 = *l.UserData
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: l.AWSCloudFormationMetadata,
	}}
}

// isASCII  not in use as we dont need to convert base64 data to string for iac scanning
// func isASCII(s string) bool {
// 	for i := 0; i < len(s); i++ {
// 		if s[i] > unicode.MaxASCII {
// 			return false
// 		}
// 	}
// 	return true
// }
