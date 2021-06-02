/*
    Copyright (C) 2021 Accurics, Inc.

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
	"github.com/awslabs/goformation/v4/cloudformation/cloudformation"
)

// CloudFormationStackConfig holds config for aws_cloudformation_stack
type CloudFormationStackConfig struct {
	Config
	TemplateURL      interface{} `json:"template_url"`
	NotificationARNs interface{} `json:"notification_arns"`
}

// GetCloudFormationStackConfig returns config for aws_cloudformation_stack
func GetCloudFormationStackConfig(s *cloudformation.Stack) []AWSResourceConfig {
	cf := CloudFormationStackConfig{
		Config: Config{
			Tags: s.Tags,
		},
	}
	if len(s.NotificationARNs) > 0 {
		cf.NotificationARNs = s.NotificationARNs
	}
	if len(s.TemplateURL) > 0 {
		cf.TemplateURL = s.TemplateURL
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}
