variable "m1projectid" {
    type = string
    default = "asdfasdf"
}

module "m2" {
    source = "../m2"
    m2versionyear = "2012" 
    m2versionmonth = "10" 
    m2versionday = "17" 
    m2bucketname = module.m3.fullbucketname
}
module "m3" {
    source = "../m3"
    m3bucketname = var.m1projectid
    m3environment = "dev"
}


resource "aws_s3_bucket" "bucket" {
  bucket = module.m3.fullbucketname
  policy = module.m2.fullbucketpolicy
}
