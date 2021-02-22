package help

import (
	"path/filepath"

	"github.com/accurics/terrascan/test/helper"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// ValidateExitCodeAndOutput validates the exit code and output of the command
func ValidateExitCodeAndOutput(session *gexec.Session, exitCode int, relFilePath string, isStdOut bool) {
	gomega.Eventually(session).Should(gexec.Exit(exitCode))
	goldenFileAbsPath, err := filepath.Abs(relFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	helper.CompareActualWithGolden(session, goldenFileAbsPath, isStdOut)
}
