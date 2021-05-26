package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/elasticloadbalancingv2"
)

// ElasticLoadBalancingV2ListenerConfig holds the config for aws_lb_listener
type ElasticLoadBalancingV2ListenerConfig struct {
	Config
	Protocol      string              `json:"protocol"`
	DefaultAction DefaultActionConfig `json:"default_action"`
}

// DefaultActionConfig holds the config for default_action attribute of aws_lb_listener
type DefaultActionConfig struct {
	RedirectConfig RedirectConfig `json:"redirect"`
}

// RedirectConfig holds the config for redirect attirbute of default_action
type RedirectConfig struct {
	Protocol string `json:"protocol"`
}

// GetElasticLoadBalancingV2ListenerConfig returns config for aws_lb_listener
func GetElasticLoadBalancingV2ListenerConfig(l *elasticloadbalancingv2.Listener) []AWSResourceConfig {
	// create a listener subresource per DefaultAction defined in cft
	// as only one default action per listener is possible in terraform
	resourceConfigs := make([]AWSResourceConfig, 0)

	for _, action := range l.DefaultActions {
		// DefaultActions are required
		cf := ElasticLoadBalancingV2ListenerConfig{
			Protocol: l.Protocol,
		}
		if action.RedirectConfig != nil {
			defaultAction := DefaultActionConfig{
				RedirectConfig: RedirectConfig{
					Protocol: action.RedirectConfig.Protocol,
				},
			}
			cf.DefaultAction = defaultAction
		}
		resourceConfigs = append(resourceConfigs, AWSResourceConfig{
			Resource: cf,
		})
	}

	return resourceConfigs
}
