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
	"github.com/awslabs/goformation/v5/cloudformation/ec2"
)

// SubnetConfig holds config for aws_subnet
type SubnetConfig struct {
	Config
	AssignIpv6AddressOnCreation bool   `json:"assign_ipv6_address_on_creation"`
	AvailabilityZone            string `json:"availability_zone"`
	CidrBlock                   string `json:"cidr_block"`
	Ipv6CidrBlock               string `json:"ipv6_cidr_block"`
	MapPublicIpOnLaunch         bool   `json:"map_public_ip_on_launch"`
	OutpostArn                  string `json:"outpost_arn"`
	VpcId                       string `json:"vpc_id"`
}

// GetSubnetConfig returns config for aws_subnet
func GetSubnetConfig(e *ec2.Subnet) []AWSResourceConfig {
	cf := SubnetConfig{
		Config: Config{
			Tags: e.Tags,
		},
		AssignIpv6AddressOnCreation: e.AssignIpv6AddressOnCreation,
		AvailabilityZone:            e.AvailabilityZone,
		CidrBlock:                   e.CidrBlock,
		Ipv6CidrBlock:               e.Ipv6CidrBlock,
		MapPublicIpOnLaunch:         e.MapPublicIpOnLaunch,
		OutpostArn:                  e.OutpostArn,
		VpcId:                       e.VpcId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
