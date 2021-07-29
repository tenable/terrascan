variable "m4bversionyear" {
    type = string
}
variable "m4bversionmonth" {
    type = string
}
variable "m4bversionday" {
    type = string
}
variable "m4bbucketname" {
    type = string
}
data "aws_iam_policy_document" "readbuckets" {
    source_json = <<EOF
{
  "Version":"${var.m4bversionyear}-${var.m4bversionmonth}-${var.m4bversionday}",
  "Statement":[
    {
      "Sid": "PublicRead",
      "Effect": "Allow",
      "Principal": "*",
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::${var.m4bbucketname}/*"]
    }
  ]
}
EOF  
}

output "fullbucketpolicy" {
    value = data.aws_iam_policy_document.readbuckets.json
} 
output "BucketARN" {
    value = var.m4bbucketname
} 

