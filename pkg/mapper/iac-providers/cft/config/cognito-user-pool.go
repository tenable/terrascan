/*
    Copyright (C) 2022 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANT IES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package config

import "github.com/awslabs/goformation/v4/cloudformation/cognito"

// PasswordPolicyBlock holds config for PasswordPolicy
type PasswordPolicyBlock struct {
	MinimumLength                 int  `json:"minimum_length"`
	RequireLowercase              bool `json:"require_lowercase"`
	RequireUppercase              bool `json:"require_uppercase"`
	RequireNumbers                bool `json:"require_numbers"`
	RequireSymbols                bool `json:"require_symbols"`
	TemporaryPasswordValidityDays int  `json:"temporary_password_validity_days"`
}

// CognitoUserPoolConfig holds config for CognitoUserPool
type CognitoUserPoolConfig struct {
	Config
	Name           string                `json:"name"`
	PasswordPolicy []PasswordPolicyBlock `json:"password_policy"`
}

// GetCognitoUserPoolConfig returns config for CognitoUserPool
func GetCognitoUserPoolConfig(u *cognito.UserPool) []AWSResourceConfig {
	var passwordPolicy []PasswordPolicyBlock
	if u.Policies != nil && u.Policies.PasswordPolicy != nil {
		passwordPolicy := make([]PasswordPolicyBlock, 1)

		passwordPolicy[0].MinimumLength = u.Policies.PasswordPolicy.MinimumLength
		passwordPolicy[0].RequireLowercase = u.Policies.PasswordPolicy.RequireLowercase
		passwordPolicy[0].RequireUppercase = u.Policies.PasswordPolicy.RequireUppercase
		passwordPolicy[0].RequireNumbers = u.Policies.PasswordPolicy.RequireNumbers
		passwordPolicy[0].RequireSymbols = u.Policies.PasswordPolicy.RequireSymbols
		passwordPolicy[0].TemporaryPasswordValidityDays = u.Policies.PasswordPolicy.TemporaryPasswordValidityDays
	}

	cf := CognitoUserPoolConfig{
		Config: Config{
			Name: u.UserPoolName,
		},
		Name:           u.UserPoolName,
		PasswordPolicy: passwordPolicy,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: u.AWSCloudFormationMetadata,
	}}
}
