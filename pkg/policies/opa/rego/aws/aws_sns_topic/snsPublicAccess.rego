package accurics

{{.prefix}}snsPublicAccess[retVal] {
	sns := input.aws_sns_topic[_]
    policy := json_unmarshal(sns.config.policy)
    statement = policy.Statement[_]
    check_role(statement, "Allow", "*") == true

    statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
    expected := object.union( policy, {"Statement": statements} )
    traverse = "policy"
    retVal := { "Id": sns.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
	check_role(statement, "Allow", "*") == true
    value := object.union(statement, { "Principal": {"AWS": "##principal##"} })
}

replace_if_needed(statement) = value {
	not check_role(statement, "Allow", "*")
    value := statement
}

check_role(s, e, p) = true {
    s.Effect == e
    s.Principal.AWS == p
}