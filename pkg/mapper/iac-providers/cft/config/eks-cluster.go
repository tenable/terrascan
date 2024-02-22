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

// EKSVPCConfigBlock holds config for EKSVPCConfig
type EKSVPCConfigBlock struct {
	SubnetIDs             []string `json:"subnet_ids"`
	SecurityGroupIDs      []string `json:"security_group_ids"`
	EndpointPrivateAccess bool     `json:"endpoint_private_access"`
	EndpointPublicAccess  bool     `json:"endpoint_public_access"`
}

// EksClusterConfig holds config for EksCluster
type EksClusterConfig struct {
	Config
	Name                   string              `json:"name"`
	RoleARN                string              `json:"role_arn"`
	VPCConfig              []EKSVPCConfigBlock `json:"vpc_config"`
	EnabledClusterLogTypes []string            `json:"enabled_cluster_log_types"`
}

// GetEksClusterConfig returns config for EksCluster
// aws_eks_cluster
func GetEksClusterConfig(c *eks.Cluster) []AWSResourceConfig {
	var vpcConfig []EKSVPCConfigBlock
	if c.ResourcesVpcConfig != nil {
		vpcConfig := make([]EKSVPCConfigBlock, 1)

		vpcConfig[0].SubnetIDs = c.ResourcesVpcConfig.SubnetIds
		vpcConfig[0].SecurityGroupIDs = c.ResourcesVpcConfig.SecurityGroupIds
		vpcConfig[0].EndpointPrivateAccess = functions.GetVal(c.ResourcesVpcConfig.EndpointPrivateAccess)
		vpcConfig[0].EndpointPublicAccess = functions.GetVal(c.ResourcesVpcConfig.EndpointPublicAccess)
	}

	cf := EksClusterConfig{
		Config: Config{
			Name: functions.GetVal(c.Name),
		},
		Name:      functions.GetVal(c.Name),
		RoleARN:   c.RoleArn,
		VPCConfig: vpcConfig,
	}

	if c.Logging != nil {
		if c.Logging.ClusterLogging != nil {
			enabledTypes := c.Logging.ClusterLogging.EnabledTypes
			if len(enabledTypes) > 0 {
				enabledClusterLogTypes := make([]string, len(enabledTypes))
				for i, enabledType := range enabledTypes {
					enabledClusterLogTypes[i] = functions.GetVal(enabledType.Type)
				}
				cf.EnabledClusterLogTypes = enabledClusterLogTypes
			}
		}
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
