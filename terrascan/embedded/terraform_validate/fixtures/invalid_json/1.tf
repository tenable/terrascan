variable "bucket_name" {
  default = "defaultvalue"
}

resource "aws_s3_bucket" "invalidjson" {
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "enforceTLS",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::examplebucketname/*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      }
POLICY
}

resource "aws_s3_bucket" "validjsonwithvariable" {
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "enforceTLS",
      "Effect": "Deny",
      "Principal": "*",
      "Action": "s3:*",
      "Resource": "arn:aws:s3:::${var.bucket_name}/*",
      "Condition": {
        "Bool": {
          "aws:SecureTransport": "false"
        }
      }
    }]
}
POLICY
}

resource "aws_s2_bucket" "blankpolicy" {
  policy = ""
}

resource "aws_s2_bucket" "nopolicy" {}
