variable "tags" {
  type = map
}

variable "name" { type = string }
variable "prefix" { type = string }

locals {
  common_tags = merge({
          CreatedBy = "Terraform"
        }, var.tags)
  name = var.prefix != "" ? "${var.prefix}-${var.name}" : var.name
}

resource "aws_iam_user" "lb" {
  name = local.name
  tags = merge(local.common_tags,
  {
    Name = local.name
  })
}
