package utils

import "strings"

// EnsureUpperCaseTrimmed make sure the string is in UPPERCASE and TRIMMED
func EnsureUpperCaseTrimmed(s string) string {
	return strings.TrimSpace(strings.ToUpper(s))
}
