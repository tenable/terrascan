package accurics


{{.prefix}}ecrRepoIsPublic[retVal] {
	repo = input.aws_ecr_repository_policy[_]
    policy := json_unmarshal(repo.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement, "Allow", "*") == true

    statements := [ item | item := replace_if_needed(policy.Statement[_]) ]
    expected := object.union(policy, {"Statement": statements})
    traverse = "policy"
    retVal := { "Id": repo.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "policy", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
	policyCheck(statement, "Allow", "*") == true
    # value := object.union(statement, { "Principal": { "AWS": "arn:aws:iam::##account_number##:root"} })
    value := statement
}

replace_if_needed(statement) = value {
	not policyCheck(statement, "Allow", "*")
    value := statement
}

policyCheck(s, e ,r) = true {
    s.Effect == e
    s.Principal == r
}