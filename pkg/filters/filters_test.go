package filters

import (
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/utils"
)

func TestRegoMetadataPreLoadFilterIsFiltered(t *testing.T) {
	testRuleID := "Rule.1"

	type fields struct {
		skipRules []string
	}
	type args struct {
		regoMetadata *policy.RegoMetadata
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no skip rules",
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
				},
			},
			want: false,
		},
		{
			name: "skip rules not matching with metadata reference id",
			fields: fields{
				skipRules: []string{"Rule.2", "Rule.3", "Rule.4"},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
				},
			},
			want: false,
		},
		{
			name: "skip rules contain a reference id matching with metadata reference id",
			fields: fields{
				skipRules: []string{"Rule.2", "Rule.3", testRuleID},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegoMetadataPreLoadFilter(nil, tt.fields.skipRules, nil, "")
			if got := r.IsFiltered(tt.args.regoMetadata); got != tt.want {
				t.Errorf("RegoMetadataPreLoadFilter.IsFiltered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegoMetadataPreLoadFilterIsAllowed(t *testing.T) {
	testRuleID := "Rule.1"
	testCategory := "Category.1"

	type fields struct {
		scanRules  []string
		categories []string
		severity   string
	}
	type args struct {
		regoMetadata *policy.RegoMetadata
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			// when no values are present, all regometadata are allowed
			name: "no scan rules, categories or severity specified",
			args: args{
				regoMetadata: &policy.RegoMetadata{},
			},
			want: true,
		},
		{
			name: "only scan rules specified, regometadata referecen id doesn't match",
			fields: fields{
				scanRules: []string{testRuleID},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: "Rule.2",
				},
			},
			want: false,
		},
		{
			name: "only scan rules specified, regometadata referecen id matches one of the scan rule id",
			fields: fields{
				scanRules: []string{testRuleID},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
				},
			},
			want: true,
		},
		{
			name: "only categories specified, regometadata category doesn't match",
			fields: fields{
				categories: []string{testCategory},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					Category: "Category.2",
				},
			},
			want: false,
		},
		{
			name: "only categories specified, regometadata category matches one of the category",
			fields: fields{
				categories: []string{testCategory},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					Category: testCategory,
				},
			},
			want: true,
		},
		{
			name: "only severity specified, regometadata severity doesn't match",
			fields: fields{
				severity: utils.HighSeverity,
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					Severity: utils.LowSeverity,
				},
			},
			want: false,
		},
		{
			name: "only severity specified, regometadata severity matches one of the severity specified",
			fields: fields{
				severity: utils.HighSeverity,
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					Severity: utils.HighSeverity,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegoMetadataPreLoadFilter(tt.fields.scanRules, nil, tt.fields.categories, tt.fields.severity)
			if got := r.IsAllowed(tt.args.regoMetadata); got != tt.want {
				t.Errorf("RegoMetadataPreLoadFilter.IsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegoDataFilter_Filter(t *testing.T) {
	testRegoDataMap := map[string]*policy.RegoData{
		"aws_s3_bucket":    {},
		"aws_ec2_instance": {},
		"kubernetes_pod":   {},
	}

	type args struct {
		rmap  map[string]*policy.RegoData
		input policy.EngineInput
	}
	tests := []struct {
		name string
		args args
		want map[string]*policy.RegoData
	}{
		{
			name: "config input doesn't have any resources",
			args: args{
				rmap: testRegoDataMap,
				input: policy.EngineInput{
					InputData: &output.AllResourceConfigs{},
				},
			},
			want: testRegoDataMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RegoDataFilter{}
			got := r.Filter(tt.args.rmap, tt.args.input)
			if len(got) != len(tt.want) {
				t.Errorf("RegoDataFilter.Filter() = got size of map %v, want size of map %v", len(got), len(tt.want))
			}
		})
	}
}
