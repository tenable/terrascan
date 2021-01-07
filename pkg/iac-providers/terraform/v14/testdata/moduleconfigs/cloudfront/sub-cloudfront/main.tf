resource "aws_kms_key" "kmsKeyDisabled" {
  description = "KMS key 2"
  is_enabled  = false
  tags = {
    Name  = "kmsKeyDisabled"
    Setup = "self-healing"
  }
}

