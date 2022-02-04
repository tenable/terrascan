/*
    Copyright (C) 2021 Accurics, Inc.

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
	"github.com/awslabs/goformation/v5/cloudformation/rds"
)

// DBSecurityGroupConfig holds config for aws_db_security_group
type DBSecurityGroupConfig struct {
	Config
	Ingress []map[string]interface{} `json:"ingress"`
}

// GetDBSecurityGroupConfig returns config for aws_db_security_group
func GetDBSecurityGroupConfig(dbsg *rds.DBSecurityGroup) []AWSResourceConfig {
	cf := DBSecurityGroupConfig{
		Config: Config{
			Tags: dbsg.Tags,
		},
	}
	for _, dbsgi := range dbsg.DBSecurityGroupIngress {
		i := make(map[string]interface{})
		i["cidr"] = dbsgi.CIDRIP
		cf.Ingress = append(cf.Ingress, i)
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: dbsg.AWSCloudFormationMetadata,
	}}
}
