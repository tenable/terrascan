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
	"github.com/awslabs/goformation/v7/cloudformation/redshift"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// RedshiftClusterConfig holds config for aws_redshift_cluster
type RedshiftClusterConfig struct {
	Config
	LoggingProperties  interface{} `json:"logging,omitempty"`
	KmsKeyID           string      `json:"kms_key_id,omitempty"`
	Encrypted          bool        `json:"encrypted"`
	PubliclyAccessible bool        `json:"publicly_accessible"`
}

// GetRedshiftClusterConfig returns config for aws_redshift_cluster
// aws_redshift_cluster
func GetRedshiftClusterConfig(c *redshift.Cluster) []AWSResourceConfig {
	cf := RedshiftClusterConfig{
		Config: Config{
			Name: c.DBName,
			Tags: functions.PatchAWSTags(c.Tags),
		},
		KmsKeyID:           functions.GetVal(c.KmsKeyId),
		Encrypted:          functions.GetVal(c.Encrypted),
		PubliclyAccessible: functions.GetVal(c.PubliclyAccessible),
	}
	if c.LoggingProperties != nil {
		// if LoggingProperties are mentioned in cft,
		// its always enabled
		logging := make(map[string]bool)
		logging["enable"] = true
		cf.LoggingProperties = []map[string]bool{logging}
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
