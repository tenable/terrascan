package scan_test

import (
	"os"
	"path/filepath"

	scanUtils "github.com/accurics/terrascan/test/e2e/scan"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	terrascanConfigEnvName      string = "TERRASCAN_CONFIG"
	severityLevelIncorrectError string = "severity level not supported"
)

var _ = Describe("Scan command with rule filtering options", func() {

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	var policyDir, iacDir string
	var err error

	iacDir, err = filepath.Abs("../test_data/iac/aws/aws_db_instance_violation")
	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	policyDir, err = filepath.Abs("../test_data/policies/")
	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("rule filtering via command line options", func() {

		Context("--scan-rules options is used", func() {
			Context("single rule is specified via --scan-rules option", func() {
				It("should scan only the rules specified in --scan-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041"
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/scan_single_rule.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
			Context("multiple rules are specified via --scan-rules option", func() {
				It("should scan only the rules specified in --scan-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041,AWS.AWS RDS.NS.High.0101,AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/scan_multiple_rules.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("--skip-rules options is used", func() {
			Context("single rule is specified via --skip-rules option", func() {
				It("should skip the rule specified in --skip-rules option", func() {
					skipRules := "AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "--skip-rules", skipRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/skip_single_rule.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
			Context("multiple rules are specified via --skip-rules option", func() {
				It("should skip the rules specified in --skip-rules option", func() {
					skipRules := "AWS.RDS.DS.High.1041,AWS.RDS.DataSecurity.High.0414,AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "--skip-rules", skipRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/skip_multiple_rules.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("both --scan-rules and --skip-rules are specified", func() {
			Context("single rule is specified via --skip-rules option", func() {
				It("should scan and skip the rules as specified with --scan-rules and --skip-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041,AWS.AWS RDS.NS.High.0101,AWS.RDS.DataSecurity.High.0577"
					skipRules := "AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "--skip-rules", skipRules, "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/scan_and_skip_rules.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("severity level is specified", func() {
			When("severity specified is invalid", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--severity", "test"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
					helper.ContainsErrorSubString(session, severityLevelIncorrectError)
				})
			})

			When("valid severity level is specified", func() {
				oldIacDir := iacDir
				JustBeforeEach(func() {
					iacDir, err = filepath.Abs("../test_data/iac/aws/aws_ami_violation")
				})

				JustAfterEach(func() {
					iacDir = oldIacDir
				})
				Context("severity leve specified is 'low'", func() {
					Context("iac file has only medium severity violations", func() {
						It("should report the violations and exit with status code 3", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--severity", "low"}
							session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
						})
					})
				})
				Context("severity leve specified is 'high'", func() {
					Context("iac files has only medium severity violations", func() {
						It("should not report any violation and exit with status code 0", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--severity", "high"}
							session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
						})
					})
				})
			})
		})
	})

	Describe("rule filtering via config file", func() {
		Context("config file is specified using -c flag", func() {
			Context("both scan and skip rules are specified", func() {
				It("should scan and skip the rules as specified with --scan-rules and --skip-rules option", func() {
					configFileAbsPath, err := filepath.Abs("config/scan_and_skip_rules.toml")
					Expect(err).NotTo(HaveOccurred())
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json", "-c", configFileAbsPath}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/rules_filtering/scan_and_skip_rules.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("config file is specified using TERRASCAN_CONFIG env variable", func() {
			Context("both scan and skip rules are specified", func() {
				JustBeforeEach(func() {
					os.Setenv(terrascanConfigEnvName, "config/scan_and_skip_rules.toml")
				})
				JustAfterEach(func() {
					os.Setenv(terrascanConfigEnvName, "")
				})
				It("should scan and skip the rules as specified with --scan-rules and --skip-rules option", func() {
					Skip("skipping this test due to https://github.com/accurics/terrascan/issues/570, should be implemented when fixed")
				})
			})

			Context("invalid severity is specified in config file", func() {
				It("should error out and exit with status code 1", func() {
					configFileAbsPath, err := filepath.Abs("config/invalid_severity.toml")
					Expect(err).NotTo(HaveOccurred())

					scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-c", configFileAbsPath}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, severityLevelIncorrectError, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Describe("resource specific rule skipping", func() {
		Context("resource skipping in tf files", func() {
			oldIacDir := iacDir
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs("../test_data/iac/resource_skipping/terraform")
			})
			JustAfterEach(func() {
				iacDir = oldIacDir
			})
			It("should display violations, skipped violations and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/resource_skipping/terraform_file_resource_skipping.txt", helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})

		Context("resource skipping in k8s files", func() {
			oldIacDir := iacDir
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs("../test_data/iac/resource_skipping/kubernetes")
			})
			JustAfterEach(func() {
				iacDir = oldIacDir
			})

			// the iac file has only one resource with one violation, which is skipped.
			// hence, the exit code is 0
			It("should display skipped violations and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-d", iacDir, "-i", "k8s", "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, "golden/resource_skipping/kubernetes_file_resource_skipping.txt", helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
	})
})
