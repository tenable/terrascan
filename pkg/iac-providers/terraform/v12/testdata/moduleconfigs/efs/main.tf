resource "aws_efs_file_system" "efsNotEncrypted" {
  creation_token = "my-product"

  tags = {
    Name = "not-encrypted"
  }
}


resource "aws_efs_file_system" "efsEncryptedFalse" {
  creation_token = "my-product"

  tags = {
    Name = "encrypted"
  }

  encrypted = false

}

resource "aws_efs_file_system" "efsEncryptedWithNoKms" {
  creation_token = "my-product"

  tags = {
    Name = "encrypted"
  }

  encrypted = true

}

