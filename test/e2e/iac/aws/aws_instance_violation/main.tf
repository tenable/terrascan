provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "instanceWithNoVpc" {
  ami           = "some-id"
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "disabled"
    http_tokens = "required"
  }

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "instanceWithPublicIp" {
  ami           = "some-id"
  instance_type = "t2.micro"

  associate_public_ip_address = true
  tags = {
    Name = "HelloWorld"
  } 
}

resource "aws_instance" "instanceWithIMDv1_emptyblock" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {}

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "instanceWithIMDv1_fullblock" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "optional"
  }

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "instanceWithIMDv1_token_not_present" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_endpoint = "enabled"
  }

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "instanceWithIMDv1_endpoint_not_present" {
  ami           = "ami-1234"
  instance_type = "t2.micro"

  metadata_options {
    http_tokens = "optional"
  }

  tags = {
    Name = "HelloWorld"
  }
}
