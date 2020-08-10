provider "aws" {
  region = "us-east-1"
}

module "cloudfront" {
  source = "./cloudfront"
}
