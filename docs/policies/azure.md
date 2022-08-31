
### azurerm_storage_container
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | azure | HIGH | Anonymous, public read access to a container and its blobs can be enabled in Azure Blob storage. This is only recommended if absolutely necessary. | accurics.azure.IAM.368 | AC_AZURE_0366 |


### azurerm_mysql_server
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | HIGH | Ensure 'Enforce SSL connection' is set to 'ENABLED' for MySQL Database Server. | accurics.azure.NS.361 | AC_AZURE_0131 |


### azurerm_sql_firewall_rule
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | MEDIUM | Restrict Azure SQL Server accessibility to a minimal address range | accurics.azure.NS.169 | AC_AZURE_0280 |
| Infrastructure Security | azure | HIGH | Ensure entire Azure infrastructure doesn't have access to Azure SQL ServerEnsure entire Azure infrastructure doesn't have access to Azure SQL Server | accurics.azure.NS.5 | AC_AZURE_0381 |
| Infrastructure Security | azure | MEDIUM | Ensure that no SQL Server allows ingress from 0.0.0.0/0 (ANY IP) | accurics.azure.NS.21 | AC_AZURE_0380 |


### azurerm_key_vault
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | azure | MEDIUM | Ensure the key vault is recoverable - enable "Soft Delete" setting for a Key Vault | accurics.azure.EKM.164 | AC_AZURE_0170 |
| Logging and Monitoring | azure | HIGH | Ensure that logging for Azure KeyVault is 'Enabled' | accurics.azure.EKM.20 | AC_AZURE_0169 |


### azurerm_resource_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | azure | LOW | Ensure that Azure Resource Group has resource lock enabled | accurics.azure.NS.272 | AC_AZURE_0389 |


### azurerm_storage_account_network_rules
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | MEDIUM | Ensure default network access rule for Storage Accounts is set to deny. | accurics.azure.NS.370 | AC_AZURE_0309 |


### azurerm_storage_account
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | HIGH | Ensure 'Trusted Microsoft Services' is enabled for Storage Account access | accurics.azure.NS.2 | AC_AZURE_0371 |
| Infrastructure Security | azure | HIGH | Ensure default network access rule for Storage Accounts is not open to public | accurics.azure.NS.4 | AC_AZURE_0370 |
| Data Protection | azure | HIGH | Ensure that 'Secure transfer required' is enabled for Storage Accounts | accurics.azure.EKM.7 | AC_AZURE_0373 |


### azurerm_sql_server
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | azure | MEDIUM | Ensure that 'Auditing' is set to 'On' for SQL servers | accurics.azure.MON.354 | AC_AZURE_0376 |
| Identity and Access Management | azure | HIGH | Ensure that Azure Active Directory Admin is configured for SQL Server | accurics.azure.IAM.10 | AC_AZURE_0378 |
| Compliance Validation | azure | MEDIUM | Avoid using names like 'Admin' for an Azure SQL Server admin account login | accurics.azure.IAM.138 | AC_AZURE_0377 |
| Compliance Validation | azure | LOW | Ensure that 'Auditing' Retention is 'greater than 90 days' for SQL servers. | accurics.azure.LOG.356 | AC_AZURE_0375 |


### azurerm_postgresql_configuration
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'log_duration' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.154 | AC_AZURE_0411 |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'log_retention_days' is greater than 3 days for PostgreSQL Database Server | accurics.azure.LOG.155 | AC_AZURE_0410 |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'log_connections' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.152 | AC_AZURE_0413 |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'log_checkpoints' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.364 | AC_AZURE_0409 |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'log_disconnections' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.153 | AC_AZURE_0412 |
| Logging and Monitoring | azure | MEDIUM | Ensure server parameter 'connection_throttling' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.151 | AC_AZURE_0414 |


### azurerm_sql_database
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | azure | MEDIUM | Ensure that 'Threat Detection' is enabled for Azure SQL Database | accurics.azure.MON.157 | AC_AZURE_0383 |


