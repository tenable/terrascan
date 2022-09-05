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

import "github.com/awslabs/goformation/v6/cloudformation/neptune"

// NeptuneClusterInstanceConfig holds config for aws_neptune_cluster_instance resource
type NeptuneClusterInstanceConfig struct {
	Config
	AutoMinorVersionUpgrade    bool   `json:"auto_minor_version_upgrade,omitempty"`
	AvailabilityZone           string `json:"availability_zone,omitempty"`
	DBClusterIdentifier        string `json:"cluster_identifier,omitempty"`
	DBInstanceClass            string `json:"instance_class,omitempty"`
	DBParameterGroupName       string `json:"neptune_parameter_group_name,omitempty"`
	DBSubnetGroupName          string `json:"neptune_subnet_group_name,omitempty"`
	PreferredMaintenanceWindow string `json:"preferred_backup_window,omitempty"`
}

// GetNeptuneClusterInstanceConfig returns config for aws_neptune_cluster_instance resource
func GetNeptuneClusterInstanceConfig(n *neptune.DBInstance) []AWSResourceConfig {
	cf := NeptuneClusterInstanceConfig{
		Config: Config{
			Tags: n.Tags,
		},
		AutoMinorVersionUpgrade:    *n.AutoMinorVersionUpgrade,
		AvailabilityZone:           *n.AvailabilityZone,
		DBClusterIdentifier:        *n.DBClusterIdentifier,
		DBInstanceClass:            n.DBInstanceClass,
		DBParameterGroupName:       *n.DBParameterGroupName,
		DBSubnetGroupName:          *n.DBSubnetGroupName,
		PreferredMaintenanceWindow: *n.PreferredMaintenanceWindow,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: n.AWSCloudFormationMetadata,
	}}
}
