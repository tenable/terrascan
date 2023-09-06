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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/pkg/utils"
	scanUtils "github.com/tenable/terrascan/test/e2e/scan"
	"github.com/tenable/terrascan/test/helper"
)

var (
	session              *gexec.Session
	terrascanBinaryPath  string
	iacRootRelPath       = filepath.Join("..", "test_data", "iac")
	awsIacRelPath        = filepath.Join(iacRootRelPath, "aws")
	k8sIacRelPath        = filepath.Join(iacRootRelPath, "k8s")
	dockerIacRelPath     = filepath.Join(iacRootRelPath, "docker")
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

	Describe("scan command is run", func() {
		Context("when no iac type is provided, terrascan scans with all iac providers", func() {
			Context("no tf files are present in the working directory", func() {
				It("scans the directory with all iac and display results", func() {
					scanArgs := []string{scanUtils.ScanCommand}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					helper.ValidateExitCode(session, scanUtils.ScanTimeout, helper.ExitCodeFour)
				})
				Context("iac loading errors would be displayed in the output, output type is json", func() {
					It("scan the directory and display results", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-o", "json", "-p", policyRootRelPath}
						// these errors would come from terraform, helm, and kustomize iac providers
						errString1 := "kustomization.y(a)ml file not found in the directory"
						errString2 := "no helm charts found in directory"
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						helper.ContainsDirScanErrorSubString(session, errString1)
						helper.ContainsDirScanErrorSubString(session, errString2)
					})
				})
				Context("iac loading errors would be displayed in the output, output type is human readable", func() {
					It("scan the directory and display results", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-p", policyRootRelPath}
						// these errors would come from terraform, helm, and kustomize iac providers
						errString1 := "kustomization.y(a)ml file not found in the directory"
						errString2 := "no helm charts found in directory"
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						helper.ContainsDirScanErrorSubString(session, errString1)
						helper.ContainsDirScanErrorSubString(session, errString2)
					})
				})
				When("iac type is terraform and --non-recursive flag is used", func() {
					It("should error out if no terraform files are present in working directory", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-i", "terraform", "--non-recursive"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
					})
				})
			})
			Context("tf files are present in the working directory", func() {
				It("should scan the directory, return results and exit with status code 3 as there would no directory scan errors", func() {
					workDir, err := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
					Expect(err).NotTo(HaveOccurred())

					scanArgs := []string{scanUtils.ScanCommand, "-i", "terraform", "--non-recursive"}
					session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeThree))
				})

				When("tf file present in the dir has no violations", func() {
					Context("when there are no violations, but has dir scan errors, terrascan exits with status code 4", func() {
						It("should scan the directory and exit with status code 0", func() {
							workDir, err := filepath.Abs(filepath.Join(awsIacRelPath, "aws_db_instance_violation"))
							Expect(err).NotTo(HaveOccurred())

							// set a policy path that doesn't have any s3 bucket policies
							policyDir, err := filepath.Abs(filepath.Join(policyRootRelPath, "k8s"))
							Expect(err).NotTo(HaveOccurred())

							scanArgs := []string{scanUtils.ScanCommand, "-p", policyDir}
							session = helper.RunCommandDir(terrascanBinaryPath, workDir, outWriter, errWriter, scanArgs...)
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeFour))
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
			Context("default iac type is all and --iac-version would be ignored", func() {
				It("should error out and exit with status code 4", func() {
					scanArgs := []string{scanUtils.ScanCommand, "--iac-version", "test"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					helper.ValidateExitCode(session, scanUtils.ScanTimeout, helper.ExitCodeFour)
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
		invalidPolicyPath := filepath.Join("path", "policy", "invalid")
		errString1 := "failed to initialize OPA policy engine"
		var errString2 string
		if utils.IsWindowsPlatform() {
			errString2 = fmt.Sprintf("%s: The system cannot find the path specified", invalidPolicyPath)
		} else {
			errString2 = fmt.Sprintf("%s: no such file or directory", invalidPolicyPath)
		}

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
				It("should error out and exit with status code 1", func() {
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

					Context("default iac type is all", func() {
						It("should scan with the policies and exit with status code 0", func() {
							scanArgs := []string{scanUtils.ScanCommand, "-p", validPolicyPath1, "-p", validPolicyPath2}
							session = helper.RunCommandDir(terrascanBinaryPath, workDirPath, outWriter, errWriter, scanArgs...)
							// exits with status code 4, because there are no iac files for violations but
							// would contain directory scan errors
							Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeFour))
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

	Describe("terrascan scan command is run with notification webhook and repo detail flags", func() {

		notificationURL := "https://httpbin.org/post"
		notificationToken := "token"

		tfGoldenRelPath := filepath.Join("golden", "terraform_scans")
		tfAwsAmiGoldenRelPath := filepath.Join(tfGoldenRelPath, "aws", "aws_ami_violations")

		iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
		policyDir, err2 := filepath.Abs(policyRootRelPath)
		It("should not error out while getting absolute path", func() {
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})
		Context("valid --webhook-url and --webhook-token flag is supplied", func() {
			It("should scan and display violations in human output format and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "--webhook-url", notificationURL, "--webhook-token", notificationToken}
				scanUtils.RunScanAndAssertGoldenOutputRegexWithTimeout(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanUtils.WebhookScanTimeout, scanArgs...)

			})
		})

		Context("only --webhook-url flag is supplied", func() {
			It("should scan and display violations in human output format and exit with status code 3", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "--webhook-url", notificationURL}
				scanUtils.RunScanAndAssertGoldenOutputRegexWithTimeout(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanUtils.WebhookScanTimeout, scanArgs...)

			})
		})
		Context("terrascan scan command is run with --repo-url and --repo-ref flag", func() {
			It("should scan and result in json output format and exit with status code 3", func() {
				scanArgs := []string{"scan", "-p", policyDir, "-i", "terraform", "-d", iacDir, "--repo-url", "https://github.com/tenable/terrascan.git", "--repo-ref", "main", "-o", "json"}
				// scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human_with_repo_detail.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeThree))
				sessionBytes := session.Wait().Out.Contents()
				sessionBytes = bytes.TrimSpace(sessionBytes)
				var sessionEngineOutput policy.EngineOutput

				err := json.Unmarshal(sessionBytes, &sessionEngineOutput)
				Expect(err).NotTo(HaveOccurred())
				Expect(sessionEngineOutput.Summary.ResourcePath).To((Equal("https://github.com/tenable/terrascan.git")))
				Expect(sessionEngineOutput.Summary.Branch).To((Equal("main")))
			})
		})
	})
	Describe("terrascan scan command run with policy-path", func() {
		Context("terrascan scan command is run with --p flag", func() {
			It("should scan and should use provided policies without downloading any from default repo", func() {
				if path, err := os.Stat(config.GetPolicyRepoPath()); err == nil && path.IsDir() {
					os.RemoveAll(config.GetPolicyRepoPath())
				}
				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				policyDir, err2 := filepath.Abs(policyRootRelPath)
				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-p", policyDir, "-i", "terraform", "-d", iacDir}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeThree))
				_, err := os.Stat(config.GetPolicyRepoPath())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})
	Describe("terrascan scan, write logs and results in file", func() {
		Context("terrascan scan command is run with --log-output-dir flag", func() {
			It("should scan and use the directory path to write the logs and result", func() {
				logOutputdir := "/tmp/terrascan"
				logFile := filepath.Join(logOutputdir, "terrascan.log")
				resultFile := filepath.Join(logOutputdir, "scan-result.txt")

				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				Expect(err1).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-d", iacDir, "--log-output-dir", logOutputdir}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeFive))
				_, err := os.Stat(logFile)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(resultFile)
				Expect(err).NotTo(HaveOccurred())

				// clean up
				os.RemoveAll(logOutputdir)
			})
			It("should scan and use the directory path to write the logs and result in json file", func() {
				logOutputdir := "/tmp/terrascan"
				logFile := filepath.Join(logOutputdir, "terrascan.log")
				resultFile := filepath.Join(logOutputdir, "scan-result.json")

				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				Expect(err1).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-d", iacDir, "--log-output-dir", logOutputdir, "-o", "json"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeFive))
				_, err := os.Stat(logFile)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(resultFile)
				Expect(err).NotTo(HaveOccurred())

				// clean up
				os.RemoveAll(logOutputdir)
			})
			It("should scan and use the directory path to write the logs and result in yaml file", func() {
				logOutputdir := "/tmp/terrascan"
				logFile := filepath.Join(logOutputdir, "terrascan.log")
				resultFile := filepath.Join(logOutputdir, "scan-result.yaml")

				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				Expect(err1).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-d", iacDir, "--log-output-dir", logOutputdir, "-o", "yaml"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeFive))
				_, err := os.Stat(logFile)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(resultFile)
				Expect(err).NotTo(HaveOccurred())

				// clean up
				os.RemoveAll(logOutputdir)
			})
			It("should scan and use the directory path to write the logs and result in xml file", func() {
				logOutputdir := "/tmp/terrascan"
				logFile := filepath.Join(logOutputdir, "terrascan.log")
				resultFile := filepath.Join(logOutputdir, "scan-result.xml")

				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				Expect(err1).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-d", iacDir, "--log-output-dir", logOutputdir, "-o", "xml"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeFive))
				_, err := os.Stat(logFile)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(resultFile)
				Expect(err).NotTo(HaveOccurred())

				// clean up
				os.RemoveAll(logOutputdir)
			})
			It("should scan and use the directory path to write the logs and result in sarif file", func() {
				logOutputdir := "/tmp/terrascan"
				logFile := filepath.Join(logOutputdir, "terrascan.log")
				resultFile := filepath.Join(logOutputdir, "scan-result.sarif")

				iacDir, err1 := filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
				Expect(err1).NotTo(HaveOccurred())
				scanArgs := []string{"scan", "-d", iacDir, "--log-output-dir", logOutputdir, "-o", "sarif"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, 3).Should(gexec.Exit(helper.ExitCodeFive))
				_, err := os.Stat(logFile)
				Expect(err).NotTo(HaveOccurred())
				_, err = os.Stat(resultFile)
				Expect(err).NotTo(HaveOccurred())

				// clean up
				os.RemoveAll(logOutputdir)
			})
		})
	})
})
