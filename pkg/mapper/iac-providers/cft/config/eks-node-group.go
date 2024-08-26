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
	"github.com/awslabs/goformation/v7/cloudformation/eks"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EksNodeGroupScalingConfigBlock holds config for EksNodeGroupScalingConfig
type EksNodeGroupScalingConfigBlock struct {
	DesiredSize int `json:"desired_size"`
	MaxSize     int `json:"max_size"`
	MinSize     int `json:"min_size"`
}

// EksNodeGroupConfig holds config for EksNodeGroup
type EksNodeGroupConfig struct {
	Config
	ClusterName   string                           `json:"cluster_name"`
	NodeGroupName string                           `json:"node_group_name"`
	SubnetIDs     []string                         `json:"subnet_ids"`
	NodeRoleARN   string                           `json:"node_role_arn"`
	ScalingConfig []EksNodeGroupScalingConfigBlock `json:"scaling_config"`
	Labels        interface{}                      `json:"labels"`
}

// GetEksNodeGroupConfig returns config for EksNodeGroup
// aws_eks_node_group
func GetEksNodeGroupConfig(g *eks.Nodegroup) []AWSResourceConfig {
	var scalingConfig []EksNodeGroupScalingConfigBlock
	if g.ScalingConfig != nil {
		scalingConfig = make([]EksNodeGroupScalingConfigBlock, 1)
		scalingConfig[0].DesiredSize = functions.GetVal(g.ScalingConfig.DesiredSize)
		scalingConfig[0].MaxSize = functions.GetVal(g.ScalingConfig.MaxSize)
		scalingConfig[0].MinSize = functions.GetVal(g.ScalingConfig.MinSize)
	}

	cf := EksNodeGroupConfig{
		Config: Config{
			Name: functions.GetVal(g.NodegroupName),
			Tags: functions.PatchAWSTags(g.Tags),
		},
		ClusterName:   g.ClusterName,
		NodeGroupName: functions.GetVal(g.NodegroupName),
		NodeRoleARN:   g.NodeRole,
		SubnetIDs:     g.Subnets,
		ScalingConfig: scalingConfig,
		Labels:        g.Labels,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: g.AWSCloudFormationMetadata,
	}}
}
