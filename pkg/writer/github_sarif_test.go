package writer

import (
	"bytes"
	"fmt"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/pkg/version"
	"strings"
	"testing"
)

const violationTemplateForGH = `{
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
                          "uri": "%s",
						  "uriBaseId": "test"
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

var expectedSarifViolationOutputGH = fmt.Sprintf(violationTemplateForGH, version.GetNumeric(), testpathForGH)

func TestGithubSarifWriter(t *testing.T) {

	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Sarif Writer for Github: Violations",
			input:          violationsInput,
			expectedOutput: expectedSarifViolationOutputGH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := GithubSarifWriter(tt.input, writer); (err != nil) != tt.expectedError {
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
