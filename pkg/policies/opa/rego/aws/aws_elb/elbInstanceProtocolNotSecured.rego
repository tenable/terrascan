package accurics

{{.prefix}}elbInstanceProtocolNotSecured[retVal] {
    elb = input.aws_elb[_]
    some i
    listener = elb.config.listener[i]
	listener.instance_protocol != "https"
    traverse := sprintf("listener[%d].instance_protocol", [i])
    retVal := { "Id": elb.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "listener.instance_protocol", "AttributeDataType": "string", "Expected": "https", "Actual": listener.instance_protocol }
}