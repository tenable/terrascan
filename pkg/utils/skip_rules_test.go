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
	testRule1 := "AWS.S3Bucket.DS.High.1041"
	testRule2 := "AWS.S3Bucket.DS.High.1042"

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
			name:  "single rule id with no comment",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041",
			expected: []output.SkipRule{
				{Rule: testRule1},
			},
		},
		{
			name:  "single rule id with comment",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041 | this rule should be skipped.",
			expected: []output.SkipRule{
				{
					Rule:    testRule1,
					Comment: "this rule should be skipped.",
				},
			},
		},
		{
			name:  "multiple comma separated no space, with comments",
			input: "#ts:skip=AWS.S3Bucket.DS.High.1041|some reason to skip. , AWS.S3Bucket.DS.High.1042| should_skip_the_rule.",
			expected: []output.SkipRule{
				{
					Rule:    testRule1,
					Comment: "some reason to skip.",
				},
				{
					Rule:    testRule2,
					Comment: "should_skip_the_rule.",
				},
			},
		},
		{
			name: "multiple comma separated random space, without comments",
			input: "#ts:skip=  AWS.S3Bucket.DS.High.1041   ,		AWS.S3Bucket.DS.High.1042",
			expected: []output.SkipRule{
				{
					Rule: testRule1,
				},
				{
					Rule: testRule2,
				},
			},
		},
		{
			name:  "sample resource config",
			input: "{\n     #ts:skip=AWS.S3Bucket.DS.High.1041 | skip the rule. \n   region        = var.region\n   #ts:skip=AWS.S3Bucket.DS.High.1042 AWS.S3Bucket.DS.High.1043\n   bucket        = local.bucket_name\n   #ts:skip=AWS.S3Bucket.DS.High.1044 | resource skipped for this rule. ,AWS.S3Bucket.DS.High.1045\n   force_destroy = true\n   #ts:skip= AWS.S3Bucket.DS.High.1046   ,   AWS.S3Bucket.DS.High.1047\n   acl           = \"public-read\"\n }",
			expected: []output.SkipRule{
				{
					Rule:    testRule1,
					Comment: "skip the rule.",
				},
				{
					Rule: testRule2,
				},
				{
					Rule:    "AWS.S3Bucket.DS.High.1044",
					Comment: "resource skipped for this rule.",
				},
				{
					Rule: "AWS.S3Bucket.DS.High.1045",
				},
				{
					Rule: "AWS.S3Bucket.DS.High.1046",
				},
				{
					Rule: "AWS.S3Bucket.DS.High.1047",
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
