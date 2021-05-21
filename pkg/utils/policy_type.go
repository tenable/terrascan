package utils

func CheckPolicyType(rulePolicyType string, desiredPolicyTypes []string) bool {
	normDesiredPolicyTypes := make(map[string]bool, len(desiredPolicyTypes))
	normRulePolicyType := EnsureUpperCaseTrimmed(rulePolicyType)

	for _, desiredPolicyType := range desiredPolicyTypes {
		desiredPolicyType = EnsureUpperCaseTrimmed(desiredPolicyType)
		normDesiredPolicyTypes[desiredPolicyType] = true
	}

	if _, ok := normDesiredPolicyTypes["ALL"]; ok {
		return true
	}

	_, ok := normDesiredPolicyTypes[normRulePolicyType]
	return ok
}
