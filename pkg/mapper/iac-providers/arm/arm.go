/*
    Copyright (C) 2020 Accurics, Inc.

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
	"encoding/json"
	"errors"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/mapper/core"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/config"
	fn "github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/functions"
	"github.com/accurics/terrascan/pkg/mapper/iac-providers/arm/types"
	"github.com/accurics/terrascan/pkg/utils"
)

const errUnsupportedDoc = "unsupported document type"

type armMapper struct {
	templateParameters map[string]interface{}
}

// Mapper returns an ARM mapper for given template schema
func Mapper() core.Mapper {
	return armMapper{}
}

// Map transforms the provider specific template to terrascan native format.
func (m armMapper) Map(doc *utils.IacDocument, params ...map[string]interface{}) (output.AllResourceConfigs, error) {
	allRC := make(map[string][]output.ResourceConfig)
	template, err := extractTemplate(doc)
	if err != nil {
		return nil, err
	}

	// set template parameters with default values if not found
	m.templateParameters = params[0]
	for key, param := range template.Parameters {
		if _, ok := m.templateParameters[key]; !ok {
			m.templateParameters[key] = param.DefaultValue
		}
	}

	// transform each resource and generate config
	for _, r := range template.Resources {
		// skip if resource does not have a mapping
		if _, ok := types.ResourceTypes[r.Type]; !ok {
			continue
		}

		rc := output.ResourceConfig{
			Name:   fn.LookUp(template.Variables, m.templateParameters, r.Name).(string),
			Source: doc.FilePath,
			Line:   doc.StartLine,
			Type:   types.ResourceTypes[r.Type],
		}

		rc.ID = rc.Type + "." + rc.Name
		fn.ResourceIDs[r.Type] = rc.ID
		rc.Config = m.mapConfigForResource(r, template.Variables)
		allRC[rc.Type] = append(allRC[rc.Type], rc)
	}
	return allRC, nil
}

func extractTemplate(doc *utils.IacDocument) (*types.Template, error) {
	if doc.Type == utils.JSONDoc {
		var t types.Template
		err := json.Unmarshal(doc.Data, &t)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}
	return nil, errors.New(errUnsupportedDoc)
}

func (m armMapper) mapConfigForResource(r types.Resource, vars map[string]interface{}) interface{} {
	switch types.ResourceTypes[r.Type] {
	case types.AzureRMKeyVault:
		return config.KeyVaultConfig(r, m.templateParameters)
	case types.AzureRMKeyVaultSecret:
		return config.KeyVaultSecretConfig(r, m.templateParameters)
	case types.AzureRMKeyVaultKey:
		return config.KeyVaultKeyConfig(r, m.templateParameters)
	case types.AzureRMApplicationGateway:
		return config.ApplicationGatewayConfig(r, m.templateParameters)
	case types.AzureRMMonitorDiagnosticSetting:
		return config.DiagnosticSettingConfig(r, vars, m.templateParameters)
	case types.AzureRMKubernetesCluster:
		return config.KubernetesClusterConfig(r, vars, m.templateParameters)
	case types.AzureRMManagedDisk:
		return config.ManagedDiskConfig(r, vars, m.templateParameters)
	case types.AzureRMCosmosDBAccount:
		return config.CosmosDBAccountConfig(r, m.templateParameters)
	case types.AzureRMContainerRegistry:
		return config.ContainerRegistryConfig(r, m.templateParameters)
	case types.AzureRMManagementLock:
		return config.ManagementLockConfig(r, vars, m.templateParameters)
	case types.AzureRMRoleAssignment:
		return config.RoleAssignmentConfig(r, vars, m.templateParameters)
	case types.AzureRMMSSQLServer:
		return config.MSSQLServerConfig(r, vars, m.templateParameters)
	case types.AzureRMMySQLServer:
		return config.MySQLServerConfig(r, vars, m.templateParameters)
	case types.AzureRMNetworkWatcherFlowLog:
		return config.NetworkWatcherFlowLogConfig(r, vars, m.templateParameters)
	case types.AzureRMResourceGroup:
		return config.ResourceGroupConfig(r, m.templateParameters)
	case types.AzureRMSecurityCenterContact:
		return config.SecurityCenterContactConfig(r, m.templateParameters)
	case types.AzureRMSecurityCenterSubscriptionPricing:
		return config.SecurityCenterSubscriptionPricingConfig(r, m.templateParameters)
	case types.AzureRMSQLActiveDirectoryAdministrator:
		return config.SQLActiveDirectoryAdministratorConfig(r, vars, m.templateParameters)
	case types.AzureRMNetworkSecurityRule:
		return config.NetworkSecurityRuleConfig(r, m.templateParameters)
	case types.AzureRMPostgreSQLConfiguration:
		return config.PostgreSQLConfigurationConfig(r, m.templateParameters)
	case types.AzureRMPostgreSQLServers:
		return config.PostgreSQLServerConfig(r, vars, m.templateParameters)
	case types.AzureRMRedisCache:
		return config.RedisCacheConfig(r, m.templateParameters)
	case types.AzureRMRedisFirewallRule:
		return config.RedisFirewallRuleConfig(r, m.templateParameters)
	case types.AzureRMStorageAccount:
		return config.StorageAccountConfig(r, vars, m.templateParameters)
	case types.AzureRMSQLFirewallRule:
		return config.SQLFirewallRuleConfig(r, m.templateParameters)
	case types.AzureRMStorageContainer:
		return config.StorageContainerConfig(r, m.templateParameters)
	case types.AzureRMVirtualMachine:
		return config.VirtualMachineConfig(r, m.templateParameters)
	case types.AzureRMVirtualNetwork:
		return config.VirtualNetworkConfig(r, vars, m.templateParameters)
	}
	return nil
}
