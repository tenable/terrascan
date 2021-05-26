package writer

import (
	"bytes"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"strings"
	"testing"
)

const expectedSarifOutput1 = `{
          "version": "2.1.0",
          "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "informationUri": "https://github.com/accurics/terrascan",
                  "rules": [
                    {
                      "id": "AWS.S3Bucket.DS.High.1043",
                      "name": "s3EnforceUserACL",
                      "shortDescription": {
                        "text": "S3 bucket Access is allowed to all AWS Account Users."
                      },
                      "properties": {
                        "category": "S3",
                        "severity": "HIGH"
                      }
                    }
                  ]
                }
              },
              "results": [
                {
                  "ruleId": "AWS.S3Bucket.DS.High.1043",
                  "level": "error",
                  "message": {
                    "text": "S3 bucket Access is allowed to all AWS Account Users."
                  },
                  "locations": [
                    {
                      "physicalLocation": {
                        "artifactLocation": {
                          "uri": "file://test/modules/m1/main.tf"
                        },
                        "region": {
                          "startLine": 19
                        }
                      },
                      "logicalLocations": [
                        {
                          "name": "bucket",
                          "kind": "aws_s3_bucket"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }`

const expectedSarifOutput2 = `{
          "version": "2.1.0",
          "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "informationUri": "https://github.com/accurics/terrascan"
                }
              },
              "results": []
            }
          ]
        }`

const expectedSarifOutput3 = `{
          "version": "2.1.0",
          "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "informationUri": "https://github.com/accurics/terrascan",
                  "rules": [
                    {
                      "id": "AWS.S3Bucket.DS.High.1043",
                      "name": "s3EnforceUserACL",
                      "shortDescription": {
                        "text": "S3 bucket Access is allowed to all AWS Account Users."
                      },
                      "properties": {
                        "category": "S3",
                        "severity": "HIGH"
                      }
                    }
                  ]
                }
              },
              "results": []
            }
          ]
        }`

func TestSarifWriter(t *testing.T) {

	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Human Readable Writer: Violations",
			input:          violationsInput,
			expectedOutput: expectedSarifOutput1,
		},
		{
			name: "Human Readable Writer: No Violations",
			input: policy.EngineOutput{
				ViolationStore: &results.ViolationStore{
					Summary: summaryWithNoViolations,
				},
			},
			expectedOutput: expectedSarifOutput2,
		},
		{
			name:           "Human Readable Writer: With PassedRules",
			input:          outputWithPassedRules,
			expectedOutput: expectedSarifOutput3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := SarifWriter(tt.input, writer); (err != nil) != tt.expectedError {
				t.Errorf("HumanReadbleWriter() error = gotErr: %v, wantErr: %v", err, tt.expectedError)
			}
			outputBytes := writer.Bytes()
			gotOutput := string(bytes.TrimSpace(outputBytes))

			if equal, _ := utils.AreEqualJSON(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)); !equal {
				t.Errorf("HumanReadbleWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
