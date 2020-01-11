variable "foo" {
    default = "bar"
}

resource "aws_instance" "foo" {
    value = "${var.missing}"
}