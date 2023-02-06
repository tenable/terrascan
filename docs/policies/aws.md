
### aws_iam_role_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | json | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AC-AW-IA-H-1189 | AC_AWS_0147 |
| Identity and Access Management | AIRP | HIGH | Ensure IAM roles do not have any policies attached that may cause privilege escalation. | AWS.AIRP.IAM.HIGH.0051 | AC_AWS_0473 |


### aws_route53_record
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | Route53HostedZone | HIGH | Route53HostedZone should have recordSets. | AWS.Route53HostedZone.DNSManagement.High.0422 | AC_AWS_0205 |


### aws_elasticsearch_domain_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | AEDP | HIGH | Ensure Elasticsearch domains do not have wildcard policies. | AWS.AEDP.IAM.HIGH.0060 | AC_AWS_0469 |


### aws_lb_target_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | ALTG | MEDIUM | Ensure Target Group use HTTPs to ensure end to end encryption | AWS.ALTG.IS.MEDIUM.0042 | AC_AWS_0492 |


### aws_api_gateway_method_settings
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | API Gateway | MEDIUM | Enable Detailed CloudWatch Metrics for APIs | AWS.APIGateway.Logging.Medium.0569 | AC_AWS_0007 |


### aws_workspaces_workspace
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | AWW | MEDIUM | Ensure user volume for Workspaces is Encrypted | AWS.AWW.DP.MEDIUM.041 | AC_AWS_0504 |
| Data Protection | AWW | MEDIUM | Ensure root volume for Workspaces is Encrypted | AWS.AWW.DP.MEDIUM.040 | AC_AWS_0503 |


### aws_vpc
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | VPC | MEDIUM | Avoid creating resources in default VPC | AWS.VPC.Logging.Medium.0471 | AC_AWS_0370 |
| Logging and Monitoring | VPC | LOW | Ensure VPC flow logging is enabled in all VPCs | AWS.VPC.Logging.Medium.0470 | AC_AWS_0369 |


### aws_iam_account_password_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | Iam | LOW | Lower case alphabet not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0454 | AC_AWS_0134 |
| Compliance Validation | Iam | MEDIUM | Setting a lengthy password increases account resiliency against brute force login attempts | AWS.Iam.IAM.Medium.0458 | AC_AWS_0142 |
| Compliance Validation | Iam | LOW | It is recommended that the password policy prevent the reuse of passwords.Preventing password reuse increases account resiliency against brute force login attempts | AWS.Iam.IAM.Low.0539 | AC_AWS_0472 |
| Compliance Validation | Iam | MEDIUM | Number not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0455 | AC_AWS_0136 |
| Compliance Validation | Iam | MEDIUM | Setting a lengthy password increases account resiliency against brute force login attempts | AWS.Iam.IAM.Medium.0495 | AC_AWS_0141 |
| Compliance Validation | Iam | MEDIUM | Special symbols not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0456 | AC_AWS_0137 |
| Compliance Validation | Iam | MEDIUM | Upper case alphabet not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0457 | AC_AWS_0135 |
| Compliance Validation | Iam | LOW | Reducing the password lifetime increases account resiliency against brute force login attempts | AWS.Iam.IAM.Low.0540 | AC_AWS_0138 |


### aws_mq_broker
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | ElasticSearch | MEDIUM | Publicly Accessible MQ Brokers | AWS.ElasticSearch.NetworkSecurity.Medium.0887 | AC_AWS_0175 |
| Logging and Monitoring | ElasticSearch | LOW | Enable AWS MQ Log Exports | AWS.ElasticSearch.Logging.Medium.0885 | AC_AWS_0174 |


### aws_db_instance
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | RDS | HIGH | Ensure Certificate used in RDS instance is updated | AWS.RDS.DS.High.1042 | AC_AWS_0057 |
| Logging and Monitoring | ADI | MEDIUM | Ensure AWS RDS instances have logging enabled. | AWS.ADI.LM.MEDIUM.0076 | AC_AWS_0454 |
| Data Protection | RDS | MEDIUM | Ensure that your RDS database has IAM Authentication enabled. | AWS.RDS.DataSecurity.High.0577 | AC_AWS_0053 |
| Infrastructure Security | RDS | HIGH | RDS Instance publicly_accessible flag is true | AWS.RDS.NS.High.0101 | AC_AWS_0054 |
| Data Protection | RDS | HIGH | RDS Instance Auto Minor Version Upgrade flag disabled | AWS.RDS.DS.High.1041 | AC_AWS_0056 |
| Data Protection | RDS | HIGH | Ensure that your RDS database instances encrypt the underlying storage. Encrypted RDS instances use the industry standard AES-256 encryption algorithm to encrypt data on the server that hosts RDS DB instances. After data is encrypted, RDS handles authentication of access and descryption of data transparently with minimal impact on performance. | AWS.RDS.DataSecurity.High.0414 | AC_AWS_0058 |


### aws_secretsmanager_secret_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | ASSP | HIGH | Ensure secrets manager do not wildcard policies attached | AWS.ASSP.IAM.HIGH.0066 | AC_AWS_0501 |


### aws_ebs_volume
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | EcsCluster | HIGH | Ensure that AWS EBS clusters are encrypted. Data encryption at rest, prevents unauthorized users from accessing sensitive data on your AWS EBS clusters and associated cache storage systems. | AWS.EcsCluster.EncryptionandKeyManagement.High.0413 | AC_AWS_0460 |
| Data Protection | EBS | HIGH | Enable AWS EBS Snapshot Encryption | AWS.EBS.EKM.Medium.0682 | AC_AWS_0459 |


### aws_api_gateway_rest_api
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | APIGateway | MEDIUM | Enable Content Encoding | AWS.APIGateway.Medium.0568 | AC_AWS_0010 |
| Infrastructure Security | APIGateway | MEDIUM | API Gateway Private Endpoints | AWS.APIGateway.NetworkSecurity.Medium.0570 | AC_AWS_0011 |


### aws_iam_role
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | json | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AC-AW-IA-H-1188 | AC_AWS_0146 |


### aws_iam_user_policy_attachment
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | AIUPA | MEDIUM | Ensure IAM permissions are not given directly to users | AWS.AIUPA.IAM.MEDIUM.0050 | AC_AWS_0476 |


### aws_ebs_encryption_by_default
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | EBS | HIGH | Ensure that the AWS EBS that hold sensitive and critical data is encrypted by default to fulfill compliance requirements for data-at-rest encryption. | AWS.EBS.DataSecurity.High.0580 | AC_AWS_0079 |


### aws_sns_topic
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | SNS | HIGH | Ensure SNS Topic is Publicly Accessible For Subscription | AWS.SNS.NS.Medium.1044 | AC_AWS_0385 |
| Data Protection | AST | MEDIUM | Ensure SNS topic is Encrypted using KMS master key | AWS.AST.DP.MEDIUM.0037 | AC_AWS_0502 |


