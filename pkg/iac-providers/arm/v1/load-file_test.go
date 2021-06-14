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

package armv1

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

func TestLoadIacFile(t *testing.T) {
	table := []struct {
		name     string
		filePath string
		armv1    ARMV1
		typeOnly bool
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			// file is skipped if no kind is specified or bad
			name:     "empty config file",
			filePath: filepath.Join(fileTestDataDir, "empty-file.json"),
			armv1:    ARMV1{},
			wantErr:  nil,
		},
		{
			name:     "key-vault",
			filePath: filepath.Join(fileTestDataDir, "key-vault.json"),
			armv1:    ARMV1{},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.armv1.LoadIacFile(tt.filePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}

func TestReadSkipRulesFromTags(t *testing.T) {
	// test data
	testRuleA := "RuleA"
	testCommentA := "RuleA can be skipped"
	testRuleB := "RuleB"
	testCommentB := "RuleB must be skipped"
	testRuleC := "RuleC"
	testCommentC := "RuleC skipped"

	testSkipRule := output.SkipRule{Rule: testRuleA}

	type args struct {
		tags       map[string]interface{}
		resourceID string
	}
	tests := []struct {
		name string
		args args
		want []output.SkipRule
	}{
		{
			name: "nil tags",
			args: args{
				tags: nil,
			},
		},
		{
			name: "tags with no terrascanSkipRules",
			args: args{
				tags: map[string]interface{}{
					"test": "test",
				},
			},
		},
		{
			name: "tags with invalid terrascanSkipRules type",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: "test",
				},
			},
			want: nil,
		},
		{
			name: "tags with invalid SkipRule object",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: []interface{}{1},
				},
			},
			want: nil,
		},
		{
			name: "tags with invalid terrascanSkipRules rule value",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: fmt.Sprintf(`{"%s":%d}`, terrascanSkipRule, 1),
				},
			},
			want: nil,
		},
		{
			name: "tags with one terrascanSkipRules",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: fmt.Sprintf(`[{"%s":"%s"}]`, terrascanSkipRule, testRuleA),
				},
			},
			want: []output.SkipRule{
				{
					Rule: testRuleA,
				},
			},
		},
		{
			name: "tags with multiple terrascanSkipRules",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: fmt.Sprintf(`[{"rule":"%s","comment":"%s"}, {"rule":"%s","comment":"%s"}, {"rule":"%s","comment":"%s"}]`, testRuleA, testCommentA, testRuleB, testCommentB, testRuleC, testCommentC),
				},
			},
			want: []output.SkipRule{
				{
					Rule:    testRuleA,
					Comment: testCommentA,
				},
				{
					Rule:    testRuleB,
					Comment: testCommentB,
				},
				{
					Rule:    testRuleC,
					Comment: testCommentC,
				},
			},
		},
		{
			name: "tags with invalid rule key in terrascanSkipRules",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: fmt.Sprintf(`[{"skip":"%s","comment":"%s"}]`, testRuleA, testCommentA),
				},
			},
			want: []output.SkipRule{{Comment: testCommentA}},
		},
		{
			name: "tags with no comment key in terrascanSkipRules",
			args: args{
				tags: map[string]interface{}{
					terrascanSkip: fmt.Sprintf(`[{"rule":"%s"}]`, testRuleA),
				},
			},
			want: []output.SkipRule{testSkipRule},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readSkipRulesFromTags(tt.args.tags, tt.args.resourceID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readSkipRulesFromTags() = got %v, want %v", got, tt.want)
			}
		})
	}
}
