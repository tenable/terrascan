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

package helper

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"gopkg.in/yaml.v3"
)

const (
	// ExitCodeZero represents command exit code 0
	ExitCodeZero = iota
	// ExitCodeOne represents command exit code 1
	ExitCodeOne
	// ExitCodeThree represents command exit code 3
	ExitCodeThree = iota + 1
	// ExitCodeFour represents command exit code 4
	ExitCodeFour
	// ExitCodeFive represents command exit code 5
	ExitCodeFive
)

var (
	// scannedAtPattern is regex for 'scanned at' attribute in violations output
	scannedAtPattern = regexp.MustCompile(`["]*[sS]canned[ _][aA]t["]*[ \t]*[:=][ \t]*["]*[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}[.][0-9]{1,9} [+-][0-9]{4} UTC["]*[,]{0,1}`)

	// fileFolderPattern is regex for 'file/folder' attribute in violations output
	fileFolderPattern = regexp.MustCompile(`["]*[fF]ile[\/_][fF]older["]*[ \t]*[:=][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// planRootPattern is regex for 'file/folder' attribute in violations output
	planRootPattern = regexp.MustCompile(`["]*[pP]lan[\/_][rR]oot["]*[ \t]*[:=][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// filePattern is regex for 'file' attribute in violations output
	filePattern = regexp.MustCompile(`["]*[fF]ile["]*[ \t]*[:=][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// sarifVersionPattern is regex for 'version' attribute in sarif violations output
	sarifVersionPattern = regexp.MustCompile(`["]*[vV]ersion["][:=][ \t]*["]([0-9])+[\.]([0-9])+[\.]([0-9])+["],`)

	// sarifUriPattern is regex for 'uri' attribute in sarif violations output
	sarifURIPattern = regexp.MustCompile(`["]*[uU]ri["][:=][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// directoryPattern is regex for 'directory' attribute in scan_errors
	directoryPattern = regexp.MustCompile(`["]*directory["]*[ \t]*[:=][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// errMsgPattern is regex for 'errMsg' attribute in scan_errors
	errMsgPattern = regexp.MustCompile(`["]*errMsg["]*[ \t]*[:=][ \t]*["]*(.+)["]*`)

	// packagePattern is regex for 'package' attribute in junit-xml output
	packagePattern = regexp.MustCompile(`package=["]*(.+)[\\\/](.+)["]*`)

	// classnamePattern is regex for 'package' attribute in junit-xml output
	classnamePattern = regexp.MustCompile(`classname=["]*(.+)[\\\/](.+)["]*`)

	// versionValuePattern is regex for 'value' attribute in junit-xml output (which is terrascan version)
	versionValuePattern = regexp.MustCompile(`value="v[0-9][\.]([0-9])+[\.]([0-9])+"`)

	// sourceRegexPattern is regex for 'file/folder' attribute in violations output
	sourceRegexPattern = regexp.MustCompile(`["]*source["]*[ \t]*[:][ \t]*["]*(.+)[\\\/](.+)["]*`)

	// version is regex for 'version' attribute in violations output
	versionRegexPattern = regexp.MustCompile(`["version":]*"[0-9][\.]([0-9])+[\.]([0-9])"`)
)

// ValidateExitCode validates the exit code of a gexec.Session
func ValidateExitCode(session *gexec.Session, timeout, exitCode int) {
	gomega.Eventually(session, timeout).Should(gexec.Exit(exitCode))
}

// ValidateDirectoryExists validates that a directory exists at the provided path
func ValidateDirectoryExists(path string) {
	_, err := os.Stat(path)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(path).To(gomega.BeADirectory())
}

// CompareActualWithGolden compares actual string with contents of golden file path passed as parameter
func CompareActualWithGolden(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)
	if utils.IsWindowsPlatform() {
		fileBytes = utils.ReplaceWinNewLineBytes(fileBytes)
	}
	gomega.Expect(string(sessionBytes)).Should(gomega.Equal(string(fileBytes)))
}

// CompareActualWithGoldenConfigOnlyRegex compares actual string with contents of golden file path passed as parameter
func CompareActualWithGoldenConfigOnlyRegex(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)
	sessionBytes = sourceRegexPattern.ReplaceAll(sessionBytes, []byte(""))
	fileBytes = sourceRegexPattern.ReplaceAll(fileBytes, []byte(""))
	gomega.Expect(string(sessionBytes)).Should(gomega.Equal(string(fileBytes)))
}

// CompareActualWithGoldenSummaryRegex compares actual string with contents of golden file passed as parameter
// ignores specified regex patterns from the actual and golden text
func CompareActualWithGoldenSummaryRegex(session *gexec.Session, goldenFileAbsPath string, isJunitXML, isStdOut bool) {
	fileData, err := os.ReadFile(goldenFileAbsPath)
	if utils.IsWindowsPlatform() {
		fileData = utils.ReplaceWinNewLineBytes(fileData)
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	var sessionOutput, fileContents string

	/*
		-There are a few attributes in our generated output which is going to vary for every machine.
		 eg: file/folder, scannedAt, file

		-These attributes needs to be removed from the actual and golden output before comparing
		-These attributes are removed based on the scan result output type
		 eg: 1. junit-xml has attributes "package", "classname", "value" that needs to be removed
		 	 2. other output formats has attributes "scannedAt", "file/folder" that needs to be removed
	*/

	if isStdOut {
		sessionOutput = string(session.Wait().Out.Contents())
	} else {
		sessionOutput = string(session.Wait().Err.Contents())
	}

	fileContents = string(fileData)

	sessionOutput = strings.TrimSpace(sessionOutput)
	fileContents = strings.TrimSpace(fileContents)

	// replace file from the output, it will cause issues for absolute paths
	sessionOutput = filePattern.ReplaceAllString(sessionOutput, "")
	fileContents = filePattern.ReplaceAllString(fileContents, "")

	if isJunitXML {
		sessionOutput = packagePattern.ReplaceAllString(sessionOutput, "")
		fileContents = packagePattern.ReplaceAllString(fileContents, "")

		sessionOutput = classnamePattern.ReplaceAllString(sessionOutput, "")
		fileContents = classnamePattern.ReplaceAllString(fileContents, "")

		sessionOutput = versionValuePattern.ReplaceAllString(sessionOutput, "")
		fileContents = versionValuePattern.ReplaceAllString(fileContents, "")
	} else {
		sessionOutput = scannedAtPattern.ReplaceAllString(sessionOutput, "")
		fileContents = scannedAtPattern.ReplaceAllString(fileContents, "")

		sessionOutput = fileFolderPattern.ReplaceAllString(sessionOutput, "")
		fileContents = fileFolderPattern.ReplaceAllString(fileContents, "")

		sessionOutput = directoryPattern.ReplaceAllString(sessionOutput, "")
		fileContents = directoryPattern.ReplaceAllString(fileContents, "")

		sessionOutput = errMsgPattern.ReplaceAllString(sessionOutput, "")
		fileContents = errMsgPattern.ReplaceAllString(fileContents, "")

		sessionOutput = planRootPattern.ReplaceAllString(sessionOutput, "")
		fileContents = planRootPattern.ReplaceAllString(fileContents, "")
	}

	gomega.Expect(sessionOutput).Should(gomega.BeIdenticalTo(fileContents))
}

// CompareActualWithGoldenJSON compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenJSON(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := json.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = json.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	CompareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenJSONString compares actual data with golden json string passed as parameter
func CompareActualWithGoldenJSONString(session *gexec.Session, golden string, isStdOut bool) {
	goldenBytes := []byte(golden)

	var sessionBytes []byte

	if isStdOut {
		sessionBytes = session.Wait().Out.Contents()
	} else {
		sessionBytes = session.Wait().Err.Contents()
	}

	sessionBytes = bytes.TrimSpace(sessionBytes)
	goldenBytes = bytes.TrimSpace(goldenBytes)

	gomega.Expect(utils.AreEqualJSONBytes(sessionBytes, goldenBytes)).To(gomega.BeTrue())
}

// CompareActualWithGoldenYAML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenYAML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := yaml.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = yaml.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	CompareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenXML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenXML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := xml.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = xml.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	CompareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// ContainsErrorSubString will assert if error string is part of error output
func ContainsErrorSubString(session *gexec.Session, errSubString string) {
	gomega.Expect(string(session.Wait().Err.Contents())).Should(gomega.ContainSubstring(errSubString))
}

// ContainsDirScanErrorSubString will assert if dir error substring is part of scan output
func ContainsDirScanErrorSubString(session *gexec.Session, errSubString string) {
	gomega.Expect(string(session.Wait().Out.Contents())).Should(gomega.ContainSubstring(errSubString))
}

// DoesNotContainsErrorSubString will assert that a string is not part of the error output
func DoesNotContainsErrorSubString(session *gexec.Session, errSubString string) {
	gomega.Expect(string(session.Wait().Err.Contents())).ShouldNot(gomega.ContainSubstring(errSubString))
}

// GetTerrascanBinaryPath returns the terrascan binary path
func GetTerrascanBinaryPath() string {
	terrascanBinaryPath := os.Getenv("TERRASCAN_BIN_PATH")
	ginkgo.Describe("terrascan binary path should be set for executing tests", func() {
		if terrascanBinaryPath == "" {
			ginkgo.Fail("ensure that TERRASCAN_BIN_PATH is set")
		}
	})
	return terrascanBinaryPath
}

// RunCommand will initialise the command to run and return session
func RunCommand(path string, outWriter, errWriter io.Writer, args ...string) *gexec.Session {
	cmd := exec.Command(path, args...)
	return getSession(cmd, outWriter, errWriter)
}

// RunCommandDir will initialise the command to run in a specific directory and return session
func RunCommandDir(path, workDir string, outWriter, errWriter io.Writer, args ...string) *gexec.Session {
	cmd := exec.Command(path, args...)
	cmd.Dir = workDir
	return getSession(cmd, outWriter, errWriter)
}

func getSession(cmd *exec.Cmd, outWriter, errWriter io.Writer) *gexec.Session {
	session, err := gexec.Start(cmd, outWriter, errWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session
}

// GetByteData is a helper function to get data in byte slice from session and golden file
func GetByteData(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) ([]byte, []byte) {
	fileBytes, err := os.ReadFile(goldenFileAbsPath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	var sessionBytes []byte

	if isStdOut {
		sessionBytes = session.Wait().Out.Contents()
	} else {
		sessionBytes = session.Wait().Err.Contents()
	}

	sessionBytes = bytes.TrimSpace(sessionBytes)
	fileBytes = bytes.TrimSpace(fileBytes)
	return sessionBytes, fileBytes
}

// CompareSummaryAndViolations is a helper function to compare actual and expected, summary and violations
func CompareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput policy.EngineOutput) {
	var sessionOutputSummary, fileDataSummary results.ScanSummary
	var actualViolations, expectedViolations violations
	var actualSkippedViolations, expectedSkippedViolations violations
	var actualPassedRules, expectedPassedRules passedRules

	actualViolations = sessionEngineOutput.ViolationStore.Violations
	expectedViolations = fileDataEngineOutput.ViolationStore.Violations

	/*
		-There are a few attributes in our generated output which is going to vary for every machine.
		-eg: file/folder, scannedAt, file
		-These attributes needs to be removed from the actual and golden output before comparing
		-Also, the violations are not in order, they need to be sorted from both actual and golden output,
		 before the comparision is made. Below are the steps:

		1. sort actual and golden violations and remove "file" attribute
		2. sort actual and golden skipped violations and remove "file" attribute
		3. remove "scannedAt" attribute, which is a timestamp from actual summary
		4. remove "scannedAt" from golden summary
		5. compare violations, skipped violations and summary in actual and golden
	*/

	// 1. sort actual and golden violations and remove "file" attribute
	sort.Sort(actualViolations)
	sort.Sort(expectedViolations)
	removeFileAndRoothFromViolations(actualViolations)
	removeFileAndRoothFromViolations(expectedViolations)

	actualSkippedViolations = sessionEngineOutput.ViolationStore.SkippedViolations
	expectedSkippedViolations = fileDataEngineOutput.ViolationStore.SkippedViolations

	// 2. sort actual and golden skipped violations and remove "file" attribute
	sort.Sort(actualSkippedViolations)
	sort.Sort(expectedSkippedViolations)
	removeFileAndRoothFromViolations(actualSkippedViolations)
	removeFileAndRoothFromViolations(expectedSkippedViolations)

	actualPassedRules = sessionEngineOutput.ViolationStore.PassedRules
	expectedPassedRules = fileDataEngineOutput.ViolationStore.PassedRules

	// 3. sort actual and golden passed rules
	sort.Sort(actualPassedRules)
	sort.Sort(expectedPassedRules)

	// 3. remove "scannedAt" attribute, which is a timestamp from actual summary
	sessionOutputSummary = sessionEngineOutput.ViolationStore.Summary
	removeTimestampAndResourcePath(&sessionOutputSummary)

	// 4. remove "scannedAt" from golden summary
	fileDataSummary = fileDataEngineOutput.ViolationStore.Summary
	removeTimestampAndResourcePath(&fileDataSummary)

	// 5. compare passed rules, violations, skipped violations and summary in actual and golden
	gomega.Expect(reflect.DeepEqual(sessionOutputSummary, fileDataSummary)).To(gomega.BeTrue())
	gomega.Expect(reflect.DeepEqual(actualPassedRules, expectedPassedRules)).To(gomega.BeTrue())
	gomega.Expect(reflect.DeepEqual(actualViolations, expectedViolations)).To(gomega.BeTrue())
	gomega.Expect(reflect.DeepEqual(actualSkippedViolations, expectedSkippedViolations)).To(gomega.BeTrue())
}

// removeTimestampAndResourcePath is helper func to make timestamp and resource path blank
func removeTimestampAndResourcePath(summary *results.ScanSummary) {
	summary.Timestamp = ""
	summary.ResourcePath = ""
}

// removeFileAndRoothFromViolations is helper func to make file in violations blank
func removeFileAndRoothFromViolations(v violations) {
	vs := []*results.Violation(v)

	for _, violation := range vs {
		violation.File = ""
		violation.PlanRoot = ""
	}
}

// GetAbsoluteFilePathForSarif helper for sarif path
func GetAbsoluteFilePathForSarif(resourcePath, filePath string) (string, error) {
	if !filepath.IsAbs(resourcePath) {
		resourcePath, err := filepath.Abs(resourcePath)
		if err != nil {
			zap.S().Errorf("unable to get absolute path for %s, error: %v", resourcePath, err)
			return "", err
		}
	}
	mode := utils.GetFileMode(resourcePath)
	if mode != nil && (*mode).IsDir() {
		return filepath.Join(resourcePath, filePath), nil
	}
	return resourcePath, nil
}

// CompareActualSarifOutputWithGoldenSummaryRegex compares actual string with contents of golden file passed as parameter
// ignores specified regex patterns from the actual and golden text
func CompareActualSarifOutputWithGoldenSummaryRegex(session *gexec.Session, goldenFileAbsPath string) {
	fileData, err := os.ReadFile(goldenFileAbsPath)
	if utils.IsWindowsPlatform() {
		fileData = utils.ReplaceWinNewLineBytes(fileData)
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	var sessionOutput, fileContents string

	sessionOutput = string(session.Wait().Out.Contents())

	fileContents = string(fileData)

	sessionOutput = strings.TrimSpace(sessionOutput)
	fileContents = strings.TrimSpace(fileContents)

	// replace uri from the output, it will cause issues for absolute paths
	sessionOutput = sarifURIPattern.ReplaceAllString(sessionOutput, "")
	fileContents = sarifURIPattern.ReplaceAllString(fileContents, "")

	sessionOutput = sarifVersionPattern.ReplaceAllString(sessionOutput, "")
	fileContents = sarifVersionPattern.ReplaceAllString(fileContents, "")

	sessionOutput = versionRegexPattern.ReplaceAllString(sessionOutput, "")
	fileContents = versionRegexPattern.ReplaceAllString(fileContents, "")

	gomega.Expect(sessionOutput).Should(gomega.BeIdenticalTo(fileContents))
}

// CheckSummaryForVulnerabilities is a helper function to check vulnerabilies exists
func CheckSummaryForVulnerabilities(session *gexec.Session, expectedCount int) {
	var sessionBytes []byte

	sessionBytes = session.Wait().Out.Contents()

	sessionBytes = bytes.TrimSpace(sessionBytes)

	var sessionEngineOutput policy.EngineOutput

	err := json.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	sessionOutputSummary := sessionEngineOutput.ViolationStore.Summary
	gomega.Expect(sessionOutputSummary).NotTo(gomega.BeNil())
	gomega.Expect(sessionOutputSummary.Vulnerabilities).NotTo(gomega.BeNil())
	gomega.Eventually(func() int {
		return *sessionOutputSummary.Vulnerabilities
	}).Should(gomega.BeNumerically(">=", expectedCount))
}
