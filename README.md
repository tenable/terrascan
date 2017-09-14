terrascan
==========
A collection of security and best practice tests for static code analysis of (terraform)[https://www.terraform.io] templates using (terraform-validate)[https://github.com/elmundio87/terraform_validate].

Resourced based tests
----------------------
These tests are performed on terraform resources where applicable:
- **Encryption**
    - Server Side Encription (SSE) enabled
    - Use of AWS Key Management Service (KMS) with Customer Managed Keys (CMK)
    - Use of SSL/TLS and proper configuration
- **Security Groups**
    - Ingress rule not open to 0.0.0.0/0 on port other than 443 or 22
    - Ingress rules on non-standard ports outside of RFC1918 IP space
- **Logging & Monitoring**
    - Access logs enabled to resources that support it
- **Public Exposure**
    - Services with public exposure other than ELBs or Gateways (NAT, VGW, IGW)

Terraform resource | Encryption | Security Groups
------------------ | ---------- | ---------------
alb_listener | :heavy_check_mark: |
ami | :heavy_check_mark: |
ami_copy | :heavy_check_mark: |
api_gateway_domain_name | :heavy_check_mark: |
aws_db_security_group | | :heavy_minus_sign:
aws_instance | :heavy_check_mark: |
aws_redshift_security_group | | :heavy_minus_sign:
aws_security_group | | :heavy_minus_sign:
aws_security_group_rule | | :heavy_minus_sign:
cloudfront_distribution | :heavy_check_mark: |
cloudtrail | :heavy_check_mark: |
codebuild_project | :heavy_check_mark: |
codepipeline | :heavy_check_mark: |
db_instance | :heavy_check_mark: |
dms_endpoint | :heavy_check_mark: |
dms_replication_instance | :heavy_check_mark: |
ebs_volume | :heavy_check_mark: |
efs_file_system | :heavy_check_mark: |
elasticache_security_group | | :heavy_minus_sign:
elastictranscoder_pipeline | :heavy_check_mark: |
elb | :heavy_check_mark: |
kinesis_firehose_delivery_stream | :heavy_check_mark: |
lambda_function | :heavy_check_mark: |
lb_ssl_negotiation_policy | :heavy_minus_sign: |
load_balancer_backend_server_policy | :heavy_minus_sign: |
load_balancer_listener_policy | :heavy_minus_sign: |
load_balancer_policy | :heavy_minus_sign: |
opsworks_application | :heavy_check_mark: |
rds_cluster | :heavy_check_mark: |
redshift_cluster | :heavy_check_mark: |
redshift_parameter_group | :heavy_minus_sign: |
s3_bucket_object | :heavy_check_mark: |
ses_receipt_rule | :heavy_minus_sign: |
sqs_queue | :heavy_check_mark: |
ssm_parameter | :heavy_check_mark: |

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
