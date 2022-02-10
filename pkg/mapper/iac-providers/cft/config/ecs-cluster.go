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

import "github.com/awslabs/goformation/v5/cloudformation/ecs"

type ClusterSettingsBlock struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CapacityProviderStrategyBlock struct {
	Base             int    `json:"base"`
	CapacityProvider string `json:"capacity_provider"`
	Weight           int    `json:"weight"`
}

type LogConfigurationBlock struct {
	CloudWatchEncryptionEnabled bool   `json:"cloud_watch_encryption_enabled"`
	CloudWatchLogGroupName      string `json:"cloud_watch_log_group_name"`
	S3BucketName                string `json:"s3_bucket_name"`
	S3EncryptionEnabled         bool   `json:"s3_bucket_encryption_enabled"`
	S3KeyPrefix                 string `json:"s3_key_prefix"`
}

type ExecuteCommandConfiguration struct {
	KmsKeyId         string                  `json:"kms_key_id"`
	Logging          string                  `json:"logging"`
	LogConfiguration []LogConfigurationBlock `json:"log_configuration"`
}

type ConfigurationBlock struct {
	ExecuteCommandConfig []ExecuteCommandConfiguration `json:"execute_command_configuration"`
}

type EcsClusterConfig struct {
	Config
	ClusterName                     string                          `json:"name"`
	ClusterSettings                 []ClusterSettingsBlock          `json:"settings"`
	DefaultCapacityProviderStrategy []CapacityProviderStrategyBlock `json:"default_capacity_provider_strategy"`
	Configuration                   []ConfigurationBlock            `json:"configuration"`
}

func GetEcsClusterConfig(e *ecs.Cluster) []AWSResourceConfig {

	clusterSettingsData := make([]ClusterSettingsBlock, len(e.ClusterSettings))
	for i := 0; i < len(e.ClusterSettings); i++ {
		clusterSettingsData[i].Name = e.ClusterSettings[i].Name
		clusterSettingsData[i].Value = e.ClusterSettings[i].Value
	}

	capacityProviderStrategyData := make([]CapacityProviderStrategyBlock, len(e.DefaultCapacityProviderStrategy))
	for i := 0; i < len(e.DefaultCapacityProviderStrategy); i++ {
		capacityProviderStrategyData[i].Base = e.DefaultCapacityProviderStrategy[i].Base
		capacityProviderStrategyData[i].CapacityProvider = e.DefaultCapacityProviderStrategy[i].CapacityProvider
		capacityProviderStrategyData[i].Weight = e.DefaultCapacityProviderStrategy[i].Weight
	}

	cf := EcsClusterConfig{
		Config: Config{
			Name: e.ClusterName,
			Tags: e.Tags,
		},
		ClusterName:                     e.ClusterName,
		ClusterSettings:                 clusterSettingsData,
		DefaultCapacityProviderStrategy: capacityProviderStrategyData,
	}

	config := setConfigurationBlock(e)
	cf.Configuration = config

	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: e.AWSCloudFormationMetadata,
	}}
}

func setConfigurationBlock(e *ecs.Cluster) []ConfigurationBlock {
	configurationData := make([]ConfigurationBlock, 1)
	configurationData[0].ExecuteCommandConfig = setExecCommandConfigBlock(e)

	return configurationData
}

func setExecCommandConfigBlock(e *ecs.Cluster) []ExecuteCommandConfiguration {

	execCommandConfigData := make([]ExecuteCommandConfiguration, 1)

	execCommandConfigData[0].KmsKeyId = e.Configuration.ExecuteCommandConfiguration.KmsKeyId
	execCommandConfigData[0].Logging = e.Configuration.ExecuteCommandConfiguration.Logging
	execCommandConfigData[0].LogConfiguration = setLogConfigurationBlock(e)

	return execCommandConfigData
}

func setLogConfigurationBlock(e *ecs.Cluster) []LogConfigurationBlock {

	logConfigData := make([]LogConfigurationBlock, 1)
	logConfigData[0].S3BucketName = e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3BucketName
	logConfigData[0].S3KeyPrefix = e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3KeyPrefix
	logConfigData[0].S3EncryptionEnabled = e.Configuration.ExecuteCommandConfiguration.LogConfiguration.S3EncryptionEnabled
	logConfigData[0].CloudWatchLogGroupName = e.Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchLogGroupName
	logConfigData[0].CloudWatchEncryptionEnabled = e.Configuration.ExecuteCommandConfiguration.LogConfiguration.CloudWatchEncryptionEnabled

	return logConfigData
}
