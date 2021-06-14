resource "aws_s3_bucket" "noS3BucketSseRules" {
  bucket = "mybucket"
  acl    = "private"

  tags = {
    Name        = "nos3BucketSseRules"
    Environment = "Dev"
  }
}
