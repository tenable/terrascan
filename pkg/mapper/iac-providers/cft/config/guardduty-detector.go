package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/guardduty"
)

// GuardDutyDetectorConfig holds config for aws_guardduty_detector
type GuardDutyDetectorConfig struct {
	Config
	Enable bool `json:"enable"`
}

// GetGuardDutyDetectorConfig returns config for aws_guardduty_detector
func GetGuardDutyDetectorConfig(d *guardduty.Detector) []AWSResourceConfig {
	cf := GuardDutyDetectorConfig{
		Config: Config{},
		Enable: d.Enable,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
