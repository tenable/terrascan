package writer

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
)

// these variables would be used as test input across the writer package
var (
	resourceConfigInput = output.AllResourceConfigs{
		"aws_s3_bucket": []output.ResourceConfig{
			{
				ID:     "aws_s3_bucket.bucket",
				Name:   "bucket",
				Source: "modules/m1/main.tf",
				Line:   20,
				Type:   "aws_s3_bucket",
				Config: map[string]string{
					"bucket": "${module.m3.fullbucketname}",
					"policy": "${module.m2.fullbucketpolicy}",
				},
			},
		},
	}

	vulnerabilitiesInput = policy.EngineOutput{
		ViolationStore: &results.ViolationStore{
			Vulnerabilities: []*results.Vulnerability{
				{
					Image:           "test",
					Container:       "test",
					VulnerabilityID: "CVE-2019-18276",
					PrimaryURL:      "http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-18276",
					Description:     "GNU Bash. Bash is the GNU Project's shell",
					Severity:        "HIGH",
				},
			},
		},
	}

	violationsInput = policy.EngineOutput{
		ViolationStore: &results.ViolationStore{
			Violations: []*results.Violation{
				{
					RuleName:     "s3EnforceUserACL",
					Description:  "S3 bucket Access is allowed to all AWS Account Users.",
					RuleID:       "AWS.S3Bucket.DS.High.1043",
					Severity:     "HIGH",
					Category:     "S3",
					ResourceName: "bucket",
					ResourceType: "aws_s3_bucket",
					File:         "modules/m1/main.tf",
					LineNumber:   20,
				},
			},
			SkippedViolations: []*results.Violation{
				{
					RuleName:     "s3EnforceUserACL",
					Description:  "S3 bucket Access is allowed to all AWS Account Users.",
					RuleID:       "AWS.S3Bucket.DS.High.1043",
					Severity:     "HIGH",
					Category:     "S3",
					Comment:      "",
					ResourceName: "bucket",
					ResourceType: "aws_s3_bucket",
					File:         "modules/m1/main.tf",
					LineNumber:   20,
				},
			},
			Summary: results.ScanSummary{
				ResourcePath:         "test",
				IacType:              "terraform",
				Timestamp:            "2020-12-12 11:21:29.902796 +0000 UTC",
				TotalPolicies:        566,
				LowCount:             0,
				MediumCount:          0,
				HighCount:            1,
				ViolatedPolicies:     1,
				ShowViolationDetails: true,
			},
		},
	}
)

const (
	configOnlyTestOutputYAML = `aws_s3_bucket:
    - id: aws_s3_bucket.bucket
      name: bucket
      source: modules/m1/main.tf
      line: 20
      type: aws_s3_bucket
      config:
        bucket: ${module.m3.fullbucketname}
        policy: ${module.m2.fullbucketpolicy}
      skip_rules: []
      maxseverity: ""
      minseverity: ""
      containerimages: []
      initcontainerimages: []
      isremotemodule: null
      terraformversion: ""
      providerversion: ""`

	scanTestOutputYAML = `results:
    violations:
        - rule_name: s3EnforceUserACL
          description: S3 bucket Access is allowed to all AWS Account Users.
          rule_id: AWS.S3Bucket.DS.High.1043
          severity: HIGH
          category: S3
          resource_name: bucket
          resource_type: aws_s3_bucket
          file: modules/m1/main.tf
          line: 20
    skipped_violations:
        - rule_name: s3EnforceUserACL
          description: S3 bucket Access is allowed to all AWS Account Users.
          rule_id: AWS.S3Bucket.DS.High.1043
          severity: HIGH
          category: S3
          resource_name: bucket
          resource_type: aws_s3_bucket
          file: modules/m1/main.tf
          line: 20
    scan_summary:
        file/folder: test
        iac_type: terraform
        scanned_at: 2020-12-12 11:21:29.902796 +0000 UTC
        policies_validated: 566
        violated_policies: 1
        low: 0
        medium: 0
        high: 1`

	vulnerabilityScanOutput = `results:
    violations: []
    skipped_violations: []
    vulnerabilities:
        - image: test
          container: test
          severity: HIGH
          description: GNU Bash. Bash is the GNU Project's shell
          vulnerability_id: CVE-2019-18276
          primary_url: http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-18276
          resource_name: ""
          resource_type: ""
    scan_summary:
        file/folder: ""
        iac_type: ""
        scanned_at: ""
        policies_validated: 0
        violated_policies: 0
        low: 0
        medium: 0
        high: 0`
)

func TestYAMLWriter(t *testing.T) {
	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedOutput string
	}{
		{
			name:           "YAML Writer: ResourceConfig",
			input:          resourceConfigInput,
			expectedOutput: configOnlyTestOutputYAML,
		},
		{
			name:           "YAML Writer: Violations",
			input:          violationsInput,
			expectedOutput: scanTestOutputYAML,
		},
		{
			name:           "YAML Writer: Vulnerabilities",
			input:          vulnerabilitiesInput,
			expectedOutput: vulnerabilityScanOutput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			YAMLWriter(tt.input, w)
			if gotOutput := bf.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("YAMLWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
