package writer

import (
	"bytes"
	"testing"

	"github.com/accurics/terrascan/pkg/results"
)

// TODO: string comparision - expected and output

func TestHumanReadbleWriter(t *testing.T) {
	type funcInput interface{}
	tests := []struct {
		name          string
		input         funcInput
		expectedError bool
	}{
		{
			name: "Human Readable Writer: Violations",
			input: results.DefaultScanResult{
				IacType:              "terraform",
				ResourcePath:         "/test",
				Timestamp:            "2020-12-12 11:21:29.902796 +0000 UTC",
				TotalPolicies:        566,
				ShowViolationDetails: true,
				ViolationStore: results.ViolationStore{
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
					Count: results.ViolationStats{
						LowCount:    0,
						MediumCount: 0,
						HighCount:   1,
						TotalCount:  1,
					},
				},
			},
		},
		{
			name: "Human Readable Writer: No Violations",
			input: results.DefaultScanResult{
				IacType:              "k8s",
				ResourcePath:         "/test",
				Timestamp:            "2020-12-12 11:21:29.902796 +0000 UTC",
				TotalPolicies:        566,
				ShowViolationDetails: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := HumanReadbleWriter(tt.input, writer); (err != nil) != tt.expectedError {
				t.Errorf("HumanReadbleWriter() error = gotErr: %v, wantErr: %v", err, tt.expectedError)
				return
			}
		})
	}
}
