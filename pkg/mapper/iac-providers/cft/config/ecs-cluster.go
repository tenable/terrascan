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
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/cft/functions"
)

// ClusterSettingsBlock holds config for settings attribute
type ClusterSettingsBlock struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CapacityProviderStrategyBlock holds config for default_capacity_provider_strategy attribute
type CapacityProviderStrategyBlock struct {
	Base             int    `json:"base"`
	CapacityProvider string `json:"capacity_provider"`
	Weight           int    `json:"weight"`
}

// LogConfigurationBlock holds config for log_configuration attribute
type LogConfigurationBlock struct {
	CloudWatchEncryptionEnabled bool   `json:"cloud_watch_encryption_enabled"`
	CloudWatchLogGroupName      string `json:"cloud_watch_log_group_name"`
	S3BucketName                string `json:"s3_bucket_name"`
	S3EncryptionEnabled         bool   `json:"s3_bucket_encryption_enabled"`
	S3KeyPrefix                 string `json:"s3_key_prefix"`
}

// ExecuteCommandConfiguration holds config for execute_command_configuration attribute
type ExecuteCommandConfiguration struct {
	KmsKeyID         string                  `json:"kms_key_id"`
	Logging          string                  `json:"logging"`
	LogConfiguration []LogConfigurationBlock `json:"log_configuration"`
}

// ConfigurationBlock holds config for configuration attribute
type ConfigurationBlock struct {
	ExecuteCommandConfig []ExecuteCommandConfiguration `json:"execute_command_configuration"`
}

// EcsClusterConfig holds config for aws_ecs_cluster resource
type EcsClusterConfig struct {
	Config
	ClusterName                     string                          `json:"name"`
	ClusterSettings                 []ClusterSettingsBlock          `json:"settings"`
	DefaultCapacityProviderStrategy []CapacityProviderStrategyBlock `json:"default_capacity_provider_strategy"`
	Configuration                   []ConfigurationBlock            `json:"configuration"`
}

// GetEcsClusterConfig returns config for aws_ecs_cluster resource
// aws_ecs_cluster
func GetEcsClusterConfig(e *ecs.Cluster) []AWSResourceConfig {
	var clusterSettingsData []ClusterSettingsBlock
	var capacityProviderStrategyData []CapacityProviderStrategyBlock

	clusterSettings := e.ClusterSettings
	clusterSettingsData = make([]ClusterSettingsBlock, len(clusterSettings))
	for i, clusterSetting := range clusterSettings {
		clusterSettingsData[i].Name = functions.GetVal(clusterSetting.Name)
		clusterSettingsData[i].Value = functions.GetVal(clusterSetting.Value)
	}

	defaultCapacityProviderStrategy := e.DefaultCapacityProviderStrategy
	capacityProviderStrategyData = make([]CapacityProviderStrategyBlock, len(defaultCapacityProviderStrategy))
	for i, defaultCapacityProviderStrategy := range defaultCapacityProviderStrategy {
		capacityProviderStrategyData[i].Base = functions.GetVal(defaultCapacityProviderStrategy.Base)
		capacityProviderStrategyData[i].CapacityProvider = functions.GetVal(defaultCapacityProviderStrategy.CapacityProvider)
		capacityProviderStrategyData[i].Weight = functions.GetVal(defaultCapacityProviderStrategy.Weight)
	}

	cf := EcsClusterConfig{
		Config: Config{
			Name: functions.GetVal(e.ClusterName),
			Tags: functions.PatchAWSTags(e.Tags),
		},
		ClusterName:                     functions.GetVal(e.ClusterName),
		ClusterSettings:                 clusterSettingsData,
		DefaultCapacityProviderStrategy: capacityProviderStrategyData,
	}

	cf.Configuration = setConfigurationBlock(e)

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}

func setConfigurationBlock(e *ecs.Cluster) []ConfigurationBlock {
	var configurationData []ConfigurationBlock

	if e.Configuration == nil {
		return configurationData
	}

	configurationData = make([]ConfigurationBlock, 1)
	configurationData[0].ExecuteCommandConfig = setExecCommandConfigBlock(e)

	return configurationData
}

func setExecCommandConfigBlock(e *ecs.Cluster) []ExecuteCommandConfiguration {
	var execCommandConfigData []ExecuteCommandConfiguration

	if e.Configuration.ExecuteCommandConfiguration == nil {
		return execCommandConfigData
	}

	execCommandConfigData = make([]ExecuteCommandConfiguration, 1)

	execCommandConfigData[0].KmsKeyID = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.KmsKeyId)
	execCommandConfigData[0].Logging = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.Logging)

	if e.Configuration.ExecuteCommandConfiguration.LogConfiguration == nil {
		return execCommandConfigData
	}

	execCommandConfigData[0].LogConfiguration = setLogConfigurationBlock(e)

	return execCommandConfigData
}

func setLogConfigurationBlock(e *ecs.Cluster) []LogConfigurationBlock {
	logConfigData := make([]LogConfigurationBlock, 1)

	logConfigData[0].S3BucketName = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3BucketName)
	logConfigData[0].S3KeyPrefix = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3KeyPrefix)
	logConfigData[0].S3EncryptionEnabled = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3EncryptionEnabled)
	logConfigData[0].CloudWatchLogGroupName = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchLogGroupName)
	logConfigData[0].CloudWatchEncryptionEnabled = functions.GetVal(e.Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchEncryptionEnabled)

	return logConfigData
}