### aws_apigatewayv2_api
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | ApiGatewayV2Api | MEDIUM | Insecure Cross-Origin Resource Sharing Configuration allowing all domains | AWS.ApiGatewayV2Api.AccessControl.0630 | AC_AWS_0441 |


### aws_efs_file_system
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | EFS | HIGH | Enable encryption of your EFS file systems in order to protect your data and metadata from breaches or unauthorized access and fulfill compliance requirements for data-at-rest encryption within your organization. | AWS.EFS.EncryptionandKeyManagement.High.0409 | AC_AWS_0097 |
| Data Protection | EFS | HIGH | Enable encryption of your EFS file systems in order to protect your data and metadata from breaches or unauthorized access and fulfill compliance requirements for data-at-rest encryption within your organization. | AWS.EFS.EncryptionandKeyManagement.High.0410 | AC_AWS_0098 |


### aws_lb_listener
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | ALL | MEDIUM | Ensure there is a listener configured on HTTPs or with a port 443 | AWS.ALL.IS.MEDIUM.0046 | AC_AWS_0491 |


### aws_sqs_queue
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | SQS | HIGH | Identify any publicly accessible SQS queues available in your AWS account and update their permissions in order to protect against unauthorized users. | AWS.SQS.NetworkSecurity.High.0569 | AC_AWS_0365 |
| Security Best Practices | SQS | HIGH | Ensure that your Amazon Simple Queue Service (SQS) queues are protecting the contents of their messages using Server-Side Encryption (SSE). The SQS service uses an AWS KMS Customer Master Key (CMK) to generate data keys required for the encryption/decryption process of SQS messages. There is no additional charge for using SQS Server-Side Encryption, however, there is a charge for using AWS KMS | AWS.SQS.NetworkSecurity.High.0570 | AC_AWS_0366 |


### aws_docdb_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ADC | MEDIUM | Ensure DocDb is encrypted at rest | AWS.ADC.DP.MEDIUM.0022 | AC_AWS_0455 |
| Logging and Monitoring | ADC | MEDIUM | Ensure DocDb clusters have log exports enabled. | AWS.ADC.LM.MEDIUM.0069 | AC_AWS_0456 |


### aws_cloudwatch_log_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | ACLG | MEDIUM | Ensure AWS Cloudwatch log group has retention policy set. | AWS.ACLG.LM.MEDIUM.0068 | AC_AWS_0452 |


### aws_instance
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | EC2 instances should disable IMDS or require IMDSv2 as this can be related to the weaponization phase of kill chain | AC-AWS-NS-IN-M-1172 | AC_AWS_0479 |
| Identity and Access Management | json | HIGH | Ensure that instance launched follows the least privilege principle as this can be related to delivery-exploitation-Installation phases of kill  chain | AC-AW-IA-LC-H-0442 | AC_AWS_0477 |
| Logging and Monitoring | AI | HIGH | Ensure that detailed monitoring is enabled for EC2 instances. | AWS.AI.LM.HIGH.0070 | AC_AWS_0480 |
| Infrastructure Security | json | HIGH | Security group attached to launch configuration is wide open to internet and this can be related to reconnaissance phase | AC-AW-IS-LC-H-0443 | AC_AWS_0478 |
| Infrastructure Security | json | MEDIUM | Ensure that your AWS application is not deployed within the default Virtual Private Cloud in order to follow security best practices | AC-AW-IS-IN-M-0144 | AC_AWS_0153 |


### aws_config
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | Config | MEDIUM | Ensure AWS Config Rule is enabled for Encrypted Volumes | AWS.Config.EncryptionandKeyManagement.Medium.0660 | AC_AWS_0048 |


### aws_cloudformation_stack
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Security Best Practices | CloudFormation | MEDIUM | AWS CloudFormation Not In Use | AWS.CloudFormation.Medium.0599 | AC_AWS_0444 |
| Security Best Practices | CloudFormation | MEDIUM | AWS CloudFormation Stack Policy | AWS.CloudFormation.Medium.0604 | AC_AWS_0445 |
| Security Best Practices | CloudFormation | MEDIUM | Enable AWS CloudFormation Stack Termination Protection | AWS.CloudFormation.Medium.0605 | AC_AWS_0022 |
| Security Best Practices | CloudFormation | MEDIUM | Enable AWS CloudFormation Stack Notifications | AWS.CloudFormation.Medium.0603 | AC_AWS_0021 |


### aws_iam_user_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | IamUser | HIGH | Ensure Hardware MFA device is enabled for the "root" account | AWS.IamUser.IAM.High.0387 | AC_AWS_0150 |
| Identity and Access Management | json | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AC-AW-IA-H-1190 | AC_AWS_0474 |
| Identity and Access Management | AIUP | MEDIUM | Ensure IAM policies are attached only to groups or roles | AWS.AIUP.IAM.MEDIUM.0049 | AC_AWS_0475 |
| Compliance Validation | IamUser | HIGH | Ensure Virtual MFA device is enabled for the "root" account | AWS.IamUser.IAM.High.0388 | AC_AWS_0149 |
| Compliance Validation | IamUser | HIGH | It is recommended that MFA be enabled for all accounts that have a console password. Enabling MFA provides increased security for console access as it requires the authenticating principal to possess a device that emits a time-sensitive key and have knowledge of a credential | AWS.IamUser.IAM.High.0389 | AC_AWS_0151 |


### aws_ecs_task_definition
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | EcsCluster | HIGH | Like any other EC2 instance it is recommended to place ECS instance within a VPC. AWS VPCs provides the controls to facilitate a formal process for approving and testing all network connections and changes to the firewall and router configurations | AWS.EcsCluster.NetworkSecurity.High.0104 | AC_AWS_0088 |
| Infrastructure Security | AETD | MEDIUM | Ensure EFS volume used for ECS task defination has in transit encryption enabled | AWS.AETD.IS.MEDIUM.0043 | AC_AWS_0463 |
| Data Protection | LaunchConfiguration | HIGH | Sensitive Information Disclosure | AWS.LaunchConfiguration.DataSecurity.High.0101 | AC_AWS_0095 |


### aws_ecr_repository_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | ECR | HIGH | Identify any exposed Amazon ECR image repositories available within your AWS account and update their permissions in order to protect against unauthorized access. Amazon Elastic Container Registry (ECR) is a managed Docker registry service that makes it easy for DevOps teams to store, manage and deploy Docker container images. An ECR repository is a collection of Docker images available on AWS cloud. | AWS.ECR.DataSecurity.High.0579 | AC_AWS_0084 |


### aws_iam_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | json | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AC-AW-IA-H-1187 | AC_AWS_0144 |


