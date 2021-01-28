package accurics

imdsv1LaunchConfig[res.id] {
    res = input.aws_launch_configuration[_]
    not res.config.metadata_options
}

imdsv1LaunchConfig[res.id] {
  res = input.aws_launch_configuration[_]
  value := res.config.metadata_options[_]
  not value.http_endpoint == "disabled"
  not value.http_tokens == "required"
}