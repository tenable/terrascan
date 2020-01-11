variable "test_variable" {
    default = "1"
}

variable "test_variable_2" {
    description = "no default value"
}

resource "aws_instance" "foo" {
    value = "${var.test_variable}"
}

resource "aws_elb" "bar" {
    value = "${var.test_variable_2}"
}