package accurics

{{.prefix}}route53LoggingDisabled[route.id] {
  route := input.aws_route53_zone[_]
  not input.aws_route53_query_log
}

{{.prefix}}route53LoggingDisabled[route.id] {
  route := input.aws_route53_query_log[_]
  logName := route.config.cloudwatch_log_group_arn
  not re_match(route.name, logName)
}