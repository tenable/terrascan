
### google_container_node_pool
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Security Best Practices | gcp | LOW | Ensure 'Automatic node upgrade' is enabled for Kubernetes Clusters. | accurics.gcp.OPS.101 |
| Compliance Validation | gcp | HIGH | Ensure Container-Optimized OS (cos) is used for Kubernetes Engine Clusters Node image. | accurics.gcp.OPS.114 |
| Security Best Practices | gcp | LOW | Ensure 'Automatic node repair' is enabled for Kubernetes Clusters. | accurics.gcp.OPS.144 |


### github_repository
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | MEDIUM | Repository is Not Private. | accurics.gcp.IAM.145 |


### google_bigquery_dataset
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | HIGH | BigQuery datasets may be anonymously or publicly accessible. | accurics.gcp.IAM.106 |


### google_compute_project_metadata
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | HIGH | Ensure oslogin is enabled for a Project | accurics.gcp.IAM.127 |


### google_compute_subnetwork
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging and Monitoring | gcp | MEDIUM | Ensure that VPC Flow Logs is enabled for every subnet in a VPC Network. | accurics.gcp.LOG.118 |


### google_project_iam_audit_config
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging and Monitoring | gcp | LOW | Ensure that Cloud Audit Logging is configured properly across all services and all users from a project. | accurics.gcp.LOG.010 |


### google_sql_database_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Resilience | gcp | HIGH | Ensure all Cloud SQL database instance have backup configuration enabled. | accurics.gcp.BDR.105 |
| Infrastructure Security | gcp | HIGH | Ensure that Cloud SQL database Instances are not open to the world. | accurics.gcp.NS.102 |
| Infrastructure Security | gcp | HIGH | Ensure that Cloud SQL database instance requires all incoming connections to use SSL | accurics.gcp.EKM.141 |


### google_compute_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | gcp | MEDIUM | Ensure IP forwarding is not enabled on Instances. | accurics.gcp.NS.130 |
| Infrastructure Security | gcp | HIGH | Ensure 'Block Project-wide SSH keys' is enabled for VM instances. | accurics.gcp.NS.126 |
| Data Protection | gcp | MEDIUM | VM disks attached to a compute instance should be encrypted with Customer Supplied Encryption Keys (CSEK) . | accurics.gcp.EKM.132 |
| Identity and Access Management | gcp | HIGH | Instances may have been configured to use the default service account with full access to all Cloud APIs | accurics.gcp.IAM.124 |
| Infrastructure Security | gcp | MEDIUM | Ensure 'Enable connecting to serial ports' is not enabled for VM instances. | accurics.gcp.NS.129 |
| Infrastructure Security | gcp | MEDIUM | Ensure Compute instances are launched with Shielded VM enabled. | accurics.gcp.NS.133 |
| Identity and Access Management | gcp | MEDIUM | Ensure that no instance in the project overrides the project setting for enabling OSLogin | accurics.gcp.IAM.128 |
| Infrastructure Security | gcp | HIGH | Instances may have been configured to use the default service account with full access to all Cloud APIs | accurics.gcp.NS.125 |


### google_storage_bucket_iam_binding
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | MEDIUM | Ensure that Cloud Storage bucket is not anonymously or publicly accessible. | accurics.gcp.IAM.121 |


