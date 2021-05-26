package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/config"
)

// AWSConfigConfigRuleConfig holds config for aws_config_config_rule
type AWSConfigConfigRuleConfig struct {
	Config
	Source interface{} `json:"source"`
}

// GetConfigConfigRuleConfig returns config for aws_config_config_rule
func GetConfigConfigRuleConfig(c *config.ConfigRule) []AWSResourceConfig {
	cf := AWSConfigConfigRuleConfig{
		Config: Config{Name: c.ConfigRuleName},
	}
	if c.Source != nil {
		sources := make([]map[string]interface{}, 0)
		source := make(map[string]interface{})
		source["source_identifier"] = c.Source.SourceIdentifier
		sources = append(sources, source)
		if len(sources) > 0 {
			cf.Source = sources
		}
	}

	return []AWSResourceConfig{{Resource: cf}}
}
