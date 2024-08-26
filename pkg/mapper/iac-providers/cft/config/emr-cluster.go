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
	"github.com/awslabs/goformation/v7/cloudformation/emr"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// KerberosAttributesBlock holds config for KerberosAttributes
type KerberosAttributesBlock struct {
	KDCAdminPassword string `json:"kdc_admin_password"`
	Realm            string `json:"realm"`
}

// EmrClusterConfig holds config for EmrCluster
type EmrClusterConfig struct {
	Config
	Name                  string                    `json:"name"`
	ReleaseLabel          string                    `json:"release_label"`
	ServiceRole           string                    `json:"service_role"`
	TerminationProtection bool                      `json:"termination_protection"`
	KerberosAttributes    []KerberosAttributesBlock `json:"kerberos_attributes"`
}

// GetEmrClusterConfig returns config for EmrCluster
// aws_emr_cluster
func GetEmrClusterConfig(c *emr.Cluster) []AWSResourceConfig {
	var kerberosAttributes []KerberosAttributesBlock
	if c.KerberosAttributes != nil {
		kerberosAttributes = make([]KerberosAttributesBlock, 1)

		kerberosAttributes[0].KDCAdminPassword = c.KerberosAttributes.KdcAdminPassword
		kerberosAttributes[0].Realm = c.KerberosAttributes.Realm
	}

	cf := EmrClusterConfig{
		Config: Config{
			Name: c.Name,
		},
		Name:                  c.Name,
		ReleaseLabel:          functions.GetVal(c.ReleaseLabel),
		ServiceRole:           c.ServiceRole,
		TerminationProtection: functions.GetVal(c.Instances.TerminationProtected),
		KerberosAttributes:    kerberosAttributes,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
