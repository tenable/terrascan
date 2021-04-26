variable bar {
  type = string
}

variable foo {
  type = string
}

locals {
  foo = lower(var.bar != null ? var.bar : var.foo)
}

resource "aws_iam_user" "lb" {
  name = local.foo
}