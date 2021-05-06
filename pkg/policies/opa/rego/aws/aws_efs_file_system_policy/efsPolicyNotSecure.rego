package accurics

{{.prefix}}efsPolicyNotSecure[efs_policy.id] {
	efs_policy := input.aws_efs_file_system_policy[_]
    policy := json_unmarshal(efs_policy.config.policy)
    statement = policy.Statement[_]
    checkStatementPrincipal(statement) == true
}

{{.prefix}}efsPolicyNotSecure[efs_policy.id] {
	efs_policy := input.aws_efs_file_system_policy[_]
    policy := json_unmarshal(efs_policy.config.policy)
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
    action == "elasticfilesystem:*"
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}