### azurerm_redis_cache
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | azure | HIGH | Ensure that Redis is updated regularly with security and operational updates.Note this feature is only available to Premium tier Redis Caches. | accurics.azure.NS.13 | AC_AZURE_0393 |
| Infrastructure Security | azure | MEDIUM | Ensure that the Redis Cache accepts only SSL connections | accurics.azure.EKM.23 | AC_AZURE_0394 |
| Infrastructure Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from other Azure sources | accurics.azure.NS.31 | AC_AZURE_0391 |
| Infrastructure Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from the Internet | accurics.azure.NS.30 | AC_AZURE_0392 |
| Infrastructure Security | azure | MEDIUM | Ensure there are no firewall rules allowing Redis Cache access for a large number of source IPs | accurics.azure.NS.166 | AC_AZURE_0390 |


### azurerm_mssql_server
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | azure | MEDIUM | Ensure that 'Auditing' is set to 'On' for MSSQL servers | accurics.azure.MON.355 | AC_AZURE_0137 |
| Logging and Monitoring | azure | MEDIUM | Ensure that 'Auditing' Retention is 'greater than 90 days' for MSSQL servers. | accurics.azure.LOG.357 | AC_AZURE_0136 |


### azurerm_kubernetes_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | MEDIUM | Ensure Kube Dashboard is disabled | accurics.azure.NS.383 | AC_AZURE_0161 |
| Infrastructure Security | azure | MEDIUM | Ensure AKS cluster has Network Policy configured. | accurics.azure.NS.382 | AC_AZURE_0158 |


### azurerm_managed_disk
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | azure | MEDIUM | Ensure that 'Unattached disks' are encrypted in Azure Managed Disk | accurics.azure.EKM.156 | AC_AZURE_0143 |


### azurerm_network_watcher_flow_log
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Resilience | azure | MEDIUM | Ensure that Network Security Group Flow Log retention period is 'greater than 90 days' for Azure Network Watcher Flow Log | accurics.azure.NS.342 | AC_AZURE_0419 |
| Logging and Monitoring | azure | HIGH | Enable Network Watcher for Azure subscriptions. Network diagnostic and visualization tools available with Network Watcher help users understand, diagnose, and gain insights to the network in Azure. | accurics.azure.NS.11 | AC_AZURE_0418 |


### azurerm_key_vault_secret
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | azure | HIGH | Ensure that all secrets have an expiration date set which must expire within 2 years | accurics.azure.EKM.26 | AC_AZURE_0163 |


### azurerm_key_vault_key
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | azure | HIGH | Ensure that the expiration date is set on all keys | accurics.azure.EKM.25 | AC_AZURE_0164 |


### azurerm_security_center_contact
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | azure | MEDIUM | Ensure that 'Send email notification for high severity alerts' is set to 'On' | accurics.azure.MON.353 | AC_AZURE_0386 |


