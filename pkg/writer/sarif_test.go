package writer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"github.com/tenable/terrascan/pkg/version"
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
                  "informationUri": "https://github.com/tenable/terrascan",
                  "name": "terrascan",
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
                  ],
                  "version": "%s"
                }
              },
              "results": [
                {
                  "ruleId": "AWS.S3Bucket.DS.High.1043",
                  "ruleIndex": 0,
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
                  "informationUri": "https://github.com/tenable/terrascan",
                  "name": "terrascan",
                  "rules": [],
                  "version": "%s"
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
                  "informationUri": "https://github.com/tenable/terrascan",
                  "name": "terrascan",
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
                  ],
                  "version": "%s"
                }
              },
              "results": []
            }
          ]
        }`, version.GetNumeric())

var expectedSarifOutput4 = fmt.Sprintf(`{
          "version": "2.1.0",
          "$schema": "https://json.schemastore.org/sarif-2.1.0-rtm.5.json",
          "runs": [
            {
              "tool": {
                "driver": {
                  "informationUri": "https://github.com/tenable/terrascan",
                  "name": "terrascan",
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
                  ],
                  "version": "%s"
                }
              },
              "invocations": [
                {
                  "executionSuccessful": true,
                  "toolExecutionNotifications": [
                    {
                      "level": "warning",
                      "message": {
                        "text": "kustomization.y(a)ml file not found in the directory test/e2e/test_data/iac/aws/aws_db_instance_violation"
                      }
                    },
                    {
                      "level": "warning",
                      "message": {
                        "text": "no helm charts found in directory test/e2e/test_data/iac/aws/aws_db_instance_violation"
                      }
                    }
                  ]
                }
              ],
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
		{
			name:           "Human Readable Writer: with directory scan error",
			input:          outputWithDirScanErrors,
			expectedOutput: expectedSarifOutput4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			if err := SarifWriter(tt.input, w); (err != nil) != tt.expectedError {
				t.Errorf("HumanReadableWriter() error = gotErr: %v, wantErr: %v", err, tt.expectedError)
			}
			outputBytes := bf.Bytes()
			gotOutput := string(bytes.TrimSpace(outputBytes))

			if equal, _ := utils.AreEqualJSON(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)); !equal {
				t.Errorf("HumanReadableWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
