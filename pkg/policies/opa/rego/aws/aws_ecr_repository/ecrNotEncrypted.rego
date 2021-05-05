package accurics

{{.prefix}}ecrNotEncrypted[ecr.id] {
    ecr := input.aws_ecr_repository[_]
    object.get(ecr.config, "encryption_configuration", "undefined") == ["undefined", []][_]
}

{{.prefix}}ecrNotEncrypted[ecr.id] {
    ecr := input.aws_ecr_repository[_]
    ecr.config.encryption_configuration[_] == {}
}