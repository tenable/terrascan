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
	// ScannedAtRegex is regex for 'scanned at' attribute in violations output
	ScannedAtRegex = regexp.MustCompile(`["]*[sS]canned[ _][aA]t["]*[ \t]*[:=][ \t]*["]*[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}[.][0-9]{3,6} [+-][0-9]{4} UTC["]*[,]{0,1}`)
	// FileFolderRegex is regex for 'file/folder' attribute in violations output
	FileFolderRegex = regexp.MustCompile(`["]*[fF]ile\/[fF]older["]*[ \t]*[:=][ \t]*["]*(.+)\/(.+)["]*`)
	// PackageRegex is regex for 'package' attribute in junit-xml output
	PackageRegex = regexp.MustCompile(`package=["]*(.+)\/(.+)["]*`)
	// VersionValueRegex is regex for 'value' attribute in junit-xml output (which is terrascan version)
	VersionValueRegex = regexp.MustCompile(`value="v[1][\.][0-9][\.][0-9]"`)
)

// CompareActualWithGolden compares actual string with contents of golden file path passed as parameter
func CompareActualWithGolden(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	fileData, err := ioutil.ReadFile(goldenFileAbsPath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	var sessionContents, fileContents string

	fileContents = string(fileData)
	if isStdOut {
		sessionContents = string(session.Wait().Out.Contents())
	} else {
		sessionContents = string(session.Wait().Err.Contents())
	}

	fileContents = strings.TrimSpace(fileContents)
	sessionContents = strings.TrimSpace(sessionContents)
	gomega.Expect(sessionContents).Should(gomega.BeIdenticalTo(fileContents))
}

// CompareActualWithGoldenRegex compares actual string with contents of golden file passed as parameter
// ignores specified regex patterns from the actual and golden text
func CompareActualWithGoldenRegex(session *gexec.Session, goldenFileAbsPath string, isJunitXML, isStdOut bool) {
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

	if isJunitXML {
		sessionOutput = PackageRegex.ReplaceAllString(sessionOutput, "")
		fileContents = PackageRegex.ReplaceAllString(fileContents, "")

		sessionOutput = VersionValueRegex.ReplaceAllString(sessionOutput, "")
		fileContents = VersionValueRegex.ReplaceAllString(fileContents, "")
	} else {
		sessionOutput = ScannedAtRegex.ReplaceAllString(sessionOutput, "")
		fileContents = ScannedAtRegex.ReplaceAllString(fileContents, "")

		sessionOutput = FileFolderRegex.ReplaceAllString(sessionOutput, "")
		fileContents = FileFolderRegex.ReplaceAllString(fileContents, "")
	}

	gomega.Expect(sessionOutput).Should(gomega.BeIdenticalTo(fileContents))
}

// CompareActualWithGoldenJSON compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenJSON(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := getByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := json.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = json.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenYAML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenYAML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := getByteData(session, goldenFileAbsPath, isStdOut)

	var sessionEngineOutput, fileDataEngineOutput policy.EngineOutput

	err := yaml.Unmarshal(sessionBytes, &sessionEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = yaml.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	compareSummaryAndViolations(sessionEngineOutput, fileDataEngineOutput)
}

// CompareActualWithGoldenXML compares actual data with contents of golden file passed as parameter
func CompareActualWithGoldenXML(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	sessionBytes, fileBytes := getByteData(session, goldenFileAbsPath, isStdOut)

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

// RunCommand will return session
func RunCommand(path string, outWriter, errWriter io.Writer, args ...string) *gexec.Session {
	cmd := exec.Command(path, args...)
	session, err := gexec.Start(cmd, outWriter, errWriter)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return session
}

// getByteData is a helper function to get data in byte slice from session and golden file
func getByteData(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) ([]byte, []byte) {
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
