package accurics

{{.name}}[api.id] {
    api := input.aws_cloudformation_stack[_]
    api.config.{{.property}} == null
}