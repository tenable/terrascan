package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

// IamPolicyConfig holds config for aws_iam_policy
type IamPolicyConfig struct {
	Config
	PolicyDocument string `json:"policy"`
	PolicyName     string `json:"name"`
}

// GetIamPolicyConfig returns config for aws_iam_policy
func GetIamPolicyConfig(r *iam.Policy) []AWSResourceConfig {
	cf := IamPolicyConfig{
		Config: Config{
			Name: r.PolicyName,
		},
		PolicyName: r.PolicyName,
	}
	policyDocument, err := json.Marshal(r.PolicyDocument)
	if err == nil {
		cf.PolicyDocument = string(policyDocument)
	}
	return []AWSResourceConfig{{Resource: cf}}
}
