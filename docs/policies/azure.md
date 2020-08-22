
### azurerm_virtual_machine
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure that at least one Network Security Group is attached to all VMs and subnets that are public | accurics.azure.NS.18 |


### azurerm_storage_container
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | azure | HIGH | Anonymous, public read access to a container and its blobs can be enabled in Azure Blob storage. This is only recommended if absolutely necessary. | accurics.azure.IAM.368 |


### azurerm_mysql_server
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure 'Enforce SSL connection' is set to 'ENABLED' for MySQL Database Server. | accurics.azure.NS.361 |


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


### azurerm_storage_account_network_rules
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | MEDIUM | Ensure default network access rule for Storage Accounts is set to deny. | accurics.azure.NS.370 |


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


### azurerm_postgresql_configuration
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | azure | MEDIUM | Ensure server parameter 'log_duration' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.154 |
| Logging | azure | MEDIUM | Ensure server parameter 'log_retention_days' is greater than 3 days for PostgreSQL Database Server | accurics.azure.LOG.155 |
| Logging | azure | MEDIUM | Ensure server parameter 'log_connections' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.152 |
| Logging | azure | MEDIUM | Ensure server parameter 'log_checkpoints' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.364 |
| Logging | azure | MEDIUM | Ensure server parameter 'log_disconnections' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.153 |
| Logging | azure | MEDIUM | Ensure server parameter 'connection_throttling' is set to 'ON' for PostgreSQL Database Server | accurics.azure.LOG.151 |


### azurerm_sql_database
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | azure | MEDIUM | Ensure that 'Threat Detection' is enabled for Azure SQL Database | accurics.azure.MON.157 |


### azurerm_redis_cache
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure that Redis is updated regularly with security and operational updates.Note this feature is only available to Premium tier Redis Caches. | accurics.azure.NS.13 |
| Encryption and Key Management | azure | MEDIUM | Ensure that the Redis Cache accepts only SSL connections | accurics.azure.EKM.23 |
| Network Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from other Azure sources | accurics.azure.NS.31 |
| Network Security | azure | HIGH | Ensure there are no firewall rules allowing unrestricted access to Redis from the Internet | accurics.azure.NS.30 |
| Network Security | azure | MEDIUM | Ensure there are no firewall rules allowing Redis Cache access for a large number of source IPs | accurics.azure.NS.166 |


### azurerm_mssql_server
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | azure | MEDIUM | Ensure that 'Auditing' is set to 'On' for MSSQL servers | accurics.azure.MON.355 |
| Monitoring | azure | MEDIUM | Ensure that 'Auditing' Retention is 'greater than 90 days' for MSSQL servers. | accurics.azure.LOG.357 |


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


### azurerm_security_center_contact
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | azure | MEDIUM | Ensure that 'Send email notification for high severity alerts' is set to 'On' | accurics.azure.MON.353 |


