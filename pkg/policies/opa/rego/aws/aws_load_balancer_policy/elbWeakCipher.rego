package accurics

{{.prefix}}elbWeakCipher[retVal] {
    lb = input.aws_load_balancer_policy[_]
    some i
    policy := lb.config.policy_attribute[i]
    name := policy.name
    contains([{{range .weak_ciphers}}{{- printf "%q" . }},{{end}}], name)
    traverse := sprintf("policy_attribute[%d].name", [i])
    retVal := { "Id": lb.id, "ReplaceType": "edit", "CodeType": "attribute", "Traverse": traverse, "Attribute": "policy_attribute.name", "AttributeDataType": "string", "Expected": "AES256-SHA256", "Actual": name }
}

contains(weak_cipher, nam) {
  weak_cipher[_] = nam
}