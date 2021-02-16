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

package writer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/version"
)

func TestJUnitXMLWriter(t *testing.T) {
	testOutput := fmt.Sprintf(`
<testsuites tests="566" name="TERRASCAN_POLICY_SUITES" failures="1" time="0">
  <testsuite tests="566" failures="1" time="0" name="TERRASCAN_POLICY_SUITE" package="test">
    <properties>
      <property name="Terrascan Version" value="%s"></property>
    </properties>
    <testcase classname="modules/m1/main.tf" name="[ERROR] resource: &#34;bucket&#34; at line: 20, violates: RULE - AWS.S3Bucket.DS.High.1043" severity="HIGH" category="S3">
      <failure message="Description: S3 bucket Access is allowed to all AWS Account Users., File: modules/m1/main.tf, Line: 20, Severity: HIGH, Rule Name: s3EnforceUserACL, Rule ID: AWS.S3Bucket.DS.High.1043, Resource Name: bucket, Resource Type: aws_s3_bucket, Category: S3" type=""></failure>
    </testcase>
    <testcase classname="modules/m1/main.tf" name="[ERROR] resource: &#34;bucket&#34; at line: 20, violates: RULE - AWS.S3Bucket.DS.High.1043" severity="HIGH" category="S3">
      <skipped message=""></skipped>
      <failure message="Description: S3 bucket Access is allowed to all AWS Account Users., File: modules/m1/main.tf, Line: 20, Severity: HIGH, Rule Name: s3EnforceUserACL, Rule ID: AWS.S3Bucket.DS.High.1043, Resource Name: bucket, Resource Type: aws_s3_bucket, Category: S3" type=""></failure>
    </testcase>
  </testsuite>
</testsuites>
  `, version.Get())

	testOutputNoViolations := fmt.Sprintf(`
<testsuites tests="550" name="TERRASCAN_POLICY_SUITES" failures="1" time="0">
  <testsuite tests="550" failures="1" time="0" name="TERRASCAN_POLICY_SUITE" package="test_resource_path">
    <properties>
      <property name="Terrascan Version" value="%s"></property>
    </properties>
  </testsuite>
</testsuites>
	`, version.Get())

	type args struct {
		data interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "incorrect input for JunitXMLWriter",
			args: args{
				// some invalid data
				data: struct{ Name string }{Name: "test"},
			},
			wantErr: true,
		},
		{
			name: "data with violations and skipped violations",
			args: args{
				data: violationsInput,
			},
			wantWriter: testOutput,
		},
		{
			name: "data with no violations",
			args: args{
				policy.EngineOutput{
					ViolationStore: &results.ViolationStore{
						Summary: results.ScanSummary{
							ResourcePath:     "test_resource_path",
							IacType:          "k8s",
							Timestamp:        "2020-12-12 11:21:29.902796 +0000 UTC",
							TotalPolicies:    550,
							LowCount:         0,
							MediumCount:      0,
							HighCount:        1,
							ViolatedPolicies: 1,
						},
					},
				},
			},
			wantWriter: testOutputNoViolations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := JUnitXMLWriter(tt.args.data, writer); (err != nil) != tt.wantErr {
				t.Errorf("JUnitXMLWriter() got error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); !strings.EqualFold(strings.TrimSpace(gotWriter), strings.TrimSpace(tt.wantWriter)) {
				t.Errorf("JUnitXMLWriter() got = %v, want = %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestGetViolationString(t *testing.T) {
	type args struct {
		v results.Violation
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "violation with all fields",
			args: args{
				v: results.Violation{
					RuleName:     "RuleA",
					Description:  "test rule",
					RuleID:       "Rule.A",
					Severity:     "MEDIUM",
					Category:     "A",
					ResourceName: "aws_resource",
					ResourceType: "some_resource_type",
					File:         "file.yaml",
					LineNumber:   1,
				},
			},
			want: "Description: test rule, File: file.yaml, Line: 1, Severity: MEDIUM, Rule Name: RuleA, Rule ID: Rule.A, Resource Name: aws_resource, Resource Type: some_resource_type, Category: A",
		},
		{
			name: "violation with all fields, blank resource name",
			args: args{
				v: results.Violation{
					RuleName:     "RuleB",
					Description:  "test rule 2",
					RuleID:       "Rule.B",
					Severity:     "HIGH",
					Category:     "B",
					ResourceType: "test_resource_type",
					File:         "file1.yaml",
					LineNumber:   2,
				},
			},
			want: `Description: test rule 2, File: file1.yaml, Line: 2, Severity: HIGH, Rule Name: RuleB, Rule ID: Rule.B, Resource Name: "", Resource Type: test_resource_type, Category: B`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getViolationString(tt.args.v); got != tt.want {
				t.Errorf("getViolationString() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
