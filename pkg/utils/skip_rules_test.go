/*
    Copyright (C) 2020 Accurics, Inc.

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
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// ---------------------- unit tests -------------------------------- //
func TestGetSkipRules(t *testing.T) {
	testRuleAWS1 := "AWS.S3Bucket.DS.High.1041"
	testRuleAWS2 := "AWS.S3Bucket.DS.High.1042"
	testRuleAWSwithHyphen := "AC-AWS-NS-IN-M-1172"
	testRuleAzure := "accurics.azure.NS.147"
	testRuleKubernetesWithHyphen := "AC-K8-DS-PO-M-0143"

	table := []struct {
		name     string
		input    string
		expected []output.SkipRule
	}{
		{
			name:  "no rules",
			input: "no rules here",
			// expected would be an empty slice of output.SkipRule
		},
		{
			name:  "rule id with no comment, aws",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041\n",
			expected: []output.SkipRule{
				{Rule: testRuleAWS1},
			},
		},
		{
			name:  "rule id with no comment, aws, with '-'",
			input: "#ts:skip=AC-AWS-NS-IN-M-1172\n",
			expected: []output.SkipRule{
				{Rule: testRuleAWSwithHyphen},
			},
		},
		{
			// gcp, kubernetes, github rules are of same format
			name:  "rule id with no comment, azure",
			input: "#ts:skip=accurics.azure.NS.147\n",
			expected: []output.SkipRule{
				{Rule: testRuleAzure},
			},
		},
		{
			name:  "rule id with no comment, kubernetes with '-'",
			input: "#ts:skip=AC-K8-DS-PO-M-0143\n",
			expected: []output.SkipRule{
				{Rule: testRuleKubernetesWithHyphen},
			},
		},
		{
			name:  "rule id with comment",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041 This rule should be skipped.\n",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS1,
					Comment: "This rule should be skipped.",
				},
			},
		},
		{
			// should match only one rule, we support single rule and comment in one line
			// everything after the first group match will be considered a comment
			name:  "multiple comma separated no space, with comments",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041 some reason to skip. , AWS.S3Bucket.DS.High.1042 should_skip_the_rule.\n",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS1,
					Comment: "some reason to skip. , AWS.S3Bucket.DS.High.1042 should_skip_the_rule.",
				},
			},
		},
		{
			name: "rule and comment with random space characters",
			input: "#ts:skip=  AWS.S3Bucket.DS.High.1041  		reason_to skip. the rule\n",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS1,
					Comment: "reason_to skip. the rule",
				},
			},
		},
		{
			name: "sample resource config",
			input: `{
			#ts:skip=AWS.S3Bucket.DS.High.1041 skip the rule.
			region        = var.region
			#ts:skip=AWS.S3Bucket.DS.High.1042 AWS.S3Bucket.DS.High.1043
			bucket        = local.bucket_name
			#ts:skip=AWS.S3Bucket.DS.High.1044 resource skipped for this rule.
			force_destroy = true
			#ts:skip= AWS.S3Bucket.DS.High.1046
			acl           = "public-read"
			}`,
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS1,
					Comment: "skip the rule.",
				},
				{
					Rule:    testRuleAWS2,
					Comment: "AWS.S3Bucket.DS.High.1043",
				},
				{
					Rule:    "AWS.S3Bucket.DS.High.1044",
					Comment: "resource skipped for this rule.",
				},
				{
					Rule: "AWS.S3Bucket.DS.High.1046",
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetSkipRules(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("rule ids got: '%v', want: '%v'", actual, tt.expected)
			}
		})
	}
}

func TestReadSkipRulesFromMap(t *testing.T) {
	testRuleAWS1 := "AWS.CloudFormation.Medium.0603"
	testRuleK8s := "accurics.kubernetes.IAM.109"

	table := []struct {
		name     string
		input    map[string]interface{}
		expected []output.SkipRule
	}{
		{
			name:  "no rules",
			input: make(map[string]interface{}),
			// expected would be empty
		},
		{
			name:  "with valid aws rule",
			input: map[string]interface{}{TerrascanSkip: "[{\"rule\":\"AWS.CloudFormation.Medium.0603\"}]"},
			expected: []output.SkipRule{
				{Rule: testRuleAWS1},
			},
		},
		{
			name:  "with valid k8s rule",
			input: map[string]interface{}{TerrascanSkip: "[{\"rule\":\"accurics.kubernetes.IAM.109\"}]"},
			expected: []output.SkipRule{
				{Rule: testRuleK8s},
			},
		},
		{
			name:  "with invalid rule format",
			input: map[string]interface{}{TerrascanSkip: "[{\"rule\"\"accurics.kubernetes.IAM.109\"}]"},
			// expected would be empty
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			actual := ReadSkipRulesFromMap(tt.input, "testID")
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("rule ids got: '%v', want: '%v'", actual, tt.expected)
			}
		})
	}
}
