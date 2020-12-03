
### aws_iam_role_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamPolicy | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AWS.IamPolicy.IAM.High.0392 |


### aws_route53_record
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| DNS Management | Route53HostedZone | HIGH | Route53HostedZone should have recordSets. | AWS.Route53HostedZone.DNSManagement.High.0422 |


### aws_api_gateway_method_settings
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | API Gateway | MEDIUM | Enable Detailed CloudWatch Metrics for APIs | AWS.API Gateway.Logging.Medium.0569 |


### aws_vpc
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | VPC | MEDIUM | Avoid creating resources in default VPC | AWS.VPC.Logging.Medium.0471 |
| Logging | VPC | MEDIUM | Ensure VPC flow logging is enabled in all VPCs | AWS.VPC.Logging.Medium.0470 |


### aws_iam_account_password_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| IAM | Iam | MEDIUM | Lower case alphabet not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0454 |
| IAM | Iam | MEDIUM | Setting a lengthy password increases account resiliency against brute force login attempts | AWS.Iam.IAM.Medium.0458 |
| IAM | Iam | LOW | It is recommended that the password policy prevent the reuse of passwords.Preventing password reuse increases account resiliency against brute force login attempts | AWS.Iam.IAM.Low.0539 |
| IAM | Iam | MEDIUM | Number not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0455 |
| IAM | Iam | MEDIUM | Setting a lengthy password increases account resiliency against brute force login attempts | AWS.Iam.IAM.Medium.0495 |
| IAM | Iam | MEDIUM | Special symbols not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0456 |
| IAM | Iam | MEDIUM | Upper case alphabet not present in the Password, Password Complexity is not high. Increased Password complexity increases resiliency against brute force attack | AWS.Iam.IAM.Medium.0457 |
| IAM | Iam | LOW | Reducing the password lifetime increases account resiliency against brute force login attempts | AWS.Iam.IAM.Low.0540 |


### aws_mq_broker
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | ElasticSearch | MEDIUM | Publicly Accessible MQ Brokers | AWS.ElasticSearch.NetworkSecurity.Medium.0887 |
| Logging | ElasticSearch | MEDIUM | Enable AWS MQ Log Exports | AWS.ElasticSearch.Logging.Medium.0885 |


### aws_db_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | AWS RDS | HIGH | RDS Instance publicly_accessible flag is true | AWS.AWS RDS.NS.High.0101 |
| Data Security | RDS | HIGH | Ensure Certificate used in RDS instance is updated | AWS.RDS.DS.High.1042 |
| Data Security | RDS | HIGH | Ensure that your RDS database has IAM Authentication enabled. | AWS.RDS.DataSecurity.High.0577 |
| Data Security | RDS | HIGH | RDS Instance Auto Minor Version Upgrade flag disabled | AWS.RDS.DS.High.1041 |
| Data Security | RDS | HIGH | Ensure that your RDS database instances have automated backups enabled for point-in-time recovery. To back up your database instances, AWS RDS take automatically a full daily snapshot of your data (with transactions logs) during the specified backup window and keeps the backups for a limited period of time (known as retention period) defined by the instance owner. | AWS.RDS.DataSecurity.High.0414 |


### aws_ebs_volume
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | EcsCluster | HIGH | Ensure that AWS EBS clusters are encrypted. Data encryption at rest, prevents unauthorized users from accessing sensitive data on your AWS EBS clusters and associated cache storage systems. | AWS.EcsCluster.EncryptionandKeyManagement.High.0413 |
| Encryption and Key Management | EBS | HIGH | Enable AWS EBS Snapshot Encryption | AWS.EBS.EKM.Medium.0682 |


### aws_api_gateway_rest_api
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
|   | APIGateway | MEDIUM | Enable Content Encoding | AWS.APIGateway.Medium.0568 |
| Network Security | APIGateway | MEDIUM | API Gateway Private Endpoints | AWS.APIGateway.Network Security.Medium.0570 |


### aws_iam_role
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamPolicy | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AWS.IamPolicy.IAM.High.0392 |


### aws_ebs_encryption_by_default
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Data Security | EBS | HIGH | Ensure that the AWS EBS that hold sensitive and critical data is encrypted by default to fulfill compliance requirements for data-at-rest encryption. | AWS.EBS.DataSecurity.High.0580 |


### aws_sns_topic
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | SNS | MEDIUM | Ensure SNS Topic is Publicly Accessible For Subscription | AWS.SNS.NS.Medium.1044 |


