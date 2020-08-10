package accurics

{{.prefix}}{{.name}}[retVal] {
    password_policy := input.aws_iam_account_password_policy[_]
    password_policy.config.{{.required_parameter}} == false
    password_policy_id := password_policy.id
    traverse = "{{.required_parameter}}"
	retVal := { "Id": password_policy.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "{{.required_parameter}}", "AttributeDataType": "bool", "Expected": true, "Actual": password_policy.config.{{.required_parameter}} }
}