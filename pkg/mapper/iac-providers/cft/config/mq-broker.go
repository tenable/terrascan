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
	"github.com/awslabs/goformation/v7/cloudformation/amazonmq"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// MqBrokerConfig holds config for aws_mq_broker
type MqBrokerConfig struct {
	Logs interface{} `json:"logs,omitempty"`
	Config
	PubliclyAccessible bool `json:"publicly_accessible"`
}

// GetMqBorkerConfig returns config for aws_mq_broker
// aws_mq_broker
func GetMqBorkerConfig(c *amazonmq.Broker) []AWSResourceConfig {
	cf := MqBrokerConfig{
		Config: Config{
			Name: c.BrokerName,
			Tags: functions.PatchAWSTags(c.Tags),
		},
		PubliclyAccessible: c.PubliclyAccessible,
	}
	if c.Logs != nil {
		log := make(map[string]bool)
		if functions.GetVal(c.Logs.Audit) {
			log["audit"] = true
		} else {
			log["audit"] = false
		}
		if functions.GetVal(c.Logs.General) {
			log["general"] = true
		} else {
			log["general"] = false
		}
		cf.Logs = []map[string]bool{log}
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: c.AWSCloudFormationMetadata,
	}}
}