### azurerm_network_security_rule
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Ports Security | azure | LOW | VNC Listener (TCP:5500) is exposed to small Private network | accurics.azure.NPS.314 |
| Network Ports Security | azure | MEDIUM | VNC Listener (TCP:5500) is exposed to small Public network | accurics.azure.NPS.251 |
| Network Ports Security | azure | MEDIUM | Cassandra OpsCenter (TCP:61621) is exposed to wide Private network | accurics.azure.NPS.178 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (TCP:11214) is exposed to small Public network | accurics.azure.NPS.197 |
| Network Ports Security | azure | MEDIUM | Microsoft-DS (TCP:445) is exposed to wide Private network | accurics.azure.NPS.206 |
| Network Ports Security | azure | HIGH | NetBIOS Name Service (TCP:137) is exposed to wide Public network | accurics.azure.NPS.76 |
| Network Ports Security | azure | HIGH | Prevalent known internal port (TCP:3000) is exposed to entire Public network | accurics.azure.NPS.99 |
| Network Ports Security | azure | MEDIUM | MySQL (TCP:3306) is exposed to wide Private network | accurics.azure.NPS.210 |
| Network Ports Security | azure | HIGH | MSSQL Server (TCP:1433) is exposed to wide Public network | accurics.azure.NPS.60 |
| Network Ports Security | azure | LOW | POP3 (TCP:110) is exposed to small Private network | accurics.azure.NPS.302 |
| Network Ports Security | azure | MEDIUM | SaltStack Master (TCP:4506) is exposed to small Public network | accurics.azure.NPS.247 |
| Network Ports Security | azure | HIGH | SSH (TCP:22) is exposed to the wide public internet | accurics.azure.NPS.37 |
| Network Ports Security | azure | MEDIUM | Hadoop Name Node (TCP:9000) is exposed to small Public network | accurics.azure.NPS.181 |
| Network Ports Security | azure | HIGH | Telnet (TCP:23) is exposed to entire Public network | accurics.azure.NPS.115 |
| Network Ports Security | azure | LOW | MSSQL Browser (UDP:1434) is exposed to small Private network | accurics.azure.NPS.284 |
| Network Ports Security | azure | LOW | Mongo Web Portal (TCP:27018) is exposed to small Private network | accurics.azure.NPS.292 |
| Network Ports Security | azure | HIGH | SMTP (TCP:25) is exposed to entire Public network | accurics.azure.NPS.103 |
| Network Ports Security | azure | MEDIUM | CIFS / SMB (TCP:3020) is exposed to wide Private network | accurics.azure.NPS.174 |
| Network Ports Security | azure | HIGH | PostgreSQL (TCP:5432) is exposed to entire Public network | accurics.azure.NPS.95 |
| Network Ports Security | azure | HIGH | NetBIOS Datagram Service (UDP:138) is exposed to entire Public network | accurics.azure.NPS.83 |
| Network Ports Security | azure | HIGH | VNC Server (TCP:5900) is exposed to entire Public network | accurics.azure.NPS.119 |
| Network Ports Security | azure | HIGH | Cassandra (TCP:7001) is exposed to wide Public network | accurics.azure.NPS.40 |
| Network Ports Security | azure | MEDIUM | PostgreSQL (TCP:5432) is exposed to wide Private network | accurics.azure.NPS.230 |
| Network Ports Security | azure | LOW | Memcached SSL (TCP:11215) is exposed to small Private network | accurics.azure.NPS.288 |
| Network Ports Security | azure | HIGH | MSSQL Browser (UDP:1434) is exposed to wide Public network | accurics.azure.NPS.56 |
| Network Ports Security | azure | MEDIUM | Oracle DB SSL (UDP:2484) is exposed to wide Private network | accurics.azure.NPS.226 |
| Network Ports Security | azure | MEDIUM | POP3 (TCP:110) is exposed to small Public network | accurics.azure.NPS.227 |
| Network Ports Security | azure | HIGH | MSSQL Browser (UDP:1434) is exposed to entire Public network | accurics.azure.NPS.57 |
| Network Ports Security | azure | LOW | Memcached SSL (UDP:11214) is exposed to small Private network | accurics.azure.NPS.289 |
| Network Ports Security | azure | MEDIUM | PostgreSQL (UDP:5432) is exposed to small Public network | accurics.azure.NPS.231 |
| Network Ports Security | azure | HIGH | VNC Server (TCP:5900) is exposed to wide Public network | accurics.azure.NPS.118 |
| Network Ports Security | azure | HIGH | Cassandra (TCP:7001) is exposed to entire Public network | accurics.azure.NPS.41 |
| Network Ports Security | azure | HIGH | NetBIOS Datagram Service (UDP:138) is exposed to wide Public network | accurics.azure.NPS.82 |
| Network Ports Security | azure | HIGH | PostgreSQL (TCP:5432) is exposed to wide Public network | accurics.azure.NPS.94 |
| Network Ports Security | azure | MEDIUM | Cassandra (TCP:7001) is exposed to small Public network | accurics.azure.NPS.175 |
| Network Ports Security | azure | HIGH | SMTP (TCP:25) is exposed to wide Public network | accurics.azure.NPS.102 |
| Network Ports Security | azure | LOW | MySQL (TCP:3306) is exposed to small Private network | accurics.azure.NPS.293 |
| Network Ports Security | azure | LOW | MSSQL Debugger (TCP:135) is exposed to small Private network | accurics.azure.NPS.285 |
| Network Ports Security | azure | HIGH | Telnet (TCP:23) is exposed to wide Public network | accurics.azure.NPS.114 |
| Network Ports Security | azure | HIGH | Remote Desktop (TCP:3389) is exposed to the wide public internet | accurics.azure.NPS.36 |
| Network Ports Security | azure | MEDIUM | DNS (UDP:53) is exposed to wide Private network | accurics.azure.NPS.180 |
| Network Ports Security | azure | LOW | PostgreSQL (TCP:5432) is exposed to small Private network | accurics.azure.NPS.303 |
| Network Ports Security | azure | MEDIUM | SaltStack Master (TCP:4505) is exposed to wide Private network | accurics.azure.NPS.246 |
| Network Ports Security | azure | HIGH | MSSQL Server (TCP:1433) is exposed to entire Public network | accurics.azure.NPS.61 |
| Network Ports Security | azure | MEDIUM | NetBIOS Name Service (TCP:137) is exposed to small Public network | accurics.azure.NPS.211 |
| Network Ports Security | azure | HIGH | NetBIOS Name Service (TCP:137) is exposed to entire Public network | accurics.azure.NPS.77 |
| Network Ports Security | azure | HIGH | Prevalent known internal port (TCP:3000) is exposed to wide Public network | accurics.azure.NPS.98 |
| Network Ports Security | azure | MEDIUM | Mongo Web Portal (TCP:27018) is exposed to small Public network | accurics.azure.NPS.207 |
| Network Ports Security | azure | MEDIUM | DNS (UDP:53) is exposed to small Public network | accurics.azure.NPS.179 |
| Network Ports Security | azure | MEDIUM | MSSQL Server (TCP:1433) is exposed to wide Private network | accurics.azure.NPS.196 |
| Network Ports Security | azure | LOW | VNC Server (TCP:5900) is exposed to small Private network | accurics.azure.NPS.315 |
| Network Ports Security | azure | MEDIUM | Telnet (TCP:23) is exposed to wide Private network | accurics.azure.NPS.250 |
| Network Ports Security | azure | HIGH | SSH (TCP:22) is exposed to the entire public internet | accurics.azure.NPS.172 |
| Network Ports Security | azure | HIGH | POP3 (TCP:110) is exposed to entire Public network | accurics.azure.NPS.93 |
| Network Ports Security | azure | HIGH | NetBIOS Session Service (TCP:139) is exposed to entire Public network | accurics.azure.NPS.85 |
| Network Ports Security | azure | LOW | SNMP (UDP:161) is exposed to small Private network | accurics.azure.NPS.308 |
| Network Ports Security | azure | HIGH | Hadoop Name Node (TCP:9000) is exposed to wide Public network | accurics.azure.NPS.46 |
| Network Ports Security | azure | MEDIUM | Puppet Master (TCP:8140) is exposed to wide Private network | accurics.azure.NPS.236 |
| Network Ports Security | azure | LOW | Cassandra OpsCenter (TCP:61621) is exposed to small Private network | accurics.azure.NPS.277 |
| Network Ports Security | azure | LOW | NetBIOS Session Service (TCP:139) is exposed to small Private network | accurics.azure.NPS.298 |
| Network Ports Security | azure | HIGH | SQL Server Analysis (TCP:2383) is exposed to entire Public network | accurics.azure.NPS.109 |
| Network Ports Security | azure | HIGH |  Known internal web port (TCP:8080) is exposed to wide Public network | accurics.azure.NPS.50 |
| Network Ports Security | azure | MEDIUM | NetBIOS Session Service (TCP:139) is exposed to wide Private network | accurics.azure.NPS.220 |
| Network Ports Security | azure | LOW | SaltStack Master (TCP:4506) is exposed to small Private network | accurics.azure.NPS.312 |
| Network Ports Security | azure | MEDIUM | MSSQL Browser (UDP:1434) is exposed to small Public network | accurics.azure.NPS.191 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (TCP:11215) is exposed to wide Private network | accurics.azure.NPS.200 |
| Network Ports Security | azure | HIGH | Microsoft-DS (TCP:445) is exposed to wide Public network | accurics.azure.NPS.70 |
| Network Ports Security | azure | MEDIUM | NetBIOS Datagram Service (TCP:138) is exposed to wide Private network | accurics.azure.NPS.216 |
| Network Ports Security | azure | HIGH | Memcached SSL (UDP:11214) is exposed to wide Public network | accurics.azure.NPS.66 |
| Network Ports Security | azure | HIGH | Oracle DB SSL (TCP:2484) is exposed to entire Public network | accurics.azure.NPS.89 |
| Network Ports Security | azure | LOW | PostgreSQL (UDP:5432) is exposed to small Private network | accurics.azure.NPS.304 |
| Network Ports Security | azure | MEDIUM | SQL Server Analysis (TCP:2382) is exposed to small Public network | accurics.azure.NPS.241 |
| Network Ports Security | azure | MEDIUM | LDAP SSL (TCP:636) is exposed to small Public network | accurics.azure.NPS.187 |
| Network Ports Security | azure | HIGH | SaltStack Master (TCP:4506) is exposed to entire Public network | accurics.azure.NPS.113 |
| Network Ports Security | azure | LOW | LDAP SSL (TCP:636) is exposed to small Private network | accurics.azure.NPS.282 |
| Network Ports Security | azure | LOW | NetBIOS Name Service (TCP:137) is exposed to small Private network | accurics.azure.NPS.294 |
| Network Ports Security | azure | HIGH | SNMP (UDP:161) is exposed to entire Public network | accurics.azure.NPS.105 |
| Network Ports Security | azure | HIGH | SNMP (UDP:161) is exposed to wide Public network | accurics.azure.NPS.104 |
| Network Ports Security | azure | LOW | NetBIOS Name Service (UDP:137) is exposed to small Private network | accurics.azure.NPS.295 |
| Network Ports Security | azure | LOW | MSSQL Admin (TCP:1434) is exposed to small Private network | accurics.azure.NPS.283 |
| Network Ports Security | azure | HIGH | SaltStack Master (TCP:4506) is exposed to wide Public network | accurics.azure.NPS.112 |
| Network Ports Security | azure | MEDIUM |  Known internal web port (TCP:8080) is exposed to wide Private network | accurics.azure.NPS.186 |
| Network Ports Security | azure | LOW | Prevalent known internal port (TCP:3000) is exposed to small Private network | accurics.azure.NPS.305 |
| Network Ports Security | azure | MEDIUM | SNMP (UDP:161) is exposed to wide Private network | accurics.azure.NPS.240 |
| Network Ports Security | azure | HIGH | Memcached SSL (UDP:11214) is exposed to entire Public network | accurics.azure.NPS.67 |
| Network Ports Security | azure | HIGH | Oracle DB SSL (TCP:2484) is exposed to wide Public network | accurics.azure.NPS.88 |
| Network Ports Security | azure | MEDIUM | NetBIOS Datagram Service (UDP:138) is exposed to small Public network | accurics.azure.NPS.217 |
| Network Ports Security | azure | HIGH | Microsoft-DS (TCP:445) is exposed to entire Public network | accurics.azure.NPS.71 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (UDP:11214) is exposed to small Public network | accurics.azure.NPS.201 |
| Network Ports Security | azure | MEDIUM | MSSQL Admin (TCP:1434) is exposed to wide Private network | accurics.azure.NPS.190 |
| Network Ports Security | azure | LOW | Telnet (TCP:23) is exposed to small Private network | accurics.azure.NPS.313 |
| Network Ports Security | azure | MEDIUM | NetBIOS Session Service (UDP:139) is exposed to small Public network | accurics.azure.NPS.221 |
| Network Ports Security | azure | HIGH | SQL Server Analysis (TCP:2383) is exposed to wide Public network | accurics.azure.NPS.108 |
| Network Ports Security | azure | HIGH |  Known internal web port (TCP:8080) is exposed to entire Public network | accurics.azure.NPS.51 |
| Network Ports Security | azure | LOW | Cassandra (TCP:7001) is exposed to small Private network | accurics.azure.NPS.276 |
| Network Ports Security | azure | LOW | NetBIOS Session Service (UDP:139) is exposed to small Private network | accurics.azure.NPS.299 |
| Network Ports Security | azure | MEDIUM | SMTP (TCP:25) is exposed to small Public network | accurics.azure.NPS.237 |
| Network Ports Security | azure | HIGH | Hadoop Name Node (TCP:9000) is exposed to entire Public network | accurics.azure.NPS.47 |
| Network Ports Security | azure | LOW | SQL Server Analysis (TCP:2382) is exposed to small Private network | accurics.azure.NPS.309 |
| Network Ports Security | azure | HIGH | NetBIOS Session Service (TCP:139) is exposed to wide Public network | accurics.azure.NPS.84 |
| Network Ports Security | azure | HIGH | POP3 (TCP:110) is exposed to wide Public network | accurics.azure.NPS.92 |
| Network Ports Security | azure | MEDIUM | CIFS / SMB (TCP:3020) is exposed to small Public network | accurics.azure.NPS.173 |
| Network Ports Security | azure | MEDIUM | NetBIOS Session Service (UDP:139) is exposed to wide Private network | accurics.azure.NPS.222 |
| Network Ports Security | azure | HIGH | LDAP SSL (TCP:636) is exposed to wide Public network | accurics.azure.NPS.52 |
| Network Ports Security | azure | LOW | CIFS / SMB (TCP:3020) is exposed to small Private network | accurics.azure.NPS.275 |
| Network Ports Security | azure | MEDIUM | Prevalent known internal port (TCP:3000) is exposed to wide Private network | accurics.azure.NPS.234 |
| Network Ports Security | azure | HIGH | DNS (UDP:53) is exposed to wide Public network | accurics.azure.NPS.44 |
| Network Ports Security | azure | MEDIUM | MSSQL Admin (TCP:1434) is exposed to small Public network | accurics.azure.NPS.189 |
| Network Ports Security | azure | MEDIUM | NetBIOS Datagram Service (UDP:138) is exposed to wide Private network | accurics.azure.NPS.218 |
| Network Ports Security | azure | HIGH | NetBIOS Session Service (UDP:139) is exposed to entire Public network | accurics.azure.NPS.87 |
| Network Ports Security | azure | HIGH | Memcached SSL (UDP:11215) is exposed to wide Public network | accurics.azure.NPS.68 |
| Network Ports Security | azure | HIGH | Oracle DB SSL (UDP:2484) is exposed to entire Public network | accurics.azure.NPS.91 |
| Network Ports Security | azure | HIGH | CiscoSecure, WebSM (TCP:9090) is exposed to the entire public internet | accurics.azure.NPS.170 |
| Network Ports Security | azure | HIGH | SQL Server Analysis (TCP:2382) is exposed to entire Public network | accurics.azure.NPS.107 |
| Network Ports Security | azure | LOW | NetBIOS Datagram Service (TCP:138) is exposed to small Private network | accurics.azure.NPS.296 |
| Network Ports Security | azure | LOW | Hadoop Name Node (TCP:9000) is exposed to small Private network | accurics.azure.NPS.279 |
| Network Ports Security | azure | LOW |  Known internal web port (TCP:8000) is exposed to small Private network | accurics.azure.NPS.280 |
| Network Ports Security | azure | HIGH |  Known internal web port (TCP:8000) is exposed to wide Public network | accurics.azure.NPS.48 |
| Network Ports Security | azure | HIGH | SaltStack Master (TCP:4505) is exposed to entire Public network | accurics.azure.NPS.111 |
| Network Ports Security | azure | MEDIUM | SMTP (TCP:25) is exposed to wide Private network | accurics.azure.NPS.238 |
| Network Ports Security | azure | MEDIUM |  Known internal web port (TCP:8080) is exposed to small Public network | accurics.azure.NPS.185 |
| Network Ports Security | azure | MEDIUM | SQL Server Analysis (TCP:2383) is exposed to small Public network | accurics.azure.NPS.243 |
| Network Ports Security | azure | LOW | Puppet Master (TCP:8140) is exposed to small Private network | accurics.azure.NPS.306 |
| Network Ports Security | azure | HIGH | Memcached SSL (TCP:11215) is exposed to wide Public network | accurics.azure.NPS.64 |
| Network Ports Security | azure | MEDIUM | NetBIOS Name Service (UDP:137) is exposed to wide Private network | accurics.azure.NPS.214 |
| Network Ports Security | azure | HIGH | Mongo Web Portal (TCP:27018) is exposed to wide Public network | accurics.azure.NPS.72 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (UDP:11214) is exposed to wide Private network | accurics.azure.NPS.202 |
| Network Ports Security | azure | MEDIUM | MSSQL Debugger (TCP:135) is exposed to small Public network | accurics.azure.NPS.193 |
| Network Ports Security | azure | LOW | SQL Server Analysis (TCP:2383) is exposed to small Private network | accurics.azure.NPS.310 |
| Network Ports Security | azure | MEDIUM | VNC Server (TCP:5900) is exposed to wide Private network | accurics.azure.NPS.254 |
| Network Ports Security | azure | LOW | SaltStack Master (TCP:4505) is exposed to small Private network | accurics.azure.NPS.311 |
| Network Ports Security | azure | MEDIUM | MSSQL Browser (UDP:1434) is exposed to wide Private network | accurics.azure.NPS.192 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (UDP:11215) is exposed to small Public network | accurics.azure.NPS.203 |
| Network Ports Security | azure | HIGH | Mongo Web Portal (TCP:27018) is exposed to entire Public network | accurics.azure.NPS.73 |
| Network Ports Security | azure | MEDIUM | NetBIOS Datagram Service (TCP:138) is exposed to small Public network | accurics.azure.NPS.215 |
| Network Ports Security | azure | HIGH | Memcached SSL (TCP:11215) is exposed to entire Public network | accurics.azure.NPS.65 |
| Network Ports Security | azure | MEDIUM | SQL Server Analysis (TCP:2382) is exposed to wide Private network | accurics.azure.NPS.242 |
| Network Ports Security | azure | LOW | SMTP (TCP:25) is exposed to small Private network | accurics.azure.NPS.307 |
| Network Ports Security | azure | MEDIUM |  Known internal web port (TCP:8000) is exposed to wide Private network | accurics.azure.NPS.184 |
| Network Ports Security | azure | MEDIUM | SNMP (UDP:161) is exposed to small Public network | accurics.azure.NPS.239 |
| Network Ports Security | azure | HIGH |  Known internal web port (TCP:8000) is exposed to entire Public network | accurics.azure.NPS.49 |
| Network Ports Security | azure | HIGH | SaltStack Master (TCP:4505) is exposed to wide Public network | accurics.azure.NPS.110 |
| Network Ports Security | azure | LOW |  Known internal web port (TCP:8080) is exposed to small Private network | accurics.azure.NPS.281 |
| Network Ports Security | azure | LOW | NetBIOS Datagram Service (UDP:138) is exposed to small Private network | accurics.azure.NPS.297 |
| Network Ports Security | azure | LOW | DNS (UDP:53) is exposed to small Private network | accurics.azure.NPS.278 |
| Network Ports Security | azure | HIGH | SQL Server Analysis (TCP:2382) is exposed to wide Public network | accurics.azure.NPS.106 |
| Network Ports Security | azure | HIGH | Remote Desktop (TCP:3389) is exposed to the entire public internet | accurics.azure.NPS.171 |
| Network Ports Security | azure | HIGH | Oracle DB SSL (UDP:2484) is exposed to wide Public network | accurics.azure.NPS.90 |
| Network Ports Security | azure | HIGH | NetBIOS Session Service (UDP:139) is exposed to wide Public network | accurics.azure.NPS.86 |
| Network Ports Security | azure | HIGH | Memcached SSL (UDP:11215) is exposed to entire Public network | accurics.azure.NPS.69 |
| Network Ports Security | azure | MEDIUM | NetBIOS Session Service (TCP:139) is exposed to small Public network | accurics.azure.NPS.219 |
| Network Ports Security | azure | MEDIUM | LDAP SSL (TCP:636) is exposed to wide Private network | accurics.azure.NPS.188 |
| Network Ports Security | azure | HIGH | DNS (UDP:53) is exposed to entire Public network | accurics.azure.NPS.45 |
| Network Ports Security | azure | MEDIUM | Puppet Master (TCP:8140) is exposed to small Public network | accurics.azure.NPS.235 |
| Network Ports Security | azure | HIGH | LDAP SSL (TCP:636) is exposed to entire Public network | accurics.azure.NPS.53 |
| Network Ports Security | azure | MEDIUM | Oracle DB SSL (TCP:2484) is exposed to small Public network | accurics.azure.NPS.223 |
| Network Ports Security | azure | HIGH | MSSQL Debugger (TCP:135) is exposed to wide Public network | accurics.azure.NPS.58 |
| Network Ports Security | azure | HIGH | Puppet Master (TCP:8140) is exposed to entire Public network | accurics.azure.NPS.101 |
| Network Ports Security | azure | MEDIUM | POP3 (TCP:110) is exposed to wide Private network | accurics.azure.NPS.228 |
| Network Ports Security | azure | LOW | Memcached SSL (UDP:11215) is exposed to small Private network | accurics.azure.NPS.290 |
| Network Ports Security | azure | LOW | MSSQL Server (TCP:1433) is exposed to small Private network | accurics.azure.NPS.286 |
| Network Ports Security | azure | HIGH | VNC Listener (TCP:5500) is exposed to entire Public network | accurics.azure.NPS.117 |
| Network Ports Security | azure | MEDIUM |  Known internal web port (TCP:8000) is exposed to small Public network | accurics.azure.NPS.183 |
| Network Ports Security | azure | HIGH | CiscoSecure, WebSM (TCP:9090) is exposed to the wide public internet | accurics.azure.NPS.35 |
| Network Ports Security | azure | MEDIUM | SaltStack Master (TCP:4505) is exposed to small Public network | accurics.azure.NPS.245 |
| Network Ports Security | azure | LOW | Oracle DB SSL (TCP:2484) is exposed to small Private network | accurics.azure.NPS.300 |
| Network Ports Security | azure | HIGH | Memcached SSL (TCP:11214) is exposed to wide Public network | accurics.azure.NPS.62 |
| Network Ports Security | azure | MEDIUM | NetBIOS Name Service (TCP:137) is exposed to wide Private network | accurics.azure.NPS.212 |
| Network Ports Security | azure | HIGH | MySQL (TCP:3306) is exposed to wide Public network | accurics.azure.NPS.74 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (UDP:11215) is exposed to wide Private network | accurics.azure.NPS.204 |
| Network Ports Security | azure | MEDIUM | MSSQL Server (TCP:1433) is exposed to small Public network | accurics.azure.NPS.195 |
| Network Ports Security | azure | MEDIUM | VNC Server (TCP:5900) is exposed to small Public network | accurics.azure.NPS.253 |
| Network Ports Security | azure | MEDIUM | Oracle DB SSL (TCP:2484) is exposed to wide Private network | accurics.azure.NPS.224 |
| Network Ports Security | azure | HIGH | MSSQL Admin (TCP:1434) is exposed to wide Public network | accurics.azure.NPS.54 |
| Network Ports Security | azure | MEDIUM | PostgreSQL (UDP:5432) is exposed to wide Private network | accurics.azure.NPS.232 |
| Network Ports Security | azure | HIGH | Cassandra OpsCenter (TCP:61621) is exposed to wide Public network | accurics.azure.NPS.42 |
| Network Ports Security | azure | MEDIUM | Telnet (TCP:23) is exposed to small Public network | accurics.azure.NPS.249 |
| Network Ports Security | azure | HIGH | CIFS / SMB (TCP:3020) is exposed to entire Public network | accurics.azure.NPS.39 |
| Network Ports Security | azure | HIGH | NetBIOS Datagram Service (TCP:138) is exposed to entire Public network | accurics.azure.NPS.81 |
| Network Ports Security | azure | MEDIUM | Mongo Web Portal (TCP:27018) is exposed to wide Private network | accurics.azure.NPS.208 |
| Network Ports Security | azure | HIGH | PostgreSQL (UDP:5432) is exposed to entire Public network | accurics.azure.NPS.97 |
| Network Ports Security | azure | HIGH | NetBIOS Name Service (UDP:137) is exposed to wide Public network | accurics.azure.NPS.78 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (TCP:11215) is exposed to small Public network | accurics.azure.NPS.199 |
| Network Ports Security | azure | MEDIUM | Cassandra (TCP:7001) is exposed to wide Private network | accurics.azure.NPS.176 |
| Network Ports Security | azure | MEDIUM | Memcached SSL (TCP:11214) is exposed to wide Private network | accurics.azure.NPS.198 |
| Network Ports Security | azure | MEDIUM | Cassandra OpsCenter (TCP:61621) is exposed to small Public network | accurics.azure.NPS.177 |
| Network Ports Security | azure | HIGH | PostgreSQL (UDP:5432) is exposed to wide Public network | accurics.azure.NPS.96 |
| Network Ports Security | azure | HIGH | NetBIOS Name Service (UDP:137) is exposed to entire Public network | accurics.azure.NPS.79 |
| Network Ports Security | azure | MEDIUM | MySQL (TCP:3306) is exposed to small Public network | accurics.azure.NPS.209 |
| Network Ports Security | azure | HIGH | NetBIOS Datagram Service (TCP:138) is exposed to wide Public network | accurics.azure.NPS.80 |
| Network Ports Security | azure | HIGH | CIFS / SMB (TCP:3020) is exposed to wide Public network | accurics.azure.NPS.38 |
| Network Ports Security | azure | MEDIUM | SaltStack Master (TCP:4506) is exposed to wide Private network | accurics.azure.NPS.248 |
| Network Ports Security | azure | HIGH | Cassandra OpsCenter (TCP:61621) is exposed to entire Public network | accurics.azure.NPS.43 |
| Network Ports Security | azure | MEDIUM | Prevalent known internal port (TCP:3000) is exposed to small Public network | accurics.azure.NPS.233 |
| Network Ports Security | azure | HIGH | MSSQL Admin (TCP:1434) is exposed to entire Public network | accurics.azure.NPS.55 |
| Network Ports Security | azure | MEDIUM | Oracle DB SSL (UDP:2484) is exposed to small Public network | accurics.azure.NPS.225 |
| Network Ports Security | azure | MEDIUM | VNC Listener (TCP:5500) is exposed to wide Private network | accurics.azure.NPS.252 |
| Network Ports Security | azure | MEDIUM | MSSQL Debugger (TCP:135) is exposed to wide Private network | accurics.azure.NPS.194 |
| Network Ports Security | azure | MEDIUM | Microsoft-DS (TCP:445) is exposed to small Public network | accurics.azure.NPS.205 |
| Network Ports Security | azure | HIGH | MySQL (TCP:3306) is exposed to entire Public network | accurics.azure.NPS.75 |
| Network Ports Security | azure | MEDIUM | NetBIOS Name Service (UDP:137) is exposed to small Public network | accurics.azure.NPS.213 |
| Network Ports Security | azure | HIGH | Memcached SSL (TCP:11214) is exposed to entire Public network | accurics.azure.NPS.63 |
| Network Ports Security | azure | MEDIUM | SQL Server Analysis (TCP:2383) is exposed to wide Private network | accurics.azure.NPS.244 |
| Network Ports Security | azure | LOW | Oracle DB SSL (UDP:2484) is exposed to small Private network | accurics.azure.NPS.301 |
| Network Ports Security | azure | MEDIUM | Hadoop Name Node (TCP:9000) is exposed to wide Private network | accurics.azure.NPS.182 |
| Network Ports Security | azure | HIGH | VNC Listener (TCP:5500) is exposed to wide Public network | accurics.azure.NPS.116 |
| Network Ports Security | azure | LOW | Memcached SSL (TCP:11214) is exposed to small Private network | accurics.azure.NPS.287 |
| Network Ports Security | azure | LOW | Microsoft-DS (TCP:445) is exposed to small Private network | accurics.azure.NPS.291 |
| Network Ports Security | azure | MEDIUM | PostgreSQL (TCP:5432) is exposed to small Public network | accurics.azure.NPS.229 |
| Network Ports Security | azure | HIGH | MSSQL Debugger (TCP:135) is exposed to entire Public network | accurics.azure.NPS.59 |
| Network Ports Security | azure | HIGH | Puppet Master (TCP:8140) is exposed to wide Public network | accurics.azure.NPS.100 |


### azurerm_cosmosdb_account
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | azure | HIGH | Ensure to filter source Ips for Cosmos DB Account | accurics.azure.NS.32 |
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


