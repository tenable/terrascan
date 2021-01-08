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
)

// ---------------------- unit tests -------------------------------- //
func TestGetSkipRules(t *testing.T) {

	table := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "no rules",
			input: "no rules here",
			// expected would be an empty slice of strings
		},
		{
			name:     "single rule id",
			input:    "#ts:skip=AWS.S3Bucket.DS.High.1041",
			expected: []string{"AWS.S3Bucket.DS.High.1041"},
		},
		{
			name:     "multiple comma separated no space",
			input:    "#ts:skip=AWS.S3Bucket.DS.High.1041,AWS.S3Bucket.DS.High.1042",
			expected: []string{"AWS.S3Bucket.DS.High.1041", "AWS.S3Bucket.DS.High.1042"},
		},
		{
			name: "multiple comma separated random space",
			input: "#ts:skip=  AWS.S3Bucket.DS.High.1041   ,		AWS.S3Bucket.DS.High.1042",
			expected: []string{"AWS.S3Bucket.DS.High.1041", "AWS.S3Bucket.DS.High.1042"},
		},
		{
			name:     "sample resource config",
			input:    "{\n     #ts:skip=AWS.S3Bucket.DS.High.1041\n   region        = var.region\n   #ts:skip=AWS.S3Bucket.DS.High.1042 AWS.S3Bucket.DS.High.1043\n   bucket        = local.bucket_name\n   #ts:skip=AWS.S3Bucket.DS.High.1044,AWS.S3Bucket.DS.High.1045\n   force_destroy = true\n   #ts:skip= AWS.S3Bucket.DS.High.1046   ,   AWS.S3Bucket.DS.High.1047\n   acl           = \"public-read\"\n }",
			expected: []string{"AWS.S3Bucket.DS.High.1041", "AWS.S3Bucket.DS.High.1042", "AWS.S3Bucket.DS.High.1044", "AWS.S3Bucket.DS.High.1045", "AWS.S3Bucket.DS.High.1046", "AWS.S3Bucket.DS.High.1047"},
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
