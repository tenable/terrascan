package accurics

{{.prefix}}kmsKeyExposedPolicy[retVal] {
	kms = input.aws_kms_key[_]
    policy := json_unmarshal(kms.config.policy)
    statement = policy.Statement[_]
   	check_role(statement, "kms:*", "*", "*") == true

   	statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
   	expected := object.union(policy, {"Statement": statements})
    traverse = "policy"
   	retVal := { "Id": kms.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
    check_role(statement, "kms:*", "*", "*") == true
    actions := [ action | action := replace_action_if_needed( statement.Action[_] ) ]
    value := object.union(statement, { "Principal": "##principal##", "Action": actions })
}

replace_if_needed(statement) = value {
	not check_role(statement, "kms:*", "*", "*")
    value := statement
}

replace_action_if_needed(action) = value {
	action == "kms:*"
    value := "kms:##kms_action##"
}

replace_action_if_needed(action) = value {
	action != "kms:*"
    value := action
}

check_role(s, a, p, r) = true {
    s.Action[_] = a
    s.Principal.AWS == p
    s.Resource == r
}