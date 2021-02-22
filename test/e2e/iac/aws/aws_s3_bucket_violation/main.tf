provider "aws" {
  region = "us-west-2"
}

resource "aws_s3_bucket" "noS3BucketSseRules" {
  bucket = "mybucket"
  acl    = "private"

  tags = {
    Name        = "nos3BucketSseRules"
    Environment = "Dev"
  }
}


resource "aws_s3_bucket" "s3BucketSseRulesWithKmsNull" {
  bucket = "mybucket"
  acl    = "private"

  tags = {
    Name        = "s3BucketSseRulesWithNoKms"
    Environment = "Dev"
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "aws:kms"
      }
    }
  }
}

resource "aws_s3_bucket" "s3BucketNoWebsiteIndexDoc" {
  bucket = "website"
  acl    = "public-read"

   server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "some-key-id"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

resource "aws_s3_bucket" "s3VersioningMfaFalse" {
  bucket = "tf-test"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "some-key-id"
        sse_algorithm     = "aws:kms"
      }
    }
  }
  
  versioning {
    enabled = false
    mfa_delete = false
  }
}

resource "aws_s3_bucket" "allUsersReadAccess" {
  bucket = "my-tf-test-bucket"
  acl    = "public-read"
}

resource "aws_s3_bucket" "authUsersReadAccess" {
  bucket = "my-tf-test-bucket"
  acl    = "authenticated-read"
}

resource "aws_s3_bucket" "allUsersWriteAccess" {
  bucket = "my-tf-test-bucket"
  acl    = "public-read-write"
}

resource "aws_s3_bucket" "allUsersReadWriteAccess" {
  bucket = "my-tf-test-bucket"
  acl    = "public-read-write"
}

resource "aws_s3_bucket" "s3VersioningOnly" {
  bucket = "tf-test"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "some-key-id"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "s3MFADeleteOnly" {
  bucket = "tf-test"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "some-key-id"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning {
    mfa_delete = true
  }
}