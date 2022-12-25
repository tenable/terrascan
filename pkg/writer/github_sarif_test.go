package writer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/tenable/terrascan/pkg/utils"
	"github.com/tenable/terrascan/pkg/version"
)

const violationTemplateForGH = `{
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

func TestGitHubSarifWriter(t *testing.T) {

	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedError  bool
		expectedOutput string
	}{
		{
			name:           "Sarif Writer for GitHub: Violations",
			input:          violationsInput,
			expectedOutput: expectedSarifViolationOutputGH,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bf bytes.Buffer
			w := []io.Writer{&bf}
			if err := GitHubSarifWriter(tt.input, w); (err != nil) != tt.expectedError {
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
