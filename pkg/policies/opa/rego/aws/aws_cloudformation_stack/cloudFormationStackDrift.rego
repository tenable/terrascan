package accurics

cloudFormationStackDrift[api.id]
{
    api := input.aws_config_config_rule[_]
    data := api.config.source[_]
    data.source_identifier == "CLOUDFORMATION_STACK_DRIFT_DETECTION_CHECK"
}