### aws_apigatewayv2_api
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| AccessControl | ApiGatewayV2Api | Medium | Insecure Cross-Origin Resource Sharing Configuration allowing all domains | AWS.ApiGatewayV2Api.AccessControl.0630 |


### aws_efs_file_system
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | EFS | HIGH | Enable encryption of your EFS file systems in order to protect your data and metadata from breaches or unauthorized access and fulfill compliance requirements for data-at-rest encryption within your organization. | AWS.EFS.EncryptionandKeyManagement.High.0409 |
| Encryption and Key Management | EFS | HIGH | Enable encryption of your EFS file systems in order to protect your data and metadata from breaches or unauthorized access and fulfill compliance requirements for data-at-rest encryption within your organization. | AWS.EFS.EncryptionandKeyManagement.High.0410 |


### aws_sqs_queue
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | SQS | HIGH | Identify any publicly accessible SQS queues available in your AWS account and update their permissions in order to protect against unauthorized users. | AWS.SQS.NetworkSecurity.High.0569 |
| Network Security | SQS | HIGH | Ensure that your Amazon Simple Queue Service (SQS) queues are protecting the contents of their messages using Server-Side Encryption (SSE). The SQS service uses an AWS KMS Customer Master Key (CMK) to generate data keys required for the encryption/decryption process of SQS messages. There is no additional charge for using SQS Server-Side Encryption, however, there is a charge for using AWS KMS | AWS.SQS.NetworkSecurity.High.0570 |


### aws_instance
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | json | MEDIUM | EC2 instances should disable IMDS or require IMDSv2 | AC-AWS-NS-IN-M-1172 |
| Network Security | Instance | MEDIUM | Instance should be configured in vpc. AWS VPCs provides the controls to facilitate a formal process for approving and testing all network connections and changes to the firewall and router configurations. | AWS.Instance.NetworkSecurity.Medium.0506 |


### aws_config
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & Key Management | Config | MEDIUM | Ensure AWS Config Rule is enabled for Encrypted Volumes | AWS.Config.Encryption&KeyManagement.Medium.0660 |


### aws_cloudformation_stack
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
|   | CloudFormation | MEDIUM | AWS CloudFormation Not In Use | AWS.CloudFormation.Medium.0599 |
|   | CloudFormation | MEDIUM | AWS CloudFormation Stack Policy | AWS.CloudFormation.Medium.0604 |
|   | CloudFormation | MEDIUM | Enable AWS CloudFormation Stack Termination Protection | AWS.CloudFormation.Medium.0605 |
|   | CloudFormation | MEDIUM | Enable AWS CloudFormation Stack Notifications | AWS.CloudFormation.Medium.0603 |


### aws_iam_user_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamPolicy | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AWS.IamPolicy.IAM.High.0392 |
| Identity and Access Management | IamUser | HIGH | Ensure Hardware MFA device is enabled for the "root" account | AWS.IamUser.IAM.High.0387 |
| Identity and Access Management | IamUser | HIGH | Ensure Virtual MFA device is enabled for the "root" account | AWS.IamUser.IAM.High.0388 |
| Identity and Access Management | IamUser | HIGH | It is recommended that MFA be enabled for all accounts that have a console password. Enabling MFA provides increased security for console access as it requires the authenticating principal to possess a device that emits a time-sensitive key and have knowledge of a credential | AWS.IamUser.IAM.High.0389 |


### aws_ecs_task_definition
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | EcsCluster | HIGH | Like any other EC2 instance it is recommended to place ECS instance within a VPC. AWS VPCs provides the controls to facilitate a formal process for approving and testing all network connections and changes to the firewall and router configurations | AWS.EcsCluster.NetworkSecurity.High.0104 |
| Data Security | LaunchConfiguration | HIGH | Sensitive Information Disclosure | AWS.LaunchConfiguration.DataSecurity.High.0101 |


### aws_ecr_repository_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Data Security | ECR | HIGH | Identify any exposed Amazon ECR image repositories available within your AWS account and update their permissions in order to protect against unauthorized access. Amazon Elastic Container Registry (ECR) is a managed Docker registry service that makes it easy for DevOps teams to store, manage and deploy Docker container images. An ECR repository is a collection of Docker images available on AWS cloud. | AWS.ECR.DataSecurity.High.0579 |


### aws_iam_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamPolicy | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AWS.IamPolicy.IAM.High.0392 |


