package accurics

{{.prefix}}{{.name}}[retVal] {
    password_policy := input.aws_iam_account_password_policy[_]
    password_policy.config.max_password_age > 90
    password_policy_id := password_policy.id
    traverse = "max_password_age"
	retVal := { "Id": password_policy.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "max_password_age", "AttributeDataType": "int", "Expected": 90, "Actual": password_policy.config.max_password_age }
}