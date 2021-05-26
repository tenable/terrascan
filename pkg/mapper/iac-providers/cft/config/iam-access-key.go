package config

import (
	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

// IamAccessKeyConfig holds config for aws_iam_access_key
type IamAccessKeyConfig struct {
	Config
	UserName string `json:"user"`
	Status   string `json:"status"`
}

// GetIamAccessKeyConfig returns config for aws_iam_access_key
func GetIamAccessKeyConfig(r *iam.AccessKey) []AWSResourceConfig {
	cf := IamAccessKeyConfig{
		Config: Config{
			Name: r.UserName,
		},
		UserName: r.UserName,
		Status:   r.Status,
	}
	return []AWSResourceConfig{{Resource: cf}}
}
