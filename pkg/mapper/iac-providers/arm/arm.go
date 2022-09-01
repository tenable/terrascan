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

package arm

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/mapper/core"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/config"
	fn "github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/tenable/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/tenable/terrascan/pkg/utils"
)

type armMapper struct{}

// Mapper returns an ARM mapper for given template schema
func Mapper() core.Mapper {
	return armMapper{}
}

// Map transforms the provider specific template to terrascan native format.
func (m armMapper) Map(resource interface{}, params ...map[string]interface{}) ([]output.ResourceConfig, error) {
	r, ok := resource.(types.Resource)
	if !ok {
		return nil, errors.New("failed to cast resource into types.Resource")
	}

	config := output.ResourceConfig{
		SkipRules: make([]output.SkipRule, 0),
	}

	// add skipRules if available
	if r.Tags != nil {
		skipRules := utils.ReadSkipRulesFromMap(r.Tags, config.ID)
		if skipRules != nil {
			config.SkipRules = append(config.SkipRules, skipRules...)
		}
	}

	variables := params[0]
	parameters := params[1]
	configName, ok := fn.LookUp(variables, parameters, r.Name).(string)
	if !ok {
		configName = uuid.NewString()
	}
	config.Name = configName
	config.Type = types.ResourceTypes[r.Type]
	config.ID = config.Type + "." + config.Name

	fn.ResourceIDs[r.Type] = config.ID
	config.Config = m.mapConfigForResource(r, variables, parameters)

	return []output.ResourceConfig{config}, nil
}

func (m armMapper) mapConfigForResource(r types.Resource, vars, params map[string]interface{}) interface{} {
	switch types.ResourceTypes[r.Type] {
	case types.AzureRMKeyVault:
		return config.KeyVaultConfig(r, params)
	case types.AzureRMKeyVaultSecret:
		return config.KeyVaultSecretConfig(r, params)
	case types.AzureRMKeyVaultKey:
		return config.KeyVaultKeyConfig(r, params)
	case types.AzureRMApplicationGateway:
		return config.ApplicationGatewayConfig(r, params)
	case types.AzureRMMonitorDiagnosticSetting:
		return config.DiagnosticSettingConfig(r, vars, params)
	case types.AzureRMKubernetesCluster:
		return config.KubernetesClusterConfig(r, vars, params)
	case types.AzureRMManagedDisk:
		return config.ManagedDiskConfig(r, vars, params)
	case types.AzureRMCosmosDBAccount:
		return config.CosmosDBAccountConfig(r, params)
	case types.AzureRMContainerRegistry:
		return config.ContainerRegistryConfig(r, params)
	case types.AzureRMManagementLock:
		return config.ManagementLockConfig(r, vars, params)
	case types.AzureRMRoleAssignment:
		return config.RoleAssignmentConfig(r, vars, params)
	case types.AzureRMMSSQLServer:
		return config.MSSQLServerConfig(r, vars, params)
	case types.AzureRMMySQLServer:
		return config.MySQLServerConfig(r, vars, params)
	case types.AzureRMNetworkWatcherFlowLog:
		return config.NetworkWatcherFlowLogConfig(r, vars, params)
	case types.AzureRMResourceGroup:
		return config.ResourceGroupConfig(r, params)
	case types.AzureRMSecurityCenterContact:
		return config.SecurityCenterContactConfig(r, params)
	case types.AzureRMSecurityCenterSubscriptionPricing:
		return config.SecurityCenterSubscriptionPricingConfig(r, params)
	case types.AzureRMSQLActiveDirectoryAdministrator:
		return config.SQLActiveDirectoryAdministratorConfig(r, vars, params)
	case types.AzureRMNetworkSecurityRule:
		return config.NetworkSecurityRuleConfig(r, params)
	case types.AzureRMPostgreSQLConfiguration:
		return config.PostgreSQLConfigurationConfig(r, params)
	case types.AzureRMPostgreSQLServers:
		return config.PostgreSQLServerConfig(r, vars, params)
	case types.AzureRMRedisCache:
		return config.RedisCacheConfig(r, params)
	case types.AzureRMRedisFirewallRule:
		return config.RedisFirewallRuleConfig(r, params)
	case types.AzureRMStorageAccount:
		return config.StorageAccountConfig(r, vars, params)
	case types.AzureRMSQLFirewallRule:
		return config.SQLFirewallRuleConfig(r, params)
	case types.AzureRMStorageContainer:
		return config.StorageContainerConfig(r, params)
	case types.AzureRMVirtualMachine:
		return config.VirtualMachineConfig(r, params)
	case types.AzureRMVirtualNetwork:
		return config.VirtualNetworkConfig(r, vars, params)
	case types.AzureRMDeployments:
		return config.DeploymentsConfig(r, vars, params)
	}
	return nil
}
