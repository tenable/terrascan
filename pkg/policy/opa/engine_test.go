package opa

import (
	"fmt"
	"testing"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

var testPolicyPath string = "test"

type args struct {
	policyPath string
	scanRules  []string
	skipRules  []string
	categories []string
	severity   string
}

func TestFilterRules(t *testing.T) {

	tests := []struct {
		name        string
		args        args
		assert      bool
		regoMapSize int
		regoDataMap map[string]*policy.RegoData
	}{
		{
			name:        "no scan and skip rules",
			args:        args{},
			regoDataMap: nil,
		},
		{
			name: "scan rules test",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.0", "Rule.1", "Rule.2", "Rule.3", "Rule.10"},
			},
			regoDataMap: getTestRegoDataMap(10),
			assert:      true,
			regoMapSize: 4,
		},
		{
			name: "scan rules not found in path",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.11", "Rule.12", "Rule.13"},
			},
			regoDataMap: getTestRegoDataMap(10),
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "skip rules test",
			args: args{
				policyPath: testPolicyPath,
				skipRules:  []string{"Rule.1"},
			},
			regoDataMap: getTestRegoDataMap(6),
			assert:      true,
			regoMapSize: 5,
		},
		{
			name: "skip rules not found in policy path",
			args: args{
				policyPath: testPolicyPath,
				skipRules:  []string{"Rule.21", "Rule.22"},
			},
			regoDataMap: getTestRegoDataMap(20),
			assert:      true,
			regoMapSize: 20,
		},
		{
			name: "both scan and skip rules supplied",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
			},
			regoDataMap: getTestRegoDataMap(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : low and rule severity : blank",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "LOW",
			},
			regoDataMap: getTestRegoDataMap(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : high and rule severity : high",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "HIGH",
			},
			regoDataMap: getTestRegoDataMapHighSeverity(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : high and rule severity : high",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "HIGH",
			},
			regoDataMap: getTestRegoDataMapLowSeverity(50),
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : high and rule severity : medium",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "high",
			},
			regoDataMap: getTestRegoDataMapMediumSeverity(50),
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : medium and rule severity : low",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "MEDIUM",
			},
			regoDataMap: getTestRegoDataMapLowSeverity(50),
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : medium and rule severity : medium",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "MEDIUM",
			},
			regoDataMap: getTestRegoDataMapMediumSeverity(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired severity : low and rule severity : low",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				severity:   "low",
			},
			regoDataMap: getTestRegoDataMapLowSeverity(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired category : COMPLIANCE VALIDATION",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				categories: []string{"COMPLIANCE VALIDATION"},
			},
			regoDataMap: getTestRegoDataMapCVCategory(50),
			assert:      true,
			regoMapSize: 7,
		},
		{
			name: "both scan and skip rules supplied, with desired category : DATA PROTECTION",
			args: args{
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.6", "Rule.7", "Rule.8", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
				categories: []string{"DATA PROTECTION"},
			},
			regoDataMap: getTestRegoDataMapDPCategory(50),
			assert:      true,
			regoMapSize: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				regoDataMap: tt.regoDataMap,
			}
			e.FilterRules(tt.args.policyPath, tt.args.scanRules, tt.args.skipRules, tt.args.categories, tt.args.severity)
			if tt.assert {
				if len(e.regoDataMap) != tt.regoMapSize {
					t.Errorf("filterRules(): expected regoDataMap size = %d, got = %d", tt.regoMapSize, len(e.regoDataMap))
				}
			}
		})
	}
}

