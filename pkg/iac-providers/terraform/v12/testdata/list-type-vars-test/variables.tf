variable "instance_count" {
  description = "Number of EC2 instances to deploy"
  type        = list(number)
  default = [1,2]
}

variable "instance_type" {
  description = "Type of EC2 instance to use"
  type        = string
  default = "blah_instance"
}

variable "subnet_ids" {
  description = "Subnet IDs for EC2 instances"
  type        = list(string)
  default = ["10.0.0.0/8"]
}

variable "security_group_ids" {
  description = "Security group IDs for EC2 instances"
  type        = list(string)
}

variable "tags" {
  description = "Tags for instances"
  type        = map
  default     = {
    "a" = "b"
    "c" = "d"
  }
}
