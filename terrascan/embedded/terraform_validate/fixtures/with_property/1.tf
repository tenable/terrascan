resource "aws_s3_bucket" "private_bucket" {
  acl = "private"

  policy = <<POLICY
{INVALID JSON {}
POLICY
}

resource "aws_s3_bucket" "public_bucket" {
  acl = "public"
}

resource "aws_s3_bucket" "tagged_bucket" {
  tags {
    Tag1      = "Tag1"
    CustomTag = "CustomValue"
    Tag2      = "Tag2"
  }

  policy = <<POLICY
{INVALID JSON {}
POLICY
}
