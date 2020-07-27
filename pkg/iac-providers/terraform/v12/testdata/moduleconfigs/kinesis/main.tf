resource "aws_kinesis_stream" "kinesisEncryptedWithKms" {
  name             = "kinesisEncryptedWithKms"
  shard_count      = 1
  retention_period = 48

  shard_level_metrics = [
    "IncomingBytes",
    "OutgoingBytes",
  ]

  encryption_type = "KMS"
  kms_key_id      = "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"

  tags = {
    Environment = "kinesisEncryptedWithKms"
  }
}