### azurerm_network_security_rule
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (Tcp:8000) is not exposed to public for Azure Network Security Rule | AC_AZURE_0528 | AC_AZURE_0528 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Name Service (Udp:137) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0482 | AC_AZURE_0482 |
| Infrastructure Security | json | LOW | Ensure Microsoft-DS (Tcp:445) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0494 | AC_AZURE_0494 |
| Infrastructure Security | json | LOW | Ensure MSSQL Debugger (Tcp:135) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0512 | AC_AZURE_0512 |
| Infrastructure Security | json | HIGH | Ensure PostgreSQL (Udp:5432) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0457 | AC_AZURE_0457 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (Tcp:11215) is not exposed to public for Azure Network Security Rule | AC_AZURE_0504 | AC_AZURE_0504 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis (Tcp:2382) is not exposed to public for Azure Network Security Rule | AC_AZURE_0441 | AC_AZURE_0441 |
| Infrastructure Security | json | LOW | Ensure POP3 (Tcp:110) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0461 | AC_AZURE_0461 |
| Infrastructure Security | json | LOW | Ensure Known internal web port (Tcp:8080) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0524 | AC_AZURE_0524 |
| Infrastructure Security | json | HIGH | Ensure SaltStack Master (Tcp:4505) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0436 | AC_AZURE_0436 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Datagram Service (Udp:138) is not exposed to public for Azure Network Security Rule | AC_AZURE_0477 | AC_AZURE_0477 |
| Infrastructure Security | json | HIGH | Ensure Hadoop Name Node (Tcp:9000) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0532 | AC_AZURE_0532 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (Udp:11215) is not exposed to public for Azure Network Security Rule | AC_AZURE_0498 | AC_AZURE_0498 |
| Infrastructure Security | json | MEDIUM | Ensure CIFS / SMB (Tcp:3020) is not exposed to public for Azure Network Security Rule | AC_AZURE_0271 | AC_AZURE_0271 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (Tcp:11214) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0508 | AC_AZURE_0508 |
| Infrastructure Security | json | LOW | Ensure MSSQL Server (Tcp:1433) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0509 | AC_AZURE_0509 |
| Infrastructure Security | json | HIGH | Ensure CIFS / SMB (Tcp:3020) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0270 | AC_AZURE_0270 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Datagram Service (Udp:138) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0476 | AC_AZURE_0476 |
| Infrastructure Security | json | LOW | Ensure DNS (Udp:53) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0533 | AC_AZURE_0533 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (Udp:11215) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0499 | AC_AZURE_0499 |
| Infrastructure Security | json | LOW | Ensure server is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0421 | AC_AZURE_0421 |
| Infrastructure Security | json | LOW | Ensure SQL Server Analysis (Tcp:2383) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0437 | AC_AZURE_0437 |
| Infrastructure Security | json | HIGH | Ensure PostgreSQL (Tcp:5432) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0460 | AC_AZURE_0460 |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (Tcp:8080) is not exposed to public for Azure Network Security Rule | AC_AZURE_0525 | AC_AZURE_0525 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (Tcp:11215) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0505 | AC_AZURE_0505 |
| Infrastructure Security | json | LOW | Ensure SQL Server Analysis (Tcp:2382) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0440 | AC_AZURE_0440 |
| Infrastructure Security | json | HIGH | Ensure SSH (Tcp:22) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0285 | AC_AZURE_0285 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Debugger (Tcp:135) is not exposed to public for Azure Network Security Rule | AC_AZURE_0513 | AC_AZURE_0513 |
| Infrastructure Security | json | MEDIUM | Ensure PostgreSQL (Udp:5432) is not exposed to public for Azure Network Security Rule | AC_AZURE_0456 | AC_AZURE_0456 |
| Infrastructure Security | json | MEDIUM | Ensure Microsoft-DS (Tcp:445) is not exposed to public for Azure Network Security Rule | AC_AZURE_0495 | AC_AZURE_0495 |
| Infrastructure Security | json | HIGH | Ensure that RDP access is restricted from the internet for Azure Network Security Rule | AC_AZURE_0342 | AC_AZURE_0342 |
| Infrastructure Security | json | HIGH | Ensure Known internal web port (Tcp:8000) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0529 | AC_AZURE_0529 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (Udp:137) is not exposed to public for Azure Network Security Rule | AC_AZURE_0483 | AC_AZURE_0483 |
| Infrastructure Security | json | LOW | Ensure Oracle DB SSL (Tcp:2484) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0467 | AC_AZURE_0467 |
| Infrastructure Security | json | LOW | Ensure MySQL (Tcp:3306) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0488 | AC_AZURE_0488 |
| Infrastructure Security | json | MEDIUM | Ensure LDAP SSL (Tcp:636) is not exposed to public for Azure Network Security Rule | AC_AZURE_0522 | AC_AZURE_0522 |
| Infrastructure Security | json | HIGH | Ensure Telnet (Tcp:23) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0430 | AC_AZURE_0430 |
| Infrastructure Security | json | MEDIUM | Ensure VNC Listener (Tcp:5500) is not exposed to public for Azure Network Security Rule | AC_AZURE_0426 | AC_AZURE_0426 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Session Service (Udp:139) is not exposed to public for Azure Network Security Rule | AC_AZURE_0471 | AC_AZURE_0471 |
| Infrastructure Security | json | MEDIUM | Ensure DNS (Udp:53) is not exposed to public for Azure Network Security Rule | AC_AZURE_0534 | AC_AZURE_0534 |
| Infrastructure Security | json | LOW | Ensure MSSQL Admin (Tcp:1434) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0518 | AC_AZURE_0518 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Name Service (Udp:137) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0484 | AC_AZURE_0484 |
| Infrastructure Security | json | MEDIUM | Ensure Mongo Web Portal (Tcp:27018) is not exposed to public for Azure Network Security Rule | AC_AZURE_0492 | AC_AZURE_0492 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Debugger (Tcp:135) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0514 | AC_AZURE_0514 |
| Infrastructure Security | json | HIGH | Ensure Puppet Master (Tcp:8140) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0451 | AC_AZURE_0451 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (Udp:11214) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0502 | AC_AZURE_0502 |
| Infrastructure Security | json | MEDIUM | Ensure SMTP (Tcp:25) is not exposed to public for Azure Network Security Rule | AC_AZURE_0447 | AC_AZURE_0447 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (Tcp:11215) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0503 | AC_AZURE_0503 |
| Infrastructure Security | json | LOW | Ensure SMTP (Tcp:25) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0446 | AC_AZURE_0446 |
| Infrastructure Security | json | LOW | Ensure MSSQL Browser (Udp:1434) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0515 | AC_AZURE_0515 |
| Infrastructure Security | json | MEDIUM | Ensure Puppet Master (Tcp:8140) is not exposed to public for Azure Network Security Rule | AC_AZURE_0450 | AC_AZURE_0450 |
| Infrastructure Security | json | HIGH | Ensure Mongo Web Portal (Tcp:27018) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0493 | AC_AZURE_0493 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Name Service (Tcp:137) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0485 | AC_AZURE_0485 |
| Infrastructure Security | json | HIGH | Ensure Cassandra OpsCenter (Tcp:61621) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0276 | AC_AZURE_0276 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Admin (Tcp:1434) is not exposed to public for Azure Network Security Rule | AC_AZURE_0519 | AC_AZURE_0519 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Session Service (Udp:139) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0470 | AC_AZURE_0470 |
| Infrastructure Security | json | HIGH | Ensure DNS (Udp:53) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0535 | AC_AZURE_0535 |
| Infrastructure Security | json | HIGH | Ensure VNC Listener (Tcp:5500) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0427 | AC_AZURE_0427 |
| Infrastructure Security | json | LOW | Ensure SaltStack Master (Tcp:4506) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0431 | AC_AZURE_0431 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB SSL (Udp:2484) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0466 | AC_AZURE_0466 |
| Infrastructure Security | json | MEDIUM | Ensure MySQL (Tcp:3306) is not exposed to public for Azure Network Security Rule | AC_AZURE_0489 | AC_AZURE_0489 |
| Infrastructure Security | json | HIGH | Ensure LDAP SSL (Tcp:636) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0523 | AC_AZURE_0523 |
| Infrastructure Security | json | LOW | Ensure Puppet Master (Tcp:8140) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0449 | AC_AZURE_0449 |
| Infrastructure Security | json | LOW | Ensure Cassandra (Tcp:7001) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0275 | AC_AZURE_0275 |
| Infrastructure Security | json | LOW | Ensure Cassandra OpsCenter (Tcp:61621) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0536 | AC_AZURE_0536 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Session Service (Tcp:139) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0473 | AC_AZURE_0473 |
| Infrastructure Security | json | HIGH | Ensure VNC Server (Tcp:5900) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0424 | AC_AZURE_0424 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (Tcp:4506) is not exposed to public for Azure Network Security Rule | AC_AZURE_0432 | AC_AZURE_0432 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Admin (Tcp:1434) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0520 | AC_AZURE_0520 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (Udp:2484) is not exposed to public for Azure Network Security Rule | AC_AZURE_0465 | AC_AZURE_0465 |
| Infrastructure Security | json | HIGH | Ensure SNMP (Udp:161) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0445 | AC_AZURE_0445 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (Udp:11214) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0500 | AC_AZURE_0500 |
| Infrastructure Security | json | MEDIUM | Ensure Prevalent known internal port (Tcp:3000) is not exposed to public for Azure Network Security Rule | AC_AZURE_0453 | AC_AZURE_0453 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Browser (Udp:1434) is not exposed to public for Azure Network Security Rule | AC_AZURE_0516 | AC_AZURE_0516 |
| Infrastructure Security | json | HIGH | Ensure MySQL (Tcp:3306) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0490 | AC_AZURE_0490 |
| Infrastructure Security | json | LOW | Ensure Telnet (Tcp:23) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0428 | AC_AZURE_0428 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (Tcp:137) is not exposed to public for Azure Network Security Rule | AC_AZURE_0486 | AC_AZURE_0486 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB SSL (Tcp:2484) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0469 | AC_AZURE_0469 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Name Service (Tcp:137) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0487 | AC_AZURE_0487 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (Tcp:2484) is not exposed to public for Azure Network Security Rule | AC_AZURE_0468 | AC_AZURE_0468 |
| Infrastructure Security | json | MEDIUM | Ensure Telnet (Tcp:23) is not exposed to public for Azure Network Security Rule | AC_AZURE_0429 | AC_AZURE_0429 |
| Infrastructure Security | json | LOW | Ensure Mongo Web Portal (Tcp:27018) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0491 | AC_AZURE_0491 |
| Infrastructure Security | json | LOW | Ensure Prevalent known internal port (Tcp:3000) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0452 | AC_AZURE_0452 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Browser (Udp:1434) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0517 | AC_AZURE_0517 |
| Infrastructure Security | json | MEDIUM | Ensure SNMP (Udp:161) is not exposed to public for Azure Network Security Rule | AC_AZURE_0444 | AC_AZURE_0444 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (Udp:11214) is not exposed to public for Azure Network Security Rule | AC_AZURE_0501 | AC_AZURE_0501 |
| Infrastructure Security | json | LOW | Ensure LDAP SSL (Tcp:636) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0521 | AC_AZURE_0521 |
| Infrastructure Security | json | LOW | Ensure Oracle DB SSL (Udp:2484) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0464 | AC_AZURE_0464 |
| Infrastructure Security | json | HIGH | Ensure SaltStack Master (Tcp:4506) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0433 | AC_AZURE_0433 |
| Infrastructure Security | json | LOW | Ensure VNC Listener (Tcp:5500) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0425 | AC_AZURE_0425 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra OpsCenter (Tcp:61621) is not exposed to public for Azure Network Security Rule | AC_AZURE_0537 | AC_AZURE_0537 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Session Service (Udp:139) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0472 | AC_AZURE_0472 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra (Tcp:7001) is not exposed to public for Azure Network Security Rule | AC_AZURE_0274 | AC_AZURE_0274 |
| Infrastructure Security | json | HIGH | Ensure SMTP (Tcp:25) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0448 | AC_AZURE_0448 |
| Infrastructure Security | json | LOW | Ensure SNMP (Udp:161) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0443 | AC_AZURE_0443 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (Tcp:11214) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0506 | AC_AZURE_0506 |
| Infrastructure Security | json | MEDIUM | Ensure SSH (Tcp:22) is not exposed to public for Azure Network Security Rule | AC_AZURE_0286 | AC_AZURE_0286 |
| Infrastructure Security | json | LOW | Ensure PostgreSQL (Udp:5432) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0455 | AC_AZURE_0455 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Server (Tcp:1433) is not exposed to public for Azure Network Security Rule | AC_AZURE_0510 | AC_AZURE_0510 |
| Infrastructure Security | json | HIGH | Ensure Microsoft-DS (Tcp:445) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0496 | AC_AZURE_0496 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Datagram Service (Tcp:138) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0479 | AC_AZURE_0479 |
| Infrastructure Security | json | HIGH | Ensure that request initiated from all ports (*) for all destination ports (*) is restricted from the internet for Azure Network Security Rule | AC_AZURE_0357 | AC_AZURE_0357 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis (Tcp:2383) is not exposed to public for Azure Network Security Rule | AC_AZURE_0438 | AC_AZURE_0438 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Datagram Service (Tcp:138) is not exposed to public for Azure Network Security Rule | AC_AZURE_0480 | AC_AZURE_0480 |
| Infrastructure Security | json | HIGH | Ensure Cassandra (Tcp:7001) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0273 | AC_AZURE_0273 |
| Infrastructure Security | json | MEDIUM | Ensure PostgreSQL (Tcp:5432) is not exposed to public for Azure Network Security Rule | AC_AZURE_0459 | AC_AZURE_0459 |
| Infrastructure Security | json | LOW | Ensure Hadoop Name Node (Tcp:9000) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0530 | AC_AZURE_0530 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Session Service (Tcp:139) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0475 | AC_AZURE_0475 |
| Infrastructure Security | json | LOW | Ensure VNC Server (Tcp:5900) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0422 | AC_AZURE_0422 |
| Infrastructure Security | json | LOW | Ensure SaltStack Master (Tcp:4505) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0434 | AC_AZURE_0434 |
| Infrastructure Security | json | HIGH | Ensure Known internal web port (Tcp:8080) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0526 | AC_AZURE_0526 |
| Infrastructure Security | json | HIGH | Ensure POP3 (Tcp:110) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0463 | AC_AZURE_0463 |
| Infrastructure Security | json | LOW | Ensure Known internal web port (Tcp:8000) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0527 | AC_AZURE_0527 |
| Infrastructure Security | json | MEDIUM | Ensure POP3 (Tcp:110) is not exposed to public for Azure Network Security Rule | AC_AZURE_0462 | AC_AZURE_0462 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (Tcp:4505) is not exposed to public for Azure Network Security Rule | AC_AZURE_0435 | AC_AZURE_0435 |
| Infrastructure Security | json | MEDIUM | Ensure VNC Server (Tcp:5900) is not exposed to public for Azure Network Security Rule | AC_AZURE_0423 | AC_AZURE_0423 |
| Infrastructure Security | json | MEDIUM | Ensure Hadoop Name Node (Tcp:9000) is not exposed to public for Azure Network Security Rule | AC_AZURE_0531 | AC_AZURE_0531 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Session Service (Tcp:139) is not exposed to public for Azure Network Security Rule | AC_AZURE_0474 | AC_AZURE_0474 |
| Infrastructure Security | json | LOW | Ensure PostgreSQL (Tcp:5432) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0458 | AC_AZURE_0458 |
| Infrastructure Security | json | LOW | Ensure CIFS / SMB (Tcp:3020) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0272 | AC_AZURE_0272 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Datagram Service (Tcp:138) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0481 | AC_AZURE_0481 |
| Infrastructure Security | json | HIGH | Ensure SQL Server Analysis (Tcp:2383) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0439 | AC_AZURE_0439 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (Udp:11215) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0497 | AC_AZURE_0497 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Datagram Service (Udp:138) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0478 | AC_AZURE_0478 |
| Infrastructure Security | json | HIGH | Ensure Prevalent known internal port (Tcp:3000) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0454 | AC_AZURE_0454 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Server (Tcp:1433) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0511 | AC_AZURE_0511 |
| Infrastructure Security | json | LOW | Ensure SSH (Tcp:22) is not exposed to private hosts more than 32 for Azure Network Security Rule | AC_AZURE_0287 | AC_AZURE_0287 |
| Infrastructure Security | json | HIGH | Ensure SQL Server Analysis (Tcp:2382) is not exposed to entire internet for Azure Network Security Rule | AC_AZURE_0442 | AC_AZURE_0442 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (Tcp:11214) is not exposed to public for Azure Network Security Rule | AC_AZURE_0507 | AC_AZURE_0507 |


