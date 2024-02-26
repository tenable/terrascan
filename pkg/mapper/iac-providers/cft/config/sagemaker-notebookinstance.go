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
	"github.com/awslabs/goformation/v7/cloudformation/sagemaker"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

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
// aws_sagemaker_notebook_instance
func GetSagemakerNotebookInstanceConfig(n *sagemaker.NotebookInstance) []AWSResourceConfig {
	cf := SagemakerNotebookInstanceConfig{
		Config: Config{
			Name: functions.GetVal(n.NotebookInstanceName),
			Tags: functions.PatchAWSTags(n.Tags),
		},
		Name:                 functions.GetVal(n.NotebookInstanceName),
		RoleARN:              n.RoleArn,
		InstanceType:         n.InstanceType,
		KMSKeyID:             functions.GetVal(n.KmsKeyId),
		DirectInternetAccess: functions.GetVal(n.DirectInternetAccess),
		RootAccess:           functions.GetVal(n.RootAccess),
	}

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: n.AWSCloudFormationMetadata,
	}}
}
