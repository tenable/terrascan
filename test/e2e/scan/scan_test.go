package scan_test

import (
	"io"
	"path/filepath"

	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	scanComand string = "scan"
)

var _ = Describe("Scan", func() {

	var session *gexec.Session
	var terrascanBinaryPath string

	var outWriter, errWriter io.Writer

	BeforeSuite(func() {
		terrascanBinaryPath = helper.GetTerrascanBinaryPath()
	})

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	Describe("scan command is run with -h flag", func() {
		It("should exit with status code 0", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-h")
			Eventually(session).Should(gexec.Exit(0))
		})

		It("should print the scan help", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-h")
			goldenFileAbsPath, err := filepath.Abs("golden/scan_help.txt")
			Expect(err).NotTo(HaveOccurred())
			helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
		})
	})
})
