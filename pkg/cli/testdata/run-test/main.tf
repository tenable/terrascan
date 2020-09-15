# Specify the provider and access details
provider "aws" {
   region = "var.aws_region"
}

provider "kubernetes" {
}

# Create a VPC to launch our instances into
resource "aws_vpc" "acme_root" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "acme_root"
  }
}

# Create an internet gateway to give our subnet access to the outside world
resource "aws_internet_gateway" "acme_root" {
  vpc_id = "aws_vpc.acme_root.id"
  tags = {
    Name = "acme_root"
  }
}

# Grant the VPC internet access on its main route table
resource "aws_route" "acme_root" {
  route_table_id         = "aws_vpc.acme_root.main_route_table_id"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "aws_internet_gateway.acme_root.id"
}

# Create a subnet to launch our instances into
resource "aws_subnet" "acme_web" {
  vpc_id                  = "aws_vpc.acme_root.id"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name = "acme_web"
  }
}

resource "aws_key_pair" "auth" {
  key_name   = "var.key_name"
  public_key = "file(var.public_key_path)"
}

# resource "aws_s3_bucket" "acme_main" {
#   bucket = "main-bucket"
#   acl    = "private"
# }

resource "aws_ecr_repository" "scanOnPushDisabled" {
  name                 = "test"

  image_scanning_configuration {
    scan_on_push = false
  }
}

resource "aws_ecr_repository_policy" "ecrRepoIsPublic" {
  repository = "some-Repo-Name"

  policy = <<EOF
{
    "Version": "2008-10-17",
    "Statement": [
        {
            "Sid": "new policy",
            "Effect": "Allow",
            "Principal": "*",
            "Action": ["*"]
        }
    ]
}
EOF
}
