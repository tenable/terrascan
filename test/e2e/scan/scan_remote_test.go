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

var _ = Describe("Scan Command using remote types", func() {
	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	Context("remote type is supplied, but remote URL is not", func() {
		errString := "empty remote url or type or destination dir path"

		When("remote type is git", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "git"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		When("remote type is s3", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "s3"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		When("remote type is gcs", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "gcs"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		When("remote type is http", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "http"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})

		When("remote type is terraform-registry", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry"}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})
	})

	Context("valid remote type is supplied with invalid remote URL ", func() {
		invalidRemoteURL := "test"
		When("remote type is git", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "git", "--remote-url", invalidRemoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is s3", func() {
			remoteURL := invalidRemoteURL
			JustBeforeEach(func() {
				invalidRemoteURL = "s3://"
			})
			JustAfterEach(func() {
				invalidRemoteURL = remoteURL
			})
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "s3", "--remote-url", invalidRemoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is gcs", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "gcs", "--remote-url", invalidRemoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is http", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "http", "--remote-url", invalidRemoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is terraform-registry", func() {
			It("should error out and exit with status code 1", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", invalidRemoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		Context("when remote type is unsupported", func() {
			It("should error out and exit with status code 1", func() {
				errString := "supplied remote type is not supported"
				scanArgs := []string{scanUtils.ScanCommand, "-r", "unsupportedType", "--remote-url", invalidRemoteURL}
				scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.RemoteScanTimeout, errString, outWriter, errWriter, scanArgs...)
			})
		})
	})

	Context("valid remote type is supplied with valid remote URL", func() {
		When("remote type is git", func() {
			remoteURL := "github.com/tenable/kaimonkey/terraform/aws"
			It("should download the resource and generate scan results", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-r", "git", "--remote-url", remoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				// exit code is 5 because iac files in directory has violations
				// and directory scan errors
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeFive))
			})

			It("should download the resource and generate scan results", func() {
				remoteURL := "https://github.com/tenable/kaimonkey.git//terraform/aws"
				scanArgs := []string{scanUtils.ScanCommand, "-r", "git", "--remote-url", remoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				// exit code is 5 because iac files in directory has violations
				// and directory scan errors
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeFive))
			})
		})

		When("remote type is s3", func() {
			It("should download the resource and generate scan results", func() {
				Skip("Skipping this test until we have a s3 url")
			})
		})

		When("remote type is gcs", func() {
			It("should download the resource and generate scan results", func() {
				Skip("Skipping this test until we have a gcs url")
			})
		})

		When("remote type is http", func() {
			It("should download the resource and generate scan results", func() {
				Skip("Skipping this test")
			})
		})

		When("remote type is terraform-registry", func() {
			remoteURL := "terraform-aws-modules/vpc/aws"
			When("terraform registry URL doesn't have version specified, it downloads the latest available version", func() {
				It("should download the resource and generate scan results", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					// has a OR condition because we don't know if there would be violations or not
					// there would be directory scan errors due to all iac type
					Eventually(session, scanUtils.RemoteScanTimeout).Should(Or(gexec.Exit(helper.ExitCodeFive), gexec.Exit(helper.ExitCodeFour)))
				})
			})

			When("terraform registry remote url has a version", func() {
				remoteURL = "terraform-aws-modules/vpc/aws:2.22.0"
				It("should download the remote registry and generate scan results", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					// has a OR condition because we don't know if there would be violations or not
					// there would be directory scan errors due to all iac type
					Eventually(session, scanUtils.RemoteScanTimeout).Should(Or(gexec.Exit(helper.ExitCodeFive), gexec.Exit(helper.ExitCodeFour)))
				})
			})

			Context("remote modules has reference to its local modules", func() {
				When("remote type is terraform registry and remote url has a subdirectory", func() {
					remoteURL := "terraform-aws-modules/security-group/aws//modules/http-80"
					It("should download the remote registry and generate scan results", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL, "-i", "terraform", "--non-recursive"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						// has a OR condition because we don't know if there would be violations or not
						Eventually(session, scanUtils.RemoteScanTimeout).Should(Or(gexec.Exit(helper.ExitCodeThree), gexec.Exit(helper.ExitCodeZero)))
					})
				})

				When("remote type is git and remote url has a subdirectory", func() {
					remoteURL := "github.com/terraform-aws-modules/terraform-aws-security-group//modules/http-80"
					It("should download the remote registry and generate scan results", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-r", "git", "--remote-url", remoteURL, "-i", "terraform", "--non-recursive"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						// has a OR condition because we don't know if there would be violations or not
						Eventually(session, scanUtils.RemoteScanTimeout).Should(Or(gexec.Exit(helper.ExitCodeThree), gexec.Exit(helper.ExitCodeZero)))
					})
				})
			})

			When("terraform registry remote url has a invalid version", func() {
				remoteURL := "terraform-aws-modules/vpc/aws:blah"
				It("should error out and exit with status code 1", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				})
			})
		})
	})
	Context("when scan is run using custom temp directory with env variable", func() {
		When("remote type is git", func() {
			remoteURL := "github.com/tenable/kaimonkey/terraform/aws"
			tmpDir, err := filepath.Abs(filepath.Join(iacRootRelPath, "temp_dir"))
			Expect(err).NotTo(HaveOccurred())
			JustBeforeEach(func() {
				err = os.Setenv("TERRASCAN_CUSTOM_TEMP_DIR", tmpDir)
				Expect(err).NotTo(HaveOccurred())
			})

			JustAfterEach(func() {
				err := os.Unsetenv("TERRASCAN_CUSTOM_TEMP_DIR")
				Expect(err).NotTo(HaveOccurred())
			})
			It("should download the resource in provided custom temp dir and generate scan results", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-o", "json", "-r", "git", "--remote-url", remoteURL}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				// exit code is 5 because iac files in directory has violations
				// and directory scan errors
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeFive))
				helper.ContainsDirScanErrorSubString(session, tmpDir)
			})
		})
	})
	Context("when scan is run on remote dir and using flag --temp-dir to set custom temp dir", func() {
		When("remote type is git", func() {
			remoteURL := "github.com/tenable/kaimonkey/terraform/aws"
			tmpDir, err := filepath.Abs(filepath.Join(iacRootRelPath, "temp_dir"))
			Expect(err).NotTo(HaveOccurred())
			It("should download the resource in provided custom temp dir and generate scan results", func() {
				scanArgs := []string{scanUtils.ScanCommand, "-o", "json", "-r", "git", "--remote-url", remoteURL, "--temp-dir", tmpDir}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
				// exit code is 5 because iac files in directory has violations
				// and directory scan errors
				Eventually(session, scanUtils.RemoteScanTimeout).Should(gexec.Exit(helper.ExitCodeFive))
				helper.ContainsDirScanErrorSubString(session, tmpDir)
			})
		})
	})
})
