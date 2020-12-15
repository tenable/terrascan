package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
)

const (
	// TODO: --config-only test for XML Writer (xml.Marshal doesn't support maps)

	scanTestOutputXML = `
<results>
  <violations>
    <violation rule_name="s3EnforceUserACL" description="S3 bucket Access is allowed to all AWS Account Users." rule_id="AWS.S3Bucket.DS.High.1043" severity="HIGH" category="S3" resource_name="bucket" resource_type="aws_s3_bucket" file="modules/m1/main.tf" line="20"></violation>
  </violations>
  <scan_summary iac_type="terraform" file/folder="test" scanned_at="2020-12-12 11:21:29.902796 +0000 UTC" policies_validated="566" violated_policies="1" low="0" medium="0" high="1"></scan_summary>
</results>
	`
)

func TestXMLWriter(t *testing.T) {
	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedOutput string
	}{
		{
			name: "XML Writer: Violations",
			input: policy.EngineOutput{
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
					Summary: results.ScanSummary{
						ResourcePath:     "test",
						IacType:          "terraform",
						Timestamp:        "2020-12-12 11:21:29.902796 +0000 UTC",
						TotalPolicies:    566,
						LowCount:         0,
						MediumCount:      0,
						HighCount:        1,
						ViolatedPolicies: 1,
					},
				},
			},
			expectedOutput: scanTestOutputXML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			XMLWriter(tt.input, writer)
			if gotOutput := writer.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("XMLWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
