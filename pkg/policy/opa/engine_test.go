package opa

import (
	"fmt"
	"testing"
)

func TestFilterRules(t *testing.T) {
	testPolicyPath := "test"

	type args struct {
		policyPath string
		scanRules  []string
		skipRules  []string
	}
	tests := []struct {
		name        string
		args        args
		assert      bool
		regoMapSize int
		regoDataMap map[string]*RegoData
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
				scanRules:  []string{"Rule.0", "Rule.1", "Rule.2", "Rule.3", "Rule.11"},
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
				scanRules:  []string{"Rule.10", "Rule.11", "Rule.12", "Rule.15", "Rule.31", "Rule.32", "Rule.40", "Rule.41", "Rule.42"},
				skipRules:  []string{"Rule.31", "Rule.32", "Rule.38"},
			},
			regoDataMap: getTestRegoDataMap(50),
			assert:      true,
			regoMapSize: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{
				regoDataMap: tt.regoDataMap,
			}
			e.FilterRules(tt.args.policyPath, tt.args.scanRules, tt.args.skipRules)
			if tt.assert {
				if len(e.regoDataMap) != tt.regoMapSize {
					t.Errorf("filterRules(): expected regoDataMap size = %d, got = %d", tt.regoMapSize, len(e.regoDataMap))
				}
			}
		})
	}
}

// helper func to generate test rego data map of given size
func getTestRegoDataMap(size int) map[string]*RegoData {
	testRegoDataMap := make(map[string]*RegoData)
	for i := 0; i < size; i++ {
		ruleID := fmt.Sprintf("Rule.%d", i)
		testRegoDataMap[ruleID] = &RegoData{}
	}
	return testRegoDataMap
}
