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
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	scanUtils "github.com/tenable/terrascan/test/e2e/scan"
	"github.com/tenable/terrascan/test/helper"
)

const (
	backwardsCompatibilityWarningMessage = "There may be a few breaking changes while working with terraform v0.12 files. For further information, refer to https://github.com/tenable/terrascan/releases/v1.3.0"
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
		iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))

		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		When("terraform iac provider is used", func() {
			It("should scan successfully and exit with status code 3", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "-i", "terraform"}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
			})
		})

		Context("default iac version for terraform is v14", func() {
			When("iac version is v12", func() {
				It("terrascan should display the warning message related to version", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "-i", "terraform", "--iac-version", "v12"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeThree, scanUtils.ScanTimeout, backwardsCompatibilityWarningMessage, outWriter, errWriter, scanArgs...)
				})
			})

			When("iac version is v13", func() {
				It("terrascan should not display the warning message related to version", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "-i", "terraform", "--iac-version", "v13"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
					helper.DoesNotContainsErrorSubString(session, backwardsCompatibilityWarningMessage)
				})
			})
		})
	})

	Describe("scan iac files violating aws policies against test policy set", func() {
		var policyDir, iacDir string
		var err error

		tfGoldenRelPath := filepath.Join("golden", "terraform_scans")
		tfAwsAmiGoldenRelPath := filepath.Join(tfGoldenRelPath, "aws", "aws_ami_violations")

		policyDir, err = filepath.Abs(policyRootRelPath)
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("iac file violates aws_ami policy", func() {
			It("should scan and display violations in human output format", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir}
				scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})

			When("-v flag is used for verbose output", func() {
				It("should display verbose output for human output format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-v"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human_verbose.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_json.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is sarif", func() {
				It("should display violations in sarif format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "sarif"}
					scanUtils.RunScanAndAssertGoldenSarifOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_sarif.txt"), helper.ExitCodeThree, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is json and no iac type is specified", func() {
				Context("when iac type is not specified and a directory is specified, it will be scanned will all iac providers", func() {
					It("should display violations in json format, and should have iac type as 'all'", func() {
						scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
						scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_json_all.txt"), helper.ExitCodeFive, false, true, outWriter, errWriter, scanArgs...)
					})
				})
			})

			When("output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_yaml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_xml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is junit-xml", func() {
				It("should display violations in junit-xml format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "junit-xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_junit_xml.txt"), helper.ExitCodeThree, true, true, outWriter, errWriter, scanArgs...)
				})
			})

			Context("policy path with only aws policies is supplied", func() {
				JustBeforeEach(func() {
					policyDir, err = filepath.Abs(filepath.Join(policyRootRelPath, "aws"))
				})
				It("should display violations", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfGoldenRelPath, "scanned_with_only_aws_policies.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			Context("policy path with no aws policies is supplied", func() {
				JustBeforeEach(func() {
					policyDir, err = filepath.Abs(filepath.Join(policyRootRelPath, "k8s"))
				})
				It("should not display any violations and exit with status code 0", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfGoldenRelPath, "scanned_with_no_aws_policies.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		tfAwsDBInstanceGoldenRelPath := filepath.Join(tfGoldenRelPath, "aws", "aws_db_instance_violations")

		Context("iac file violates aws_db_instance policy", func() {
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_db_instance_violation"))
				policyDir, err = filepath.Abs(policyRootRelPath)
			})

			When("output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(tfAwsDBInstanceGoldenRelPath, "aws_db_instance_json.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanAndAssertYAMLOutput(terrascanBinaryPath, filepath.Join(tfAwsDBInstanceGoldenRelPath, "aws_db_instance_yaml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanAndAssertXMLOutput(terrascanBinaryPath, filepath.Join(tfAwsDBInstanceGoldenRelPath, "aws_db_instance_xml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("--show-passed option is used", func() {
				It("should display passed rules in the output", func() {
					scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "-o", "json", "--show-passed"}
					scanUtils.RunScanAndAssertJSONOutput(terrascanBinaryPath, filepath.Join(tfAwsDBInstanceGoldenRelPath, "aws_db_instance_json_show_passed.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("when --non-recursive flag is not used, all sub folders will be scanned in the specified directory", func() {
			When("output type is json", func() {
				It("should display violations in json format", func() {
					iacDir := filepath.Join(iacRootRelPath, "terraform_recursive")
					scanArgs := []string{"-i", "terraform", "-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_json_recursive.txt"), helper.ExitCodeFive, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
		Context("when --use-terraform-cache flag is used, all remote modules are refered from terraform cache", func() {
			When("when --use-terraform-cache is set with output format json", func() {
				It("should not display any violations", func() {
					iacDir := filepath.Join(iacRootRelPath, "terraform_cache_use_in_scan")
					scanArgs := []string{"-i", "terraform", "-p", policyDir, "-d", iacDir, "-o", "json", "--use-terraform-cache"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfGoldenRelPath, "terraform_cache_use_in_scan", "terraform_cache_use_in_scan_result.txt"), helper.ExitCodeZero, false, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})
})