### google_container_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | Medium | GKE Control Plane is exposed to few public IP addresses using master-authorized-network-config | AC-GC-IS-CC-M-0367 |
| Logging and Monitoring | gcp | HIGH | Ensure Stackdriver Monitoring is enabled on Kubernetes Engine Clusters. | accurics.gcp.MON.143 |
| Infrastructure Security | gcp | HIGH | Ensure Kubernetes Cluster is created with Private cluster enabled. | accurics.gcp.NS.117 |
| Compliance Validation | gcp | HIGH | Ensure PodSecurityPolicy controller is enabled on the Kubernetes Engine Clusters. | accurics.gcp.OPS.116 |
| Identity and Access Management | gcp | HIGH | Ensure GKE basic auth is disabled. | accurics.gcp.IAM.110 |
| Infrastructure Security | gcp | HIGH | Ensure Master Authentication is set to enabled on Kubernetes Engine Clusters. | accurics.gcp.NS.112 |
| Compliance Validation | gcp | HIGH | Ensure Kubernetes Cluster is created with Alias IP ranges enabled | accurics.gcp.OPS.115 |
| Infrastructure Security | gcp | HIGH | Ensure GKE Control Plane is not public. | accurics.gcp.NS.109 |
| Identity and Access Management | gcp | MEDIUM | Ensure Kubernetes Cluster is created with Client Certificate disabled. | accurics.gcp.IAM.104 |
| Compliance Validation | gcp | HIGH | Ensure Kubernetes Clusters are configured with Labels. | accurics.gcp.OPS.113 |
| Identity and Access Management | gcp | HIGH | Ensure Legacy Authorization is set to disabled on Kubernetes Engine Clusters. | accurics.gcp.IAM.142 |
| Logging and Monitoring | gcp | HIGH | Ensure Stackdriver Logging is enabled on Kubernetes Engine Clusters. | accurics.gcp.LOG.100 |
| Infrastructure Security | gcp | HIGH | Ensure Network policy is enabled on Kubernetes Engine Clusters. | accurics.gcp.NS.103 |


### google_project
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | gcp | MEDIUM | Ensure that the default network does not exist in a project. | accurics.gcp.NS.119 |


