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

// NatGatewayConfig holds config for aws_nat_gateway
type NatGatewayConfig struct {
	Config
	AllocationID     string `json:"allocation_id"`
	ConnectivityType string `json:"connectivity_type"`
	SubnetID         string `json:"subnet_id"`
}

// GetNatGatewayConfig returns config for aws_nat_gateway
// aws_nat_gateway
func GetNatGatewayConfig(e *ec2.NatGateway) []AWSResourceConfig {
	cf := NatGatewayConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(e.Tags),
		},
		AllocationID:     functions.GetVal(e.AllocationId),
		ConnectivityType: functions.GetVal(e.ConnectivityType),
		SubnetID:         e.SubnetId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
