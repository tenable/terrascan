package writer

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

const (
	// TODO: --config-only test for XML Writer (xml.Marshal doesn't support maps)

	scanTestOutputXML = `
<results>
  <scan_errors></scan_errors>
  <passed_rules></passed_rules>
  <violations>
    <violation rule_name="s3EnforceUserACL" description="S3 bucket Access is allowed to all AWS Account Users." rule_id="AWS.S3Bucket.DS.High.1043" severity="HIGH" category="S3" resource_name="bucket" resource_type="aws_s3_bucket" file="modules/m1/main.tf" line="20"></violation>
  </violations>
  <skipped_violations>
    <violation rule_name="s3EnforceUserACL" description="S3 bucket Access is allowed to all AWS Account Users." rule_id="AWS.S3Bucket.DS.High.1043" severity="HIGH" category="S3" resource_name="bucket" resource_type="aws_s3_bucket" file="modules/m1/main.tf" line="20"></violation>
  </skipped_violations>
  <scan_summary file_folder="test" iac_type="terraform" scanned_at="2020-12-12 11:21:29.902796 +0000 UTC" policies_validated="566" violated_policies="1" low="0" medium="0" high="1"></scan_summary>
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
			name:           "XML Writer: Violations",
			input:          violationsInput,
			expectedOutput: scanTestOutputXML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			XMLWriter(tt.input, w)
			if gotOutput := bf.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("XMLWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
