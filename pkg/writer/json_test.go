package writer

import (
	"bytes"
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
      "min_severity": "",
      "container_images": null,
      "init_container_images": null
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			JSONWriter(tt.input, writer)
			if gotOutput := writer.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("JSONWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
