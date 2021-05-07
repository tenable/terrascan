package accurics

{{.prefix}}kmsKeyPolicyMissingPrincipal[kms_policy.id] {
    kms_policy = input.aws_kms_key[_]
    policy := json_unmarshal(kms_policy.config.policy)
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
    object.get(statement, "Principal", "undefined") == "undefined"
}