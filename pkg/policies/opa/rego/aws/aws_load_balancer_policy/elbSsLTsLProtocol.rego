package accurics

{{.prefix}}elbSsLTsLProtocol[retVal] {
    lb = input.aws_load_balancer_policy[_]
    some i
    policy := lb.config.policy_attribute[i]
    name := policy.name
    contains([{{range .security_protocols}}{{- printf "%q" . }},{{end}}], name)
    traverse := sprintf("policy_attribute[%d].name", [i])
    retVal := { "Id": lb.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "policy_attribute.name", "AttributeDataType": "string", "Expected": "Protocol-TLSv1.2", "Actual": name }

}

contains(security_protocol, nam) {
  security_protocol[_] = nam
}