### google_compute_firewall
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | json | MEDIUM | Ensure Puppet Master (TCP:8140) is not exposed to public for Google Compute Firewall | AC_GCP_0049 |
| Infrastructure Security | json | HIGH | Ensure Remote Desktop (TCP:3389) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0225 |
| Infrastructure Security | json | HIGH | Ensure LDAP SSL (TCP:636) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0161 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (TCP:4506) is not exposed to public for Google Compute Firewall | AC_GCP_0073 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra (TCP:7001) is not exposed to public for Google Compute Firewall | AC_GCP_0136 |
| Infrastructure Security | json | HIGH | Ensure VNC Listener (TCP:5500) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0065 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (UDP:11215) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0120 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB (TCP:1521) is not exposed to public for Google Compute Firewall | AC_GCP_0209 |
| Infrastructure Security | json | HIGH | Ensure Cassandra Internode Communication (TCP:7000) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0198 |
| Infrastructure Security | json | LOW | Ensure Elastic Search (TCP:9300) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0177 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Datagram Service (TCP:138) is not exposed to public for Google Compute Firewall | AC_GCP_0100 |
| Infrastructure Security | json | LOW | Ensure Mongo Web Portal (TCP:27018) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0045 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Server (TCP:1433) is not exposed to public for Google Compute Firewall | AC_GCP_0157 |
| Infrastructure Security | json | LOW | Ensure Postgres SQL (TCP:5432) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0141 |
| Infrastructure Security | json | HIGH | Ensure Microsoft-DS (TCP:445) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0116 |
| Infrastructure Security | json | HIGH | Ensure SQL Server Analysis Service browser (TCP:2382) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0053 |
| Infrastructure Security | json | HIGH | Ensure Elastic Search (TCP:9200) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0182 |
| Infrastructure Security | json | HIGH | Ensure LDAP (UDP:389) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0213 |
| Infrastructure Security | json | LOW | Ensure NetBios Session Service (UDP:139) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0090 |
| Infrastructure Security | json | LOW | Ensure Oracle DB (TCP:2483) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0205 |
| Infrastructure Security | json | LOW | Ensure Known internal web port (TCP:8000) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0069 |
| Infrastructure Security | json | HIGH | Ensure DNS (UDP:53) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0086 |
| Infrastructure Security | json | HIGH | Ensure Cassandra Monitoring (TCP:7199) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0194 |
| Infrastructure Security | json | HIGH | Ensure Known internal web port (TCP:8080) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0068 |
| Infrastructure Security | json | LOW | Ensure SNMP (UDP:161) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0087 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB (UDP:2483) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0204 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Session Service (UDP:139) is not exposed to public for Google Compute Firewall | AC_GCP_0091 |
| Infrastructure Security | json | MEDIUM | Ensure LDAP (UDP:389) is not exposed to public for Google Compute Firewall | AC_GCP_0212 |
| Infrastructure Security | json | LOW | Ensure Cassandra Thrift (TCP:9160) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0183 |
| Infrastructure Security | json | LOW | Ensure Telnet (TCP:23) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0117 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis Service browser (TCP:2382) is not exposed to public for Google Compute Firewall | AC_GCP_0052 |
| Infrastructure Security | json | HIGH | Ensure Postgres SQL (UDP:5432) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0140 |
| Infrastructure Security | json | LOW | Ensure MSSQL Server (TCP:1433) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0156 |
| Infrastructure Security | json | HIGH | Ensure NetBios Datagram Service (TCP:138) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0101 |
| Infrastructure Security | json | HIGH | Ensure Cassandra OpsCenter agent (TCP:61621) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0044 |
| Infrastructure Security | json | HIGH | Ensure SSH (TCP:20) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0228 |
| Infrastructure Security | json | LOW | Ensure Redis (TCP:6379) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0199 |
| Infrastructure Security | json | HIGH | Ensure Unencrypted Memcached Instances (TCP:11211) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0176 |
| Infrastructure Security | json | LOW | Ensure Oracle DB (TCP:1521) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0208 |
| Infrastructure Security | json | MEDIUM | Ensure VNC Listener (TCP:5500) is not exposed to public for Google Compute Firewall | AC_GCP_0064 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (UDP:11215) is not exposed to public for Google Compute Firewall | AC_GCP_0121 |
| Infrastructure Security | json | LOW | Ensure SaltStack Master (TCP:4506) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0072 |
| Infrastructure Security | json | HIGH | Ensure Cassandra (TCP:7001) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0137 |
| Infrastructure Security | json | MEDIUM | Ensure LDAP SSL (TCP:636) is not exposed to public for Google Compute Firewall | AC_GCP_0160 |
| Infrastructure Security | json | MEDIUM | Ensure Remote Desktop (TCP:3389) is not exposed to public for Google Compute Firewall | AC_GCP_0224 |
| Infrastructure Security | json | LOW | Ensure Puppet Master (TCP:8140) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0048 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (TCP:137) is not exposed to public for Google Compute Firewall | AC_GCP_0106 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra OpsCenter agent (TCP:61621) is not exposed to public for Google Compute Firewall | AC_GCP_0043 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (TCP:2484) is not exposed to public for Google Compute Firewall | AC_GCP_0151 |
| Infrastructure Security | json | LOW | Ensure Oracle DB SSL (UDP:2484) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0147 |
| Infrastructure Security | json | HIGH | Ensure POP3 (TCP:110) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0110 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Browser Service (UDP:1434) is not exposed to public for Google Compute Firewall | AC_GCP_0055 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra Thrift (TCP:9160) is not exposed to public for Google Compute Firewall | AC_GCP_0184 |
| Infrastructure Security | json | MEDIUM | Ensure LDAP (TCP:389) is not exposed to public for Google Compute Firewall | AC_GCP_0215 |
| Infrastructure Security | json | MEDIUM | Ensure CIFS / SMB (TCP:3020) is not exposed to public for Google Compute Firewall | AC_GCP_0079 |
| Infrastructure Security | json | LOW | Ensure NetBios Datagram Service (TCP:138) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0096 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB (UDP:2483) is not exposed to public for Google Compute Firewall | AC_GCP_0203 |
| Infrastructure Security | json | HIGH | Ensure CIFS / SMB (TCP:3020) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0080 |
| Infrastructure Security | json | LOW | Ensure Cassandra Monitoring (TCP:7199) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0192 |
| Infrastructure Security | json | LOW | Ensure Remote Desktop (TCP:3389) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0223 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Admin (TCP:1434) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0059 |
| Infrastructure Security | json | HIGH | Ensure Cassandra Client (TCP:9042) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0188 |
| Infrastructure Security | json | HIGH | Ensure Cassandra OpsCenter Monitoring (TCP:61620) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0167 |
| Infrastructure Security | json | LOW | Ensure SaltStack Master (TCP:4505) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0075 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (TCP:11214) is not exposed to public for Google Compute Firewall | AC_GCP_0130 |
| Infrastructure Security | json | HIGH | Ensure CiscoSecure, Websm (TCP:9090) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0219 |
| Infrastructure Security | json | LOW | Ensure VNC Listener (TCP:5500) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0063 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (UDP:11214) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0126 |
| Infrastructure Security | json | LOW | Ensure Unencrypted Memcached Instances (UDP:11211) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0171 |
| Infrastructure Security | json | HIGH | Ensure Unencrypted Mongo Instances (TCP:27017) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0170 |
| Infrastructure Security | json | HIGH | Ensure VNC Server (TCP:5900) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0062 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (UDP:11214) is not exposed to public for Google Compute Firewall | AC_GCP_0127 |
| Infrastructure Security | json | MEDIUM | Ensure CiscoSecure, Websm (TCP:9090) is not exposed to public for Google Compute Firewall | AC_GCP_0218 |
| Infrastructure Security | json | HIGH | Ensure SaltStack Master (TCP:4506) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0074 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (TCP:11214) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0131 |
| Infrastructure Security | json | LOW | Ensure Cassandra OpsCenter Website (TCP:8888) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0189 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra OpsCenter Monitoring (TCP:61620) is not exposed to public for Google Compute Firewall | AC_GCP_0166 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Admin (TCP:1434) is not exposed to public for Google Compute Firewall | AC_GCP_0058 |
| Infrastructure Security | json | HIGH | Ensure Hadoop Name Node (TCP:9000) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0222 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra Monitoring (TCP:7199) is not exposed to public for Google Compute Firewall | AC_GCP_0193 |
| Infrastructure Security | json | LOW | Ensure Prevalent known internal port (TCP:3000) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0081 |
| Infrastructure Security | json | LOW | Ensure Oracle DB (UDP:2483) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0202 |
| Infrastructure Security | json | LOW | Ensure CIFS / SMB (TCP:3020) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0078 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Datagram Service (TCP:138) is not exposed to public for Google Compute Firewall | AC_GCP_0097 |
| Infrastructure Security | json | LOW | Ensure LDAP (TCP:389) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0214 |
| Infrastructure Security | json | HIGH | Ensure Cassandra Thrift (TCP:9160) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0185 |
| Infrastructure Security | json | LOW | Ensure SMTP (TCP:25) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0111 |
| Infrastructure Security | json | LOW | Ensure MSSQL Browser Service (UDP:1434) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0054 |
| Infrastructure Security | json | HIGH | Ensure MySQL (TCP:3306) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0146 |
| Infrastructure Security | json | LOW | Ensure Oracle DB SSL (TCP:2484) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0150 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Name Service (TCP:137) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0107 |
| Infrastructure Security | json | LOW | Ensure Cassandra OpsCenter agent (TCP:61621) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0042 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra OpsCenter Website (TCP:8888) is not exposed to public for Google Compute Firewall | AC_GCP_0190 |
| Infrastructure Security | json | MEDIUM | Ensure Prevalent known internal port (TCP:3000) is not exposed to public for Google Compute Firewall | AC_GCP_0082 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (UDP:11214) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0128 |
| Infrastructure Security | json | HIGH | Ensure Redis (TCP:6379) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0201 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Session Service (TCP:139) is not exposed to public for Google Compute Firewall | AC_GCP_0094 |
| Infrastructure Security | json | LOW | Ensure CiscoSecure, Websm (TCP:9090) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0217 |
| Infrastructure Security | json | MEDIUM | Ensure Unencrypted Mongo Instances (TCP:27017) is not exposed to public for Google Compute Firewall | AC_GCP_0169 |
| Infrastructure Security | json | LOW | Ensure Cassandra Client (TCP:9042) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0186 |
| Infrastructure Security | json | LOW | Ensure MSSQL Admin (TCP:1434) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0057 |
| Infrastructure Security | json | MEDIUM | Ensure SMTP (TCP:25) is not exposed to public for Google Compute Firewall | AC_GCP_0112 |
| Infrastructure Security | json | MEDIUM | Ensure MySQL (TCP:3306) is not exposed to public for Google Compute Firewall | AC_GCP_0145 |
| Infrastructure Security | json | LOW | Ensure SQL Server Analysis Services (TCP:2383) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0153 |
| Infrastructure Security | json | HIGH | Ensure NetBIOS Name Service (UDP:137) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0104 |
| Infrastructure Security | json | HIGH | Ensure Unencrypted Memcached Instances (UDP:11211) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0173 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (TCP:11215) is not exposed to public for Google Compute Firewall | AC_GCP_0124 |
| Infrastructure Security | json | MEDIUM | Ensure VNC Server (TCP:5900) is not exposed to public for Google Compute Firewall | AC_GCP_0061 |
| Infrastructure Security | json | HIGH | Ensure NetBios Datagram Service (TCP:138) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0098 |
| Infrastructure Security | json | LOW | Ensure Remote Desktop (TCP:3389) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0132 |
| Infrastructure Security | json | HIGH | Ensure SaltStack Master (TCP:4505) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0077 |
| Infrastructure Security | json | LOW | Ensure Cassandra OpsCenter Monitoring (TCP:61620) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0165 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB SSL (UDP:2484) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0149 |
| Infrastructure Security | json | MEDIUM | Ensure Hadoop Name Node (TCP:9000) is not exposed to public for Google Compute Firewall | AC_GCP_0221 |
| Infrastructure Security | json | LOW | Ensure POP3 (TCP:110) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0108 |
| Infrastructure Security | json | MEDIUM | Ensure POP3 (TCP:110) is not exposed to public for Google Compute Firewall | AC_GCP_0109 |
| Infrastructure Security | json | LOW | Ensure Hadoop Name Node (TCP:9000) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0220 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (UDP:2484) is not exposed to public for Google Compute Firewall | AC_GCP_0148 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Debugger (TCP:135) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0164 |
| Infrastructure Security | json | LOW | Ensure NetBios Datagram Service (TCP:138) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0099 |
| Infrastructure Security | json | MEDIUM | Ensure Remote Desktop (TCP:3389) is not exposed to public for Google Compute Firewall | AC_GCP_0133 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (TCP:4505) is not exposed to public for Google Compute Firewall | AC_GCP_0076 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (TCP:11215) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0125 |
| Infrastructure Security | json | LOW | Ensure VNC Server (TCP:5900) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0060 |
| Infrastructure Security | json | MEDIUM | Ensure Unencrypted Memcached Instances (UDP:11211) is not exposed to public for Google Compute Firewall | AC_GCP_0172 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Name Service (TCP:137) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0105 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB SSL (TCP:2484) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0152 |
| Infrastructure Security | json | LOW | Ensure MySQL (TCP:3306) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0144 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Browser Service (UDP:1434) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0056 |
| Infrastructure Security | json | HIGH | Ensure SMTP (TCP:25) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0113 |
| Infrastructure Security | json | LOW | Ensure Unencrypted Mongo Instances (TCP:27017) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0168 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra Client (TCP:9042) is not exposed to public for Google Compute Firewall | AC_GCP_0187 |
| Infrastructure Security | json | HIGH | Ensure LDAP (TCP:389) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0216 |
| Infrastructure Security | json | HIGH | Ensure NetBios Session Service (TCP:139) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0095 |
| Infrastructure Security | json | MEDIUM | Ensure Redis (TCP:6379) is not exposed to public for Google Compute Firewall | AC_GCP_0200 |
| Infrastructure Security | json | HIGH | Ensure Prevalent known internal port (TCP:3000) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0083 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (TCP:11214) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0129 |
| Infrastructure Security | json | HIGH | Ensure Cassandra OpsCenter Website (TCP:8888) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0191 |
| Infrastructure Security | json | MEDIUM | Ensure Unencrypted Memcached Instances (TCP:11211) is not exposed to public for Google Compute Firewall | AC_GCP_0175 |
| Infrastructure Security | json | HIGH | Ensure Memcached SSL (UDP:11215) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0122 |
| Infrastructure Security | json | MEDIUM | Ensure SNMP (UDP:161) is not exposed to public for Google Compute Firewall | AC_GCP_0088 |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (TCP:8080) is not exposed to public for Google Compute Firewall | AC_GCP_0067 |
| Infrastructure Security | json | HIGH | Ensure Remote Desktop (TCP:3389) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0134 |
| Infrastructure Security | json | HIGH | Ensure Known internal web port (TCP:8000) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0071 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Debugger (TCP:135) is not exposed to public for Google Compute Firewall | AC_GCP_0163 |
| Infrastructure Security | json | MEDIUM | Ensure Telnet (TCP:23) is not exposed to public for Google Compute Firewall | AC_GCP_0118 |
| Infrastructure Security | json | LOW | Ensure LDAP SSL (TCP:636) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0159 |
| Infrastructure Security | json | MEDIUM | Ensure SSH (TCP:20) is not exposed to public for Google Compute Firewall | AC_GCP_0227 |
| Infrastructure Security | json | HIGH | Ensure Elastic Search (TCP:9300) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0179 |
| Infrastructure Security | json | LOW | Ensure Cassandra Internode Communication (TCP:7000) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0196 |
| Infrastructure Security | json | LOW | Ensure DNS (UDP:53) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0084 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB (TCP:2483) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0207 |
| Infrastructure Security | json | LOW | Ensure Postgres SQL (UDP:5432) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0138 |
| Infrastructure Security | json | HIGH | Ensure NetBios Session Service (UDP:139) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0092 |
| Infrastructure Security | json | LOW | Ensure LDAP (UDP:389) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0211 |
| Infrastructure Security | json | LOW | Ensure Elastic Search (TCP:9200) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0180 |
| Infrastructure Security | json | LOW | Ensure SQL Server Analysis Service browser (TCP:2382) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0051 |
| Infrastructure Security | json | LOW | Ensure Microsoft-DS (TCP:445) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0114 |
| Infrastructure Security | json | HIGH | Ensure Postgres SQL (TCP:5432) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0143 |
| Infrastructure Security | json | HIGH | Ensure SQL Server Analysis Services (TCP:2383) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0155 |
| Infrastructure Security | json | HIGH | Ensure Mongo Web Portal (TCP:27018) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0047 |
| Infrastructure Security | json | LOW | Ensure NetBIOS Name Service (UDP:137) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0102 |
| Infrastructure Security | json | MEDIUM | Ensure Mongo Web Portal (TCP:27018) is not exposed to public for Google Compute Firewall | AC_GCP_0046 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (UDP:137) is not exposed to public for Google Compute Firewall | AC_GCP_0103 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis Services (TCP:2383) is not exposed to public for Google Compute Firewall | AC_GCP_0154 |
| Infrastructure Security | json | MEDIUM | Ensure Postgres SQL (TCP:5432) is not exposed to public for Google Compute Firewall | AC_GCP_0142 |
| Infrastructure Security | json | HIGH | Ensure Puppet Master (TCP:8140) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0050 |
| Infrastructure Security | json | MEDIUM | Ensure Microsoft-DS (TCP:445) is not exposed to public for Google Compute Firewall | AC_GCP_0115 |
| Infrastructure Security | json | MEDIUM | Ensure Elastic Search (TCP:9200) is not exposed to public for Google Compute Firewall | AC_GCP_0181 |
| Infrastructure Security | json | HIGH | Ensure Oracle DB (TCP:1521) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0210 |
| Infrastructure Security | json | MEDIUM | Ensure Postgres SQL (UDP:5432) is not exposed to public for Google Compute Firewall | AC_GCP_0139 |
| Infrastructure Security | json | LOW | Ensure NetBios Session Service (TCP:139) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0093 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB (TCP:2483) is not exposed to public for Google Compute Firewall | AC_GCP_0206 |
| Infrastructure Security | json | MEDIUM | Ensure DNS (UDP:53) is not exposed to public for Google Compute Firewall | AC_GCP_0085 |
| Infrastructure Security | json | MEDIUM | Ensure Elastic Search (TCP:9300) is not exposed to public for Google Compute Firewall | AC_GCP_0178 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra Internode Communication (TCP:7000) is not exposed to public for Google Compute Firewall | AC_GCP_0197 |
| Infrastructure Security | json | LOW | Ensure SSH (TCP:20) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0226 |
| Infrastructure Security | json | HIGH | Ensure MSSQL Server (TCP:1433) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0158 |
| Infrastructure Security | json | HIGH | Ensure Telnet (TCP:23) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0119 |
| Infrastructure Security | json | LOW | Ensure MSSQL Debugger (TCP:135) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0162 |
| Infrastructure Security | json | LOW | Ensure Cassandra (TCP:7001) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0135 |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (TCP:8000) is not exposed to public for Google Compute Firewall | AC_GCP_0070 |
| Infrastructure Security | json | LOW | Ensure Memcached SSL (TCP:11215) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0123 |
| Infrastructure Security | json | HIGH | Ensure SNMP (UDP:161) is not exposed to entire internet for Google Compute Firewall | AC_GCP_0089 |
| Infrastructure Security | json | LOW | Ensure Known internal web port (TCP:8080) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0066 |
| Infrastructure Security | json | LOW | Ensure Unencrypted Memcached Instances (TCP:11211) is not exposed to private hosts more than 32 for Google Compute Firewall | AC_GCP_0174 |


