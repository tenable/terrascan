###
# Our default security group to access
# the instances over SSH and HTTP
resource "aws_security_group" "acme_web" {
  name        = "acme_web"
  description = "Used in the terraform"
  vpc_id      = "aws_vpc.acme_root.id"

  tags = {
    Name = "acme_web"
  }

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }


  # HTTP access from the VPC
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
}

resource "aws_instance" "acem_web" {
  # The connection block tells our provisioner how to
  # communicate with the resource (instance)
  # connection {
  #   # The default username for our AMI
  #   user = "ubuntu"
  #   host = "acme"
  #   # The connection will use the local SSH agent for authentication.
  # }

  tags = {
    Name = "acem_web"
  }

  instance_type = "t2.micro"

  # Lookup the correct AMI based on the region
  # we specified
  ami = "lookup(var.aws_amis, var.aws_region)"

  # The name of our SSH keypair we created above.
  key_name = "aws_key_pair.auth.id"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["aws_security_group.acme_web.id"]

  # We're going to launch into the same subnet as our ELB. In a production
  # environment it's more common to have a separate private subnet for
  # backend instances.
  subnet_id = "aws_subnet.acme_web.id"

  # We run a remote provisioner on the instance after creating it.
  # In this case, we just install nginx and start it. By default,
  # this should be on port 80
  # provisioner "remote-exec" {
  #   inline = [
  #     "sudo apt-get -y update",
  #     "sudo apt-get -y install nginx",
  #     "sudo service nginx start",
  #   ]
  # }
}
