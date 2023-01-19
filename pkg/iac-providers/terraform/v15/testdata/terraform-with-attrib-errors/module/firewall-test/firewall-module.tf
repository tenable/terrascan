resource "aws_route53_resolver_firewall_domain_list" "dlist" {
  name    = "some name"


  lifecycle {
    ignore_changes       = [domains]
    replace_triggered_by = [null_resource.someteam]
  }
}

resource "null_resource" "someteam" {
  triggers = {
    domains =  ["example.com"]
  }
}

resource "aws_route53_resolver_firewall_domain_list" "team_allow" {
  name    = "domain list"
  domains =  ["example.com"]

 

  lifecycle {
    ignore_changes       = [domains]
    replace_triggered_by = [null_resource.team_allow]
  }
}

resource "null_resource" "team_allow" {
  triggers = {
    domains = ["example.com"]
  }
}

resource "aws_vpc" "example" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}


resource "aws_route53_resolver_firewall_config" "firewall" {
  firewall_fail_open = "ENABLED"
  resource_id        = aws_vpc.example.id
}