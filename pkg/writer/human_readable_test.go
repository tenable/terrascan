package writer

import (
	"bytes"
	"testing"

	"github.com/accurics/terrascan/pkg/policy"
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
			name:  "Human Readable Writer: Violations",
			input: violationsInput,
		},
		{
			name: "Human Readable Writer: No Violations",
			input: policy.EngineOutput{
				ViolationStore: &results.ViolationStore{
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
