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
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/cft/functions"
	"github.com/awslabs/goformation/v5/cloudformation/cloudformation"
)

// CloudFormationStackConfig holds config for aws_cloudformation_stack
type CloudFormationStackConfig struct {
	Config
	TemplateURL      string            `json:"template_url"`
	NotificationARNs interface{}       `json:"notification_arns"`
	Parameters       map[string]string `json:"-"`
	TemplateData     []byte            `json:"-"`
}

// GetCloudFormationStackConfig returns config for aws_cloudformation_stack
func GetCloudFormationStackConfig(s *cloudformation.Stack) []AWSResourceConfig {
	cf := CloudFormationStackConfig{
		Config:           Config{Tags: s.Tags},
		TemplateURL:      "",
		NotificationARNs: nil,
		TemplateData:     []byte{},
	}

	if len(s.NotificationARNs) > 0 {
		cf.NotificationARNs = s.NotificationARNs
	}

	// Add and resolve template URL
	if len(s.TemplateURL) > 0 {
		cf.TemplateURL = s.TemplateURL

		templateData, err := fn.DownloadBucketObj(s.TemplateURL)
		if err == nil {
			cf.TemplateData = templateData
		}
	}

	// Add Parameters for propogation to the nested Stack
	if s.Parameters != nil {
		cf.Parameters = s.Parameters
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: s.AWSCloudFormationMetadata,
	}}
}
