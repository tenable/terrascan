package utils

const (
	// HighSeverity high
	HighSeverity = "HIGH"
	// MediumSeverity medium
	MediumSeverity = "MEDIUM"
	// LowSeverity low
	LowSeverity = "LOW"
)

// ValidateSeverityInput validates input for --severity flag
func ValidateSeverityInput(severity string) bool {
	severity = EnsureUpperCaseTrimmed(severity)
	return severity == LowSeverity || severity == MediumSeverity || severity == HighSeverity
}

// CheckSeverity validates if the severity of policy rule is equal or above the desired severity
func CheckSeverity(ruleSeverity, desiredSeverity string) bool {
	ruleSeverity = EnsureUpperCaseTrimmed(ruleSeverity)
	desiredSeverity = EnsureUpperCaseTrimmed(desiredSeverity)

	if desiredSeverity == LowSeverity {
		return true
	}

	if desiredSeverity == MediumSeverity {
		return ruleSeverity == MediumSeverity || ruleSeverity == HighSeverity
	}

	return ruleSeverity == HighSeverity
}
