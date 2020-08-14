resource "aws_efs_file_system" "efsNotEncrypted" {
  creation_token = "my-product"

  tags = {
    Name = "not-encrypted"
  }
}
