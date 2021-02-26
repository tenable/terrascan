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

var _ = Describe("Scan is run for k8s files", func() {

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	Context("scan iac files violating k8s policies", func() {
		var policyDir, iacDir string
		var err error

		policyDir, err = filepath.Abs("../test_data/policies/")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		iacDir, err = filepath.Abs("../test_data/iac/k8s/kubernetes_ingress_violation")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("iac type k8s is not default", func() {
			When("k8s files are scanned but iac type is not specified", func() {
				It("should print error related to terraform files not being present", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-d", iacDir)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
					helper.ContainsErrorSubString(session, "has no terraform config files")
				})
			})
		})

		Context("iac type is specified as k8s", func() {
			It("should scan and display violations in human output format", func() {
				scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir}
				scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_human.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})

			When("-v flag is used for verbose output", func() {
				It("should display verbose output for human output format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-v"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_human_verbose.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_json.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_yaml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_xml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is junit-xml", func() {
				It("should display violations in junit-xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "junit-xml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/k8s_scans/k8s/kubernetes_ingress_violations/kubernetes_ingress_junit_xml.txt", helper.ExitCodeThree, true, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})
})
