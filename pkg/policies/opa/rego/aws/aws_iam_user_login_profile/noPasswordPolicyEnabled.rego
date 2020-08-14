package accurics

{{.prefix}}noPasswordPolicyEnabled[result.retVal] {
    policy := input.aws_iam_user_login_profile[_]
    result := checkPassword(policy.id, policy.config)
	result != null
}

checkPassword(id, c) = { "retVal": retVal } {
    c.password_length < 14
    c.password_reset_required == false
    traverse = ""
    rc := "ewogICJwYXNzd29yZF9sZW5ndGgiOiAxNiwKICAicGFzc3dvcmRfcmVzZXRfcmVxdWlyZWQiOiB0cnVlCn0="
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "block", "Traverse": traverse, "Attribute": "", "AttributeDataType": "block", "Expected": rc, "Actual": null }
}

checkPassword(id, c) = { "retVal": retVal } {
    c.password_length >= 14
    c.password_reset_required == false
    traverse = "password_reset_required"
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "password_reset_required", "AttributeDataType": "boolean", "Expected": true, "Actual": c.password_reset_required }
}

checkPassword(id, c) = { "retVal": retVal } {
    not c.password_length
    c.password_reset_required == false
    traverse = "password_reset_required"
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "password_reset_required", "AttributeDataType": "boolean", "Expected": true, "Actual": c.password_reset_required }
}

checkPassword(id, c) = { "retVal": retVal } {
    c.password_length < 14
    c.password_reset_required == true
    traverse = "password_length"
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "password_length", "AttributeDataType": "int", "Expected": 14, "Actual": c.password_length }
}

checkPassword(id, c) = { "retVal": retVal } {
    c.password_length < 14
    not c.password_reset_required
    traverse = "password_length"
	retVal := { "Id": id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "password_length", "AttributeDataType": "int", "Expected": 14, "Actual": c.password_length }
}