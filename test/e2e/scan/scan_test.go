/*
    Copyright (C) 2020 Accurics, Inc.

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
	"fmt"
	"io"
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
	session              *gexec.Session
	terrascanBinaryPath  string
	iacRootRelPath       = filepath.Join("..", "test_data", "iac")
	awsIacRelPath        = filepath.Join(iacRootRelPath, "aws")
	k8sIacRelPath        = filepath.Join(iacRootRelPath, "k8s")
	policyRootRelPath    = filepath.Join("..", "test_data", "policies")
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
			scanArgs := []string{scanUtils.ScanCommand, "-h"}
			scanUtils.RunScanAndAssertGoldenOutput(terrascanBinaryPath, filepath.Join("..", "help", "golden", "help_scan.txt"), helper.ExitCodeZero, true, outWriter, errWriter, scanArgs...)
		})
	})

	Describe("typo in the scan command, eg: scna", func() {
		scanTypo := "scna"

		It("should print scan command suggestion and exit with status code 1", func() {
			scanArgs := []string{scanTypo}
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
			Eventually(session).Should(gexec.Exit(helper.ExitCodeOne))
			goldenFileAbsPath, err := filepath.Abs(filepath.Join("golden", "scan_typo_help.txt"))
			Expect(err).NotTo(HaveOccurred())
			helper.CompareActualWithGolden(session, goldenFileAbsPath, false)
		})
	})

	Describe("scan command is run without any flags", func() {
		Context("by default, terrascan will scan for terraform files in the working directory", func() {
			Context("no tf files are present in the working directory", func() {
				It("should error out as no terraform files are present in working directory", func() {
					scanArgs := []string{scanUtils.ScanCommand}
					errString := "has no terraform config files"
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})

			Context("tf files are present in the working directory", func() {
				It("should scan the directory, return results and exit with status code 3", func() {
					workDir, err := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
					Expect(err).NotTo(HaveOccurred())

					scanArgs := []string{scanUtils.ScanCommand}
					session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
				})

				When("tf file present in the dir has no violations", func() {
					Context("when there are no violations, terrascan exits with status code 0", func() {
						It("should scan the directory and exit with status code 0", func() {
							workDir, err := filepath.Abs(filepath.Join(awsIacRelPath, "aws_db_instance_violation"))
							Expect(err).NotTo(HaveOccurred())

							// set a policy path that doesn't have any s3 bucket policies
							policyDir, err := filepath.Abs(filepath.Join(policyRootRelPath, "k8s"))
							Expect(err).NotTo(HaveOccurred())

							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir}
							session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
						})
					})
				})
			})
		})
	})

	Describe("terrascan scan command is run with -d and -f flag", func() {
		workDir, err := os.Getwd()
		It("should not get an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		Context("invalid path is supplied to -d or -f flag", func() {

			invalidPath := "invalid/path"
			When("supplied with -d flag", func() {
				It("should error out and exit with status code 1", func() {
					errString := fmt.Sprintf("directory '%s' does not exist", filepath.Join(workDir, invalidPath))
					scanArgs := []string{scanUtils.ScanCommand, "-d", invalidPath}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})

			When("supplied with -f flag", func() {
				It("should error out and exit with status code 1", func() {
					errString := fmt.Sprintf("file '%s' does not exist", filepath.Join(workDir, invalidPath))
					scanArgs := []string{scanUtils.ScanCommand, "-f", invalidPath}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("-d flag is supplied with a valid file path", func() {
			It("should error out and exit with status code 1", func() {
				validAbsFilePath, err := filepath.Abs(filepath.Join("golden", "scan_typo_help.txt"))
				Expect(err).NotTo(HaveOccurred())

				errString := fmt.Sprintf("input path '%s' is not a valid directory", validAbsFilePath)
				scanArgs := []string{scanUtils.ScanCommand, "-d", validAbsFilePath}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		Context("-f flag is supplied with a valid dir path", func() {
			It("should error out and exit with status code 1", func() {
				validAbsDirPath, err := filepath.Abs("golden")
				Expect(err).NotTo(HaveOccurred())

				errString := fmt.Sprintf("input path '%s' is not a valid file", validAbsDirPath)
				scanArgs := []string{scanUtils.ScanCommand, "-f", validAbsDirPath}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})
	})

	Describe("scan is run with unsupported iac type or version", func() {
		errString := "iac type or version not supported"
		When("-i flag is supplied with unsupported iac type", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-i", "test"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		When("--iac-version flag is supplied invalid version", func() {
			Context("default iac type is terraform", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "--iac-version", "test"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})

		Context("iac type is valid but iac version isn't", func() {
			When("iac type is k8s and supplied version is invalid", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-i", "k8s", "--iac-version", "test"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})

			When("iac type is helm and supplied version is invalid", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-i", "helm", "--iac-version", "test"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})

			When("iac type is kustomize and supplied version is invalid", func() {
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-i", "kustomize", "--iac-version", "test"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Describe("scan is run with -p (policy path) flag", func() {
		invalidPolicyPath := "path/policy/invalid"
		errString1 := "failed to initialize OPA policy engine"
		errString2 := fmt.Sprintf("%s: no such file or directory", invalidPolicyPath)

		When("supplied policy path doesn't exist", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-p", invalidPolicyPath}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString1)
				helper.ContainsErrorSubString(session, errString2)
			})
		})

		Context("multiple policy paths can be supplied", func() {
			When("one of the supplied policy path is invalid", func() {
				It("should error out and exit with staus code 1", func() {
					validPolicyPath, err := filepath.Abs(filepath.Join("..", "test_data"))
					Expect(err).NotTo(HaveOccurred())
					scanArgs := []string{scanUtils.ScanCommand, "-p", validPolicyPath, "-p", invalidPolicyPath}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
					helper.ContainsErrorSubString(session, errString1)
					helper.ContainsErrorSubString(session, errString2)
				})
			})

			When("multiple valid policy paths are supplied", func() {
				Context("workdir contains k8s files", func() {
					awsPolicyRelPath := filepath.Join(policyRootRelPath, "aws")
					validPolicyPath1, err1 := filepath.Abs(filepath.Join(awsPolicyRelPath, "aws_ami"))
					validPolicyPath2, err2 := filepath.Abs(filepath.Join(awsPolicyRelPath, "aws_db_instance"))
					workDirPath, err3 := filepath.Abs(filepath.Join(policyRootRelPath, "k8s"))

					It("should not error out", func() {
						Expect(err1).NotTo(HaveOccurred())
						Expect(err2).NotTo(HaveOccurred())
						Expect(err3).NotTo(HaveOccurred())
					})

					Context("default iac type is terraform", func() {
						It("should scan with the policies and exit with status code 1", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", validPolicyPath1, "-p", validPolicyPath2}
							session = helper.RunCommandDir(terrascanBinaryPath, workDirPath, outWriter, errWriter, scanArgs...)
							// exits with status code 1, because the work dir has k8s iac file and default iac type is terraform
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
						})
					})

					Context("iac type is k8s", func() {
						It("should scan with the policies and display results", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", validPolicyPath1, "-p", validPolicyPath2, "-i", "k8s"}
							session = helper.RunCommandDir(terrascanBinaryPath, workDirPath, outWriter, errWriter, scanArgs...)
							// exits with status code 0, because no violations would be reported,
							// the work dir has k8s iac file and supplied policies are for tf files
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
						})
					})
				})
			})
		})
	})
})
