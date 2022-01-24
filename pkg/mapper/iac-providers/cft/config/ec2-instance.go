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
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strconv"

	"github.com/awslabs/goformation/v4/cloudformation/ec2"
)

// NetworkInterfaceBlock holds config for NetworkInterface
type NetworkInterfaceBlock struct {
	NetworkInterfaceID  string `json:"network_interface_id"`
	DeviceIndex         int    `json:"device_index"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

// TagBlock holds config for Tag
type TagBlock struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
	AssociatePublicIPAddress bool                    `json:"associate_public_ip_address"`
	NetworkInterface         []NetworkInterfaceBlock `json:"network_interface"`
	Tags                     []TagBlock              `json:"tags"`
}

// GetEC2InstanceConfig returns config for EC2Instance
func GetEC2InstanceConfig(i *ec2.Instance) []AWSResourceConfig {
	name := fmt.Sprintf("aws_instance_%s", getname(10))

	var publicIp bool
	nics := make([]NetworkInterfaceBlock, len(i.NetworkInterfaces))
	for index := range i.NetworkInterfaces {
		nics[index].NetworkInterfaceID = i.NetworkInterfaces[index].NetworkInterfaceId
		nics[index].DeviceIndex, _ = strconv.Atoi(i.NetworkInterfaces[index].DeviceIndex)
		nics[index].DeleteOnTermination = i.NetworkInterfaces[index].DeleteOnTermination

		publicIp = i.NetworkInterfaces[index].AssociatePublicIpAddress
	}

	tags := make([]TagBlock, len(i.Tags))
	for index := range i.Tags {
		tags[index].Key = i.Tags[index].Key
		tags[index].Value = i.Tags[index].Value
	}

	cf := EC2InstanceConfig{
		Config: Config{
			Name: name,
		},
		AMI:                      i.ImageId,
		InstanceType:             i.InstanceType,
		EBSOptimized:             i.EbsOptimized,
		Hibernation:              i.HibernationOptions.Configured,
		Monitoring:               i.Monitoring,
		IAMInstanceProfile:       i.IamInstanceProfile,
		VPCSecurityGroupIDs:      i.SecurityGroupIds,
		AssociatePublicIPAddress: publicIp,
		NetworkInterface:         nics,
		Tags:                     tags,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: i.AWSCloudFormationMetadata,
	}}
}

func getname(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
