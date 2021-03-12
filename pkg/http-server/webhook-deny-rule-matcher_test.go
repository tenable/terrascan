package httpserver

import (
	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/results"
	"testing"
)

func TestDenyRuleMatcher(t *testing.T) {
	testMediumSeverity := "MEDIUM"
	testCategory := "Identity and Access Management"
	testRuleName := "My Amazing Rule"

	table := []struct {
		name           string
		ruleSeverity   string
		ruleCategory   string
		ruleName       string
		k8sDenyRules   config.K8sDenyRules
		expectedResult bool
	}{
		{
			name:           "no deny rules",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			expectedResult: false,
		},
		{
			name:           "matched severity",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{DeniedSeverity: testMediumSeverity},
			expectedResult: true,
		},

		{
			name:           "lower severity",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{DeniedSeverity: "LOW"},
			expectedResult: true,
		},
		{
			name:           "higher severity",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{DeniedSeverity: "High"},
			expectedResult: false,
		},
		{
			name:           "not matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{Categories: []string{"WRONG!"}},
			expectedResult: false,
		},

		{
			name:           "matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{Categories: []string{"WRONG!", testCategory}},
			expectedResult: true,
		},
		{
			name:           "incorrect severity by matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sDenyRules{Categories: []string{"WRONG!", testCategory}, DeniedSeverity: "HIGH"},
			expectedResult: true,
		},
	}

	var denyRuleMatcher = webhookDenyRuleMatcher{}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			violation := results.Violation{
				RuleName: tt.ruleName,
				Severity: tt.ruleSeverity,
				Category: tt.ruleCategory,
			}

			result := denyRuleMatcher.match(violation, tt.k8sDenyRules)
			if result != tt.expectedResult {
				t.Errorf("Expected: %v, Got: %v", tt.expectedResult, result)
			}
		})
	}
}
