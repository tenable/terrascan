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
		severity   string
	}
	tests := []struct {
		name          string
		fields        fields
		wantErr       bool
		assert        bool
		lenScanRules  int
		lenSkipRules  int
		lenCategories int
		severity      string
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
				severity:   "low",
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
			if err := e.initRuleSetFromConfigFile(); (err != nil) != tt.wantErr {
				t.Errorf("Executor.initRulesAndSeverity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assert {
				if len(e.scanRules) != tt.lenScanRules && len(e.skipRules) != tt.lenSkipRules && e.severity != tt.severity {
					t.Errorf("Executor.initRulesAndSeverity() expected scanRules: %d , skipRules: %d & severity : %s, got scanRules: %d , skipRules: %d and severity: %s", tt.lenScanRules, tt.lenSkipRules, tt.severity, len(e.scanRules), len(e.skipRules), e.severity)
				}
			}
		})
	}
}