### aws_apigatewayv2_stage
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | ApiGatewayV2Stage | Low | AWS API Gateway V2 Stage is missing access logs | AWS.ApiGatewayV2Stage.Logging.0630 |


### aws_ecr_repository
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Data Security | ECR | MEDIUM | Unscanned images may contain vulnerabilities | AWS.ECR.DataSecurity.High.0578 |


### aws_cloudfront_distribution
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | CloudFront | HIGH | Secure ciphers are not used in CloudFront distribution | AWS.CloudFront.EncryptionandKeyManagement.High.0408 |
| Network Security | CloudFront | LOW | Ensure that geo restriction is enabled for your Amazon CloudFront CDN distribution to whitelist or blacklist a country in order to allow or restrict users in specific locations from accessing web application content. | AWS.CloudFront.Network Security.Low.0568 |
| Encryption and Key Management | CloudFront | HIGH | Use encrypted connection between CloudFront and origin server | AWS.CloudFront.EncryptionandKeyManagement.High.0407 |
| Logging | CloudFront | MEDIUM | Ensure that your AWS Cloudfront distributions have the Logging feature enabled in order to track all viewer requests for the content delivered through the Content Delivery Network (CDN). | AWS.CloudFront.Logging.Medium.0567 |


### aws_cloudwatch
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | CloudWatch | MEDIUM | App-Tier CloudWatch Log Group Retention Period | AWS.CloudWatch.Logging.Medium.0631 |
| Encryption and Key Management | CloudWatch | HIGH | AWS CloudWatch log group is not encrypted with a KMS CMK | AWS.CloudWatch.EncryptionandKeyManagement.High.0632 |


### aws_ami_launch_permission
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | AMI | MEDIUM | Limit access to AWS AMIs | AWS.AMI.NS.Medium.1040 |


### aws_launch_configuration
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | EcsCluster | HIGH | Ensure that AWS ECS clusters are encrypted. Data encryption at rest, prevents unauthorized users from accessing sensitive data on your AWS ECS clusters and associated cache storage systems. | AWS.EcsCluster.EncryptionandKeyManagement.High.0413 |
| Data Security | LaunchConfiguration | HIGH | Avoid using base64 encoded private keys as part of config | AWS.LaunchConfiguration.DataSecurity.High.0102 |
| Data Security | LaunchConfiguration | HIGH | Avoid using base64 encoded shell script as part of config | AWS.LaunchConfiguration.DataSecurity.High.0101 |


### aws_api_gateway_stage
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | API Gateway | MEDIUM | Enable SSL Client Certificate | AWS.API Gateway.Network Security.Medium.0565 |
| Logging | API Gateway | MEDIUM | Ensure that AWS CloudWatch logs are enabled for all your APIs created with Amazon API Gateway service in order to track and analyze execution behavior at the API stage level. | AWS.API Gateway.Logging.Medium.0572 |
| Logging | API Gateway | MEDIUM | Enable Active Tracing | AWS.API Gateway.Logging.Medium.0571 |
| Logging | API Gateway | MEDIUM | Enable AWS CloudWatch Logs for APIs | AWS.API Gateway.Logging.Medium.0567 |


### aws_elasticsearch_domain
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | ElasticSearch | MEDIUM | Enable AWS ElasticSearch Encryption At Rest | AWS.ElasticSearch.EKM.Medium.0778 |
| Encryption and Key Management | ElasticSearch | MEDIUM | ElasticSearch Domain Encrypted with KMS CMKs | AWS.ElasticSearch.EKM.Medium.0768 |
| Logging | Elasticsearch | MEDIUM | Ensure that your AWS Elasticsearch clusters have enabled the support for publishing slow logs to AWS CloudWatch Logs. This feature enables you to publish slow logs from the indexing and search operations performed on your ES clusters and gain full insight into the performance of these operations. | AWS.Elasticsearch.Logging.Medium.0573 |


### aws_iam_user_login_profile
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | Iam | HIGH | Password policies are, in part, used to enforce password complexity requirements. IAM password policies can be used to ensure password are comprised of different character sets, have minimal length, rotation and history restrictions | AWS.Iam.IAM.High.0391 |


### aws_iam_group_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamPolicy | HIGH | It is recommended and considered a standard security advice to grant least privileges that is, granting only the permissions required to perform a task. IAM policies are the means by which privileges are granted to users, groups, or roles. Determine what users need to do and then craft policies for them that let the users perform only those tasks, instead of granting full administrative privileges. | AWS.IamPolicy.IAM.High.0392 |


