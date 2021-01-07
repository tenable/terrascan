resource "aws_security_group" "acme_web" {
  name        = "acme_web"
  description = "Used in the terraform"
  vpc_id      = "some_dummy_vpc"

  tags = {
    Name = "acme_web"
  }

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "19.16.0.0/24"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }
}
