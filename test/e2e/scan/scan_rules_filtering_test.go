/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package scan_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	scanUtils "github.com/tenable/terrascan/test/e2e/scan"
	"github.com/tenable/terrascan/test/helper"
)

var (
	terrascanConfigEnvName      string = "TERRASCAN_CONFIG"
	severityLevelIncorrectError string = "severity level not supported"
	categoryIncorrectError      string = "category not supported"
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

	iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_db_instance_violation"))
	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	policyDir, err = filepath.Abs(policyRootRelPath)
	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	ruleFilterGoldenRelPath := filepath.Join("golden", "rules_filtering")
	Describe("rule filtering via command line options", func() {

		Context("--scan-rules options is used", func() {
			Context("single rule is specified via --scan-rules option", func() {
				It("should scan only the rules specified in --scan-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041"
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "scan_single_rule.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
			Context("multiple rules are specified via --scan-rules option", func() {
				It("should scan only the rules specified in --scan-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041,AWS.AWS RDS.NS.High.0101,AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "scan_multiple_rules.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("--skip-rules options is used", func() {
			Context("single rule is specified via --skip-rules option", func() {
				It("should skip the rule specified in --skip-rules option", func() {
					skipRules := "AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--skip-rules", skipRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "skip_single_rule.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
			Context("multiple rules are specified via --skip-rules option", func() {
				It("should skip the rules specified in --skip-rules option", func() {
					skipRules := "AWS.RDS.DS.High.1041,AWS.RDS.DataSecurity.High.0414,AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--skip-rules", skipRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "skip_multiple_rules.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("both --scan-rules and --skip-rules are specified", func() {
			Context("single rule is specified via --skip-rules option", func() {
				It("should scan and skip the rules as specified with --scan-rules and --skip-rules option", func() {
					scanRules := "AWS.RDS.DS.High.1041,AWS.AWS RDS.NS.High.0101,AWS.RDS.DataSecurity.High.0577"
					skipRules := "AWS.RDS.DataSecurity.High.0577"
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--skip-rules", skipRules, "--scan-rules", scanRules}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "scan_and_skip_rules.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
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
					iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				})

				JustAfterEach(func() {
					iacDir = oldIacDir
				})
				Context("severity leve specified is 'low'", func() {
					Context("iac file has only medium severity violations", func() {
						It("should report the violations and exit with status code 3", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--severity", "low", "-i", "terraform"}
							session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
						})
					})
				})
				Context("severity level specified is 'high'", func() {
					Context("iac files has only medium severity violations", func() {
						// there would not no violations but directory scan errors would be present due to all iac scan
						It("should not report any violation and exit with status code 4", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--severity", "high"}
							session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeFour))
						})
					})
				})
			})
		})

		Context("category is specified", func() {
			When("category specified is invalid", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--categories", "test"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
					helper.ContainsErrorSubString(session, categoryIncorrectError)
				})
			})

			When("valid category is specified", func() {
				oldIacDir := iacDir
				JustAfterEach(func() {
					iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				})

				JustAfterEach(func() {
					iacDir = oldIacDir
				})
				Context("category specified is 'COMPLIANCE VALIDATION'", func() {
					Context("iac file has violations with only 'DATA PROTECTION' category", func() {
						It("should not report any violation and exit with status code 0", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-o", "json", "--categories", "COMPLIANCE VALIDATION"}
							session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
							// summary would contain directory scan errors due to all iac scan
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeFour))
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
					configFileAbsPath, err := filepath.Abs(filepath.Join("config", "scan_and_skip_rules.toml"))
					Expect(err).NotTo(HaveOccurred())
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "-c", configFileAbsPath}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "scan_and_skip_rules.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("config file is specified using TERRASCAN_CONFIG env variable", func() {
			Context("both scan and skip rules are specified", func() {
				JustBeforeEach(func() {
					os.Setenv(terrascanConfigEnvName, filepath.Join("config", "scan_and_skip_rules.toml"))
				})
				JustAfterEach(func() {
					os.Setenv(terrascanConfigEnvName, "")
				})
				It("should scan and skip the rules as specified with --scan-rules and --skip-rules option", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(ruleFilterGoldenRelPath, "scan_and_skip_rules.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			Context("invalid severity is specified in config file", func() {
				It("should error out and exit with status code 1", func() {
					configFileAbsPath, err := filepath.Abs(filepath.Join("config", "invalid_severity.toml"))
					Expect(err).NotTo(HaveOccurred())

					scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-c", configFileAbsPath}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, severityLevelIncorrectError, outWriter, errWriter, scanArgs...)
				})
			})

			Context("invalid category is specified in config file", func() {
				It("should error out and exit with status code 3", func() {
					configFileAbsPath, err := filepath.Abs(filepath.Join("config", "invalid_category.toml"))
					Expect(err).NotTo(HaveOccurred())

					scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir, "-d", iacDir, "-c", configFileAbsPath}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, categoryIncorrectError, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Describe("resource specific rule skipping", func() {
		resourceSkipGoldenRelPath := filepath.Join("golden", "resource_skipping")
		resourceSkipIacRelPath := filepath.Join(iacRootRelPath, "resource_skipping")
		Context("resource skipping in tf files", func() {
			oldIacDir := iacDir
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs(filepath.Join(resourceSkipIacRelPath, "terraform"))
			})
			JustAfterEach(func() {
				iacDir = oldIacDir
			})
			It("should display violations, skipped violations and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourceSkipGoldenRelPath, "terraform_file_resource_skipping.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})

		Context("resource skipping in k8s files", func() {
			oldIacDir := iacDir
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs(filepath.Join(resourceSkipIacRelPath, "kubernetes"))
			})
			JustAfterEach(func() {
				iacDir = oldIacDir
			})

			// the iac file has only one resource with one violation, which is skipped.
			// hence, the exit code is 0
			It("should display skipped violations and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-d", iacDir, "-i", "k8s", "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourceSkipGoldenRelPath, "kubernetes_file_resource_skipping.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})

		Context("resource skipping in docker files", func() {
			iacScanDir := filepath.Join(resourceSkipIacRelPath, "docker")
			// the iac file has only one resource with one violation, which is skipped.
			// hence, the exit code is 0
			It("should display skipped violations and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-d", iacScanDir, "-i", "docker", "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourceSkipGoldenRelPath, "dockerfile_resource_skipping.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
	})
	Describe("resource specific rule prioritising", func() {
		resourcePrioritisingGoldenRelPath := filepath.Join("golden", "resource_prioritising")
		resourcePrioritisingIacRelPath := filepath.Join(iacRootRelPath, "resource_prioritising")
		Context("resource max severity set to Low in tf files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set", "terraform")
			It("should display violations with change priority to Low for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set", "terraform", "terraform_file_setting_max_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High in tf files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_severity_set", "terraform")
			It("should display violations with change priority to High for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "terraform", "terraform_file_setting_min_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource max severity set to none in tf files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set_none", "terraform")
			It("should skip all violations for the resource and exit with status code 0 since only one resource is in tf file", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set_none", "terraform", "terraform_file_setting_max_severity_none.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High and max severity set to Low in tf files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_max_both_severity_set", "terraform")
			It("should display violations with change priority to High as specified by min severity for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "terraform", "terraform_file_setting_min_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High and skip rule provided in tf files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_severity_with_skip_rule", "terraform")
			It("should display skipped violations and violations with change priority to High as specified by min severity for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_with_skip_rule", "terraform", "terraform_file_setting_min_severity_with_skip_rule.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})

		// k8s file tests
		Context("resource max severity set to Low in k8s files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set", "k8s")
			It("should display violations with change priority to Low for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "k8s", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set", "k8s", "k8s_file_setting_max_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High in k8s files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_severity_set", "k8s")
			It("should display violations with change priority to High for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "k8s", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "k8s", "k8s_file_setting_min_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource max severity set to none in k8s files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set_none", "k8s")
			It("should skip all violations for the resource and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-i", "k8s", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set_none", "k8s", "k8s_file_setting_max_severity_none.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High and max severity set to Low k8s files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_max_both_severity_set", "k8s")
			It("should display violations with change priority to High as specified by min severity for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "k8s", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "k8s", "k8s_file_setting_min_severity.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High and skip rule provided in k8s files", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_severity_with_skip_rule", "k8s")
			// since only one violation is present for resource
			It("should display skipped violations with change priority to High as specified by min severity for the resource and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-i", "k8s", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_with_skip_rule", "k8s", "k8s_file_setting_min_severity_with_skip_rule.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource max severity set to Low in dockerfile", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set", "docker")
			It("should give violations for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "docker", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set", "docker", "dockerfile_max_severity_low.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High in dockerfile", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_severity_set", "docker")
			It("should give violations for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "docker", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "docker", "dockerfile_min_severity_high.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource max severity set to None in dockerfile", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "max_severity_set_none", "docker")
			It("should skip violations for the resource and exit with status code 0", func() {
				scanArgs := []string{"-p", policyDir, "-i", "docker", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "max_severity_set_none", "docker", "dockerfile_max_severity_none.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
			})
		})
		Context("resource min severity set to High and max severity set to Low in dockerfile", func() {
			iacDir := filepath.Join(resourcePrioritisingIacRelPath, "min_max_both_severity_set", "docker")
			It("should give violations for the resource and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "docker", "-d", iacDir, "-o", "json"}
				scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(resourcePrioritisingGoldenRelPath, "min_severity_set", "docker", "dockerfile_min_severity_high.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})
		})
	})
})
