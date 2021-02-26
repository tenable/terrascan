package helper

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"gopkg.in/yaml.v3"
)

var (
	// ScannedAt is regex for 'scanned at' attribute in violations output
	ScannedAt = regexp.MustCompile(`["]*[sS]canned[ _][aA]t["]*[ \t]*[:=][ \t]*["]*[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}[.][0-9]{1,9} [+-][0-9]{4} UTC["]*[,]{0,1}`)

	// FileFolder is regex for 'file/folder' attribute in violations output
	FileFolder = regexp.MustCompile(`["]*[fF]ile[\/_][fF]older["]*[ \t]*[:=][ \t]*["]*(.+)\/(.+)["]*`)

	// File is regex for 'file' attribute in violations output
	File = regexp.MustCompile(`["]*[fF]ile["]*[ \t]*[:=][ \t]*["]*(.+)\/(.+)["]*`)

	// Package is regex for 'package' attribute in junit-xml output
	Package = regexp.MustCompile(`package=["]*(.+)\/(.+)["]*`)

	// Classname is regex for 'package' attribute in junit-xml output
	Classname = regexp.MustCompile(`classname=["]*(.+)\/(.+)["]*`)

	// VersionValue is regex for 'value' attribute in junit-xml output (which is terrascan version)
	VersionValue = regexp.MustCompile(`value="v[1][\.][0-9][\.][0-9]"`)

	// SourceRegex is regex for 'file/folder' attribute in violations output
	SourceRegex = regexp.MustCompile(`["]*source["]*[ \t]*[:][ \t]*["]*(.+)\/(.+)["]*`)
)

// CompareActualWithGolden compares actual string with contents of golden file path passed as parameter
func CompareActualWithGolden(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)
	gomega.Expect(string(sessionBytes)).Should(gomega.Equal(string(fileBytes)))
}

// CompareActualWithGoldenConfigOnlyRegex compares actual string with contents of golden file path passed as parameter
func CompareActualWithGoldenConfigOnlyRegex(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)
	sessionBytes = SourceRegex.ReplaceAll(sessionBytes, []byte(""))
	fileBytes = SourceRegex.ReplaceAll(fileBytes, []byte(""))
	gomega.Expect(string(sessionBytes)).Should(gomega.Equal(string(fileBytes)))
}

// CompareActualWithGoldenSummaryRegex compares actual string with contents of golden file passed as parameter
// ignores specified regex patterns from the actual and golden text
func CompareActualWithGoldenSummaryRegex(session *gexec.Session, goldenFileAbsPath string, isJunitXML, isStdOut bool) {
	fileData, err := ioutil.ReadFile(goldenFileAbsPath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	var sessionOutput, fileContents string

	if isStdOut {
		sessionOutput = string(session.Wait().Out.Contents())
	} else {
		sessionOutput = string(session.Wait().Err.Contents())
	}

	fileContents = string(fileData)

	sessionOutput = strings.TrimSpace(sessionOutput)
	fileContents = strings.TrimSpace(fileContents)

	// replace file from the output, it will cause issues for absolute paths
	sessionOutput = File.ReplaceAllString(sessionOutput, "")
	fileContents = File.ReplaceAllString(fileContents, "")

	if isJunitXML {
		sessionOutput = Package.ReplaceAllString(sessionOutput, "")
		fileContents = Package.ReplaceAllString(fileContents, "")

		sessionOutput = Classname.ReplaceAllString(sessionOutput, "")
		fileContents = Classname.ReplaceAllString(fileContents, "")

		sessionOutput = VersionValue.ReplaceAllString(sessionOutput, "")
		fileContents = VersionValue.ReplaceAllString(fileContents, "")
	} else {
		sessionOutput = ScannedAt.ReplaceAllString(sessionOutput, "")
		fileContents = ScannedAt.ReplaceAllString(fileContents, "")

		sessionOutput = FileFolder.ReplaceAllString(sessionOutput, "")
		fileContents = FileFolder.ReplaceAllString(fileContents, "")
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

	compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenYAML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenYAML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := yaml.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = yaml.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenXML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenXML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := GetByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := xml.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = xml.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// ContainsErrorSubString will assert if error string is part of error output
func ContainsErrorSubString(session *gexec.Session, errSubString string) {
	gomega.Expect(string(session.Wait().Err.Contents())).Should(gomega.ContainSubstring(errSubString))
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
	fileBytes, err := ioutil.ReadFile(goldenFileAbsPath)
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

// compareSummaryAndViolations is a helper function to compare actual and expected, summary and violations
func compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput policy.EngineOutput) {
	var sessionOutputSummary, fileDataSummary results.ScanSummary
	var actualViolations, expectedViolations violations

	actualViolations = sessionEngineOutput.ViolationStore.Violations
	expectedViolations = fileDataEngineOutput.ViolationStore.Violations

	sort.Sort(actualViolations)
	sort.Sort(expectedViolations)

	sessionOutputSummary = sessionEngineOutput.ViolationStore.Summary
	removeTimestampAndResourcePath(&sessionOutputSummary)

	fileDataSummary = fileDataEngineOutput.ViolationStore.Summary
	removeTimestampAndResourcePath(&fileDataSummary)

	gomega.Expect(reflect.DeepEqual(sessionOutputSummary, fileDataSummary)).To(gomega.BeTrue())
	gomega.Expect(reflect.DeepEqual(actualViolations, expectedViolations)).To(gomega.BeTrue())
}

// removeTimestampAndResourcePath is helper func to make timestamp and resource path blank
func removeTimestampAndResourcePath(summary *results.ScanSummary) {
	summary.Timestamp = ""
	summary.ResourcePath = ""
}
