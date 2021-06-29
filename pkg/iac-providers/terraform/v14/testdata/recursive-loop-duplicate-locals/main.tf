locals {
  common_tags = {
    this = "that"
  }
}

module "dummy" {
  source = "./dummy"
  
  tags = local.common_tags  
  name   = "fred"
  prefix = "joe"
}



