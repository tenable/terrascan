package writer

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

const (
	configOnlyTestOutputJSON = `{
  "aws_s3_bucket": [
    {
      "id": "aws_s3_bucket.bucket",
      "name": "bucket",
      "source": "modules/m1/main.tf",
      "line": 20,
      "type": "aws_s3_bucket",
      "config": {
        "bucket": "${module.m3.fullbucketname}",
        "policy": "${module.m2.fullbucketpolicy}"
      },
      "skip_rules": null,
      "max_severity": "",
      "min_severity": ""
    }
  ]
}`

	scanTestOutputJSON = `{
  "results": {
    "violations": [
      {
        "rule_name": "s3EnforceUserACL",
        "description": "S3 bucket Access is allowed to all AWS Account Users.",
        "rule_id": "AWS.S3Bucket.DS.High.1043",
        "severity": "HIGH",
        "category": "S3",
        "resource_name": "bucket",
        "resource_type": "aws_s3_bucket",
        "file": "modules/m1/main.tf",
        "line": 20
      }
    ],
    "skipped_violations": [
      {
        "rule_name": "s3EnforceUserACL",
        "description": "S3 bucket Access is allowed to all AWS Account Users.",
        "rule_id": "AWS.S3Bucket.DS.High.1043",
        "severity": "HIGH",
        "category": "S3",
        "resource_name": "bucket",
        "resource_type": "aws_s3_bucket",
        "file": "modules/m1/main.tf",
        "line": 20
      }
    ],
    "scan_summary": {
      "file/folder": "test",
      "iac_type": "terraform",
      "scanned_at": "2020-12-12 11:21:29.902796 +0000 UTC",
      "policies_validated": 566,
      "violated_policies": 1,
      "low": 0,
      "medium": 0,
      "high": 1
    }
  }
}`

	vulnerabilityScanOutputJSON = `{
  "results": {
    "violations": null,
    "skipped_violations": null,
    "vulnerabilities": [
      {
        "image": "test",
        "container": "test",
        "severity": "HIGH",
        "cvss_score": {},
        "description": "GNU Bash. Bash is the GNU Project's shell",
        "vulnerability_id": "CVE-2019-18276",
        "primary_url": "http://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-18276",
        "resource_name": "",
        "resource_type": ""
      }
    ],
    "scan_summary": {
      "file/folder": "",
      "iac_type": "",
      "scanned_at": "",
      "policies_validated": 0,
      "violated_policies": 0,
      "low": 0,
      "medium": 0,
      "high": 0
    }
  }
}`
)

func TestJSONWriter(t *testing.T) {
	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedOutput string
	}{
		{
			name:           "JSON Writer: ResourceConfig",
			input:          resourceConfigInput,
			expectedOutput: configOnlyTestOutputJSON,
		},
		{
			name:           "JSON Writer: Violations",
			input:          violationsInput,
			expectedOutput: scanTestOutputJSON,
		},
		{
			name:           "JSON Writer: Vulnerabilities",
			input:          vulnerabilitiesInput,
			expectedOutput: vulnerabilityScanOutputJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			JSONWriter(tt.input, w)
			if gotOutput := bf.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("JSONWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
