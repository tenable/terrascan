package accurics

{{.prefix}}instanceNotInVpc[retVal] {
    instance := input.aws_ecs_task_definition[_]
    instance.config.network_mode != "awsvpc"
    traverse = "network_mode"
    retVal := { "Id": instance.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "network_mode", "AttributeDataType": "string", "Expected": "awsvpc", "Actual": instance.config.network_mode}
}