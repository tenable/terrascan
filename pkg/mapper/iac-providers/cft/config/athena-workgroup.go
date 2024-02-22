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
	"github.com/awslabs/goformation/v7/cloudformation/athena"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// EncryptionConfigurationBlock holds config for encryption_configuration attribute
type EncryptionConfigurationBlock struct {
	EncryptionOption string `json:"encryption_option"`
	KmsKeyArn        string `json:"kms_key_arn"`
}

// ResultConfigurationBlock holds config for result_configuration attribute
type ResultConfigurationBlock struct {
	EncryptionConfiguration []EncryptionConfigurationBlock `json:"encryption_configuration"`
	OutputLocation          string                         `json:"output_location"`
}

// EngineerVersionBlock holds config for engine_version attribute
type EngineerVersionBlock struct {
	SelectedEngineVersion string `json:"selected_version"`
}

// WorkgroupConfigurationBlock holds config for configuration attribute
type WorkgroupConfigurationBlock struct {
	BytesScannedCutoffPerQuery      int                        `json:"bytes_scanned_cutoff_per_query,omitempty"`
	EnforceWorkgroupConfiguration   bool                       `json:"enforce_workgroup_configuration"`
	RequesterPaysEnabled            bool                       `json:"requester_pays_enabled"`
	PublishCloudwatchMetricsEnabled bool                       `json:"publish_cloudwatch_metrics_enabled"`
	EngineVersion                   []EngineerVersionBlock     `json:"engine_version"`
	ResultConfiguration             []ResultConfigurationBlock `json:"result_configuration"`
}

// AthenaWorkGroupConfig holds config for aws_athena_workgroup resource
type AthenaWorkGroupConfig struct {
	Config
	Name          string                        `json:"name"`
	Configuration []WorkgroupConfigurationBlock `json:"configuration"`
}

// GetAthenaWorkGroupConfig returns config for aws_athena_workgroup resource
// aws_athena_workgroup
func GetAthenaWorkGroupConfig(w *athena.WorkGroup) []AWSResourceConfig {
	var workGroupConfig []WorkgroupConfigurationBlock

	if w.WorkGroupConfiguration != nil {
		workGroupConfig = make([]WorkgroupConfigurationBlock, 1)

		workGroupConfig[0].BytesScannedCutoffPerQuery = functions.GetVal(w.WorkGroupConfiguration.BytesScannedCutoffPerQuery)
		workGroupConfig[0].EnforceWorkgroupConfiguration = functions.GetVal(w.WorkGroupConfiguration.EnforceWorkGroupConfiguration)
		workGroupConfig[0].RequesterPaysEnabled = functions.GetVal(w.WorkGroupConfiguration.RequesterPaysEnabled)
		workGroupConfig[0].PublishCloudwatchMetricsEnabled = functions.GetVal(w.WorkGroupConfiguration.PublishCloudWatchMetricsEnabled)

		if w.WorkGroupConfiguration.EngineVersion != nil {
			engineConfig := make([]EngineerVersionBlock, 1)
			engineConfig[0].SelectedEngineVersion = functions.GetVal(w.WorkGroupConfiguration.EngineVersion.SelectedEngineVersion)
			workGroupConfig[0].EngineVersion = engineConfig
		}

		if w.WorkGroupConfiguration.ResultConfiguration != nil {
			resultConfig := make([]ResultConfigurationBlock, 1)
			resultConfig[0].OutputLocation = functions.GetVal(w.WorkGroupConfiguration.ResultConfiguration.OutputLocation)

			if w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration != nil {
				encryptionConfig := make([]EncryptionConfigurationBlock, 1)
				encryptionConfig[0].EncryptionOption = w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration.EncryptionOption
				encryptionConfig[0].KmsKeyArn = functions.GetVal(w.WorkGroupConfiguration.ResultConfiguration.EncryptionConfiguration.KmsKey)

				resultConfig[0].EncryptionConfiguration = encryptionConfig
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
