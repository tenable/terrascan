locals {
  foo = "bar"
}

module "dummy" {
  source = "./dummy"
  bar = local.foo
  foo = "bar"
}