### aws_load_balancer_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | ELB | HIGH | Using insecure ciphers for your ELB Predefined or Custom Security Policy, could make the SSL connection between the client and the load balancer vulnerable to exploits. TLS 1.0 was recommended to be disabled by PCI Council after June 30, 2016 | AWS.ELB.EncryptionandKeyManagement.High.0401 |
| Encryption and Key Management | ELB | HIGH | Remove insecure ciphers for your ELB Predefined or Custom Security Policy, to reduce the risk of the SSL connection between the client and the load balancer being exploited. | AWS.ELB.EncryptionandKeyManagement.High.0403 |


### aws_s3_bucket
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| IAM | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0377 |
| Network Security | S3Bucket | HIGH | Ensure that there are not any static websites being hosted on buckets you aren't aware of | AWS.S3Bucket.NetworkSecurity.High.0417 |
| IAM | S3Bucket | HIGH | Enabling S3 versioning will enable easy recovery from both unintended user actions, like deletes and overwrites | AWS.S3Bucket.IAM.High.0370 |
| S3 | S3Bucket | HIGH | S3 bucket Access is allowed to all AWS Account Users. | AWS.S3Bucket.DS.High.1043 |
| IAM | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0379 |
| Encryption and Key Management | S3Bucket | HIGH | Ensure that S3 Buckets have server side encryption at rest enabled with KMS key to protect sensitive data. | AWS.S3Bucket.EncryptionandKeyManagement.High.0405 |
| IAM | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0378 |
| IAM | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0381 |


### aws_elb
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Ports Security | ELB | LOW | AWS ELB incoming traffic not encrypted | AWS.ELB.NetworkPortsSecurity.Low.0563 |


### aws_redshift_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | Redshift | HIGH | Ensure Redshift clusters are not publicly accessible to minimise security risks. | AWS.Redshift.NetworkSecurity.HIGH.0564 |
| Logging | Redshift | MEDIUM | Ensure audit logging is enabled for Redshift clusters for security and troubleshooting purposes. | AWS.Redshift.Logging.Medium.0565 |
| Encryption and Key Management | Redshift | HIGH | Use customer-managed KMS keys instead of AWS-managed keys, to have granular control over encrypting and encrypting data. Encrypt Redshift clusters with a Customer-managed KMS key. This is a recommended best practice. | AWS.Redshift.EncryptionandKeyManagement.High.0415 |


### aws_kinesis_stream
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | Kinesis | HIGH | Kinesis Streams and metadata are not protected | AWS.Kinesis.EncryptionandKeyManagement.High.0412 |


### aws_config_configuration_aggregator
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | Config | HIGH | Ensure AWS Config is enabled in all regions | AWS.Config.Logging.HIGH.0590 |


### aws_organizations_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| IAM | Organizations | MEDIUM | Ensure that All Features is enabled within your Amazon Organizations to achieve full control over the use of AWS services and actions across multiple AWS accounts using Service Control Policies (SCPs). | AWS.Organizations.IAM.MEDIUM.0590 |


### aws_route53_query_log
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | Route53 query logs | MEDIUM | Ensure CloudWatch logging is enabled for Route53 hosted zones. | AWS.Route53 query logs.Logging.Medium.0574 |


### aws_iam_access_key
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | IamUser | HIGH | The root account is the most privileged user in an AWS account. AWS Access Keys provide programmatic access to a given AWS account. It is recommended that all access keys associated with the root account be removed. Removing access keys associated with the root account limits vectors by which the account can be compromised. Additionally, removing the root access keys encourages the creation and use of role based accounts that are least privileged. | AWS.IamUser.IAM.High.0390 |
| Identity and Access Management | IamUser | MEDIUM | Ensure that there are no exposed Amazon IAM access keys in order to protect your AWS resources against unapproved access | AWS.IamUser.IAM.High.0391 |


### aws_guardduty_detector
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | GuardDuty Enabled | MEDIUM | Ensure that Amazon GuardDuty service is currently enabled in all regions in order to protect your AWS environment and infrastructure (AWS accounts and resources, IAM credentials, guest operating systems, applications, etc) against security threats. AWS GuardDuty is a managed threat detection service that continuously monitors your VPC flow logs, AWS CloudTrail event logs and DNS logs for malicious or unauthorized behavior. The service monitors for activity such as unusual API calls, potentially compromised EC2 instances or potentially unauthorized deployments that indicate a possible AWS account compromise. AWS GuardDuty operates entirely on Amazon Web Services infrastructure and does not affect the performance or reliability of your applications. The service does not require any software agents, sensors or network appliances. | AWS.GuardDuty Enabled.Security.Medium.0575 |


