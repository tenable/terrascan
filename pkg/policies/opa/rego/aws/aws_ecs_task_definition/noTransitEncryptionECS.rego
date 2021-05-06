package accurics

{{.prefix}}noTransitEncryptionECS[ecs.id]{
	ecs := input.aws_ecs_task_definition[_]
    efs := ecs.config.volume[_].efs_volume_configuration[_]
    object.get(efs, "transit_encryption", "undefined") == ["undefined", false, ""][_]
}