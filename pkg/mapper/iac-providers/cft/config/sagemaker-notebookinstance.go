/*
    Copyright (C) 2022 Accurics, Inc.

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

import "github.com/awslabs/goformation/v4/cloudformation/sagemaker"

// SagemakerNotebookInstanceConfig holds config for SagemakerNotebookInstance
type SagemakerNotebookInstanceConfig struct {
	Config
	Name                 string `json:"name"`
	RoleARN              string `json:"role_arn"`
	InstanceType         string `json:"instance_type"`
	KMSKeyID             string `json:"kms_key_id"`
	DirectInternetAccess string `json:"direct_internet_access"`
	RootAccess           string `json:"root_access"`
}

// GetSagemakerNotebookInstanceConfig returns config for SagemakerNotebookInstance
func GetSagemakerNotebookInstanceConfig(n *sagemaker.NotebookInstance) []AWSResourceConfig {
	cf := SagemakerNotebookInstanceConfig{
		Config: Config{
			Name: n.NotebookInstanceName,
			Tags: n.Tags,
		},
		Name:                 n.NotebookInstanceName,
		RoleARN:              n.RoleArn,
		InstanceType:         n.InstanceType,
		KMSKeyID:             n.KmsKeyId,
		DirectInternetAccess: n.DirectInternetAccess,
		RootAccess:           n.RootAccess,
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: n.AWSCloudFormationMetadata,
	}}
}
