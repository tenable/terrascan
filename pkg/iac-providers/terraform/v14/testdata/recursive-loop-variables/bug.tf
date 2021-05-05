variable "filename" {
  type = string
}

module "dummy" {
  source = "./dummy"

  filename = "${path.module}/${var.filename}"
}