package accurics

ec2UsingIMDSv1[api.id] {
  api := input.aws_instance[_]
  not api.config.metadata_options
}

ec2UsingIMDSv1[api.id] {
  api := input.aws_instance[_]
  value := api.config.metadata_options[_]
  not value.http_endpoint == "disabled"
  not value.http_tokens == "required"
}