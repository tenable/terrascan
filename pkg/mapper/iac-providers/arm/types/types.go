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

package types

// ResourceTypes holds mapping for ARM resource types to TF types
var ResourceTypes = map[string]string{
	"Microsoft.KeyVault/vaults":                                 AzureRMKeyVault,
	"Microsoft.KeyVault/vaults/keys":                            AzureRMKeyVaultKey,
	"Microsoft.KeyVault/vaults/secrets":                         AzureRMKeyVaultSecret,
	"Microsoft.Network/applicationGateways":                     AzureRMApplicationGateway,
	"Microsoft.Insights/diagnosticsettings":                     AzureRMMonitorDiagnosticSetting,
	"Microsoft.ContainerService/managedClusters":                AzureRMKubernetesCluster,
	"Microsoft.Compute/disks":                                   AzureRMManagedDisk,
	"Microsoft.DocumentDB/databaseAccounts":                     AzureRMCosmosDBAccount,
	"Microsoft.ContainerRegistry/registries":                    AzureRMContainerRegistry,
	"Microsoft.Authorization/locks":                             AzureRMManagementLock,
	"Microsoft.Authorization/roleAssignments":                   AzureRMRoleAssignment,
	"Microsoft.Sql/servers":                                     AzureRMMSSQLServer,
	"Microsoft.DBforMySQL/servers":                              AzureRMMySQLServer,
	"Microsoft.Network/networkWatchers/flowLogs":                AzureRMNetworkWatcherFlowLog,
	"Microsoft.Resources/resourceGroups":                        AzureRMResourceGroup,
	"Microsoft.Security/securityContacts":                       AzureRMSecurityCenterContact,
	"Microsoft.Security/pricings":                               AzureRMSecurityCenterSubscriptionPricing,
	"Microsoft.Sql/servers/administrators":                      AzureRMSQLActiveDirectoryAdministrator,
	"Microsoft.Network/networkSecurityGroups/securityRules":     AzureRMNetworkSecurityRule,
	"Microsoft.DBforPostgreSQL/servers/configurations":          AzureRMPostgreSQLConfiguration,
	"Microsoft.DBforPostgreSQL/servers":                         AzureRMPostgreSQLServers,
	"Microsoft.Cache/redis":                                     AzureRMRedisCache,
	"Microsoft.Cache/redis/firewallRules":                       AzureRMRedisFirewallRule,
	"Microsoft.Storage/storageAccounts":                         AzureRMStorageAccount,
	"Microsoft.Sql/servers/firewallRules":                       AzureRMSQLFirewallRule,
	"Microsoft.Storage/storageAccounts/blobServices/containers": AzureRMStorageContainer,
	"Microsoft.Compute/virtualMachines":                         AzureRMVirtualMachine,
	"Microsoft.Network/virtualNetworks":                         AzureRMVirtualNetwork,
	"Microsoft.Sql/servers/auditingSettings":                    AzureRMMSSQLDBAuditingPolicy,
	"Microsoft.Cache/redis/patchSchedules":                      AzureRMPatchSchedule,
	"Microsoft.Resources/deployments":                           AzureRMDeployments,
}

// ARM equivalent TF resource types
const (
	AzureRMKeyVaultSecret                    = "azurerm_key_vault_secret"
	AzureRMKeyVault                          = "azurerm_key_vault"
	AzureRMKeyVaultKey                       = "azurerm_key_vault_key"
	AzureRMApplicationGateway                = "azurerm_application_gateway"
	AzureRMMonitorDiagnosticSetting          = "azurerm_monitor_diagnostic_setting"
	AzureRMKubernetesCluster                 = "azurerm_kubernetes_cluster"
	AzureRMManagedDisk                       = "azurerm_managed_disk"
	AzureRMCosmosDBAccount                   = "azurerm_cosmosdb_account"
	AzureRMContainerRegistry                 = "azurerm_container_registry"
	AzureRMManagementLock                    = "azurerm_management_lock"
	AzureRMRoleAssignment                    = "azurerm_role_assignment"
	AzureRMMSSQLServer                       = "azurerm_mssql_server"
	AzureRMMySQLServer                       = "azurerm_mysql_server"
	AzureRMNetworkWatcherFlowLog             = "azurerm_network_watcher_flow_log"
	AzureRMResourceGroup                     = "azurerm_resource_group"
	AzureRMSecurityCenterContact             = "azurerm_security_center_contact"
	AzureRMSecurityCenterSubscriptionPricing = "azurerm_security_center_subscription_pricing"
	AzureRMSQLActiveDirectoryAdministrator   = "azurerm_sql_active_directory_administrator"
	AzureRMNetworkSecurityRule               = "azurerm_network_security_rule"
	AzureRMPostgreSQLConfiguration           = "azurerm_postgresql_configuration"
	AzureRMPostgreSQLServers                 = "azurerm_postgresql_server"
	AzureRMRedisCache                        = "azurerm_redis_cache"
	AzureRMRedisFirewallRule                 = "azurerm_redis_firewall_rule"
	AzureRMStorageAccount                    = "azurerm_storage_account"
	AzureRMSQLFirewallRule                   = "azurerm_sql_firewall_rule"
	AzureRMStorageContainer                  = "azurerm_storage_container"
	AzureRMVirtualMachine                    = "azurerm_virtual_machine"
	AzureRMVirtualNetwork                    = "azurerm_virtual_network"
	AzureRMMSSQLDBAuditingPolicy             = "extended_auditing_policy"
	AzureRMPatchSchedule                     = "patch_schedule"
	AzureRMDeployments                       = "azurerm_resource_group_template_deployment"
)
