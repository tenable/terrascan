package accurics

ecsServiceAdmin[ecs.id] {
    ecs := input.aws_ecs_service[_]
    contains(ecs.config.iam_role, "admin")
}