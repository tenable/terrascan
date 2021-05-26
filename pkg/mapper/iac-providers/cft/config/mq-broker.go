package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/amazonmq"
)

// MqBrokerConfig holds config for aws_mq_broker
type MqBrokerConfig struct {
	Logs interface{} `json:"logs,omitempty"`
	Config
	PubliclyAccessible bool `json:"publicly_accessible"`
}

// GetMqBorkerConfig returns config for aws_mq_broker
func GetMqBorkerConfig(c *amazonmq.Broker) []AWSResourceConfig {
	cf := MqBrokerConfig{
		Config: Config{
			Name: c.BrokerName,
		},
		PubliclyAccessible: c.PubliclyAccessible,
	}
	if c.Logs != nil {
		log := make(map[string]bool)
		if c.Logs.Audit {
			log["audit"] = true
		} else {
			log["audit"] = false
		}
		if c.Logs.General {
			log["general"] = true
		} else {
			log["general"] = false
		}
		cf.Logs = []map[string]bool{log}
	}

	return []AWSResourceConfig{{Resource: cf}}
}
