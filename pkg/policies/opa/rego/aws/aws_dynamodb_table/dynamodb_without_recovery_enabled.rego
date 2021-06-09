package accurics

{{.prefix}}dynamoderecovery_enabled[policy.id] {
    policy := input.aws_dynamodb_table[_]
    object.get(policy.config, "point_in_time_recovery", "undefined") == "undefined"
}

{{.prefix}}dynamoderecovery_enabled[policy.id] {
    policy := input.aws_dynamodb_table[_]
    pitr := policy.config.point_in_time_recovery[i]
    pitr.enabled == false
}