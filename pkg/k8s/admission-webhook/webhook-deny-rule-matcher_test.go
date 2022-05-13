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

package admissionwebhook

import (
	"testing"

	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/results"
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
		k8sDenyRules   config.K8sAdmissionControl
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
			k8sDenyRules:   config.K8sAdmissionControl{DeniedSeverity: testMediumSeverity},
			expectedResult: true,
		},

		{
			name:           "lower severity",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sAdmissionControl{DeniedSeverity: "LOW"},
			expectedResult: true,
		},
		{
			name:           "higher severity",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sAdmissionControl{DeniedSeverity: "High"},
			expectedResult: false,
		},
		{
			name:           "not matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sAdmissionControl{Categories: []string{"WRONG!"}},
			expectedResult: false,
		},

		{
			name:           "matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sAdmissionControl{Categories: []string{"WRONG!", testCategory}},
			expectedResult: true,
		},
		{
			name:           "incorrect severity by matching category",
			ruleSeverity:   testMediumSeverity,
			ruleCategory:   testCategory,
			ruleName:       testRuleName,
			k8sDenyRules:   config.K8sAdmissionControl{Categories: []string{"WRONG!", testCategory}, DeniedSeverity: "HIGH"},
			expectedResult: true,
		},
	}

	var denyRuleMatcher = WebhookDenyRuleMatcher{}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {

			violation := results.Violation{
				RuleName: tt.ruleName,
				Severity: tt.ruleSeverity,
				Category: tt.ruleCategory,
			}

			result := denyRuleMatcher.Match(violation, tt.k8sDenyRules)
			if result != tt.expectedResult {
				t.Errorf("Expected: %v, Got: %v", tt.expectedResult, result)
			}
		})
	}
}