### aws_dynamodb_table
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Resilience | DynamoDb | MEDIUM | Ensure Point In Time Recovery is enabled for DynamoDB Tables | AWS.DynamoDb.Logging.Medium.007 | AC_AWS_0458 |
| Data Protection | ADT | MEDIUM | Ensure DynamoDb is encrypted at rest | AWS.ADT.DP.MEDIUM.0025 | AC_AWS_0457 |


### aws_apigatewayv2_stage
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | ApiGatewayV2Stage | LOW | AWS API Gateway V2 Stage is missing access logs | AWS.ApiGatewayV2Stage.Logging.0630 | AC_AWS_0442 |


### aws_ecr_repository
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | AER | MEDIUM | Ensure ECR repository is encrypted at rest | AWS.AER.DP.MEDIUM.0026 | AC_AWS_0461 |
| Configuration and Vulnerability Analysis | ECR | MEDIUM | Unscanned images may contain vulnerabilities | AWS.ECR.DataSecurity.High.0578 | AC_AWS_0083 |
| Identity and Access Management | AER | MEDIUM | Ensure ECR repository has policy attached. | AWS.AER.DP.MEDIUM.0058 | AC_AWS_0462 |


### aws_cloudfront_distribution
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | CloudFront | HIGH | Secure ciphers are not used in CloudFront distribution | AWS.CloudFront.EncryptionandKeyManagement.High.0408 | AC_AWS_0023 |
| Infrastructure Security | json | LOW | Ensure that geo restriction is enabled for your Amazon CloudFront CDN distribution to whitelist or blacklist a country in order to allow or restrict users in specific locations from accessing web application content. | AC-AW-IS-CD-M-0026 | AC_AWS_0026 |
| Infrastructure Security | json | MEDIUM | Ensure that cloud-front has web application firewall enabled | AC-AW-IS-CD-M-1186 | AC_AWS_0032 |
| Data Protection | CloudFront | HIGH | Use encrypted connection between CloudFront and origin server | AWS.CloudFront.EncryptionandKeyManagement.High.0407 | AC_AWS_0024 |
| Logging and Monitoring | CloudFront | MEDIUM | Ensure that your AWS Cloudfront distributions have the Logging feature enabled in order to track all viewer requests for the content delivered through the Content Delivery Network (CDN). | AWS.CloudFront.Logging.Medium.0567 | AC_AWS_0025 |


### aws_cloudwatch
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | CloudWatch | MEDIUM | App-Tier CloudWatch Log Group Retention Period | AWS.CloudWatch.Logging.Medium.0631 | AC_AWS_0041 |
| Data Protection | CloudWatch | HIGH | AWS CloudWatch log group is not encrypted with a KMS CMK | AWS.CloudWatch.EncryptionandKeyManagement.High.0632 | AC_AWS_0451 |


### aws_ami_launch_permission
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | AMI | MEDIUM | Limit access to AWS AMIs | AWS.AMI.NS.Medium.1040 | AC_AWS_0006 |


### aws_launch_configuration
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | EcsCluster | HIGH | Ensure that AWS ECS clusters are encrypted. Data encryption at rest, prevents unauthorized users from accessing sensitive data on your AWS ECS clusters and associated cache storage systems. | AWS.EcsCluster.EncryptionandKeyManagement.High.0413 | AC_AWS_0167 |
| Identity and Access Management | json | HIGH | Ensure that launch configuration launched follows the least privilege principle | AC-AW-IA-LC-H-0441 | AC_AWS_0488 |
| Data Protection | LaunchConfiguration | HIGH | Avoid using base64 encoded private keys as part of config | AWS.LaunchConfiguration.DataSecurity.High.0102 | AC_AWS_0168 |
| Data Protection | LaunchConfiguration | HIGH | Avoid using base64 encoded shell script as part of config | AWS.LaunchConfiguration.DataSecurity.High.0101 | AC_AWS_0170 |
| Logging and Monitoring | json | MEDIUM | It is important to enable cloudWatch monitoring incase monitoring the activity | AC-AW-LM-LC-M-0440 | AC_AWS_0490 |
| Configuration and Vulnerability Analysis | json | HIGH | Launch configuration uses IMDSv1 which vulnerable to SSRF | AC-AW-CA-LC-H-0439 | AC_AWS_0487 |
| Infrastructure Security | json | HIGH | Security group attached to launch configuration is wide open to internet  | AC-AW-IS-LC-H-0438 | AC_AWS_0489 |


### aws_athena_database
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ADB | MEDIUM | Ensure Athena Database is encrypted at rest | AWS.ADB.DP.MEDIUM.016 | AC_AWS_0443 |


### aws_api_gateway_stage
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | API Gateway | MEDIUM | Enable SSL Client Certificate | AWS.APIGateway.NetworkSecurity.Medium.0565 | AC_AWS_0013 |
| Logging and Monitoring | API Gateway | MEDIUM | Ensure that AWS CloudWatch logs are enabled for all your APIs created with Amazon API Gateway service in order to track and analyze execution behavior at the API stage level. | AWS.APIGateway.Logging.Medium.0572 | AC_AWS_0012 |
| Logging and Monitoring | API Gateway | LOW | Ensure AWS API Gateway has active xray tracing enabled | AWS.APIGateway.Logging.Medium.0571 | AC_AWS_0015 |
| Logging and Monitoring | API Gateway | MEDIUM | Enable AWS CloudWatch Logs for APIs | AWS.APIGateway.Logging.Medium.0567 | AC_AWS_0014 |


### aws_elasticsearch_domain
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ElasticSearch | HIGH | Enable AWS ElasticSearch Encryption At Rest | AWS.ElasticSearch.EKM.Medium.0778 | AC_AWS_0112 |
| Infrastructure Security | ElasticSearch | MEDIUM | Ensure Elasticsearch domains being created are set to be encrypted node-to-node | AWS.ElasticSearch.IS.MEDIUM.0045 | AC_AWS_0468 |
| Data Protection | ElasticSearch | MEDIUM | ElasticSearch Domain Encrypted with KMS CMKs | AWS.ElasticSearch.EKM.Medium.0768 | AC_AWS_0111 |
| Compliance Validation | Elasticsearch | MEDIUM | Ensure that your AWS Elasticsearch clusters have enabled the support for publishing slow logs to AWS CloudWatch Logs. This feature enables you to publish slow logs from the indexing and search operations performed on your ES clusters and gain full insight into the performance of these operations. | AWS.Elasticsearch.Logging.Medium.0573 | AC_AWS_0105 |


### aws_iam_user_login_profile
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | Iam | HIGH | Password policies are, in part, used to enforce password complexity requirements. IAM password policies can be used to ensure password are comprised of different character sets, have minimal length, rotation and history restrictions | AWS.Iam.IAM.High.0391 | AC_AWS_0148 |


