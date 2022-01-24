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
	"unicode"

	"github.com/awslabs/goformation/v4/cloudformation/autoscaling"
)

// EbsBlockDeviceBlock hold config for EbsBlockDevice
type EbsBlockDeviceBlock struct {
	DeviceName          string `json:"device_name"`
	Encrypted           bool   `json:"encrypted"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

// MetadataOptionsBlock hold config for MetadataOptions
type MetadataOptionsBlock struct {
	HTTPEndpoint string `json:"http_endpoint"`
	HTTPTokens   string `json:"http_tokens"`
}

// AutoScalingLaunchConfigurationConfig hold config for AutoScalingLaunchConfiguration
type AutoScalingLaunchConfigurationConfig struct {
	Config
	EnableMonitoring bool                  `json:"enable_monitoring"`
	UserDataBase64   string                `json:"user_data_base64"`
	UserData         string                `json:"user_data"`
	MetadataOptions  MetadataOptionsBlock  `json:"metadata_options"`
	EbsBlockDevice   []EbsBlockDeviceBlock `json:"ebs_block_device"`
}

// GetAutoScalingLaunchConfigurationConfig returns config for AutoScalingLaunchConfiguration
func GetAutoScalingLaunchConfigurationConfig(l *autoscaling.LaunchConfiguration) []AWSResourceConfig {
	ebsBlockDevice := make([]EbsBlockDeviceBlock, len(l.BlockDeviceMappings))

	for i := range l.BlockDeviceMappings {
		if l.BlockDeviceMappings[i].Ebs != nil {
			ebsBlockDevice[i].Encrypted = l.BlockDeviceMappings[i].Ebs.Encrypted
			ebsBlockDevice[i].DeleteOnTermination = l.BlockDeviceMappings[i].Ebs.DeleteOnTermination
		}
		ebsBlockDevice[i].DeviceName = l.BlockDeviceMappings[i].DeviceName
	}

	var metadataOptions MetadataOptionsBlock
	if l.MetadataOptions != nil {
		metadataOptions.HTTPEndpoint = l.MetadataOptions.HttpEndpoint
		metadataOptions.HTTPTokens = l.MetadataOptions.HttpTokens
	}

	cf := AutoScalingLaunchConfigurationConfig{
		Config: Config{
			Name: l.LaunchConfigurationName,
		},
		EnableMonitoring: l.InstanceMonitoring,
		MetadataOptions:  metadataOptions,
		EbsBlockDevice:   ebsBlockDevice,
	}

	data, err := base64.StdEncoding.Strict().DecodeString(l.UserData)
	datastr := string(data)

	if isASCII(datastr) && err == nil {
		cf.UserDataBase64 = l.UserData
	} else {
		cf.UserData = l.UserData
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: l.AWSCloudFormationMetadata,
	}}
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
