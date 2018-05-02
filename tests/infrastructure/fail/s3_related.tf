resource "aws_s3_bucket_object" "examplebucket_object" {
  server_side_encryption = "AES256"

  # kms_key_id             = "${var.kms_key_arn}"
}

resource "aws_s3_bucket" "public_read" {
  acl = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"

    routing_rules = <<EOF
[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
}]
EOF
  }
}

resource "aws_s3_bucket" "public_read_write" {
  acl = "public-read-write"
}

resource "aws_s3_bucket" "authenticated_read" {
  acl = "authenticated-read"
}

resource "aws_emr_cluster" "emr-test-cluster" {
  # log_uri = "s3bucket/test"
}
