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
			goldenFileAbsPath, err := filepath.Abs("golden/scan_help.txt")
			Expect(err).NotTo(HaveOccurred())
			helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
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
		When("scan command is run without any flags, terrascan with scan for terraform files in working directory", func() {
			Context("no tf files are present in the working directory", func() {
				It("should error out as no terraform files are present in working directory", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand)
					Eventually(session, 2).Should(gexec.Exit(1))
					helper.ContainsErrorSubString(session, "has no terraform config files")
				})
			})

			Context("tf files are present in the working directory", func() {
				It("should scan the directory, return results and exit with status code 3", func() {
					workDir, err := filepath.Abs("../test_data/iac/aws/aws_ami_violation")
					Expect(err).NotTo(HaveOccurred())
					session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanComand)
					Eventually(session, 2).Should(gexec.Exit(3))
				})

				When("tf file present in the dir has no violations", func() {
					Context("when there are no violations, terrascan exits with status code 0", func() {
						It("should scan the directory and exit with status code 0", func() {
							workDir, err := filepath.Abs("../test_data/iac/aws/aws_db_instance_violation")
							Expect(err).NotTo(HaveOccurred())

							// set a policy path that doesn't have any s3 bucket policies
							policyDir, err := filepath.Abs("../test_data/policies/k8s")
							Expect(err).NotTo(HaveOccurred())

							session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanComand, "-p", policyDir)
							Eventually(session, 2).Should(gexec.Exit(0))
						})
					})
				})
			})
		})
	})

})
