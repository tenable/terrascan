package accurics

{{.prefix}}iamPolicyWithFullAdminControl[retVal] {
	iamUserMfa = input.aws_iam_user_policy[_]
    policy := json_unmarshal(iamUserMfa.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement, "*", "Allow", "*") == true

    statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
    expected := object.union(policy, {"Statement": statements})

    traverse = "policy"
    retVal := { "Id": iamUserMfa.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}

replace_if_needed(statement) = value {
    policyCheck(statement, "*", "Allow", "*") == true
    actions := [ action | action := replace_action_if_needed( statement.Action[_] ) ]
    value := object.union(statement, { "Action": actions })
}

replace_if_needed(statement) = value {
	not policyCheck(statement, "*", "Allow", "*")
    value := statement
}

replace_action_if_needed(action) = value {
	action == "*"
    value := "##resource:action##"
}

replace_action_if_needed(action) = value {
	action != "*"
    value := action
}

policyCheck(s, a, e ,r) = true {
    s.Action[_] = a
    s.Effect == e
    s.Resource == r
}
