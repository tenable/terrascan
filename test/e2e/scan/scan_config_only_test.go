package scan_test

import (
	"path/filepath"

	scanUtils "github.com/accurics/terrascan/test/e2e/scan"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Scan With Config Only Flag", func() {

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	var iacDir string
	var err error
	iacDir, err = filepath.Abs("../test_data/iac/aws/aws_ami_violation")

	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("scan command is run using the --config-only flag for unsupported output types", func() {
		When("output type is human readable format", func() {
			Context("it doesn't support --config-only flag", func() {
				Context("human readable output format is the default output format", func() {
					It("should result in an error and exit with status code 1", func() {
						errString := "please use yaml or json output format when using --config-only flag"
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only"}
						scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
					})
				})
			})
		})

		When("output type is xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					errString := "failed to write XML output. error: 'xml: unsupported type: output.AllResourceConfigs'"
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "xml"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})

		When("output type is junit-xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					errString := "incorrect input for JunitXML writer, supported type is policy.EngineOutput"
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "junit-xml"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Describe("scan command is run using the --config-only flag for unsupported output types", func() {
		Context("for terraform files", func() {
			When("output type is json", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "json"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "yaml"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					})
				})
			})
		})

		Context("for yaml files", func() {
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs("../test_data/iac/k8s/kubernetes_ingress_violation")
			})
			When("output type is json", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "json", "-i", "k8s"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "yaml", "-i", "k8s"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					})
				})
			})
		})
	})
})
