provider "aws" {
  region = "us-east-1"
}

module "cloudfront" {
  source = "./cloudfront"
}

module "cloudtrail" {
  source = "./cloudtrail"
}

module "ecs" {
  source = "./ecs"
}

module "efs" {
  source = "./efs"
}

module "elb" {
  source = "./elb"
}

module "guardduty" {
  source = "./guardduty"
}

module "iam" {
  source = "./iam"
}

module "kinesis" {
  source = "./kinesis"
}

module "s3" {
  source = "./s3"
}

module "sg" {
  source = "./sg"
}

module "sqs" {
  source = "./sqs"
}

module "elasticcache" {
  source = "../relative-moduleconfigs/elasticcache"
}
