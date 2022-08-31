/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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
	// this case will never arise, added for safe check
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

	return false
}

// MaxSeverityApplicable verifies if the severity of policy rule need to be changed to the maximum severity level
func MaxSeverityApplicable(ruleSeverity, maxSeverity string) bool {
	// this case will never arise, added for safe check
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

	return false
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