func TestFilterRulesWithEquallyDividedSeverity(t *testing.T) {
	tests := []struct {
		name        string
		args        args
		assert      bool
		regoMapSize int
		regoDataMap map[string]*policy.RegoData
	}{
		{
			name: "no scan and skip rules supplied, with desired severity: low and rules severities: low, medium, high",
			args: args{
				policyPath: testPolicyPath,
				severity:   "low",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM", "HIGH"}),
			assert:      true,
			regoMapSize: 30,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: medium and rules severities: low, medium, high",
			args: args{
				policyPath: testPolicyPath,
				severity:   "medium",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM", "HIGH"}),
			assert:      true,
			regoMapSize: 20,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: high and rules severities: low, medium, high",
			args: args{
				policyPath: testPolicyPath,
				severity:   "high",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM", "HIGH"}),
			assert:      true,
			regoMapSize: 10,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: high and rules severities: low, medium",
			args: args{
				policyPath: testPolicyPath,
				severity:   "high",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM"}),
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: high and rules severities: low, medium",
			args: args{
				policyPath: testPolicyPath,
				severity:   "MEDIUM",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM"}),
			assert:      true,
			regoMapSize: 15,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: high and rules severities: low, medium",
			args: args{
				policyPath: testPolicyPath,
				severity:   "LOW",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(30, []string{"LOW", "MEDIUM"}),
			assert:      true,
			regoMapSize: 30,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: low and rules severities: low",
			args: args{
				policyPath: testPolicyPath,
				severity:   "LOW",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(10, []string{"MEDIUM"}),
			assert:      true,
			regoMapSize: 10,
		},
		{
			name: "no scan and skip rules supplied, with desired severity : medium and rules severities: medium",
			args: args{
				policyPath: testPolicyPath,
				severity:   "MEDIUM",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(10, []string{"MEDIUM"}),
			assert:      true,
			regoMapSize: 10,
		},
		{
			name: "no scan and skip rules supplied, with desired severity: high and rules severities: medium",
			args: args{
				policyPath: testPolicyPath,
				severity:   "HIGH",
			},
			regoDataMap: getTestRegoDataMapWithEquallyDividedSeverity(10, []string{"MEDIUM"}),
			assert:      true,
			regoMapSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				regoDataMap: tt.regoDataMap,
			}
			e.FilterRules(tt.args.policyPath, tt.args.scanRules, tt.args.skipRules, tt.args.categories, tt.args.severity)
			if tt.assert {
				if len(e.regoDataMap) != tt.regoMapSize {
					t.Errorf("filterRules(): expected regoDataMap size = %d, got = %d", tt.regoMapSize, len(e.regoDataMap))
				}
			}
		})
	}
}

// helper func to generate test rego data map of given size
func getTestRegoDataMap(size int) map[string]*policy.RegoData {
	testRegoDataMap := make(map[string]*policy.RegoData)
	for i := 0; i < size; i++ {
		ruleID := fmt.Sprintf("Rule.%d", i)
		testRegoDataMap[ruleID] = &policy.RegoData{}
	}
	return testRegoDataMap
}

// helper func to generate test rego data map of given size with severity levels attached in equal proportion
func getTestRegoDataMapWithEquallyDividedSeverity(size int, severities []string) map[string]*policy.RegoData {

	severitytypes := len(severities)

	if severitytypes == 0 {
		return getTestRegoDataMap(size)
	}

	severitypartitionsizes := size / severitytypes
	severitycounter := 0
	testRegoDataMap := make(map[string]*policy.RegoData)

	for i := 0; i < size; i++ {
		ruleID := fmt.Sprintf("Rule.%d", i)
		testRegoDataMap[ruleID] = &policy.RegoData{Metadata: policy.RegoMetadata{Severity: severities[severitycounter]}}
		if (i+1)%severitypartitionsizes == 0 {
			severitycounter = severitycounter + 1
		}
	}

	return testRegoDataMap
}

// helper func to generate test rego data map of given size with required severity
func getTestRegoDataMapWithSeverity(size int, severity string) map[string]*policy.RegoData {
	testRegoDataMap := getTestRegoDataMap(size)

	for _, regoData := range testRegoDataMap {
		regoData.Metadata.Severity = severity
	}
	return testRegoDataMap
}

func getTestRegoDataMapWithCategory(size int, category string) map[string]*policy.RegoData {
	testRegoDataMap := getTestRegoDataMap(size)

	for _, regoData := range testRegoDataMap {
		regoData.Metadata.Category = category
	}
	return testRegoDataMap
}

// helper func to generate test rego data map of given size with high severity
func getTestRegoDataMapHighSeverity(size int) map[string]*policy.RegoData {
	return getTestRegoDataMapWithSeverity(size, utils.HighSeverity)
}

// helper func to generate test rego data map of given size with medium severity
func getTestRegoDataMapMediumSeverity(size int) map[string]*policy.RegoData {
	return getTestRegoDataMapWithSeverity(size, utils.MediumSeverity)
}

// helper func to generate test rego data map of given size with low severity
func getTestRegoDataMapLowSeverity(size int) map[string]*policy.RegoData {
	return getTestRegoDataMapWithSeverity(size, utils.LowSeverity)
}

func getTestRegoDataMapCVCategory(size int) map[string]*policy.RegoData {
	return getTestRegoDataMapWithCategory(size, utils.AcceptedCategories[1])
}

func getTestRegoDataMapDPCategory(size int) map[string]*policy.RegoData {
	return getTestRegoDataMapWithCategory(size, utils.AcceptedCategories[7])
}
