package accurics

{{.prefix}}userWithPassNotContainMfaActive[iamUserMfa_id] {
    iamUserMfa = input.aws_iam_user_policy[_]
    policy := json_unmarshal(iamUserMfa.config.policy)
    statement := policy.Statement[_]
    statement.Condition.Bool[_] = false
    iamUserMfa_id = iamUserMfa.id
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}
