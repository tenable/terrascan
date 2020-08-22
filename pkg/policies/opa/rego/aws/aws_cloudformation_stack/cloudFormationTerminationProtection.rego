package accurics

cloudFormationTerminationProtection[api.id] {
    api := input.aws_cloudformation_stack_set_instance[_]
    api.config.retain_stack == false
}