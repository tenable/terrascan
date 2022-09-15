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

package utils

import (
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

// ---------------------- unit tests -------------------------------- //
func TestGetSkipRules(t *testing.T) {
	testRuleAWS1 := "AWS.S3Bucket.DS.High.1041"
	testRuleAWS2 := "AWS.S3Bucket.DS.High.1042"
	testRuleAWS3 := "AWS.S3 Bucket.DS.High.1041"
	testRuleAWS4 := "AWS.S3 Bucket DS.High.1041"
	testRuleAWS5 := "AWS.S3 Bucket DS .High.1041"
	testRuleAWS6 := "AC_AWS_1111"
	testRuleAZURE1 := "AC_AZURE_1111"
	testRuleGCP1 := "AC_GCP_1111"
	testRuleK8S1 := "AC_K8S_1111"
	testRuleGITHUB1 := "AC_GITHUB_1111"
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
			name:  "rule and comment with random space characters",
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
		{
			// Rule with single space should get skipped
			name:  "rule with space in between, aws",
			input: "#ts:skip=AWS.S3 Bucket.DS.High.1041",
			expected: []output.SkipRule{
				{Rule: testRuleAWS3},
			},
		},
		{
			// Rule with two spaces should get skipped
			name:  "rule with two spaces in between, aws",
			input: "#ts:skip=AWS.S3 Bucket DS.High.1041",
			expected: []output.SkipRule{
				{Rule: testRuleAWS4},
			},
		},
		{
			// Rule with multiple spaces should get skipped
			name:  "rule with multiple spaces in between, aws",
			input: "#ts:skip=AWS.S3 Bucket DS .High.1041",
			expected: []output.SkipRule{
				{Rule: testRuleAWS5},
			},
		},
		{
			// Rule with space and comment should get skipped
			name:  "rule with spaces in between and comment, aws",
			input: "#ts:skip=AWS.S3 Bucket.DS.High.1041 skip rule with spaces",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS3,
					Comment: "skip rule with spaces",
				},
			},
		},
		{
			// Rule with multiple spaces and comment should get skipped
			name:  "rule with multiple spaces in between, aws",
			input: "#ts:skip=AWS.S3 Bucket DS .High.1041 skip rule with multiple spaces",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS5,
					Comment: "skip rule with multiple spaces",
				},
			},
		},
		{
			// skipping rule by ID field
			name:  "skipping rule using ID field",
			input: "#ts:skip=AC_AWS_1111",
			expected: []output.SkipRule{
				{Rule: testRuleAWS6},
			},
		},
		{
			// skipping rule by ID field and comment
			name:  "skipping rule using ID field and comment",
			input: "#ts:skip=AC_AWS_1111 skip rule by ID",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAWS6,
					Comment: "skip rule by ID",
				},
			},
		},
		{
			// skipping AZURE rule by ID field
			name:  "skipping AZURE rule using ID field",
			input: "#ts:skip=AC_AZURE_1111",
			expected: []output.SkipRule{
				{Rule: testRuleAZURE1},
			},
		},
		{
			// skipping AZURE rule by ID field and comment
			name:  "skipping AZURE rule using ID field and comment",
			input: "#ts:skip=AC_AZURE_1111 skip rule by ID",
			expected: []output.SkipRule{
				{
					Rule:    testRuleAZURE1,
					Comment: "skip rule by ID",
				},
			},
		},
		{
			// skipping GCP rule by ID field
			name:  "skipping GCP rule using ID field",
			input: "#ts:skip=AC_GCP_1111",
			expected: []output.SkipRule{
				{Rule: testRuleGCP1},
			},
		},
		{
			// skipping GCP rule by ID field and comment
			name:  "skipping GCP rule using ID field and comment",
			input: "#ts:skip=AC_GCP_1111 skip rule by ID",
			expected: []output.SkipRule{
				{
					Rule:    testRuleGCP1,
					Comment: "skip rule by ID",
				},
			},
		},
		{
			// skipping K8S rule by ID field
			name:  "skipping K8S rule using ID field ",
			input: "#ts:skip=AC_K8S_1111",
			expected: []output.SkipRule{
				{Rule: testRuleK8S1},
			},
		},
		{
			// skipping K8S rule by ID field and comment
			name:  "skipping K8S rule using ID field and comment",
			input: "#ts:skip=AC_K8S_1111 skip rule by ID",
			expected: []output.SkipRule{
				{
					Rule:    testRuleK8S1,
					Comment: "skip rule by ID",
				},
			},
		},
		{
			// skipping GITHUB rule by ID field
			name:  "skipping GITHUB rule using ID field ",
			input: "#ts:skip=AC_GITHUB_1111",
			expected: []output.SkipRule{
				{Rule: testRuleGITHUB1},
			},
		},
		{
			// skipping K8S rule by ID field and comment
			name:  "skipping GITHUB rule using ID field and comment",
			input: "#ts:skip=AC_GITHUB_1111 skip rule by ID",
			expected: []output.SkipRule{
				{
					Rule:    testRuleGITHUB1,
					Comment: "skip rule by ID",
				},
			},
		},
		{
			// skipping rule by ID field and comment
			name:  "skipping rule using ID field and comment repeated name",
			input: "#ts:skip=AC_AWS_1111AC_AWS_1111 skip rule by ID",
			expected: []output.SkipRule{
				{Rule: testRuleAWS6},
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
