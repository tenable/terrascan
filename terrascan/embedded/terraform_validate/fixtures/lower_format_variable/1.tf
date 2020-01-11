resource "aws_instance" "foo" {
    value = "${upper(lower(var.bar))}"
}

resource "aws_instance2" "foo2" {
    value = "${lower(var.bar)}${upper(var.bizz)}"
}


variable "bar" {
    default = "aBC"
}

variable "bizz" {
    default = "deF"
}