provider "aws" {
  region = var.aws_region
}


resource "aws_vpc" "vpc_playground" {
  cidr_block           = var.cidr_vpc
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Environment = "${var.environment_tag}"
  }
}

resource "aws_internet_gateway" "igw_playground" {
  vpc_id = aws_vpc.vpc_playground.id
  tags = {
    Environment = "${var.environment_tag}"
  }
}

resource "aws_subnet" "subnet_public_playground" {
  vpc_id                  = aws_vpc.vpc_playground.id
  cidr_block              = var.cidr_subnet
  map_public_ip_on_launch = "true"
  tags = {
    Environment = "${var.environment_tag}"
  }
}

resource "aws_route_table" "rtb_public_playground" {
  vpc_id = aws_vpc.vpc_playground.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw_playground.id
  }
  tags = {
    Environment = "${var.environment_tag}"
  }
}

resource "aws_route_table_association" "rta_subnet_public_playground" {
  subnet_id      = aws_subnet.subnet_public_playground.id
  route_table_id = aws_route_table.rtb_public_playground.id
}

resource "aws_security_group" "sg_playground" {
  name   = "sg"
  vpc_id = aws_vpc.vpc_playground.id
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = {
    Environment = "${var.environment_tag}"
  }
}

resource "aws_key_pair" "ec2key_playground" {
  key_name   = "testKey"
  public_key = file(var.public_key_path)
}

resource "aws_instance" "instance_playground" {
  ami                    = lookup(var.aws_amis, var.aws_region)
  instance_type          = var.instance_type
  subnet_id              = aws_subnet.subnet_public_playground.id
  vpc_security_group_ids = [aws_security_group.sg_playground.id]
  key_name               = aws_key_pair.ec2key_playground.key_name
  tags = {
    Environment = "${var.environment_tag}"
  }
  provisioner "remote-exec" {

    inline = [
      "sudo apt-get -y update",
      "sudo apt-get -y install nginx",
      "sudo touch /var/www/html/index.html",
      "echo HelloWorld | sudo tee -a /var/www/html/index.html",
      "sudo service nginx start",
    ]
    connection {
      host        = self.public_ip
      type        = "ssh"
      user        = "ubuntu"
      private_key = file(var.privateKey)
    }
  }
}
