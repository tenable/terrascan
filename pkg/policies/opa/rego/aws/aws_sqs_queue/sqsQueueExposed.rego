package accurics

{{.prefix}}sqsQueueExposed[retVal] {
	sqs := input.aws_sqs_queue[_]
    policy := json_unmarshal(sqs.config.policy)
    statement = policy.Statement[_]
    check_role(statement, "*") == true

    statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
    expected := object.union(policy, {"Statement": statements})
    traverse = "policy"
    retVal := { "Id": sqs.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
	check_role(statement, "*") == true
    value := object.union(statement, { "Principal": "##principal###" })
}

replace_if_needed(statement) = value {
	not check_role(statement, "*")
    value := statement
}

check_role(s, p) = true {
    s.Principal == p
}