### google_dns_managed_zone
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | gcp | HIGH | Ensure that RSASHA1 is not used for the zone-signing and key-signing keys in Cloud DNS DNSSEC. | accurics.gcp.EKM.108 |
| Infrastructure Security | gcp | LOW | Ensure that DNSSEC is enabled for Cloud DNS. | accurics.gcp.NS.107 |


### google_compute_disk
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Data Protection | gcp | MEDIUM | Ensure VM disks for critical VMs are encrypted with Customer Supplied Encryption Keys (CSEK) . | accurics.gcp.EKM.131 |


### google_project_iam_member
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | HIGH | Ensure that IAM users are not assigned the Service Account User or Service Account Token Creator roles at project level. | accurics.gcp.IAM.137 |
| Identity and Access Management | gcp | HIGH | Ensure that Service Account has no Admin privileges. | accurics.gcp.IAM.138 |


### google_storage_bucket_iam_member
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | HIGH | Ensure that Cloud Storage bucket is not anonymously or publicly Accessible. | accurics.gcp.IAM.120 |


### google_compute_ssl_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Infrastructure Security | gcp | MEDIUM | Ensure no HTTPS or SSL proxy load balancers permit SSL policies with weak cipher suites. | accurics.gcp.EKM.134 |


### google_storage_bucket
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging and Monitoring | gcp | HIGH | Ensure that logging is enabled for Cloud storage buckets. | accurics.gcp.LOG.147 |
| Logging and Monitoring | gcp | HIGH | Ensure that object versioning is enabled on log-buckets. | accurics.gcp.LOG.146 |
| Identity and Access Management | gcp | MEDIUM | Ensure that Cloud Storage buckets have uniform bucket-level access enabled. | accurics.gcp.IAM.122 |


### google_kms_crypto_key
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Security Best Practices | gcp | MEDIUM | Ensure Encryption keys are rotated within a period of 90 days. | accurics.gcp.EKM.139 |
| Security Best Practices | gcp | HIGH | Ensure Encryption keys are rotated within a period of 365 days. | accurics.gcp.EKM.007 |


### google_project_iam_binding
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | gcp | HIGH | Ensure that IAM users are not assigned the Service Account User or Service Account Token Creator roles at project level. | accurics.gcp.IAM.136 |
| Identity and Access Management | gcp | MEDIUM | Ensure that corporate login credentials are used instead of Gmail accounts. | accurics.gcp.IAM.150 |


