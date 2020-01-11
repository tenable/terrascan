resource "aws_instance" "TEST_RESOURCE" {
    default = "1"
}

resource "aws_foo" "test_resource" {
    default = "1"
}

resource "aws_s3_bucket" "badResource" {
    default = "1"
}