### aws_iam_group_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | json | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AC-AW-IA-H-0392 | AC_AWS_0143 |


### aws_load_balancer_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | ELB | HIGH | Using insecure ciphers for your ELB Predefined or Custom Security Policy, could make the SSL connection between the client and the load balancer vulnerable to exploits. TLS 1.0 was recommended to be disabled by PCI Council after June 30, 2016 | AWS.ELB.EncryptionandKeyManagement.High.0401 | AC_AWS_0172 |
| Infrastructure Security | ELB | HIGH | Remove insecure ciphers for your ELB Predefined or Custom Security Policy, to reduce the risk of the SSL connection between the client and the load balancer being exploited. | AWS.ELB.EncryptionandKeyManagement.High.0403 | AC_AWS_0171 |


### aws_s3_bucket
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0377 | AC_AWS_0210 |
| Identity and Access Management | S3Bucket | HIGH | Ensure S3 buckets do not have, a both public ACL on the bucket and a public access block. | AWS.S3Bucket.IAM.HIGH.0065 | AC_AWS_0496 |
| Logging and Monitoring | S3Bucket | MEDIUM | Ensure S3 buckets have access logging enabled. | AWS.S3Bucket.LM.MEDIUM.0078 | AC_AWS_0497 |
| Infrastructure Security | S3Bucket | LOW | Ensure that there are not any static websites being hosted on buckets you aren't aware of | AWS.S3Bucket.NetworkSecurity.High.0417 | AC_AWS_0208 |
| Resilience | S3Bucket | HIGH | Enabling S3 versioning will enable easy recovery from both unintended user actions, like deletes and overwrites | AWS.S3Bucket.IAM.High.0370 | AC_AWS_0214 |
| Identity and Access Management | S3Bucket | HIGH | S3 bucket Access is allowed to all AWS Account Users. | AWS.S3Bucket.DS.High.1043 | AC_AWS_0215 |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0379 | AC_AWS_0212 |
| Data Protection | S3Bucket | HIGH | Ensure that S3 Buckets have server side encryption at rest enabled with KMS key to protect sensitive data. | AWS.S3Bucket.EncryptionandKeyManagement.High.0405 | AC_AWS_0207 |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0378 | AC_AWS_0211 |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0381 | AC_AWS_0213 |


### aws_eks_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | AEC | MEDIUM | Ensure EKS clusters have control plane logging enabled. | AWS.AEC.LM.MEDIUM.0071 | AC_AWS_0465 |


### aws_elb
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | ELB | MEDIUM | Ensure AWS ELB has access logging enabled. | AWS.ELB.LM.MEDIUM.0072 | AC_AWS_0470 |
| Infrastructure Security | ELB | LOW | AWS ELB incoming traffic not encrypted | AWS.ELB.NetworkPortsSecurity.Low.0563 | AC_AWS_0120 |


### aws_redshift_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | Redshift | HIGH | Ensure Redshift clusters are not publicly accessible to minimize security risks. | AWS.Redshift.NetworkSecurity.HIGH.0564 | AC_AWS_0199 |
| Logging and Monitoring | Redshift | LOW | Ensure AWS Redshift cluster instances have logging enabled. | AWS.Redshift.Logging.Medium.0565 | AC_AWS_0200 |
| Data Protection | Redshift | MEDIUM | Use customer-managed KMS keys instead of AWS-managed keys, to have granular control over encrypting and encrypting data. Encrypt Redshift clusters with a Customer-managed KMS key. This is a recommended best practice. | AWS.Redshift.EncryptionandKeyManagement.High.0415 | AC_AWS_0198 |


### aws_elasticcache_replication_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | AERG | MEDIUM | Ensure Elastic Cache Replication Group is encrypted at rest | AWS.AERG.DP.MEDIUM.0027 | AC_AWS_0466 |
| Data Protection | AERG | MEDIUM | Ensure Elastic Cache Replication Group is encrypted in transit | AWS.AERG.DP.MEDIUM.0044 | AC_AWS_0467 |


### aws_kinesis_stream
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | Kinesis | HIGH | Ensure Kinesis Stream is encrypted | AWS.Kinesis.EncryptionandKeyManagement.High.0412 | AC_AWS_0157 |


### aws_config_configuration_aggregator
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | Config | HIGH | Ensure AWS Config is enabled in all regions | AWS.Config.Logging.HIGH.0590 | AC_AWS_0049 |


### aws_s3_bucket_object
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ASBO | MEDIUM | Ensure S3 object is Encrypted | AWS.ASBO.DP.MEDIUM.0034 | AC_AWS_0498 |


### aws_route53_query_log
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | Route53 query logs | MEDIUM | Ensure CloudWatch logging is enabled for Route53 hosted zones. | AWS.Route53querylogs.Logging.Medium.0574 | AC_AWS_0204 |


### aws_secretsmanager_secret
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | SecretsManagerSecret | MEDIUM | Ensure SecretsManager Secrets are Encrypted using KMS key | AWS.SecretsManagerSecret.DP.MEDIUM.0036 | AC_AWS_0500 |


### aws_iam_access_key
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | IamUser | HIGH | The root account is the most privileged user in an AWS account. AWS Access Keys provide programmatic access to a given AWS account. It is recommended that all access keys associated with the root account be removed. Removing access keys associated with the root account limits vectors by which the account can be compromised. Additionally, removing the root access keys encourages the creation and use of role based accounts that are least privileged. | AWS.IamUser.IAM.High.0390 | AC_AWS_0132 |
| Identity and Access Management | IamUser | MEDIUM | Ensure that there are no exposed Amazon IAM access keys in order to protect your AWS resources against unapproved access | AWS.IamUser.IAM.High.0391 | AC_AWS_0133 |


### aws_neptune_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ANC | MEDIUM | Ensure Neptune Cluster is Encrypted | AWS.ANC.DP.MEDIUM.0030 | AC_AWS_0493 |
| Logging and Monitoring | ANC | MEDIUM | Ensure AWS Neptune clusters have logging enabled. | AWS.ANC.LM.MEDIUM.0075 | AC_AWS_0494 |


### aws_dax_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ADC | MEDIUM | Ensure server-side encryption is enabled for AWS DAX Cluster | AWS.ADC.DP.MEDIUM.0021 | AC_AWS_0375 |


