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
	"github.com/awslabs/goformation/v5/cloudformation/athena"
)

// EncryptionConfigurationBlock holds config for encryption_configuration attribute
type EncryptionConfigurationBlock struct {
	EncryptionOption string `json:"encryption_option"`
	KmsKeyArn        string `json:"kms_key_arn"`
}

// ResultConfigurationBlock holds config for ResultConfigurationBlock attribute
type ResultConfigurationBlock struct {
	EncryptionConfiguration []EncryptionConfigurationBlock `json:"encryption_configuration"`
	OutputLocation          string                         `json:"output_location"`
}

type EnginerVersionBlock struct {
	SelectedEngineVersion string `json:"selected_version"`
}

// WorkgroupConfigurationBlock holds config for configuration attribute
type WorkgroupConfigurationBlock struct {
	BytesScannedCutoffPerQuery      int                        `json:"bytes_scanned_cutoff_per_query,omitempty"`
	EnforceWorkgroupConfiguration   bool                       `json:"enforce_workgroup_configuration"`
	RequesterPaysEnabled            bool                       `json:"requester_pays_enabled"`
	PublishCloudwatchMetricsEnabled bool                       `json:"publish_cloudwatch_metrics_enabled"`
	EngineVersion                   []EnginerVersionBlock      `json:"engine_version"`
	ResultConfiguration             []ResultConfigurationBlock `json:"result_configuration"`
}

// AthenaWorkGroupConfig holds config for aws_athena_workgroup resource
type AthenaWorkGroupConfig struct {
	Config
	Name          string                        `json:"name"`
	Configuration []WorkgroupConfigurationBlock `json:"configuration"`
}

// GetAthenaWorkGroupConfig returns config for aws_athena_workgroup resource
func GetAthenaWorkGroupConfig(w *athena.WorkGroup) []AWSResourceConfig {
	var workGroupConfig []WorkgroupConfigurationBlock

	if w.WorkGroupConfiguration != nil {
		workGroupConfig = make([]WorkgroupConfigurationBlock, 1)

		workGroupConfig[0].BytesScannedCutoffPerQuery = w.WorkGroupConfiguration.BytesScannedCutoffPerQuery
		workGroupConfig[0].EnforceWorkgroupConfiguration = w.WorkGroupConfiguration.EnforceWorkGroupConfiguration
		workGroupConfig[0].RequesterPaysEnabled = w.WorkGroupConfiguration.RequesterPaysEnabled
		workGroupConfig[0].PublishCloudwatchMetricsEnabled = w.WorkGroupConfiguration.PublishCloudWatchMetricsEnabled

		if w.WorkGroupConfiguration.EngineVersion != nil {
			engineConfig := make([]EnginerVersionBlock, 1)
			engineConfig[0].SelectedEngineVersion = w.WorkGroupConfiguration.EngineVersion.SelectedEngineVersion
			workGroupConfig[0].EngineVersion = engineConfig
		}

		if w.WorkGroupConfiguration.ResultConfiguration != nil {
			resultConfig := make([]ResultConfigurationBlock, 1)
			resultConfig[0].OutputLocation = w.WorkGroupConfiguration.ResultConfiguration.OutputLocation

			if w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration != nil {
				encryptionCofig := make([]EncryptionConfigurationBlock, 1)
				encryptionCofig[0].EncryptionOption = w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration.EncryptionOption
				encryptionCofig[0].KmsKeyArn = w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration.KmsKey

				resultConfig[0].EncryptionConfiguration = encryptionCofig
			}

			workGroupConfig[0].ResultConfiguration = resultConfig
		}
	}

	cf := AthenaWorkGroupConfig{
		Config: Config{
			Name: w.Name,
			Tags: w.Tags,
		},
		Name: w.Name,
	}

	cf.Configuration = workGroupConfig

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: w.AWSCloudFormationMetadata,
	}}
}
