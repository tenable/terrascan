package config

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation/iam"
)

const (
	// IamUserLoginProfile represents the subresource aws_iam_user_login_profile for attribute LoginProfile
	IamUserLoginProfile = "LoginProfile"
	// IamUserPolicy represents the subresource aws_iam_user_policy for the attribute policy
	IamUserPolicy = "Policy"
)

// IamUserLoginProfileConfig holds the config for aws_iam_user_login_profile
type IamUserLoginProfileConfig struct {
	Config
	PasswordResetRequired bool `json:"password_reset_required"`
}

// IamUserPolicyConfig holds the config for aws_iam_user_policy
type IamUserPolicyConfig struct {
	Config
	PolicyName     string `json:"name"`
	PolicyDocument string `json:"policy"`
}

// IamUserConfig holds the config for aws_iam_user
type IamUserConfig struct {
	Config
	UserName string `json:"name"`
}

// GetIamUserConfig returns the config for aws_iam_user, aws_iam_user_policy, aws_iam_user_login_profile
func GetIamUserConfig(i *iam.User) []AWSResourceConfig {

	resourceConfigs := make([]AWSResourceConfig, 0)

	// add aws_iam_user
	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Resource: IamUserConfig{
			UserName: i.UserName,
		},
	})

	// add aws_iam_user_login_profile
	resourceConfigs = append(resourceConfigs, AWSResourceConfig{
		Type: IamUserLoginProfile,
		Name: i.UserName,
		Resource: IamUserLoginProfileConfig{
			PasswordResetRequired: i.LoginProfile.PasswordResetRequired,
		},
	})

	// add aws_iam_user_policy
	if i.Policies != nil {
		for _, policy := range i.Policies {
			pc := IamUserPolicyConfig{
				Config: Config{
					Name: policy.PolicyName,
				},
				PolicyName: policy.PolicyName,
			}
			policyDocument, err := json.Marshal(policy.PolicyDocument)
			if err == nil {
				pc.PolicyDocument = string(policyDocument)
			}
			resourceConfigs = append(resourceConfigs, AWSResourceConfig{
				Type:     IamUserPolicy,
				Name:     policy.PolicyName,
				Resource: pc,
			})
		}
	}

	return resourceConfigs
}
