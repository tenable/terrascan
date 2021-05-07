package accurics

{{.prefix}}globalAccFlowLogsDisabled[global_acc.id] {
    global_acc := input.aws_globalaccelerator_accelerator[_]
    object.get(global_acc.config, "attributes", "undefined") ==  "undefined"
}

{{.prefix}}globalAccFlowLogsDisabled[global_acc.id] {
    global_acc := input.aws_globalaccelerator_accelerator[_]
    global_acc.config.attributes == [{}]
}

{{.prefix}}globalAccFlowLogsDisabled[global_acc.id] {
    global_acc := input.aws_globalaccelerator_accelerator[_]
    global_acc.config.attributes == []
}

{{.prefix}}globalAccFlowLogsDisabled[global_acc.id] {
    global_acc := input.aws_globalaccelerator_accelerator[_]
    attributes := global_acc.config.attributes[_]
    attributes.flow_logs_enabled != true
}