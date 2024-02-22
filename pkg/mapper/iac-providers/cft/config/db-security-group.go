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
	"github.com/awslabs/goformation/v7/cloudformation/rds"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// DBIngress holds config for ingress block
type DBIngress struct {
	CIDR              string `json:"cidr"`
	SecurityGroupName string `json:"security_group_name"`
}

// DBSecurityGroupConfig holds config for aws_db_security_group
type DBSecurityGroupConfig struct {
	Config
	Ingress []DBIngress `json:"ingress"`
}

// GetDBSecurityGroupConfig returns config for aws_db_security_group
// aws_db_security_group
func GetDBSecurityGroupConfig(dbsg *rds.DBSecurityGroup) []AWSResourceConfig {
	cf := DBSecurityGroupConfig{
		Config: Config{
			Tags: dbsg.Tags,
		},
	}

	if dbsg.DBSecurityGroupIngress != nil {
		cf.Ingress = make([]DBIngress, len(dbsg.DBSecurityGroupIngress))
		for i, dbsgi := range dbsg.DBSecurityGroupIngress {
			cf.Ingress[i].CIDR = functions.GetVal(dbsgi.CIDRIP)
			cf.Ingress[i].SecurityGroupName = functions.GetVal(dbsgi.EC2SecurityGroupName)
		}
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: dbsg.AWSCloudFormationMetadata,
	}}
}
