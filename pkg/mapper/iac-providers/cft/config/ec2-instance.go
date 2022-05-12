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
	"fmt"
	"strconv"

	"github.com/awslabs/goformation/v5/cloudformation/ec2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/store"
)

// GetNetworkInterface represents subresource aws_network_interface for NetworkInterface attribute
const (
	GetNetworkInterface = "NetworkInterface"
)

// AttachmentBlock holds config for Attachment
type AttachmentBlock struct {
	Instance    string `json:"instance"`
	DeviceIndex int    `json:"device_index"`
}

// NetworkInterfaceConfig holds config for NetworkInterface
type NetworkInterfaceConfig struct {
	Config
	SubnetID   string            `json:"subnet_id"`
	PrivateIPs []string          `json:"private_ips"`
	Attachment []AttachmentBlock `json:"attachment"`
}

// NetworkInterfaceBlock holds config for NetworkInterface
type NetworkInterfaceBlock struct {
	NetworkInterfaceID  string `json:"network_interface_id"`
	DeviceIndex         int    `json:"device_index"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

// EC2InstanceConfig holds config for EC2Instance
type EC2InstanceConfig struct {
	Config
	AMI                 string                  `json:"ami"`
	InstanceType        string                  `json:"instance_type"`
	EBSOptimized        bool                    `json:"ebs_optimized"`
	Hibernation         bool                    `json:"hibernation"`
	Monitoring          bool                    `json:"monitoring"`
	IAMInstanceProfile  string                  `json:"iam_instance_profile"`
	VPCSecurityGroupIDs []string                `json:"vpc_security_group_ids"`
	NetworkInterface    []NetworkInterfaceBlock `json:"network_interface"`
}

// GetEC2InstanceConfig returns config for EC2Instance
func GetEC2InstanceConfig(i *ec2.Instance, instanceName string) []AWSResourceConfig {
	nics := make([]NetworkInterfaceBlock, len(i.NetworkInterfaces))
	niconfigs := make([]NetworkInterfaceConfig, len(i.NetworkInterfaces))
	awsconfig := make([]AWSResourceConfig, len(i.NetworkInterfaces))

	for index := range i.NetworkInterfaces {
		nics[index].NetworkInterfaceID = i.NetworkInterfaces[index].NetworkInterfaceId
		nics[index].DeleteOnTermination = i.NetworkInterfaces[index].DeleteOnTermination
		var devindex int
		devindex, err := strconv.Atoi(i.NetworkInterfaces[index].DeviceIndex)
		if err != nil {
			devindex = 0
		}
		nics[index].DeviceIndex = devindex

		// create aws_network_interface resource on the fly for every network interface used in aws_instance
		niconfigs[index].SubnetID = i.NetworkInterfaces[index].SubnetId
		if i.NetworkInterfaces[index].PrivateIpAddress != "" {
			niconfigs[index].PrivateIPs = []string{i.NetworkInterfaces[index].PrivateIpAddress}
		}

		nicname := fmt.Sprintf("%s%d", instanceName, index)
		niconfigs[index].Attachment = make([]AttachmentBlock, 1)
		niconfigs[index].Attachment[0].DeviceIndex = devindex
		niconfigs[index].Attachment[0].Instance = store.AwsEc2Instance + "." + instanceName
		niconfigs[index].Config.Name = nicname

		awsconfig[index].Type = GetNetworkInterface
		awsconfig[index].Name = nicname
		awsconfig[index].Resource = niconfigs[index]
		awsconfig[index].Metadata = i.AWSCloudFormationMetadata
	}

	ec2Config := EC2InstanceConfig{
		Config: Config{
			Tags: i.Tags,
			Name: instanceName,
		},
		AMI:                 i.ImageId,
		InstanceType:        i.InstanceType,
		EBSOptimized:        i.EbsOptimized,
		Monitoring:          i.Monitoring,
		IAMInstanceProfile:  i.IamInstanceProfile,
		VPCSecurityGroupIDs: i.SecurityGroupIds,
		NetworkInterface:    nics,
	}

	if i.HibernationOptions != nil {
		ec2Config.Hibernation = i.HibernationOptions.Configured
	}

	var awsconfigec2 AWSResourceConfig
	awsconfigec2.Resource = ec2Config
	awsconfigec2.Metadata = i.AWSCloudFormationMetadata
	awsconfig = append(awsconfig, awsconfigec2)

	return awsconfig
}