### azurerm_cosmosdb_account
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | HIGH | Ensure to filter source Ips for Cosmos DB Account | accurics.azure.NS.32 | AC_AZURE_0184 |
| Compliance Validation | azure | MEDIUM | Ensure that Cosmos DB Account has an associated tag | accurics.azure.CAM.162 | AC_AZURE_0277 |


### azurerm_security_center_subscription_pricing
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | azure | MEDIUM | Ensure that standard pricing tiers are selected | accurics.azure.OPS.349 | AC_AZURE_0385 |


### azurerm_sql_active_directory_administrator
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | azure | MEDIUM | Avoid using names like 'Admin' for an Azure SQL Server Active Directory Administrator account | accurics.azure.IAM.137 | AC_AZURE_0384 |


### azurerm_container_registry
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | azure | MEDIUM | Ensure that admin user is disabled for Container Registry | accurics.azure.EKM.164 | AC_AZURE_0186 |
| Resilience | azure | HIGH | Ensure Container Registry has locks | accurics.azure.AKS.3 | AC_AZURE_0185 |


### azurerm_virtual_network
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | MEDIUM | Ensure that Azure Virtual Network subnet is configured with a Network Security Group | accurics.azure.NS.161 | AC_AZURE_0356 |


### azurerm_role_assignment
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | azure | HIGH | Ensure that there are no guest users | accurics.azure.IAM.388 | AC_AZURE_0388 |


### azurerm_application_gateway
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | MEDIUM | Ensure Azure Application Gateway Web application firewall (WAF) is enabled | accurics.azure.NS.147 | AC_AZURE_0189 |


### azurerm_postgresql_server
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | azure | HIGH | Ensure 'Enforce SSL connection' is set to 'ENABLED' for PostgreSQL Database Server | accurics.azure.EKM.1 | AC_AZURE_0408 |
| Resilience | azure | MEDIUM | Ensure that Geo Redundant Backups is enabled on PostgreSQL | accurics.azure.BDR.163 | AC_AZURE_0407 |


