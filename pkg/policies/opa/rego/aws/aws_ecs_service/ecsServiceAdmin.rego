package accurics

ecsServiceAdmin[data.id] {
    data := input.aws_ecs_service[_]
    contains(data.config.iam_role, "admin")
}