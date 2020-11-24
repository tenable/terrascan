
### google_container_node_pool
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Operational Efficiency | gcp | HIGH | Ensure 'Automatic node upgrade' is enabled for Kubernetes Clusters. | accurics.gcp.OPS.101 |
| Operational Efficiency | gcp | HIGH | Ensure Container-Optimized OS (cos) is used for Kubernetes Engine Clusters Node image. | accurics.gcp.OPS.114 |
| Operational Efficiency | gcp | MEDIUM | Ensure 'Automatic node repair' is enabled for Kubernetes Clusters. | accurics.gcp.OPS.144 |


### github_repository
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | MEDIUM | Repository is Not Private. | accurics.gcp.IAM.145 |


### google_bigquery_dataset
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | HIGH | BigQuery datasets may be anonymously or publicly accessible. | accurics.gcp.IAM.106 |


### google_compute_project_metadata
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Access Control | gcp | HIGH | Ensure oslogin is enabled for a Project | accurics.gcp.IAM.127 |


### google_compute_subnetwork
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging  | gcp | MEDIUM | Ensure that VPC Flow Logs is enabled for every subnet in a VPC Network. | accurics.gcp.LOG.118 |


### google_project_iam_audit_config
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | gcp | HIGH | Ensure that Cloud Audit Logging is configured properly across all services and all users from a project. | accurics.gcp.LOG.010 |


### google_sql_database_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Backup & Disaster Recovery | gcp | HIGH | Ensure all Cloud SQL database instance have backup configuration enabled. | accurics.gcp.BDR.105 |
| Network Security | gcp | HIGH | Ensure that Cloud SQL database Instances are not open to the world. | accurics.gcp.NS.102 |
| Encryption & Key Management | gcp | HIGH | Ensure that Cloud SQL database instance requires all incoming connections to use SSL | accurics.gcp.EKM.141 |


### google_compute_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | gcp | MEDIUM | Ensure IP forwarding is not enabled on Instances. | accurics.gcp.NS.130 |
| Network Security | gcp | MEDIUM | Ensure 'Block Project-wide SSH keys' is enabled for VM instances. | accurics.gcp.NS.126 |
| Encryption & Key Management | gcp | MEDIUM | Ensure VM disks for critical VMs are encrypted with Customer Supplied Encryption Keys (CSEK) . | accurics.gcp.EKM.132 |
| Identity & Access Management | gcp | MEDIUM | Instances may have been configured to use the default service account with full access to all Cloud APIs | accurics.gcp.IAM.124 |
| Network Security | gcp | MEDIUM | Ensure 'Enable connecting to serial ports' is not enabled for VM instances. | accurics.gcp.NS.129 |
| Network Security  | gcp | MEDIUM | Ensure Compute instances are launched with Shielded VM enabled. | accurics.gcp.NS.133 |
| Identity & Access Management | gcp | MEDIUM | Ensure that no instance in the project overrides the project setting for enabling OSLogin | accurics.gcp.IAM.128 |
| Access Control | gcp | HIGH | Instances may have been configured to use the default service account with full access to all Cloud APIs | accurics.gcp.NS.125 |


### google_storage_bucket_iam_binding
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | MEDIUM | Ensure that Cloud Storage bucket is not anonymously or publicly accessible. | accurics.gcp.IAM.121 |


