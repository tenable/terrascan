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

import "github.com/awslabs/goformation/v5/cloudformation/ec2"

// Ec2VpcConfig holds config for Ec2Vpc
type Ec2VpcConfig struct {
	Config
	CIDRBlock          string `json:"cidr_block"`
	EnableDNSSupport   bool   `json:"enable_dns_support"`
	EnableDNSHostnames bool   `json:"enable_dns_hostnames"`
	InstanceTenancy    string `json:"instance_tenancy"`
}

// GetEc2VpcConfig returns config for Ec2Vpc
func GetEc2VpcConfig(v *ec2.VPC) []AWSResourceConfig {
	cf := Ec2VpcConfig{
		Config: Config{
			Tags: v.Tags,
		},
		CIDRBlock:          v.CidrBlock,
		EnableDNSSupport:   v.EnableDnsSupport,
		EnableDNSHostnames: v.EnableDnsHostnames,
		InstanceTenancy:    v.InstanceTenancy,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: v.AWSCloudFormationMetadata,
	}}
}
