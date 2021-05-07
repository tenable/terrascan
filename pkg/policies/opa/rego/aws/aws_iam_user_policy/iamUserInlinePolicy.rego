package accurics

{{.prefix}}iamUserInlinePolicy[iam_user_policy.id] {
	iam_user_policy := input.aws_iam_user_policy[_]
    checkUserHasPolicyAttached(iam_user_policy.config.user, input.aws_iam_user)
}

{{.prefix}}iamUserInlinePolicy[iam_user_policy.id] {
	iam_user_policy := input.aws_iam_user_policy[_]
    checkUserHasPolicy(iam_user_policy.config.user, input.aws_iam_user)
}

checkUserHasPolicyAttached(policy_user, iam_user) = true {
    iam_username = iam_user[_].name
    split_val := split(policy_user, ".")
    policy_user_name =split_val[1]
    policy_user_name == iam_username
}

checkUserHasPolicy(policy_user, iam_user) = true {
    iam_username = iam_user[_].name
    policy_user == iam_username
}