### aws_guardduty_detector
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | GuardDuty Enabled | MEDIUM | Ensure that Amazon GuardDuty service is currently enabled in all regions in order to protect your AWS environment and infrastructure (AWS accounts and resources, IAM credentials, guest operating systems, applications, etc) against security threats. AWS GuardDuty is a managed threat detection service that continuously monitors your VPC flow logs, AWS CloudTrail event logs and DNS logs for malicious or unauthorized behavior. The service monitors for activity such as unusual API calls, potentially compromised EC2 instances or potentially unauthorized deployments that indicate a possible AWS account compromise. AWS GuardDuty operates entirely on Amazon Web Services infrastructure and does not affect the performance or reliability of your applications. The service does not require any software agents, sensors or network appliances. | AWS.GuardDutyEnabled.Security.Medium.0575 | AC_AWS_0131 |


### aws_db_security_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | RDS | HIGH | RDS should not be defined with public interface. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0101 | AC_AWS_0066 |
| Infrastructure Security | RDS | HIGH | RDS should not be open to a large scope. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0103 | AC_AWS_0065 |
| Infrastructure Security | RDS | HIGH | RDS should not be open to a public scope. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0102 | AC_AWS_0067 |


### aws_s3_bucket_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0371 | AC_AWS_0217 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0376 | AC_AWS_0224 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0375 | AC_AWS_0221 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0374 | AC_AWS_0220 |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0372 | AC_AWS_0218 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0373 | AC_AWS_0219 |


### aws_ami
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | EC2 | MEDIUM | Enable AWS AMI Encryption | AWS.EC2.EncryptionandKeyManagement.Medium.0688 | AC_AWS_0005 |


### aws_elasticache_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Compliance Validation | ElastiCache | HIGH | ElastiCache for Redis version is not compliant with AWS PCI DSS requirements | AWS.ElastiCache.DataSecurity.High.0425 | AC_AWS_0102 |
| Compliance Validation | ElastiCache | HIGH | ElastiCache for Memcached is not in use in AWS PCI DSS environments | AWS.ElastiCache.DataSecurity.High.0424 | AC_AWS_0103 |
| Resilience | ElastiCache | MEDIUM | AWS ElastiCache Multi-AZ | AWS.ElastiCache.HighAvailability.Medium.0757 | AC_AWS_0104 |


### aws_kinesis_firehose_delivery_stream
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | Kinesis | HIGH | AWS Kinesis Server data at rest has server side encryption (SSE) | AWS.Kinesis.EncryptionandKeyManagement.High.0411 | AC_AWS_0156 |


### aws_rds_cluster
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Resilience | RDS | MEDIUM | Ensure backup retention period is set for rds cluster | AWS.RDS.RE.MEDIUM.0013 | AC_AWS_0495 |
| Data Protection | RDS | HIGH | Encrypt Amazon RDS instances and snapshots at rest, by enabling the encryption option for your Amazon RDS DB instance | AWS.RDS.EncryptionandKeyManagement.High.0414 | AC_AWS_0186 |


### aws_cloudtrail
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | CloudTrail | HIGH | Ensure CloudTrail logs are encrypted using KMS | AWS.CloudTrail.Logging.High.0399 | AC_AWS_0033 |
| Logging and Monitoring | CloudTrail | MEDIUM | Cloud Trail Multi Region not enabled | AWS.CloudTrail.Logging.Medium.004 | AC_AWS_0448 |
| Security Best Practices | CloudTrail | MEDIUM | Ensure that EC2 is EBS optimized | AWS.CloudTrail.Logging.Medium.008 | AC_AWS_0449 |
| Security Best Practices | CloudTrail | LOW | ECR should have an image tag be immutable | AWS.CloudTrail.Logging.Low.009 | AC_AWS_0447 |
| Logging and Monitoring | CloudTrail | MEDIUM | Cloud Trail Multi Region not enabled | AWS.CloudTrail.Logging.Medium.0460 | AC_AWS_0034 |
| Logging and Monitoring | CloudTrail | MEDIUM | Ensure CloudTrail has log file validation enabled. | AWS.CloudTrail.LM.MEDIUM.0087 | AC_AWS_0446 |
| Logging and Monitoring | CloudTrail | MEDIUM | Ensure appropriate subscribers to each SNS topic | AWS.CloudTrail.Logging.Low.0559 | AC_AWS_0035 |
| Logging and Monitoring | Config | MEDIUM | Ensure AWS Config is enabled in all regions | AWS.Config.Logging.Medium.0590 | AC_AWS_0450 |


### aws_sagemaker_notebook_instance
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Data Protection | ASNI | MEDIUM | Ensure SageMaker Instance is Encrypted | AWS.ASNI.DP.MEDIUM.0035 | AC_AWS_0499 |


### aws_lambda_function
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | LambdaFunction | MEDIUM | Ensure AWS Lambda function has policy attached. | AWS.LambdaFunction.LM.MEIDUM.0063 | AC_AWS_0484 |
| Infrastructure Security | LambdaFunction | LOW | Lambda function does not include a VPC configuration. | AWS.LambdaFunction.Logging.0472 | AC_AWS_0486 |
| Logging and Monitoring | LambdaFunction | LOW | Lambda tracing is not enabled. | AWS.LambdaFunction.Logging.0470 | AC_AWS_0485 |
| Data Protection | LambdaFunction | HIGH | Lambda does not use KMS CMK key to protect environment variables. | AWS.LambdaFunction.EncryptionandKeyManagement.0471 | AC_AWS_0483 |
| Logging and Monitoring | LambdaFunction | LOW | Lambda tracing is not enabled. | AWS.LambdaFunction.Logging.0470 | AC_AWS_0163 |


### aws_kms_key
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | AKK | HIGH | Ensure IAM policies do not have 'Principal' element missing from the policy statement. | AWS.AKK.IAM.HIGH.0012 | AC_AWS_0481 |
| Identity and Access Management | KMS | HIGH | Identify any publicly accessible AWS Key Management Service master keys and update their access policy in order to stop any unsigned requests made to these resources. | AWS.KMS.NetworkSecurity.High.0566 | AC_AWS_0162 |
| Data Protection | AKK | HIGH | Ensure rotation for customer created CMKs is enabled | AWS.AKK.DP.HIGH.0012 | AC_AWS_0160 |
| Identity and Access Management | AKK | HIGH | Ensure KMS key policy does not have wildcard policies attached. | AWS.AKK.IAM.HIGH.0082 | AC_AWS_0482 |


