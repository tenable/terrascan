package accurics

# this is still buggy, logic is still unstable

{{.prefix}}scpFullAccess[retVal]{
    org_policy = input.aws_organizations_policy[_]
    org_policy.config.type == "SERVICE_CONTROL_POLICY"
    content := json_unmarshal(org_policy.config.content)
    # policyCheck(content, "*", "Allow", "*") == true

    statements := [ content | content := replace_if_needed(content.Statement) ]
    expected := object.union(content, {"Statement": statements})
    traverse = "content"
    retVal := { "Id": org_policy.id, "ReplaceType": "edit", "CodeType": "document", "Traverse": traverse, "Attribute": "content", "AttributeDataType": "base64", "Expected": base64.encode(json.marshal(expected))}
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
    not policyCheck(statement, "*", "Allow", "*") == true
    value := object.union(statement, { "Resource": "*" })
}

replace_if_needed(statement) = value {
	policyCheck(statement, "*", "Allow", "*")
    value := statement
}

policyCheck(s, r, e, a) = true {
    s.Effect == e
	action := is_array(s.Action)
    s.Action[_] == a
	resource := is_array(s.Resource)
    s.Resource[_] == r
}

policyCheck(s, r, e, a) = true {
    s.Effect == e
	action := is_string(s.Action)
    s.Action == a
	resource := is_array(s.Resource)
    s.Resource[_] == r
}

policyCheck(s, r, e, a) = true {
    s.Effect == e
	action := is_array(s.Action)
    s.Action[_] == a
	resource := is_string(s.Resource)
    s.Resource == r
}

policyCheck(s, r, e, a) = true {
    s.Effect == e
	action := is_string(s.Action)
    s.Action == a
	resource := is_string(s.Resource)
    s.Resource == r
}