### aws_db_security_group
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Security | RDS | HIGH | RDS should not be defined with public interface. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0101 |
| Network Security | RDS | HIGH | RDS should not be open to a large scope. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0103 |
| Network Security | RDS | HIGH | RDS should not be open to a public scope. Firewall and router configurations should be used to restrict connections between untrusted networks and any system components in the cloud environment. | AWS.RDS.NetworkSecurity.High.0102 |


### aws_s3_bucket_policy
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0371 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0376 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0375 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0374 |
| Identity and Access Management | S3Bucket | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.S3Bucket.IAM.High.0372 |
| Identity and Access Management | IamPolicy | HIGH | Misconfigured S3 buckets can leak private information to the entire internet or allow unauthorized data tampering / deletion | AWS.IamPolicy.IAM.High.0373 |


### aws_ami
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption & KeyManagement | EC2 | MEDIUM | Enable AWS AMI Encryption | AWS.EC2.Encryption&KeyManagement.Medium.0688 |


### aws_elasticache_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Data Security | ElastiCache | HIGH | ElastiCache for Redis version is not compliant with AWS PCI DSS requirements | AWS.ElastiCache.DataSecurity.High.0425 |
| Data Security | ElastiCache | HIGH | ElastiCache for Memcached is not in use in AWS PCI DSS environments | AWS.ElastiCache.DataSecurity.High.0424 |
| High Availability | ElastiCache | MEDIUM | AWS ElastiCache Multi-AZ | AWS.ElastiCache.HighAvailability.Medium.0757 |


### aws_kinesis_firehose_delivery_stream
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | Kinesis | HIGH | AWS Kinesis Server data at rest has server side encryption (SSE) | AWS.Kinesis.EncryptionandKeyManagement.High.0411 |


### aws_rds_cluster
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Encryption and Key Management | RDS | HIGH | Encrypt Amazon RDS instances and snapshots at rest, by enabling the encryption option for your Amazon RDS DB instance | AWS.RDS.EncryptionandKeyManagement.High.0414 |


### aws_cloudtrail
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | CloudTrail | HIGH | Cloud Trail Log Not Enabled | AWS.CloudTrail.Logging.High.0399 |
| Logging | CloudTrail | MEDIUM | Cloud Trail Multi Region not enabled | AWS.CloudTrail.Logging.Medium.0460 |
| Logging | CloudTrail | MEDIUM | Ensure appropriate subscribers to each SNS topic | AWS.CloudTrail.Logging.Low.0559 |


### aws_lambda_function
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | LambdaFunction | Low | Lambda function doesn't not include a VPC configuration. | AWS.LambdaFunction.Logging.0472 |
| Logging | LambdaFunction | LOW | Lambda tracing is not enabled. | AWS.LambdaFunction.Logging.0470 |
| Encryption and Key Management | LambdaFunction | High | Lambda does not use KMS CMK key to protect environment variables. | AWS.LambdaFunction.EncryptionandKeyManagement.0471 |


### aws_kms_key
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Logging | KMS | HIGH | Ensure rotation for customer created CMKs is enabled | AWS.KMS.Logging.High.0400 |
| Network Security | KMS | HIGH | Identify any publicly accessible AWS Key Management Service master keys and update their access policy in order to stop any unsigned requests made to these resources. | AWS.KMS.NetworkSecurity.High.0566 |


