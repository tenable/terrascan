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

// RouteTableConfig holds config for aws_route_table
type RouteTableConfig struct {
	Config
	VpcID string `json:"vpc_id"`
}

// GetRouteTableConfig returns config for aws_route_table
// aws_route_table
func GetRouteTableConfig(e *ec2.RouteTable) []AWSResourceConfig {
	cf := RouteTableConfig{
		Config: Config{
			Tags: functions.PatchAWSTags(e.Tags),
		},
		VpcID: e.VpcId,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}
