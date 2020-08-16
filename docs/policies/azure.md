
### azurerm_virtual_machine
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure that at least one Network Security Group is attached to all VMs and subnets that are public | accurics.azure.NS.18 |


### azurerm_storage_container
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | azure | HIGH | Anonymous, public read access to a container and its blobs can be enabled in Azure Blob storage. This is only recommended if absolutely necessary. | accurics.azure.IAM.368 |


### azurerm_sql_firewall_rule
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Restrict Azure SQL Server accessibility to a minimal address range | accurics.azure.NS.169 |
| Network Security | azure | HIGH | Ensure entire Azure infrastructure doesn't have access to Azure SQL ServerEnsure entire Azure infrastructure doesn't have access to Azure SQL Server | accurics.azure.NS.5 |
| Network Security | azure | HIGH | Ensure that no SQL Server allows ingress from 0.0.0.0/0 (ANY IP) | accurics.azure.NS.21 |


### azurerm_key_vault
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | azure | MEDIUM | Ensure the key vault is recoverable - enable "Soft Delete" setting for a Key Vault | accurics.azure.EKM.164 |
| Encryption and Key Management | azure | HIGH | Ensure that logging for Azure KeyVault is 'Enabled' | accurics.azure.EKM.20 |


### azurerm_resource_group
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | LOW | Ensure that Azure Resource Group has resource lock enabled | accurics.azure.NS.272 |


### azurerm_storage_account
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure 'Trusted Microsoft Services' is enabled for Storage Account access | accurics.azure.NS.2 |
| Network Security | azure | HIGH | Ensure default network access rule for Storage Accounts is not open to public | accurics.azure.NS.4 |
| Encryption and Key Management | azure | HIGH | Ensure that 'Secure transfer required' is enabled for Storage Accounts | accurics.azure.EKM.7 |


### azurerm_sql_server
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | azure | MEDIUM | Ensure that 'Auditing' is set to 'On' for SQL servers | accurics.azure.MON.354 |
| Identity and Access Management | azure | HIGH | Ensure that Azure Active Directory Admin is configured for SQL Server | accurics.azure.IAM.10 |
| Identity and Access Management | azure | MEDIUM | Avoid using names like 'Admin' for an Azure SQL Server admin account login | accurics.azure.IAM.138 |
| Logging | azure | MEDIUM | Ensure that 'Auditing' Retention is 'greater than 90 days' for SQL servers. | accurics.azure.LOG.356 |


### azurerm_sql_database
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | azure | MEDIUM | Ensure that 'Threat Detection' is enabled for Azure SQL Database | accurics.azure.MON.157 |


### azurerm_redis_cache
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure that Redis is updated regularly with security and operational updates. Note this feature is only available to Premium tier Redis Caches. | accurics.azure.NS.13 |
| Encryption and Key Management | azure | MEDIUM | Ensure that the Redis Cache accepts only SSL connections | accurics.azure.EKM.23 |
| Network Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from other Azure sources | accurics.azure.NS.31 |
| Network Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from the Internet | accurics.azure.NS.30 |
| Network Security | azure | MEDIUM | Ensure there are no firewall rules allowing Redis Cache access for a large number of source IPs | accurics.azure.NS.166 |


### azurerm_kubernetes_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Ensure Kube Dashboard is disabled | accurics.azure.NS.383 |
| Network Security | azure | MEDIUM | Ensure AKS cluster has Network Policy configured. | accurics.azure.NS.382 |


### azurerm_managed_disk
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | azure | MEDIUM | Ensure that 'OS disk' are encrypted | accurics.azure.EKM.156 |


### azurerm_network_watcher_flow_log
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Network Security Group Flow Logs should be enabled and the retention period is set to greater than or equal to 90 days. Flow logs enable capturing information about IP traffic flowing in and out of network security groups. Logs can be used to check for anomalies and give insight into suspected breaches. | accurics.azure.NS.342 |
| Network Security | azure | HIGH | Enable Network Watcher for Azure subscriptions. Network diagnostic and visualization tools available with Network Watcher help users understand, diagnose, and gain insights to the network in Azure. | accurics.azure.NS.11 |


### azurerm_key_vault_secret
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Key Management | azure | HIGH | Ensure that the expiration date is set on all secrets | accurics.azure.EKM.26 |


### azurerm_key_vault_key
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Key Management | azure | HIGH | Ensure that the expiration date is set on all keys | accurics.azure.EKM.25 |


### azurerm_network_security_rule
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Ports Security | azure | HIGH | SSH (TCP:22) is exposed to the wide public internet | accurics.azure.NPS.37 |
| Network Ports Security | azure | HIGH | Remote Desktop (TCP:3389) is exposed to the wide public internet | accurics.azure.NPS.36 |
| Network Ports Security | azure | HIGH | SSH (TCP:22) is exposed to the entire public internet | accurics.azure.NPS.172 |
| Network Ports Security | azure | HIGH | Remote Desktop (TCP:3389) is exposed to the entire public internet | accurics.azure.NPS.171 |
| Network Ports Security | azure | HIGH | CiscoSecure, WebSM (TCP:9090) is exposed to the wide public internet | accurics.azure.NPS.35 |


### azurerm_cosmosdb_account
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Cloud Assets Management | azure | MEDIUM | Ensure that Cosmos DB Account has an associated tag | accurics.azure.CAM.162 |


### azurerm_security_center_subscription_pricing
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Operational Efficiency | azure | MEDIUM | Ensure that standard pricing tiers are selected | accurics.azure.OPS.349 |


### azurerm_sql_active_directory_administrator
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | azure | MEDIUM | Avoid using names like 'Admin' for an Azure SQL Server Active Directory Administrator account | accurics.azure.IAM.137 |


### azurerm_container_registry
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | azure | MEDIUM | Ensure that admin user is disabled for Container Registry | accurics.azure.EKM.164 |
| Azure Container Services | azure | HIGH | Ensure Container Registry has locks | accurics.azure.AKS.3 |


### azurerm_virtual_network
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Ensure that Azure Virtual Network subnet is configured with a Network Security Group | accurics.azure.NS.161 |


### azurerm_role_assignment
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | azure | HIGH | Ensure that there are no guest users | accurics.azure.IAM.388 |


### azurerm_application_gateway
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Ensure Azure Application Gateway Web application firewall (WAF) is enabled | accurics.azure.NS.147 |


### azurerm_postgresql_server
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | azure | HIGH | Ensure 'Enforce SSL connection' is set to 'ENABLED' for PostgreSQL Database Server | accurics.azure.EKM.1 |
| Backup and Disaster Recovery | azure | HIGH | Ensure that Geo Redundant Backups is enabled on PostgreSQL | accurics.azure.BDR.163 |


