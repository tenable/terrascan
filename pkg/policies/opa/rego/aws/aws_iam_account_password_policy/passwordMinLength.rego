package accurics

{{.prefix}}{{.name}}[retVal] {
    password_policy := input.aws_iam_account_password_policy[_]
    check_validity(password_policy.config, {{.value}}) == true
    traverse = "{{.parameter}}"
	retVal := { "Id": password_policy.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "{{.parameter}}", "AttributeDataType": "int", "Expected": {{.value}}, "Actual": password_policy.config.{{.parameter}} }
}

check_validity(p, v) = true {
    p.{{.parameter}} < v
}