package opa

import (
	"fmt"
	"testing"
)

func TestFilterRules(t *testing.T) {
	testPolicyPath := "test"

	type args struct {
		e          *Engine
		policyPath string
		scanRules  []string
		skipRules  []string
	}
	tests := []struct {
		name        string
		args        args
		assert      bool
		regoMapSize int
	}{
		{
			name: "no scan and skip rules",
			args: args{
				e: &Engine{},
			},
		},
		{
			name: "scan rules test",
			args: args{
				e: &Engine{
					regoDataMap: getTestRegoDataMap(10),
				},
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.0", "Rule.1", "Rule.2", "Rule.3", "Rule.11"},
			},
			assert:      true,
			regoMapSize: 4,
		},
		{
			name: "scan rules not found in path",
			args: args{
				e: &Engine{
					regoDataMap: getTestRegoDataMap(10),
				},
				policyPath: testPolicyPath,
				scanRules:  []string{"Rule.11", "Rule.12", "Rule.13"},
			},
			assert:      true,
			regoMapSize: 0,
		},
		{
			name: "skip rules test",
			args: args{
				e: &Engine{
					regoDataMap: getTestRegoDataMap(6),
				},
				policyPath: testPolicyPath,
				skipRules:  []string{"Rule.1"},
			},
			assert:      true,
			regoMapSize: 5,
		},
		{
			name: "skip rules not found in policy path",
			args: args{
				e: &Engine{
					regoDataMap: getTestRegoDataMap(20),
				},
				policyPath: testPolicyPath,
				skipRules:  []string{"Rule.21", "Rule.22"},
			},
			assert:      true,
			regoMapSize: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterRules(tt.args.e, tt.args.policyPath, tt.args.scanRules, tt.args.skipRules)
			if tt.assert {
				if len(tt.args.e.regoDataMap) != tt.regoMapSize {
					t.Errorf("filterRules(): expected regoDataMap size = %d, got = %d", tt.regoMapSize, len(tt.args.e.regoDataMap))
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