### google_container_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Monitoring | gcp | HIGH | Ensure Stackdriver Monitoring is enabled on Kubernetes Engine Clusters. | accurics.gcp.MON.143 |
| Network Security | gcp | HIGH | Ensure Kubernetes Cluster is created with Private cluster enabled. | accurics.gcp.NS.117 |
| Operational Efficiency | gcp | HIGH | Ensure PodSecurityPolicy controller is enabled on the Kubernetes Engine Clusters. | accurics.gcp.OPS.116 |
| Identity & Access Management | gcp | HIGH | Ensure GKE basic auth is disabled. | accurics.gcp.IAM.110 |
| Network Security | gcp | HIGH | Ensure Master Authentication is set to enabled on Kubernetes Engine Clusters. | accurics.gcp.NS.112 |
| Operational Efficiency | gcp | HIGH | Ensure Kubernetes Cluster is created with Alias IP ranges enabled | accurics.gcp.OPS.115 |
| Network Security | gcp | HIGH | Ensure GKE Control Plane is not public. | accurics.gcp.NS.109 |
| Identity & Access Management | gcp | HIGH | Ensure Kubernetes Cluster is created with Client Certificate disabled. | accurics.gcp.IAM.104 |
| Operational Efficiency | gcp | HIGH | Ensure Kubernetes Clusters are configured with Labels. | accurics.gcp.OPS.113 |
| Identity & Access Management | gcp | HIGH | Ensure Legacy Authorization is set to disabled on Kubernetes Engine Clusters. | accurics.gcp.IAM.142 |
| Logging | gcp | HIGH | Ensure Stackdriver Logging is enabled on Kubernetes Engine Clusters. | accurics.gcp.LOG.100 |
| Network Security | gcp | HIGH | Ensure Network policy is enabled on Kubernetes Engine Clusters. | accurics.gcp.NS.103 |


### google_project
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | gcp | MEDIUM | Ensure that the default network does not exist in a project. | accurics.gcp.NS.119 |


### google_compute_firewall
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | gcp | MEDIUM | Ensure that SSH access is restricted from the internet | accurics.gcp.NS.149 |
| Network Security | gcp | MEDIUM | Ensure that SSH access is restricted from the internet | accurics.gcp.NS.148 |
| Network Security  | gcp | MEDIUM | Ensure Google compute firewall ingress does not allow unrestricted rdp access. | accurics.gcp.NS.123 |


### google_dns_managed_zone
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & Key Management | gcp | HIGH | Ensure that RSASHA1 is not used for the zone-signing and key-signing keys in Cloud DNS DNSSEC. | accurics.gcp.EKM.108 |
| Network Security | gcp | HIGH | Ensure that DNSSEC is enabled for Cloud DNS. | accurics.gcp.NS.107 |


### google_compute_disk
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & Key Management | gcp | MEDIUM | Ensure VM disks for critical VMs are encrypted with Customer Supplied Encryption Keys (CSEK) . | accurics.gcp.EKM.131 |


### google_project_iam_member
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | MEDIUM | Ensure that IAM users are not assigned the Service Account User or Service Account Token Creator roles at project level. | accurics.gcp.IAM.137 |
| Identity & Access Management | gcp | MEDIUM | Ensure that Service Account has no Admin privileges. | accurics.gcp.IAM.138 |


### google_storage_bucket_iam_member
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | HIGH | Ensure that Cloud Storage bucket is not anonymously or publicly Accessible. | accurics.gcp.IAM.120 |


### google_compute_ssl_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & Key Management | gcp | MEDIUM | Ensure no HTTPS or SSL proxy load balancers permit SSL policies with weak cipher suites. | accurics.gcp.EKM.134 |


### google_storage_bucket
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | gcp | HIGH | Ensure that logging is enabled for Cloud storage buckets. | accurics.gcp.LOG.147 |
| Logging | gcp | HIGH | Ensure that object versioning is enabled on log-buckets. | accurics.gcp.LOG.146 |
| Identity & Access Management | gcp | MEDIUM | Ensure that Cloud Storage buckets have uniform bucket-level access enabled. | accurics.gcp.IAM.122 |


### google_kms_crypto_key
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & Key Management | gcp | MEDIUM | Ensure Encryption keys are rotated within a period of 90 days. | accurics.gcp.EKM.139 |
| Encryption & Key Management | gcp | HIGH | Ensure Encryption keys are rotated within a period of 365 days. | accurics.gcp.EKM.007 |


### google_project_iam_binding
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity & Access Management | gcp | MEDIUM | Ensure that IAM users are not assigned the Service Account User or Service Account Token Creator roles at project level. | accurics.gcp.IAM.136 |
| Identity and Access Management | gcp | HIGH | Ensure that corporate login credentials are used instead of Gmail accounts. | accurics.gcp.IAM.150 |


