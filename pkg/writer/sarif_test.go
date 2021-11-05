package writer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/pkg/version"
)

var abstestpath, _ = getAbsoluteFilePath(violationsInput.Summary.ResourcePath, violationsInput.Violations[0].File)
var testpath = fmt.Sprintf("file://%s", abstestpath)
var testpathForGH = violationsInput.Violations[0].File

const violationTemplate = `{
          "version": "2.1.0",
          "$schema": "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "version": "%s",
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
                          "uri": "%s"
                        },
                        "region": {
                          "startLine": 20
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

var expectedSarifOutput1 = fmt.Sprintf(violationTemplate, version.GetNumeric(), testpath)

var expectedSarifOutput2 = fmt.Sprintf(`{
          "version": "2.1.0",
          "$schema": "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "version": "%s",
                  "informationUri": "https://github.com/accurics/terrascan"
                }
              },
              "results": []
            }
          ]
        }`, version.GetNumeric())

var expectedSarifOutput3 = fmt.Sprintf(`{
          "version": "2.1.0",
          "$schema": "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "name": "terrascan",
                  "version": "%s",
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
        }`, version.GetNumeric())

func TestSarifWriter(t *testing.T) {

	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Sarif Writer: Violations",
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
