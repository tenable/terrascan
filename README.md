terrascan
==========
A collection of security and best practice tests for static code analysis of (terraform)[https://www.terraform.io] templates using (terraform-validate)[https://github.com/elmundio87/terraform_validate].

Encryption
----------
- Verifies server side encription, KMS, and SSL certificates are configured as applicable on the following resources:
    - [x] alb_listener
    - [x] ami
    - [x] ami_copy
    - [x] api_gateway_domain_name
    - [x] aws_instance
    - [x] cloudfront_distribution
    - [x] cloudtrail
    - [x] codebuild_project
    - [x] codepipeline
    - [x] db_instance
    - [x] dms_endpoint
    - [x] dms_replication_instance
    - [x] ebs_volume
    - [x] efs_file_system
    - [x] elastictranscoder_pipeline
    - [x] elb
    - [x] kinesis_firehose_delivery_stream
    - [ ] lambda_function
    - [ ] lb_ssl_negotiation_policy
    - [ ] load_balancer_backend_server_policy
    - [ ] load_balancer_listener_policy
    - [ ] load_balancer_policy
    - [ ] opsworks_application
    - [ ] proxy_protocol_policy
    - [ ] rds_cluster
    - [ ] rds_cluster_instance
    - [ ] redshift_cluster
    - [ ] redshift_parameter_group
    - [ ] s3_bucket_object
    - [ ] ses_receipt_rule
    - [ ] sqs_queue
    - [ ] ssm_parameter

Identity and access management
------------------------------
Checks for overly permissive permissions and bad practices.
Verifies that:
- For each of these types of policies that there are no NotActions:
    - IAM policy
    - IAM role trust relationship
    - S3 bucket policy
    - SNS topic policy
    - SQS queue policy
    - KMS policy
- For each of these types of policies that there are no NotPrincipals:
    - IAM role trust relationship
    - S3 bucket policy
    - SNS topic policy
    - SQS queue policy
    - KMS policy
- For each of these types of policies that there are no wildcard actions:
    - IAM policy
    - IAM role trust relationship
    - S3 bucket policy
    - SQS queue policy
    - KMS policy
- For each of these types of policies that there are no wildcard principals:
    - Lambda permission
    - S3 bucket policy
    - SNS topic policy
    - SQS queue policy
    - KMS policy
- No policies attached to IAM users
- No inline policies on:
    - IAM users
    - IAM roles
- S3 bucket no public-read ACL
- S3 bucket no public-read-write ACL
- S3 bucket no authenticated-read ACL
- The AWS administrator managed policy shouldn't be attached to any resources
- AWS Managed policies can't be scanned
- No creation of IAM API keys


Security Groups
---------------
Checks security groups rules for overly permissive configuration and bad practices.
Verifies that:
- Security group ingress rules are not:
     - Open to 0.0.0.0/0 on ports other than 443 or 22
     - Contain IP addresses outside of RFC1918 IP space


Logging and monitoring
----------------------
Checks if access logs and monitoring are enabled/configured
Verifies that:
- Access logs are enabled on the following resources:
    - CloudFront
    - ELB/ALB/NLB
    - S3

Public exposure
---------------
Checks if any resource is going to be publicly exposed without authentication.
- Verifies that the following resources are not publicly exposed:
    - EBS Snapshots
    - AMIs
    - Public IPs attached to EC2 instances
    - Private ELBs
    - RDS DBs
    - RDS Snapshots
    - Redshift
    - S3 website

Governance best practices
-------------------------
Checks against general governance best practices.
Verifies that:
- A specified number of tags are applied to all resources when supported.
- Autoscaling lifecycle actions are enabled to reduce uneccessary cost on unused resources
- There are no EC2 instance types provisioned for which AWS doesn't allow penetration testing: m1.small, t1.micro, or t2.nano
- There are no RDS instance types provisiones for which AWS doesn't allow penetration testing: small, micro
- Only approved AMIs are provisioned
- No S3 bucket names larger than 63 characters
- There are no hardcoded credentials in terraform templates
