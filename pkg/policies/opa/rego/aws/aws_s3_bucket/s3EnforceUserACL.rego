package accurics

{{.prefix}}s3EnforceUserACL[bucket.id] {
    bucket := input.aws_s3_bucket[_]
    checkAcl(bucket.config)
    checkPolicy(input, bucket.config)
}

checkAcl(bucket_config) {
    lower(bucket_config.acl) != "private"
}

checkAcl(bucket_config) {
    object.get(bucket_config, "acl", "undefined") == [[], null, "undefined"][_]
}

checkPolicy(inobj, bucket_config) {
    object.get(bucket_config, "policy", "undefined") != [[], null, "undefined"][_]
    policy_object := json_unmarshal(bucket_config.policy)

    checkPrincipals(policy_object)
}

checkPolicy(inobj, bucket_config) {
    object.get(bucket_config, "policy", "undefined") == [[], null, "undefined"][_]

    bucket_policy := inobj.aws_s3_bucket_policy[_]
    object.get(bucket_policy.config, "policy", "undefined") != [[], null, "undefined"][_]
    policy_object := json_unmarshal(bucket_policy.policy)

    checkPrincipals(policy_object)
}

checkPrincipals(policy) {
    statement := policy.statement[_]
    principal := statement.principals[_]
    identifier := principal.identifiers[_]
    identifier == "*"
}

# remove all id related prefix and suffix characters generated by terrascan
getCleanID(id) = cleanID {
    v1 := trim_left(id, "$")
    v2 := trim_left(v1, "{")
    v3 := trim_right(v2, "}")
    cleanID = cleanEnd(v3)
}

cleanEnd(idv3) = cleanID {
    endswith(idv3, ".id")
    cleanID = trim_right(idv3, ".id")
}

cleanEnd(idv3) = cleanID {
    endswith(idv3, ".bucket")
    cleanID = trim_right(idv3, ".bucket")
}

json_unmarshal(s) = result {
	s == null
	result := json.unmarshal("{}")
}

json_unmarshal(s) = result {
	s != null
	result := json.unmarshal(s)
}