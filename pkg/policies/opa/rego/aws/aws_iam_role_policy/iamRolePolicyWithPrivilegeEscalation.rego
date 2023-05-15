package accurics

{{.prefix}}iamRolePolicyWithPrivilegeEscalation[iamUser.id] {
    iamUser = input.aws_iam_role_policy[_]
    policy := json_unmarshal(iamUser.config.policy)
    statement = policy.Statement[_]
    policyCheck(statement) == true
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}

policyCheck(statement) = true {
    disallowed_actions := ["iam:passrole", "lambda:createfunction", "lambda:invokefunc*"]
    statement.Action[_] == disallowed_actions[_]   
}