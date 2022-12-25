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

package filters

import (
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/utils"
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
			r := NewRegoMetadataPreLoadFilter(nil, tt.fields.skipRules, nil, nil, "")
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
		scanRules   []string
		categories  []string
		policyTypes []string
		severity    string
	}
	type args struct {
		regoMetadata *policy.RegoMetadata
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          bool
		noFilterSpecs bool
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
			name: "only scan rules specified, regometadata reference id doesn't match",
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
			name: "only scan rules specified, regometadata reference id matches one of the scan rule id",
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
		{
			name: "only policyTypes specified, regometadata policy type doesn't match",
			fields: fields{
				policyTypes: []string{"k8s"},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					PolicyType: "aws",
				},
			},
			want: false,
		},
		{
			name: "only policyTypes specified, regometadata policy matches one of the policy specified",
			fields: fields{
				policyTypes: []string{"azure"},
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					PolicyType: "azure",
				},
			},
			want: true,
		},
		{
			name: "all fields specified, regometadata matches all the values specified",
			fields: fields{
				scanRules:   []string{testRuleID, "Rule.2"},
				categories:  []string{testCategory, "Category.2"},
				policyTypes: []string{"k8s", "aws"},
				severity:    utils.HighSeverity,
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
					Category:    testCategory,
					PolicyType:  "aws",
					Severity:    utils.HighSeverity,
				},
			},
			want: true,
		},
		{
			name: "all fields specified, regometadata doesn't match with one of the values specified",
			fields: fields{
				scanRules:   []string{testRuleID, "Rule.2"},
				categories:  []string{testCategory, "Category.2"},
				policyTypes: []string{"k8s", "aws"},
				severity:    utils.HighSeverity,
			},
			args: args{
				regoMetadata: &policy.RegoMetadata{
					ReferenceID: testRuleID,
					Category:    testCategory,
					PolicyType:  "gcp",
					Severity:    utils.HighSeverity,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegoMetadataPreLoadFilter(tt.fields.scanRules, nil, tt.fields.categories, tt.fields.policyTypes, tt.fields.severity)
			if got := r.IsAllowed(tt.args.regoMetadata); got != tt.want {
				t.Errorf("RegoMetadataPreLoadFilter.IsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegoDataFilter_Filter(t *testing.T) {
	testRegoDataMap := map[string]*policy.RegoData{
		"Rule.1": {},
		"Rule.2": {},
		"Rule.3": {},
	}

	testRegoDataMapWithResourceType := map[string]*policy.RegoData{
		"Rule.1": {
			Metadata: policy.RegoMetadata{
				ResourceType: "kubernetes_pod",
			},
		},
		"Rule.2": {
			Metadata: policy.RegoMetadata{
				ResourceType: "ec2_instance",
			},
		},
		"Rule.3": {
			Metadata: policy.RegoMetadata{
				ResourceType: "kubernetes_pod",
			},
		},
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
		{
			name: "config input has resources but regometadata doesn't have resource type set",
			args: args{
				rmap: testRegoDataMap,
				input: policy.EngineInput{
					InputData: &output.AllResourceConfigs{
						"pod": []output.ResourceConfig{},
					},
				},
			},
			want: testRegoDataMap,
		},
		{
			name: "config input has resources and there are policies matching the type",
			args: args{
				rmap: testRegoDataMapWithResourceType,
				input: policy.EngineInput{
					InputData: &output.AllResourceConfigs{
						"kubernetes_pod": []output.ResourceConfig{},
					},
				},
			},
			want: map[string]*policy.RegoData{
				"Rule.1": {
					Metadata: policy.RegoMetadata{
						ResourceType: "kubernetes_pod",
					},
				},
				"Rule.3": {
					Metadata: policy.RegoMetadata{
						ResourceType: "kubernetes_pod",
					},
				},
			},
		},
		{
			name: "config input has resources but there are no policies matching the type",
			args: args{
				rmap: testRegoDataMapWithResourceType,
				input: policy.EngineInput{
					InputData: &output.AllResourceConfigs{
						"kubernetes_deployment": []output.ResourceConfig{},
					},
				},
			},
			want: nil,
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
