package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/elasticloadbalancingv2"
)

// ElasticLoadBalancingV2TargetGroupConfig holds the config for aws_lb_target_group
type ElasticLoadBalancingV2TargetGroupConfig struct {
	Config
	Protocol string `json:"protocol"`
}

// GetElasticLoadBalancingV2TargetGroupConfig returns config for aws_lb_target_group
func GetElasticLoadBalancingV2TargetGroupConfig(l *elasticloadbalancingv2.TargetGroup) []AWSResourceConfig {
	// create a listener subresource per DefaultAction defined in cft
	// as only one default action per listener is possible in terraform
	cf := ElasticLoadBalancingV2TargetGroupConfig{
		Config: Config{
			Name: l.Name,
			Tags: l.Tags,
		},
		Protocol: l.Protocol,
	}

	return []AWSResourceConfig{{Resource: cf}}
}
