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

import "github.com/awslabs/goformation/v5/cloudformation/dms"

// DmsReplicationInstanceConfig holds config for DmsReplicationInstance
type DmsReplicationInstanceConfig struct {
	Config
	AlocatedStorage            int      `json:"allocated_storage"`
	AutoMinorVersionUpgrade    bool     `jons:"auto_minor_version_upgrade"`
	AvailabilityZone           string   `json:"availability_zone"`
	EngineVersion              string   `json:"engine_version"`
	KMSKeyARN                  string   `json:"kms_key_arn"`
	MultiAZ                    bool     `json:"multi_az"`
	PreferredMaintenanceWindow string   `json:"preferred_maintenance_window"`
	PubliclyAccessible         bool     `json:"publicly_accessible"`
	ReplicationInstanceClass   string   `json:"replication_instance_class"`
	ReplicationInstanceID      string   `json:"replication_instance_id"`
	ReplicationSubnetGroupID   string   `json:"replication_subnet_group_id"`
	VPCSecurityGroupIDs        []string `json:"vpc_security_group_ids"`
}

// GetDmsReplicationInstanceConfig returns config for DmsReplicationInstance
func GetDmsReplicationInstanceConfig(r *dms.ReplicationInstance) []AWSResourceConfig {
	cf := DmsReplicationInstanceConfig{
		Config: Config{
			Tags: r.Tags,
		},

		AlocatedStorage:            r.AllocatedStorage,
		AutoMinorVersionUpgrade:    r.AutoMinorVersionUpgrade,
		AvailabilityZone:           r.AvailabilityZone,
		EngineVersion:              r.EngineVersion,
		KMSKeyARN:                  r.KmsKeyId,
		MultiAZ:                    r.MultiAZ,
		PreferredMaintenanceWindow: r.PreferredMaintenanceWindow,
		PubliclyAccessible:         r.PubliclyAccessible,
		ReplicationInstanceClass:   r.ReplicationInstanceClass,
		ReplicationInstanceID:      r.ReplicationInstanceIdentifier,
		ReplicationSubnetGroupID:   r.ReplicationSubnetGroupIdentifier,
		VPCSecurityGroupIDs:        r.VpcSecurityGroupIds,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
