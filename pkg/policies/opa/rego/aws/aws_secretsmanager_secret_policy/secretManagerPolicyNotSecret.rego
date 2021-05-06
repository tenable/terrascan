package accurics

{{.prefix}}secretManagerPolicyNotSecret[secret_policy.id] {
    secret_policy := input.aws_secretsmanager_secret_policy[_]
    policy := json_unmarshal(secret_policy.config.policy)
    statement = policy.Statement[_]
    policyPrincipalCheck(statement) == true
}

{{.prefix}}secretManagerPolicyNotSecret[secret_policy.id] {
    secret_policy := input.aws_secretsmanager_secret_policy[_]
    policy := json_unmarshal(secret_policy.config.policy)
    statement = policy.Statement[_]
    policyActionCheck(statement) == true
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}

policyPrincipalCheck(statement) = true {
    statement.Principal == "*"
}

policyPrincipalCheck(statement) = true {
    statement.Principal.AWS == "*"
}

policyActionCheck(statement) = true {
    disallowed_actions := ["secretsmanager:*"]
    statement.Action == disallowed_actions[_]   
}