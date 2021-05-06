package accurics

{{.prefix}}apiGatewayRestApiPolicyNotSecure[rest_api.id] {
    rest_api := input.aws_api_gateway_rest_api_policy[_]
    checkRestApiPolicyAttached(rest_api.config.rest_api_id, input.aws_api_gateway_rest_api)

    policy := json_unmarshal(rest_api.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement) == true

}

checkRestApiPolicyAttached(api_id, api_name) = true {
    id_name := api_name[_].name
    attached_id := split(api_id, ".")
    id_name == attached_id[1]
}

policyCheck(statement) = true {
    disallowed_actions := ["execute-api:*"]
    act := disallowed_actions[_]
    act == statement.Action
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}