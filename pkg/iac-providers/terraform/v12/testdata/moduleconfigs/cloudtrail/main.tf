resource "aws_cloudtrail" "missing-multi-region" {
  name                          = "tf-trail-foobar"
  s3_bucket_name                = "some-s3-bucket"
  s3_key_prefix                 = "prefix"
  include_global_service_events = false
}
