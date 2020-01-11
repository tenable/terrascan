resource "aws_instance" "foo" {
    value = "${var.missing}"
}