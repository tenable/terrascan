provider "aws" {
  region = "us-west-2"
}

resource "aws_ebs_volume" "example" {
  availability_zone = "us-west-2a"
  size              = 8
}

resource "aws_ebs_snapshot" "example_snapshot" {
  volume_id = aws_ebs_volume.example.id
}
resource "aws_ami" "tesami" {
  name                = "ptshavaa1"
  virtualization_type = "hvm"
  root_device_name    = "/dev/xvda"

  ebs_block_device {
    device_name = "/dev/xvda"
    snapshot_id = aws_ebs_snapshot.example_snapshot.id
    volume_size = 8
  }
}
resource "aws_default_vpc" "main" {
  tags = {
    Name = "Default VPC"
  }
}

resource "aws_security_group" "allow_tls" {
  name        = "ptshavasg1"
  description = "Allow TLS inbound traffic"
  vpc_id      = aws_default_vpc.main.id

  ingress {
    description = "TLS from VPC"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [aws_default_vpc.main.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "instanceWithNoVpc" {
  ami           = aws_ami.tesami.id
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "disabled"
    http_tokens   = "required"
  }
}

resource "aws_instance" "instanceWithPublicIp" {
  #ts:skip=AC-AWS-NS-IN-M-1172 should skip this rule  
  ami           = aws_ami.tesami.id
  instance_type = "t2.micro"
  hibernation   = false

  associate_public_ip_address = true
}

resource "aws_instance" "instanceWithIMDv1_emptyblock" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {}
}

resource "aws_instance" "instanceWithIMDv1_fullblock" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "optional"
  }
}

resource "aws_instance" "instanceWithIMDv1_token_not_present" {
  #ts:skip=AC-AWS-NS-IN-M-1172 need to skip this rule
  #ts:skip=AC-AW-IS-IN-M-0144 can skip this rule 
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "enabled"
  }
}

resource "aws_instance" "instanceWithIMDv1_endpoint_not_present" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_tokens = "optional"
  }
}