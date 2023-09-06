/*
    Copyright (C) 2022 Tenable, Inc.

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
	"io"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/test/helper"
)

const (
	// ScanCommand is terrascan's scan command
	ScanCommand string = "scan"

	// ScanTimeout is default scan command execution timeout
	ScanTimeout int = 3

	// WebhookScanTimeout is default scan command webhook execution timeout
	WebhookScanTimeout int = 30

	// RemoteScanTimeout is default scan command remote execution timeout
	RemoteScanTimeout int = 30
)

// RunScanAndAssertGoldenOutputRegexWithTimeout runs the scan command with supplied paramters and compares actual and golden output
// it replaces variable parts in output with regex eg: timestamp, file path
// added to provide extra option for specifying timeout.
func RunScanAndAssertGoldenOutputRegexWithTimeout(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, scanTimeout int, args ...string) {
	session, goldenFileAbsPath := RunScanCommandWithTimeOut(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, scanTimeout, args...)
	helper.CompareActualWithGoldenSummaryRegex(session, goldenFileAbsPath, isJunitXML, isStdOut)
}

// RunScanAndAssertGoldenOutputRegex runs the scan command with supplied parameters and compares actual and golden output
// it replaces variable parts in output with regex eg: timestamp, file path
func RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenSummaryRegex(session, goldenFileAbsPath, isJunitXML, isStdOut)
}

// RunScanAndAssertGoldenOutput runs the scan command with supplied parameters and compares actual and golden output
func RunScanAndAssertGoldenOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGolden(session, goldenFileAbsPath, isStdOut)
}

// RunScanCommandWithTimeOut executes the scan command, validates exit code
// added to provide extra option for specifying timeout.
func RunScanCommandWithTimeOut(terrascanBinaryPath, relGoldenFilePath string, exitCode int, outWriter, errWriter io.Writer, scanTimeout int, args ...string) (*gexec.Session, string) {
	argList := []string{ScanCommand}
	argList = append(argList, args...)
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, argList...)
	gomega.Eventually(session, scanTimeout).Should(gexec.Exit(exitCode))
	goldenFileAbsPath, err := filepath.Abs(relGoldenFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session, goldenFileAbsPath
}

// RunScanAndAssertJSONOutput runs the scan command with supplied parameters and compares actual and golden output
func RunScanAndAssertJSONOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenJSON(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertJSONOutputString runs the scan command with supplied parameters and compares actual and golden output
func RunScanAndAssertJSONOutputString(terrascanBinaryPath, goldenString string, exitCode int, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	argList := []string{ScanCommand}
	argList = append(argList, args...)
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, argList...)
	gomega.Eventually(session, ScanTimeout).Should(gexec.Exit(exitCode))
	helper.CompareActualWithGoldenJSONString(session, goldenString, isStdOut)
}

// RunScanAndAssertYAMLOutput runs the scan command with supplied parameters and compares actual and golden output
func RunScanAndAssertYAMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenYAML(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertXMLOutput runs the scan command with supplied parameters and compares actual and golden output
func RunScanAndAssertXMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenXML(session, goldenFileAbsPath, isStdOut)
}

// RunScanAndAssertErrorMessage runs the scan command with supplied parameters and checks of error string is present
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

// RunScanAndAssertGoldenSarifOutputRegex runs the scan command with supplied parameters and compares actual and golden output
// it replaces variable parts in output with regex eg: uri, version path
func RunScanAndAssertGoldenSarifOutputRegex(terrascanBinaryPath, relGoldenFilePath string, exitCode int, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualSarifOutputWithGoldenSummaryRegex(session, goldenFileAbsPath)
}
