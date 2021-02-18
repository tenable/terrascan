package helper

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// CompareActualWithGolden compares
func CompareActualWithGolden(session *gexec.Session, goldenFileAbsPath string, isStdOut bool) {
	fileData, err := ioutil.ReadFile(goldenFileAbsPath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	if isStdOut {
		gomega.Expect(string(session.Wait().Out.Contents())).Should(gomega.BeIdenticalTo(string(fileData)))
	} else {
		gomega.Expect(string(session.Wait().Err.Contents())).Should(gomega.BeIdenticalTo(string(fileData)))
	}
}

// ContainsErrorSubString will assert if error string is part of error output
func ContainsErrorSubString(session *gexec.Session, errSubString string) {
	gomega.Expect(string(session.Wait().Err.Contents())).Should(gomega.ContainSubstring(errSubString))
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
