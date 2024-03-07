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
	"fmt"
	"strconv"

	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
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
	AMI                      string                  `json:"ami"`
	InstanceType             string                  `json:"instance_type"`
	EBSOptimized             bool                    `json:"ebs_optimized"`
	Hibernation              bool                    `json:"hibernation"`
	Monitoring               bool                    `json:"monitoring"`
	IAMInstanceProfile       string                  `json:"iam_instance_profile"`
	VPCSecurityGroupIDs      []string                `json:"vpc_security_group_ids"`
	NetworkInterface         []NetworkInterfaceBlock `json:"network_interface"`
	UserData                 string                  `json:"user_data_base64"`
	AssociatePublicIPAddress bool                    `json:"associate_public_ip_address"`
	SubnetID                 string                  `json:"subnet_id"`
}

// GetEC2InstanceConfig returns config for EC2Instance
// aws_instance
func GetEC2InstanceConfig(i *ec2.Instance, instanceName string) []AWSResourceConfig {
	networkInterfaces := i.NetworkInterfaces

	nics := make([]NetworkInterfaceBlock, len(networkInterfaces))
	niconfigs := make([]NetworkInterfaceConfig, len(networkInterfaces))
	awsconfig := make([]AWSResourceConfig, len(networkInterfaces))
	associatePublicAddress := false
	var vpcSecurityID []string
	for index, networkInterface := range networkInterfaces {
		nics[index].NetworkInterfaceID = functions.GetVal(networkInterface.NetworkInterfaceId)
		nics[index].DeleteOnTermination = functions.GetVal(networkInterface.DeleteOnTermination)
		var devindex int
		devindex, err := strconv.Atoi(networkInterface.DeviceIndex)
		if err != nil {
			devindex = 0
		}
		nics[index].DeviceIndex = devindex

		// create aws_network_interface resource on the fly for every network interface used in aws_instance
		niconfigs[index].SubnetID = functions.GetVal(networkInterface.SubnetId)
		if networkInterface.PrivateIpAddress != nil {
			niconfigs[index].PrivateIPs = []string{functions.GetVal(networkInterface.PrivateIpAddress)}
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
		if networkInterface.AssociatePublicIpAddress != nil {
			associatePublicAddress = associatePublicAddress || *networkInterface.AssociatePublicIpAddress
		}
		if networkInterface.GroupSet != nil {
			for _, groupName := range networkInterface.GroupSet {
				if groupName != "" {
					vpcSecurityID = append(vpcSecurityID, groupName)
				}
			}

		}
	}
	subnetID := ""
	if i.SubnetId != nil {
		subnetID = *i.SubnetId
	}

	if i.SecurityGroupIds != nil {
		vpcSecurityID = append(vpcSecurityID, i.SecurityGroupIds...)
	}

	ec2Config := EC2InstanceConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(i.Tags),
			Name: instanceName,
		},
		AMI:                      functions.GetVal(i.ImageId),
		InstanceType:             functions.GetVal(i.InstanceType),
		EBSOptimized:             functions.GetVal(i.EbsOptimized),
		Monitoring:               functions.GetVal(i.Monitoring),
		IAMInstanceProfile:       functions.GetVal(i.IamInstanceProfile),
		VPCSecurityGroupIDs:      vpcSecurityID,
		NetworkInterface:         nics,
		AssociatePublicIPAddress: associatePublicAddress,
		SubnetID:                 subnetID,
	}

	if i.UserData != nil {
		ec2Config.UserData = *i.UserData
	}

	if i.HibernationOptions != nil {
		ec2Config.Hibernation = functions.GetVal(i.HibernationOptions.Configured)
	}

	var awsconfigec2 AWSResourceConfig
	awsconfigec2.Resource = ec2Config
	awsconfigec2.Metadata = i.AWSCloudFormationMetadata
	awsconfig = append(awsconfig, awsconfigec2)

	return awsconfig
}
