resource "aws_sqs_queue" "sqsQueueExposed" {
  name                              = "terraform-example-queue"
  kms_master_key_id                 = "alias/aws/sqs"
  kms_data_key_reuse_period_seconds = 300
  policy                            = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Sid":"Queue1_AnonymousAccess_AllActions_WhitelistIP",
    "Effect": "Allow",
    "Principal": "*",
    "Action": "sqs:*",
    "Resource": "arn:aws:sqs:*:111122223333:queue1"
  }] 
}
EOF
}
