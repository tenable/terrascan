#file added for fix for #418

# source https://registry.terraform.io/

module "network" {
  source  = "Azure/network/azurerm"
  version = "3.2.1"
}

module "eks" {
  source = "terraform-aws-modules/eks/aws"
}

## contains local modules
module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "2.20.0"
}