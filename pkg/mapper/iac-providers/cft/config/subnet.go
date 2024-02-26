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
	"github.com/awslabs/goformation/v7/cloudformation/ec2"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// SubnetConfig holds config for aws_subnet
type SubnetConfig struct {
	Config
	AssignIpv6AddressOnCreation bool   `json:"assign_ipv6_address_on_creation"`
	AvailabilityZone            string `json:"availability_zone"`
	CidrBlock                   string `json:"cidr_block"`
	Ipv6CidrBlock               string `json:"ipv6_cidr_block"`
	MapPublicIPOnLaunch         bool   `json:"map_public_ip_on_launch"`
	OutpostArn                  string `json:"outpost_arn"`
	VpcID                       string `json:"vpc_id"`
}

// GetSubnetConfig returns config for aws_subnet
// aws_subnet
func GetSubnetConfig(e *ec2.Subnet) []AWSResourceConfig {
	cf := SubnetConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(e.Tags),
		},
		AssignIpv6AddressOnCreation: functions.GetVal(e.AssignIpv6AddressOnCreation),
		AvailabilityZone:            functions.GetVal(e.AvailabilityZone),
		CidrBlock:                   functions.GetVal(e.CidrBlock),
		Ipv6CidrBlock:               functions.GetVal(e.Ipv6CidrBlock),
		MapPublicIPOnLaunch:         functions.GetVal(e.MapPublicIpOnLaunch),
		OutpostArn:                  functions.GetVal(e.OutpostArn),
		VpcID:                       e.VpcId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
