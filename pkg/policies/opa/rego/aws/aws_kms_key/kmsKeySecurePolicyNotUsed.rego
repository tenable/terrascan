package accurics

{{.prefix}}kmsKeySecurePolicyNotUsed[kms_policy.id] {
    kms_policy = input.aws_kms_key[_]
    policy := json_unmarshal(kms_policy.config.policy)
    statement = policy.Statement[_]
    policyCheckPrincipal(statement) == true
}

{{.prefix}}kmsKeySecurePolicyNotUsed[kms_policy.id] {
    kms_policy = input.aws_kms_key[_]
    policy := json_unmarshal(kms_policy.config.policy)
    statement = policy.Statement[_]
    policyCheckAction(statement) == true
}

policyCheckAction(statement) = true {
    action := statement.Action[_]
    disallowed_actions := ["*", "kms:*"]
    action == disallowed_actions[_]
}

policyCheckPrincipal(statement) = true {
    statement.Principal == "*"
}

policyCheckPrincipal(statement) = true {
    statement.Principal.AWS == "*"
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}