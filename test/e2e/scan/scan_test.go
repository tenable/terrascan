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
	scanComand           string = "scan"
	session              *gexec.Session
	terrascanBinaryPath  string
	outWriter, errWriter io.Writer
)

var _ = Describe("Scan", func() {

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
		It("should print help for scan and exit with status code 0", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-h")
			Eventually(session).Should(gexec.Exit(0))
			sessionBytes, fileBytes := helper.GetByteData(session, "golden/scan_help.txt", true)
			sessionBytes = helper.TerraformIacVersion.ReplaceAll(sessionBytes, []byte(""))
			fileBytes = helper.TerraformIacVersion.ReplaceAll(fileBytes, []byte(""))
			Expect(string(sessionBytes)).Should(Equal(string(fileBytes)))
		})
	})

	Describe("typo in the scan command, eg: scna", func() {
		It("should exit with status code 1", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "scna")
			Eventually(session).Should(gexec.Exit(1))
		})

		It("should print scan command suggestion", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "scna")
			goldenFileAbsPath, err := filepath.Abs("golden/scan_typo_help.txt")
			Expect(err).NotTo(HaveOccurred())
			helper.CompareActualWithGolden(session, goldenFileAbsPath, false)
		})
	})

	Describe("scan command is run without any flags", func() {
		Context("when scan command is run without any flags, terrascan with scan for terraform files in working directory", func() {
			It("should error out as no terraform files are present in working directory", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand)
				Eventually(session, 2).Should(gexec.Exit(1))
				helper.ContainsErrorSubString(session, "has no terraform config files")
			})
		})
	})
})
