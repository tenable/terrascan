package scan

import (
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
	ScanTimeout int = 2
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
