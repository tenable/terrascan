/*
    Copyright (C) 2022 Tenable, Inc.

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
	"github.com/awslabs/goformation/v7/cloudformation/iam"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
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
		Status:   functions.GetVal(r.Status),
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: r.AWSCloudFormationMetadata,
	}}
}
