package accurics

{{.prefix}}rootUserNotContainMfaTypeVirtual[iamUserMfa_id] {
	iamUserMfa = input.aws_iam_user_policy[_]
    policy := json_unmarshal(iamUserMfa.config.policy)
    statement = policy.Statement[_]
    check_role(statement, "sts:AssumeRole", "Allow") == true
    root_check(iamUserMfa.config.user, "root") == true
    mfa_check(statement.Principal.AWS, ":mfa/") == true
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

root_check(s, v) = true {
    re_match(v, s)
}

mfa_check(s, v) = true {
    not re_match(v, s)
}

check_role(s, a, e) = true {
    s.Action == a
    s.Effect == e
}
