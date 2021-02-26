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

// RunScanCommandAndAssertTextOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanCommandAndAssertTextOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenSummaryRegex(session, goldenFileAbsPath, isJunitXML, isStdOut)
}

// RunScanCommandAndAssertJSONOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanCommandAndAssertJSONOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenJSON(session, goldenFileAbsPath, isStdOut)
}

// RunScanCommandAndAssertYAMLOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanCommandAndAssertYAMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenYAML(session, goldenFileAbsPath, isStdOut)
}

// RunScanCommandAndAssertXMLOutput runs the scan command with supplied paramters and compares actual and golden output
func RunScanCommandAndAssertXMLOutput(terrascanBinaryPath, relGoldenFilePath string, exitCode int, isJunitXML, isStdOut bool, outWriter, errWriter io.Writer, args ...string) {
	session, goldenFileAbsPath := RunScanCommand(terrascanBinaryPath, relGoldenFilePath, exitCode, outWriter, errWriter, args...)
	helper.CompareActualWithGoldenXML(session, goldenFileAbsPath, isStdOut)
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
