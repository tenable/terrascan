package accurics

{{.prefix}}allowActionsFromAllPrincipals[retVal] {
	s3bucket = input.aws_s3_bucket_policy[_]
    policy := json_unmarshal(s3bucket.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement, "*", "Allow", "*") == true

    statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
    expected := object.union(policy, {"Statement": statements})
    traverse = "policy"
    retVal := { "Id": s3bucket.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
}

{{.prefix}}allowActionsFromAllPrincipals[retVal] {
	s3bucket = input.aws_s3_bucket[_]
    policy := json_unmarshal(s3bucket.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement, "*", "Allow", "*") == true

    statements := [ statement | statement := replace_if_needed(policy.Statement[_]) ]
    expected := object.union(policy, {"Statement": statements})
    traverse = "policy"
    retVal := { "Id": s3bucket.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
    value := object.union(statement, { "Principal": "##principal##", "Action": "##s3:action##" })
}

replace_if_needed(statement) = value {
	not policyCheck(statement, "*", "Allow", "*")
    value := statement
}

policyCheck(s, a, e ,p) = true {
    s.Action == a
    s.Effect == e
    s.Principal == p
}
