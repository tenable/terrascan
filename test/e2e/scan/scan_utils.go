/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package scan

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/accurics/terrascan/test/helper"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	// ScanCommand is terrascan's scan command
	ScanCommand string = "scan"

	// ScanTimeout is default scan command execution timeout
	ScanTimeout int = 3

	// RemoteScanTimeout is default scan command remote execution timeout
	RemoteScanTimeout int = 30
)

// RunScanAndAssertGoldenOutputRegex runs the scan command with supplied paramters and compares actual and golden output
// it replaces variable parts in output with regex eg: timestamp, file path
func RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenSummaryRegex(session, goldenFileAbsPath, isJunitXML, isStdOut)
}

// RunScanAndAssertGoldenOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanAndAssertGoldenOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGolden(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertJSONOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanAndAssertJSONOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenJSON(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertJSONOutputString runs the scan command with supplied paramters and compares actual and golden output
func RunScanAndAssertJSONOutputString(terrascanBinaryPath, goldenString string, exitCode int, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	argList := []string{ScanCommand}
	argList = append(argList, args...)
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, argList...)
	gomega.Eventually(session, ScanTimeout).Should(gexec.Exit(exitCode))
	helper.CompareActualWithGoldenJSONString(session, goldenString, isStdOut)
}

// RunScanAndAssertYAMLOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanAndAssertYAMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenYAML(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertXMLOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanAndAssertXMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenXML(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertErrorMessage runs the scan command with supplied paramters and checks of error string is present
func RunScanAndAssertErrorMessage(terrascanBinaryPath string, exitCode, timeOut int, errString string, outWriter, errWriter io.Writer, args ...string) {
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
	gomega.Eventually(session, timeOut).Should(gexec.Exit(exitCode))
	helper.ContainsErrorSubString(session, errString)
}

// RunScanCommand executes the scan command, validates exit code
func RunScanCommand(terrascanBinaryPath, relGoldenFilePath string, exitCode int, outWriter, errWriter io.Writer, args ...string) (*gexec.Session, string) {
	argList := []string{ScanCommand}
	argList = append(argList, args...)
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, argList...)
	gomega.Eventually(session, ScanTimeout).Should(gexec.Exit(exitCode))
	goldenFileAbsPath, err := filepath.Abs(relGoldenFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session, goldenFileAbsPath
}

// GetSarifGoldenString gives out golden sarif format string
func GetSarifGoldenString(template, terrascanVersion string, absfilepath string) string {
	return fmt.Sprintf(template, terrascanVersion, fmt.Sprintf("file://%s", absfilepath))
}

// SarifTemplateAWSAMIViolation string
const SarifTemplateAWSAMIViolation = `{
  "version": "2.1.0",
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "terrascan",
          "version": "%s",
          "informationUri": "https://github.com/accurics/terrascan",
          "rules": [
            {
              "id": "AWS.EC2.Encryption\u0026KeyManagement.Medium.0688",
              "name": "amiNotEncrypted",
              "shortDescription": {
                "text": "Enable AWS AMI Encryption"
              },
              "properties": {
                "category": "Encryption \u0026 KeyManagement",
                "severity": "MEDIUM"
              }
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "AWS.EC2.Encryption\u0026KeyManagement.Medium.0688",
          "level": "warning",
          "message": {
            "text": "Enable AWS AMI Encryption"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "%s"
                },
                "region": {
                  "startLine": 5
                }
              },
              "logicalLocations": [
                {
                  "name": "awsAmiEncrypted",
                  "kind": "aws_ami"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`

// SarifTemplateK8sTLSViolation string template
const SarifTemplateK8sTLSViolation = `{
  "version": "2.1.0",
  "$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
  "runs": [
    {
      "tool": {
        "driver": {
          "name": "terrascan",
          "version": "%s",
          "informationUri": "https://github.com/accurics/terrascan",
          "rules": [
            {
              "id": "AC-K8-NS-IN-H-0020",
              "name": "noHttps",
              "shortDescription": {
                "text": "TLS disabled can affect the confidentiality of the data in transit"
              },
              "properties": {
                "category": "Network Security",
                "severity": "HIGH"
              }
            }
          ]
        }
      },
      "results": [
        {
          "ruleId": "AC-K8-NS-IN-H-0020",
          "level": "error",
          "message": {
            "text": "TLS disabled can affect the confidentiality of the data in transit"
          },
          "locations": [
            {
              "physicalLocation": {
                "artifactLocation": {
                  "uri": "%s"
                },
                "region": {
                  "startLine": 1
                }
              },
              "logicalLocations": [
                {
                  "name": "ingress-demo-disallowed",
                  "kind": "kubernetes_ingress"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
