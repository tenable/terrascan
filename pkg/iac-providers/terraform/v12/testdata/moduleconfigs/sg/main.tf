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

  # HTTP access from the VPC
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  # HTTPS access from the VPC
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    from_port   = 4505
    to_port     = 4505
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 4506
    from_port   = 4506
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 3020
    from_port   = 3020
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 61621
    from_port   = 61621
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 7001
    from_port   = 7001
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9000
    from_port   = 9000
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 8000
    from_port   = 8000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 8080
    from_port   = 8080
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 636
    from_port   = 636
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 1434
    from_port   = 1434
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 1434
    from_port   = 1434
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 135
    from_port   = 135
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 1433
    from_port   = 1433
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11214
    from_port   = 11214
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11214
    from_port   = 11214
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11215
    from_port   = 11215
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11215
    from_port   = 11215
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 27018
    from_port   = 27018
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 3306
    from_port   = 3306
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 137
    from_port   = 137
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 137
    from_port   = 137
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 138
    from_port   = 138
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 138
    from_port   = 138
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 139
    from_port   = 139
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 139
    from_port   = 139
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2484
    from_port   = 2484
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2484
    from_port   = 2484
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 5432
    from_port   = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 5432
    from_port   = 5432
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 3000
    from_port   = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 8140
    from_port   = 8140
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 161
    from_port   = 161
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2382
    from_port   = 2382
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2383
    from_port   = 2383
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9090
    from_port   = 9090
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 3389
    from_port   = 3389
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9042
    from_port   = 9042
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 7000
    from_port   = 7000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 7199
    from_port   = 7199
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 61620
    from_port   = 61620
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 8888
    from_port   = 8888
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9160
    from_port   = 9160
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9200
    from_port   = 9200
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 9300
    from_port   = 9300
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 389
    from_port   = 389
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 389
    from_port   = 389
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11211
    from_port   = 11211
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 11211
    from_port   = 11211
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 27017
    from_port   = 27017
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 1521
    from_port   = 1521
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2483
    from_port   = 2483
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 2483
    from_port   = 2483
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 6379
    from_port   = 6379
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 0
    from_port   = 6379
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  ingress {
    to_port     = 0
    from_port   = 4506
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0", "192.164.0.0/24"]
  }
}

resource "aws_security_group" "defaultSGNotRestrictsAllTraffic" {
  name        = "default"
  description = "Used in the terraform"
  vpc_id      = "some_dummy_vpc"

  tags = {
    Name = "default"
  }
}
