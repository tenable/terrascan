To Do
======
The following tests are planned to be implemented

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
