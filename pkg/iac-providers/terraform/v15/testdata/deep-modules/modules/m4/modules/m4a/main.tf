variable "m4aprojectid" {
    type = string
    default = "asdfasdf"
}

module "m4b" {
    source = "../m4b"
    m4bversionyear = "2012" 
    m4bversionmonth = "10" 
    m4bversionday = "17" 
    m4bbucketname = module.m4c.fullbucketname
}
module "m4c" {
    source = "../m4c"
    m4cbucketname = var.m4aprojectid
    m4cenvironment = "dev"
}


resource "aws_s3_bucket" "bucket4a" {
  bucket = module.m4c.fullbucketname
  policy = module.m4b.fullbucketpolicy
}
