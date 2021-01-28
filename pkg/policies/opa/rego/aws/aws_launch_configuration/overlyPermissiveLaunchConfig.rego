
package accurics

overlyPermissiveLaunchConfig[res.id] {
    res = input.aws_launch_configuration[_]
    ins_profile_name := split(res.config.iam_instance_profile, ".")[1]
    iam_instance_profile := input.aws_iam_instance_profile[_]
    ins_profile_name == iam_instance_profile.name

    role_name := split(iam_instance_profile.config.role, ".")[1]
    role_policy_attachment := input.aws_iam_role_policy_attachment[_]
    role_name == split(role_policy_attachment.config.role, ".")[1]
    policy_name := split(role_policy_attachment.config.policy_arn, ".")[1]

    iam_policy := input.aws_iam_policy[_]
    policy_name == iam_policy.name
    policy := json_unmarshal(iam_policy.config.policy)
    statement = policy.Statement[_]
    ac := statement.Action[_]
    action := split(ac, ":")[0]
    policyCheck(statement, "*", "Allow", "*")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}

policyCheck(s, a, e ,r) {
    split(s.Action[_], ":")[1] == a
    s.Effect == e
    s.Resource == r
}

policyCheck(s, a, e ,r) {
    s.Action[_] == a
    s.Effect == e
    s.Resource == r
}
