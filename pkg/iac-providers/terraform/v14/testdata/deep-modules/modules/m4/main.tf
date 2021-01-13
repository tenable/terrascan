variable "m4projectid" {
    type = string
    default = "asdfasdf"
}

module "m4a" {
    source = "./modules/m4a"
    m4aprojectid = var.m4projectid
}

resource "aws_s3_bucket" "bucket" {
  bucket = var.m4projectid
  policy = module.m4a.fullbucketpolicy
}
