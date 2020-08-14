package accurics

{{.prefix}}elbLbProtocolNotSecured[retVal] {
    elb = input.aws_elb[_]
    some i
    listener = elb.config.listener[i]
	listener.lb_protocol != "https"
    traverse := sprintf("listener[%d].lb_protocol", [i])
    retVal := { "Id": elb.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "listener.lb_protocol", "AttributeDataType": "string", "Expected": "https", "Actual": listener.lb_protocol }
}