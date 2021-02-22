provider "aws" {
  region = "us-west-2"
}

resource "aws_ami" "awsAmiEncrypted" {
  name                = "some-name"

  ebs_block_device {
    device_name = "dev-name"
    encrypted = "false"
  }
}