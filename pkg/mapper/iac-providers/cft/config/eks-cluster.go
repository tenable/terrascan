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

import "github.com/awslabs/goformation/v5/cloudformation/eks"

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
func GetEksClusterConfig(c *eks.Cluster) []AWSResourceConfig {
	var vpcConfig []EKSVPCConfigBlock
	if c.ResourcesVpcConfig != nil {
		vpcConfig := make([]EKSVPCConfigBlock, 1)

		vpcConfig[0].SubnetIDs = c.ResourcesVpcConfig.SubnetIds
		vpcConfig[0].SecurityGroupIDs = c.ResourcesVpcConfig.SecurityGroupIds
		vpcConfig[0].EndpointPrivateAccess = c.ResourcesVpcConfig.EndpointPrivateAccess
		vpcConfig[0].EndpointPublicAccess = c.ResourcesVpcConfig.EndpointPublicAccess
	}

	enabledClusterLogTypes := make([]string, len(c.Logging.ClusterLogging.EnabledTypes))
	for i := range c.Logging.ClusterLogging.EnabledTypes {
		enabledClusterLogTypes[i] = c.Logging.ClusterLogging.EnabledTypes[i].Type
	}

	cf := EksClusterConfig{
		Config: Config{
			Name: c.Name,
		},
		Name:                   c.Name,
		RoleARN:                c.RoleArn,
		VPCConfig:              vpcConfig,
		EnabledClusterLogTypes: enabledClusterLogTypes,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
