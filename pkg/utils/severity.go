package utils

import (
	"regexp"
	"strings"
)

const (
	// HighSeverity high
	HighSeverity = "HIGH"
	// MediumSeverity medium
	MediumSeverity = "MEDIUM"
	// LowSeverity low
	LowSeverity = "LOW"
)

var (
	maxSeverityPattern = regexp.MustCompile(`(#ts:maxseverity=(?i)None|#ts:maxseverity=(?i)High|#ts:maxseverity=(?i)Low|#ts:maxseverity=(?i)Medium)`)
	maxSeverityPrefix  = "#ts:maxseverity="
	minSeverityPattern = regexp.MustCompile(`(#ts:minseverity=(?i)High|#ts:minseverity=(?i)Low|#ts:minseverity=(?i)Medium)`)
	minSeverityPrefix  = "#ts:minseverity="
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

// MinSeverityApplicable verifies if the severity of policy rule need to be changed to the minimum severity level
func MinSeverityApplicable(ruleSeverity, minSeverity string) bool {
	if !ValidateSeverityInput(minSeverity) {
		return false
	}
	ruleSeverity = EnsureUpperCaseTrimmed(ruleSeverity)
	minSeverity = EnsureUpperCaseTrimmed(minSeverity)

	if minSeverity == HighSeverity {
		return ruleSeverity == MediumSeverity || ruleSeverity == LowSeverity
	}

	if minSeverity == MediumSeverity {
		return ruleSeverity == LowSeverity
	}

	return !(minSeverity == LowSeverity)

}

// MaxSeverityApplicable verifies if the severity of policy rule need to be changed to the maximum severity level
func MaxSeverityApplicable(ruleSeverity, maxSeverity string) bool {
	if !ValidateSeverityInput(maxSeverity) {
		return false
	}
	ruleSeverity = EnsureUpperCaseTrimmed(ruleSeverity)
	maxSeverity = EnsureUpperCaseTrimmed(maxSeverity)
	if maxSeverity == LowSeverity {
		return ruleSeverity == HighSeverity || ruleSeverity == MediumSeverity
	}

	if maxSeverity == MediumSeverity {
		return ruleSeverity == HighSeverity
	}

	return !(maxSeverity == HighSeverity)
}

// GetMinMaxSeverity returns the min and max severity to be applied to resources.
// can be set in terraform resource config with the following patterns
// #ts:minseverity = "High" --> any violation for this resource will be high
// #ts:maxseverity = "None" --> any violation for this resource will be ignored
// only one value will be considered
func GetMinMaxSeverity(body string) (minSeverity string, maxSeverity string) {
	if maxSeverityPattern.MatchString(body) {
		maxSeverityComment := maxSeverityPattern.FindString(body)
		maxSeverity = strings.TrimPrefix(maxSeverityComment, maxSeverityPrefix)
	}

	if minSeverityPattern.MatchString(body) {
		minSeverityComment := minSeverityPattern.FindString(body)
		minSeverity = strings.TrimPrefix(minSeverityComment, minSeverityPrefix)
	}
	return
}
