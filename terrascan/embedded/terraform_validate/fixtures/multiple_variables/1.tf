variable "test_variable" {
    default = "1"
}

variable "test_variable2" {
    default = "2"
}

resource "aws_instance" "foo" {
    value = "${var.test_variable}${var.test_variable2}"
    value_block = {
        value = "${var.test_variable2}${var.test_variable}"
    }
}