package accurics

{{.prefix}}esDomainInsecurePolicy[es_domain.id] {
    es_domain := input.aws_elasticsearch_domain_policy[_]
    policy := json_unmarshal(es_domain.config.access_policies)
    statement = policy.Statement[_]
    checkStatementPrincipal(statement) == true
}

{{.prefix}}esDomainInsecurePolicy[es_domain.id] {
	es_domain := input.aws_elasticsearch_domain_policy[_]
    policy := json_unmarshal(es_domain.config.access_policies)
    statement = policy.Statement[_]
    checkDisallowedActions(statement) == true
}

checkStatementPrincipal(statement) {
    statement.Principal.AWS == "*"
}

checkStatementPrincipal(statement) {
    statement.Principal == "*"
}

checkDisallowedActions(statement) {
    action := statement.Action[_]
    disallowed_actions := ["es:*"]
    action == disallowed_actions[_]
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}