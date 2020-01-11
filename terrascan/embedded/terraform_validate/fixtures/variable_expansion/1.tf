variable "bar" {
    default = "1"
}

resource "aws_instance" "foo" {
    value = "${var.bar}"
}