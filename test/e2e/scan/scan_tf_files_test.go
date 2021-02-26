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

const (
	backwardsCompatibilityWarningMessage = "There may be a few breaking changes while working with terraform v0.12 files. For further information, refer to https://github.com/accurics/terrascan/releases/v1.3.0"
)

var _ = Describe("Scan is run for terraform files", func() {

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	Context("terraform is the default iac type", func() {

		var iacDir string
		var err error
		iacDir, err = filepath.Abs("../test_data/iac/aws/aws_ami_violation")

		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		When("iac type is not specified, terraform iac provider is used", func() {
			It("should scan successfully and exit with status code 3", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-d", iacDir)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
			})
		})

		Context("default iac version for terraform is v14", func() {
			When("iac version is v12", func() {
				It("terrascan should display the warning message related to version", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-d", iacDir, "--iac-version", "v12")
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					helper.ContainsErrorSubString(session, backwardsCompatibilityWarningMessage)
				})
			})

			When("iac version is v13", func() {
				It("terrascan should not display the warning message related to version", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-d", iacDir, "--iac-version", "v13")
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					helper.DoesNotContainsErrorSubString(session, backwardsCompatibilityWarningMessage)
				})
			})
		})
	})

	Describe("scan iac files violating aws policies against test policy set", func() {
		var policyDir, iacDir string
		var err error

		policyDir, err = filepath.Abs("../test_data/policies/")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		iacDir, err = filepath.Abs("../test_data/iac/aws/aws_ami_violation")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("iac file violates aws_ami policy", func() {
			It("should scan and display violations in human output format", func() {
				scanArgs := []string{"-p", policyDir, "-d", iacDir}
				scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_human.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})

			When("-v flag is used for verbose output", func() {
				It("should display verbose output for human output format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-v"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_human_verbose.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_json.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_yaml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_xml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is junit-xml", func() {
				It("should display violations in junit-xml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "junit-xml"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_ami_violations/aws_ami_violation_junit_xml.txt", helper.ExitCodeThree, true, true, outWriter, errWriter, scanArgs...)
				})
			})

			Context("policy path with only aws policies is supplied", func() {
				JustBeforeEach(func() {
					policyDir, err = filepath.Abs("../test_data/policies/aws")
				})
				It("should display violations", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/scanned_with_only_aws_policies.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			Context("policy path with no aws policies is supplied", func() {
				JustBeforeEach(func() {
					policyDir, err = filepath.Abs("../test_data/policies/azure")
				})
				It("should not display any violations and exit with status code 0", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir}
					scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/scanned_with_no_aws_policies.txt", helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("iac file violates aws_db_instance policy", func() {
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs("../test_data/iac/aws/aws_db_instance_violation")
				policyDir, err = filepath.Abs("../test_data/policies/")
			})

			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanCommandAndAssertJSONOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_db_instance_violations/aws_db_instance_json.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanCommandAndAssertYAMLOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_db_instance_violations/aws_db_instance_yaml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanCommandAndAssertXMLOutput(terrascanBinaryPath, "golden/terraform_scans/aws/aws_db_instance_violations/aws_db_instance_xml.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})
})
