variable "m3bucketname" {
    type = string
}
variable "m3environment" {
    type = string
}

output "fullbucketname" {
    value = "${var.m3bucketname}-${var.m3environment}"
}
output "sourcebucketname" {
    value = var.m3bucketname
}
output "sourceenvironment" {
    value = var.m3environment
}

