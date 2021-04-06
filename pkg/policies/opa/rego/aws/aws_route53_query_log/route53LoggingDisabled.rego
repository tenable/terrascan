package accurics

{{.prefix}}route53LoggingDisabled[route.id] {
  route := input.aws_route53_zone[_]
  # From https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_query_log
  # There are restrictions on the configuration of query logging.
  # Notably, the CloudWatch log group must be in the us-east-1 region,
  # a permissive CloudWatch log resource policy must be in place,
  # and the Route53 hosted zone must be public.  <== NOTE
  # See Configuring Logging for DNS Queries for additional details.
  # https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/query-logs.html?console_help=true#query-logs-configuring

  # if it has a VPC associated, it's a private DNS zone, so this rule cannot apply because it would require
  # configuring an invalid logging resources (given the above)
  not vpc_exists
  not input.aws_route53_query_log
}

# Look for statically-defined VPC
vpc_exists = ns {
    ns := input.aws_route53_zone[_].config.vpc
}

# Also check for VPCs defined in dynamic blocks
vpc_exists = ns {
    ns := input.aws_route53_zone[_].config.dynamic[_].vpc
}

{{.prefix}}route53LoggingDisabled[route.id] {
  route := input.aws_route53_query_log[_]
  logName := route.config.cloudwatch_log_group_arn
  not re_match(route.name, logName)
}