### aws_security_group
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Network Ports Security | ALB | MEDIUM | 'MSSQL Debugger' (TCP:135) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0236 |
| Network Ports Security | ALB | MEDIUM | 'Cassandra OpsCenter agent port' (TCP:61621) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0220 |
| Network Ports Security | SecurityGroup | HIGH | remote desktop port open to internet | AWS.SecurityGroup.NetworkPortsSecurity.Low.0562 |
| Network Ports Security | ALB | MEDIUM | 'Memcached SSL' (TCP:11214) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0240 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Datagram Service' (TCP:138) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0256 |
| Network Ports Security | ALB | MEDIUM | 'SNMP' (UDP:161) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0276 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Session Service' (TCP:139) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0260 |
| Network Ports Security | ALB | MEDIUM | 'LDAP SSL ' (TCP:636) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0230 |
| Network Ports Security | ALB | MEDIUM | 'Known internal web port' (TCP:8000) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0226 |
| Network Ports Security | ALB | MEDIUM | 'Postgres SQL' (UDP:5432) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0270 |
| Network Ports Security | SecurityGroup | HIGH | It is recommended that no security group allows unrestricted ingress access | AWS.SecurityGroup.NPS.High.1045 |
| Network Ports Security | ALB | MEDIUM | 'Oracle DB SSL' (UDP:2484) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0266 |
| Network Ports Security | SecurityGroup | HIGH | A VPC comes with a default security group whose initial settings deny all inbound traffic, allow all outbound traffic, and allow all traffic between instances assigned to the security group. If you don't specify a security group when you launch an instance, the instance is automatically assigned to this default security group. Security groups provide stateful filtering of ingress/egress network traffic to AWS resources. It is recommended that the default security group restrict all traffic. Configuring the default security group to restrict all traffic will encourage least privilege security group development and mindful placement of AWS resource into security groups which will in-turn reduce the exposure of those resources. | AWS.SecurityGroup.NetworkSecurity.High.0097 |
| Network Ports Security | ALB | MEDIUM | 'Memcached SSL' (UDP:11215) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0246 |
| Network Ports Security | ALB | MEDIUM | 'SaltStack Master' (TCP:4505) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0196 |
| Network Ports Security | ALB | MEDIUM | 'MySQL' (TCP:3306) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0250 |
| Network Ports Security | ALB | MEDIUM | 'Known internal web port' (TCP:8080) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0228 |
| Network Ports Security | SecurityGroup | HIGH |  It is recommended that no security group allows unrestricted ingress access | AWS.SecurityGroup.NetworkSecurity.High.0094 |
| Network Ports Security | ALB | MEDIUM | 'Hadoop Name Node' (TCP:9000) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0224 |
| Network Ports Security | SecurityGroup | HIGH | Unknown Port is exposed to the entire internet | AWS.SecurityGroup.NPS.High.1046 |
| Network Ports Security | ALB | MEDIUM | 'MSSQL Admin' (TCP:1434) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0232 |
| Network Ports Security | ALB | MEDIUM | 'Mongo Web Portal' (TCP:27018) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0248 |
| Network Ports Security | ALB | MEDIUM | 'Oracle DB SSL' (TCP:2484) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0264 |
| Network Ports Security | ALB | MEDIUM | 'Prevalent known internal port' (TCP:3000) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0272 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Name Service' (TCP:137) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0252 |
| Network Ports Security | ALB | HIGH | 'SSH' (TCP:22) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0194 |
| Network Ports Security | ALB | MEDIUM | 'Memcached SSL' (UDP:11214) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0244 |
| Network Ports Security | ALB | MEDIUM | 'Postgres SQL' (TCP:5432) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0268 |
| Network Ports Security | ALB | MEDIUM | 'Cassandra' (TCP:7001) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0222 |
| Network Ports Security | ALB | MEDIUM | 'MSSQL Browser Service' (UDP:1434) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0234 |
| Network Ports Security | ALB | MEDIUM | 'CIFS / SMB' (TCP:3020) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0218 |
| Network Ports Security | ALB | MEDIUM | 'SQL Server Analysis Services' (TCP:2383) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0280 |
| Network Ports Security | ALB | MEDIUM | 'MSSQL Server' (TCP:1433) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0238 |
| Network Ports Security | SecurityGroup | HIGH | ssh port open to internet | AWS.SecurityGroup.NetworkPortsSecurity.Low.0560 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Name Service' (UDP:137) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0254 |
| Network Ports Security | SecurityGroup | HIGH | http port open to internet | AWS.SecurityGroup.NetworkPortsSecurity.Low.0561 |
| Network Ports Security | ALB | MEDIUM | 'Memcached SSL' (TCP:11215) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0242 |
| Network Ports Security | ALB | MEDIUM | 'SQL Server Analysis Service browser' (TCP:2382) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0278 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Datagram Service' (UDP:138) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0258 |
| Network Ports Security | ALB | MEDIUM | 'NetBIOS Session Service' (UDP:139) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0262 |
| Network Ports Security | ALB | MEDIUM | 'Puppet Master' (TCP:8140) is accessible by a CIDR block range | AWS.ALB.NetworkPortsSecurity.High.0274 |


### aws_ecs_service
| Category | Resource | Severity | Description | Reference ID |
| -------- | -------- | -------- | ----------- | ------------ |
| Identity and Access Management | ECS | HIGH | Ensure there are no ECS services Admin roles | AWS.ECS.High.0436 |


