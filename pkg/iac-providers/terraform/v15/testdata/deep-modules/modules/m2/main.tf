variable "m2versionyear" {
    type = string
}
variable "m2versionmonth" {
    type = string
}
variable "m2versionday" {
    type = string
}
variable "m2bucketname" {
    type = string
}
data "aws_iam_policy_document" "readbuckets" {
    source_json = <<EOF
{
  "Version":"${var.m2versionyear}-${var.m2versionmonth}-${var.m2versionday}",
  "Statement":[
    {
      "Sid": "PublicRead",
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::${var.m2bucketname}/*"]
    }
  ]
}
EOF  
}

output "fullbucketpolicy" {
    value = data.aws_iam_policy_document.readbuckets.json
} 
output "BucketARN" {
    value = var.m2bucketname
} 

