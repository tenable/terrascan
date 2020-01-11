resource "aws_instance" "foo" {
    value = "${unimplemented(var.bar))}"
}

variable "bar" {
    default = "aBC"
}
