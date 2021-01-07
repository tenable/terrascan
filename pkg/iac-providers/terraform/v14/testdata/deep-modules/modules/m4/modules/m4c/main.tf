variable "m4cbucketname" {
    type = string
}
variable "m4cenvironment" {
    type = string
}

output "fullbucketname" {
    value = "${var.m4cbucketname}-${var.m4cenvironment}"
}
output "sourcebucketname" {
    value = var.m4cbucketname
}
output "sourceenvironment" {
    value = var.m4cenvironment
}

