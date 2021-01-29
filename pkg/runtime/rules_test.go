package runtime

import (
	"testing"
)

func TestExecutorInitRulesAndSeverity(t *testing.T) {
	type fields struct {
		configFile string
		scanRules  []string
		skipRules  []string
		categories []string
	}
	tests := []struct {
		name          string
		fields        fields
		wantErr       bool
		assert        bool
		lenScanRules  int
		lenSkipRules  int
		lenCategories int
	}{
		{
			name:   "no config file",
			fields: fields{},
		},
		{
			name: "config file doesn't exist",
			fields: fields{
				configFile: "testdata/test.toml",
			},
			wantErr: true,
		},
		{
			name: "empty config file",
			fields: fields{
				configFile: "testdata/empty.toml",
			},
		},
		{
			name: "config file with empty rules",
			fields: fields{
				configFile: "testdata/webhook.toml",
			},
		},
		{
			name: "valid config file with scan and skip rules",
			fields: fields{
				configFile: "testdata/scan-skip-rules.toml",
				scanRules:  []string{"testRuleA", "testRuleB"},
				skipRules:  []string{"testRuleC"},
			},
			assert:       true,
			lenScanRules: 4,
			lenSkipRules: 5,
		},
		{
			name: "valid config file with scan and skip rules with low severity",
			fields: fields{
				configFile: "testdata/scan-skip-rules-low-severity.toml",
				scanRules:  []string{"testRuleA", "testRuleB"},
				skipRules:  []string{"testRuleC"},
				categories: []string{},
			},
			assert:       true,
			lenScanRules: 4,
			lenSkipRules: 5,
		},
		{
			name: "valid config file with invalid scan rules",
			fields: fields{
				configFile: "testdata/invalid-scan-skip-rules.toml",
			},
			wantErr: true,
		},
		{
			name: "valid config file with invalid skip rules",
			fields: fields{
				configFile: "testdata/invalid-skip-rules.toml",
			},
			wantErr: true,
		},
		{
			name: "valid config file with invalid severity",
			fields: fields{
				configFile: "testdata/invalid-severity.toml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Executor{
				configFile: tt.fields.configFile,
				scanRules:  tt.fields.scanRules,
				skipRules:  tt.fields.skipRules,
			}
			if err := e.initRulesAndCategories(); (err != nil) != tt.wantErr {
				t.Errorf("Executor.initRulesAndCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert {
				if len(e.scanRules) != tt.lenScanRules && len(e.skipRules) != tt.lenSkipRules && e.categories != tt.lenCategories {
					t.Errorf("Executor.initRulesAndCategories() expected scanRules: %d , skipRules: %d & categories : %d, got scanRules: %d , skipRules: %d and categories: %d", tt.lenScanRules, tt.lenSkipRules, tt.lenCategories, len(e.scanRules), len(e.skipRules), len(e.categories))
				}
			}
		})
	}
}