### aws_security_group
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (UDP,11214) is not accessible by a public CIDR block range | AC_AWS_0292 | AC_AWS_0292 |
| Infrastructure Security | json | LOW | Ensure Cassandra' (TCP,7001) is not exposed to private hosts more than 32 | AC_AWS_0338 | AC_AWS_0338 |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (TCP,8080) is not accessible by a CIDR block range | AC_AWS_0284 | AC_AWS_0284 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (TCP,2484) is not accessible by a public CIDR block range | AC_AWS_0302 | AC_AWS_0302 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MSSQL Server (TCP,1433) | AC_AWS_0247 | AC_AWS_0247 |
| Infrastructure Security | json | LOW | Ensure SNMP' (UDP,161) is not exposed to private hosts more than 32 | AC_AWS_0355 | AC_AWS_0355 |
| Infrastructure Security | json | LOW | Ensure NetBIOSNameService' (TCP,137) is not exposed to private hosts more than 32 | AC_AWS_0343 | AC_AWS_0343 |
| Infrastructure Security | json | HIGH | Ensure SMTP (TCP,25) is not accessible by a public CIDR block range | AC_AWS_0314 | AC_AWS_0314 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Memcached SSL (UDP,11215) | AC_AWS_0251 | AC_AWS_0251 |
| Infrastructure Security | json | LOW | Ensure Elasticsearch' (TCP,9300) is not exposed to private hosts more than 32 | AC_AWS_0363 | AC_AWS_0363 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Telnet (TCP,23) | AC_AWS_0271 | AC_AWS_0271 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MemcachedSSL (UDP,11214) is not exposed to  private hosts more than 32 | AC_AWS_0334 | AC_AWS_0334 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SQL Server Analysis Service browser (TCP,2382) | AC_AWS_0267 | AC_AWS_0267 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Debugger (TCP,135) is not accessible by a public CIDR block range | AC_AWS_0288 | AC_AWS_0288 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports https (TCP,443) is not exposed to private hosts more than 32 | AC_AWS_0322 | AC_AWS_0322 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - remote desktop port (TCP,3389) | AC_AWS_0230 | AC_AWS_0230 |
| Infrastructure Security | json | LOW | Ensure Telnet' (TCP,23) is not exposed to private hosts more than 32 | AC_AWS_0359 | AC_AWS_0359 |
| Infrastructure Security | json | HIGH | Ensure Elasticsearch (TCP,9300) is not accessible by a public CIDR block range | AC_AWS_0318 | AC_AWS_0318 |
| Infrastructure Security | json | LOW | Ensure SSH (TCP,22) is not accessible by a public CIDR block range | AC_AWS_0319 | AC_AWS_0319 |
| Infrastructure Security | json | LOW | Ensure OracleDatabaseServer' (TCP,521) is not exposed to private hosts more than 32 | AC_AWS_0358 | AC_AWS_0358 |
| Infrastructure Security | json | HIGH | Ensure no security groups allow ingress from 0.0.0.0/0 to ALL ports and protocols | AC_AWS_0231 | AC_AWS_0231 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SNMP (UDP,161) | AC_AWS_0266 | AC_AWS_0266 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Server (TCP,1433) is not accessible by a public CIDR block range | AC_AWS_0289 | AC_AWS_0289 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports remote desktop port (TCP,3389) is not exposed to private hosts more than 32 | AC_AWS_0323 | AC_AWS_0323 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Oracle Database Server (TCP,1521) | AC_AWS_0270 | AC_AWS_0270 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MemcachedSSL (UDP,11215) is not exposed to  private hosts more than 32 | AC_AWS_0335 | AC_AWS_0335 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - (SSH,22) | AC_AWS_0227 | AC_AWS_0227 |
| Infrastructure Security | json | LOW | Ensure MongoDB' (TCP,27017) is not exposed to private hosts more than 32 | AC_AWS_0362 | AC_AWS_0362 |
| Infrastructure Security | json | HIGH | Ensure CIFS for file/printer (TCP,445) is not accessible by a public CIDR block range | AC_AWS_0315 | AC_AWS_0315 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Memcached SSL (UDP,11214) | AC_AWS_0250 | AC_AWS_0250 |
| Infrastructure Security | json | LOW | Ensure MongoWebPortal' (TCP,27018) is not exposed to private hosts more than 32 | AC_AWS_0342 | AC_AWS_0342 |
| Infrastructure Security | json | LOW | Ensure PuppetMaster' (TCP,8140) is not exposed to private hosts more than 32 | AC_AWS_0354 | AC_AWS_0354 |
| Infrastructure Security | json | MEDIUM | Ensure Oracle DB SSL (UDP,2484) is not accessible by a public CIDR block range | AC_AWS_0303 | AC_AWS_0303 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MSSQL Debugger (TCP,135) | AC_AWS_0246 | AC_AWS_0246 |
| Infrastructure Security | json | MEDIUM | Ensure LDAP SSL (TCP,636) is not accessible by a public CIDR block range | AC_AWS_0285 | AC_AWS_0285 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (UDP,11215) is not accessible by a public CIDR block range | AC_AWS_0293 | AC_AWS_0293 |
| Infrastructure Security | json | LOW | Ensure HadoopNameNode' (TCP,9000) is not exposed to private hosts more than 32 | AC_AWS_0339 | AC_AWS_0339 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (TCP,4505) is not accessible by a public CIDR block range | AC_AWS_0277 | AC_AWS_0277 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MemcachedSSL (TCP,11214) is not exposed to private hosts more than 32 | AC_AWS_0332 | AC_AWS_0332 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Datagram Service (TCP,138) is not accessible by a public CIDR block range | AC_AWS_0298 | AC_AWS_0298 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Oracle DB SSL (UDP,2484) | AC_AWS_0261 | AC_AWS_0261 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports Elasticsearch (TCP,9200) is not exposed to  private hosts more than 32 | AC_AWS_0324 | AC_AWS_0324 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SaltStack Master (TCP,4506) | AC_AWS_0236 | AC_AWS_0236 |
| Infrastructure Security | json | MEDIUM | Ensure SNMP (UDP,161) is not accessible by a public CIDR block range | AC_AWS_0308 | AC_AWS_0308 |
| Infrastructure Security | json | LOW | Ensure OracleDBSSL' (TCP,2484) is not exposed to private hosts more than 32 | AC_AWS_0349 | AC_AWS_0349 |
| Infrastructure Security | json | MEDIUM | Ensure Mongo Web Portal (TCP,27018) is not accessible by a public CIDR block range | AC_AWS_0294 | AC_AWS_0294 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MSSQLAdmin (TCP,1434) is not exposed to private hosts more than 32 | AC_AWS_0328 | AC_AWS_0328 |
| Infrastructure Security | json | MEDIUM | Ensure Hadoop Name Node (TCP,9000) is not accessible by a public CIDR block range | AC_AWS_0282 | AC_AWS_0282 |
| Infrastructure Security | json | MEDIUM | Ensure Postgres SQL (TCP,5432) is not accessible by a public CIDR block range | AC_AWS_0304 | AC_AWS_0304 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Known internal web port (TCP,8000) | AC_AWS_0241 | AC_AWS_0241 |
| Infrastructure Security | json | LOW | Ensure Prevalentknowninternalport' (TCP,3000) is not exposed to private hosts more than 32 | AC_AWS_0353 | AC_AWS_0353 |
| Infrastructure Security | json | LOW | Ensure NetBIOSNameService' (UDP,137) is not exposed to private hosts more than 32 | AC_AWS_0345 | AC_AWS_0345 |
| Infrastructure Security | json | HIGH | Ensure Oracle Database Server (TCP,1521) is not accessible by a public CIDR block range | AC_AWS_0312 | AC_AWS_0312 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Datagram Service (UDP,138) | AC_AWS_0257 | AC_AWS_0257 |
| Infrastructure Security | json | HIGH | Ensure Telnet (TCP,23) is not accessible by a public CIDR block range | AC_AWS_0313 | AC_AWS_0313 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Datagram Service (TCP,138) | AC_AWS_0256 | AC_AWS_0256 |
| Infrastructure Security | json | LOW | Ensure NetBIOSNameService' (UDP,137) is not exposed to private hosts more than 32 | AC_AWS_0344 | AC_AWS_0344 |
| Infrastructure Security | json | LOW | Ensure PostgresSQL' (UDP,5432) is not exposed to private hosts more than 32 | AC_AWS_0352 | AC_AWS_0352 |
| Infrastructure Security | json | MEDIUM | Ensure Postgres SQL (UDP,5432) is not accessible by a CIDR block range | AC_AWS_0305 | AC_AWS_0305 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Hadoop Name Node (TCP,9000) | AC_AWS_0240 | AC_AWS_0240 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MSSQLBrowserService (UDP,1434) is not exposed to private hosts more than 32 | AC_AWS_0329 | AC_AWS_0329 |
| Infrastructure Security | json | MEDIUM | Ensure Known internal web port (TCP,8000) is not accessible by a public CIDR block range | AC_AWS_0283 | AC_AWS_0283 |
| Infrastructure Security | json | MEDIUM | Ensure MySQL (TCP,3306) is not accessible by a public CIDR block range | AC_AWS_0295 | AC_AWS_0295 |
| Infrastructure Security | json | LOW | Ensure NetBIOSSessionService' (UDP,139) is not exposed to private hosts more than 32 | AC_AWS_0348 | AC_AWS_0348 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis Service browser (TCP,2382) is not accessible by a public CIDR block range | AC_AWS_0309 | AC_AWS_0309 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - CIFS / SMB (TCP,3020) | AC_AWS_0237 | AC_AWS_0237 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Oracle DB SSL (TCP,2484) | AC_AWS_0260 | AC_AWS_0260 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports SaltStackMaster (TCP,4506) is not exposed to private hosts more than 32 | AC_AWS_0325 | AC_AWS_0325 |
| Infrastructure Security | json | HIGH | Ensure Unknown Port is not exposed to the entire internet | AC_AWS_0276 | AC_AWS_0276 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MemcachedSSL (TCP,11215) is not exposed to private hosts more than 32 | AC_AWS_0333 | AC_AWS_0333 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Datagram Service (UDP,138) is not accessible by a public CIDR block range | AC_AWS_0299 | AC_AWS_0299 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Session Service (UDP,139) | AC_AWS_0259 | AC_AWS_0259 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Elasticsearch (TCP,9200) | AC_AWS_0234 | AC_AWS_0234 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports CassandraOpsCenteragent (TCP,61621) is not exposed to private hosts more than 32 | AC_AWS_0326 | AC_AWS_0326 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Postgres SQL (UDP,5432) | AC_AWS_0263 | AC_AWS_0263 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MSSQLDebugger (TCP,135) is not exposed to private hosts more than 32 | AC_AWS_0330 | AC_AWS_0330 |
| Infrastructure Security | json | HIGH | Ensure no security groups is wide open to public, that is, allows traffic from 0.0.0.0/0 to ALL ports and protocols | AC_AWS_0275 | AC_AWS_0275 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Name Service (UDP,137) | AC_AWS_0255 | AC_AWS_0255 |
| Infrastructure Security | json | MEDIUM | Ensure SQL Server Analysis Services (TCP,2383) is not accessible by a public CIDR block range | AC_AWS_0310 | AC_AWS_0310 |
| Infrastructure Security | json | LOW | Ensure NetBIOSSessionService' (TCP,139) is not exposed to private hosts more than 32 | AC_AWS_0347 | AC_AWS_0347 |
| Infrastructure Security | json | LOW | Ensure PostgresSQL' (TCP,5432) is not exposed to private hosts more than 32 | AC_AWS_0351 | AC_AWS_0351 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - LDAP SSL (TCP,636) | AC_AWS_0243 | AC_AWS_0243 |
| Infrastructure Security | json | MEDIUM | Ensure Prevalent known internal port (TCP,3000) is not accessible by a public CIDR block range | AC_AWS_0306 | AC_AWS_0306 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Cassandra OpsCenter agent (TCP,61621) | AC_AWS_0238 | AC_AWS_0238 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra OpsCenter agent port (TCP,61621) is not accessible by a public CIDR block range | AC_AWS_0280 | AC_AWS_0280 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (TCP,137) is not accessible by a public CIDR block range | AC_AWS_0296 | AC_AWS_0296 |
| Infrastructure Security | json | MEDIUM | Ensure CIFS / SMB (TCP,3020) is not accessible by a public CIDR block range | AC_AWS_0279 | AC_AWS_0279 |
| Infrastructure Security | json | MEDIUM | Ensure NetBIOS Name Service (UDP,137) is not accessible by a public CIDR block range | AC_AWS_0297 | AC_AWS_0297 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (TCP,4506) is not accessible by a public CIDR block range | AC_AWS_0278 | AC_AWS_0278 |
| Infrastructure Security | json | MEDIUM | Ensure Cassandra (TCP,7001) is not accessible by a public CIDR block range | AC_AWS_0281 | AC_AWS_0281 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Cassandra (TCP,7001) | AC_AWS_0239 | AC_AWS_0239 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Known internal web port (TCP,8080) | AC_AWS_0242 | AC_AWS_0242 |
| Infrastructure Security | json | MEDIUM | Ensure Puppet Master (TCP:8140) is not accessible by a public CIDR block range | AC_AWS_0307 | AC_AWS_0307 |
| Infrastructure Security | json | LOW | Ensure OracleDBSSL' (UDP,2484) is not exposed to private hosts more than 32 | AC_AWS_0350 | AC_AWS_0350 |
| Infrastructure Security | json | LOW | Ensure NetBIOSDatagramService' (UDP,138) is not exposed to private hosts more than 32 | AC_AWS_0346 | AC_AWS_0346 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Name Service (TCP,137) | AC_AWS_0254 | AC_AWS_0254 |
| Infrastructure Security | json | MEDIUM | Ensure SaltStack Master (TCP,4505) is not accessible by a public CIDR block range | AC_AWS_0311 | AC_AWS_0311 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MSSQLServer (TCP,1433) is not exposed to private hosts more than 32 | AC_AWS_0331 | AC_AWS_0331 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MongoDB (TCP,27017) | AC_AWS_0274 | AC_AWS_0274 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports Knowninternalwebport (TCP,8080) is not exposed to private hosts more than 32 | AC_AWS_0327 | AC_AWS_0327 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Postgres SQL (TCP,5432) | AC_AWS_0262 | AC_AWS_0262 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Elasticsearch (TCP,9300) | AC_AWS_0235 | AC_AWS_0235 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - NetBIOS Session Service (TCP,139) | AC_AWS_0258 | AC_AWS_0258 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MySQL (TCP,3306) | AC_AWS_0253 | AC_AWS_0253 |
| Infrastructure Security | json | HIGH | Ensure MongoDB (TCP,27017) is not accessible by a public CIDR block range | AC_AWS_0316 | AC_AWS_0316 |
| Infrastructure Security | json | LOW | Ensure LDAPSSL' (TCP,636) is not exposed to private hosts more than 32 | AC_AWS_0341 | AC_AWS_0341 |
| Infrastructure Security | json | LOW | Ensure SQLServerAnalysisServices' (TCP,2383) is not exposed to private hosts more than 32 | AC_AWS_0357 | AC_AWS_0357 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MSSQL Browser Service (UDP,1434) | AC_AWS_0245 | AC_AWS_0245 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Session Service (TCP,139) is not accessible by a CIDR block range | AC_AWS_0300 | AC_AWS_0300 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Admin (TCP,1434) is not accessible by a public CIDR block range | AC_AWS_0286 | AC_AWS_0286 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SaltStack Master (TCP,4505) | AC_AWS_0269 | AC_AWS_0269 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (TCP,11214) is not accessible by a public CIDR block range | AC_AWS_0290 | AC_AWS_0290 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - (HTTP,80) | AC_AWS_0228 | AC_AWS_0228 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Memcached SSL (TCP,11215) | AC_AWS_0249 | AC_AWS_0249 |
| Infrastructure Security | json | HIGH | Ensure no default security groups are used as they allow ingress from 0.0.0.0/0 to ALL ports and protocols | AC_AWS_0232 | AC_AWS_0232 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports SSH (TCP,22) is not exposed to private hosts more than 32 | AC_AWS_0320 | AC_AWS_0320 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Puppet Master (TCP,8140) | AC_AWS_0265 | AC_AWS_0265 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports MySQL (TCP,3306) is not exposed to private hosts more than 32 | AC_AWS_0336 | AC_AWS_0336 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - CIFS for file/printer (TCP,445) | AC_AWS_0273 | AC_AWS_0273 |
| Infrastructure Security | json | LOW | Ensure CIFSforfile/printer' (TCP,445) is not exposed to private hosts more than 32 | AC_AWS_0361 | AC_AWS_0361 |
| Infrastructure Security | json | LOW | Ensure SMTP' (TCP,25) is not exposed to private hosts more than 32 | AC_AWS_0360 | AC_AWS_0360 |
| Infrastructure Security | json | LOW | Ensure CIFS/SMB' (TCP,3020) is not exposed to private hosts more than 32 | AC_AWS_0337 | AC_AWS_0337 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SMTP (TCP,25) | AC_AWS_0272 | AC_AWS_0272 |
| Infrastructure Security | json | LOW | Ensure Security Groups Unrestricted Specific Ports http (TCP,80) is not exposed to private hosts more than 32 | AC_AWS_0321 | AC_AWS_0321 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Prevalent known internal port (TCP,3000) | AC_AWS_0264 | AC_AWS_0264 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SaltStack Master (TCP,4505) | AC_AWS_0233 | AC_AWS_0233 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Memcached SSL (TCP,11214) | AC_AWS_0248 | AC_AWS_0248 |
| Infrastructure Security | json | LOW | Security Groups - Unrestricted Specific Ports - (HTTPS,443) | AC_AWS_0229 | AC_AWS_0229 |
| Infrastructure Security | json | MEDIUM | Ensure Memcached SSL (TCP,11215) is not accessible by a public CIDR block range | AC_AWS_0291 | AC_AWS_0291 |
| Infrastructure Security | json | MEDIUM | Ensure MSSQL Browser Service (UDP,1434) is not accessible by a public CIDR block range | AC_AWS_0287 | AC_AWS_0287 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - SQL Server Analysis Services (TCP,2383) | AC_AWS_0268 | AC_AWS_0268 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - MSSQL Admin (TCP,1434) | AC_AWS_0244 | AC_AWS_0244 |
| Infrastructure Security | json | MEDIUM | Ensure NetBios Session Service (UDP,139) is not accessible by a CIDR block range | AC_AWS_0301 | AC_AWS_0301 |
| Infrastructure Security | json | LOW | Ensure SQLServerAnalysisServicebrowser' (TCP,2382) is not exposed to private hosts more than 32 | AC_AWS_0356 | AC_AWS_0356 |
| Infrastructure Security | json | LOW | Ensure Knowninternalwebport' (TCP,8000) not exposed to private hosts more than 32 | AC_AWS_0340 | AC_AWS_0340 |
| Infrastructure Security | json | HIGH | Security Groups - Unrestricted Specific Ports - Mongo Web Portal (TCP,27018) | AC_AWS_0252 | AC_AWS_0252 |
| Infrastructure Security | json | HIGH | Ensure Elasticsearch (TCP,9200) is not accessible by a public CIDR block range | AC_AWS_0317 | AC_AWS_0317 |


### aws_api_gateway_method
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Infrastructure Security | APGM | LOW | Ensure there is no open access to back-end resources through API | AWS.APGM.IS.LOW.0056 | AC_AWS_0439 |


### aws_efs_file_system_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | AEFSP | HIGH | Ensure EFS file system does not use insecure wildcard policies. | AWS.AEFSP.IAM.HIGH.0059 | AC_AWS_0464 |


### aws_ecs_service
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | ECS | HIGH | Ensure there are no ECS services Admin roles | AWS.ECS.High.0436 | AC_AWS_0087 |


### aws_globalaccelerator_accelerator
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Logging and Monitoring | AGA | LOW | Ensure Global Accelerator accelerator has flow logs enabled. | AWS.AGA.LM.LOW.0073 | AC_AWS_0471 |


### aws_api_gateway_rest_api_policy
| Category | Resource | Severity | Description | Reference ID | ID |
| -------- | -------- | -------- | ----------- | ------------ | -- |
| Identity and Access Management | APGRAP | HIGH | Ensure use of API Gateway endpoint policy, and no action wildcards are being used. | AWS.APGRAP.IAM.HIGH.0064 | AC_AWS_0440 |


