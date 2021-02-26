package scan_test

import (
	scanUtils "github.com/accurics/terrascan/test/e2e/scan"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
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
		errString := "empty remote url or type or desitnation dir path"

		When("remote type is git", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "git")
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString)
			})
		})

		When("remote type is s3", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "s3")
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString)
			})
		})

		When("remote type is gcs", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "gcs")
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString)
			})
		})

		When("remote type is http", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "http")
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString)
			})
		})

		When("remote type is terraform-registry", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "terraform-registry")
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, errString)
			})
		})
	})

	Context("valid remote type is supplied with invalid remote URL ", func() {
		invalidRemoteURL := "test"
		When("remote type is git", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "git", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
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
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "s3", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is gcs", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "gcs", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is http", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "http", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		When("remote type is terraform-registry", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})

		Context("when remote type is unsupported", func() {
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "unsupportedType", "--remote-url", invalidRemoteURL)
				Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				helper.ContainsErrorSubString(session, "supplied remote type is not supported")
			})
		})
	})

	Context("valid remote type is supplied with valid remote URL", func() {
		When("remote type is git", func() {
			remoteURL := "github.com/accurics/KaiMonkey/terraform/aws"
			It("should download the resource and generate scan results", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "git", "--remote-url", remoteURL)
				Eventually(session, 10).Should(gexec.Exit(helper.ExitCodeThree))
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
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL)
					// has a OR condition because we don't know if there would be violations or not
					Eventually(session, 10).Should(Or(gexec.Exit(helper.ExitCodeThree), gexec.Exit(helper.ExitCodeZero)))
				})
			})

			When("terraform registry remote url has a version", func() {
				oldRemoteURL := remoteURL
				JustBeforeEach(func() {
					remoteURL = "terraform-aws-modules/vpc/aws:2.22.0"
				})
				JustAfterEach(func() {
					remoteURL = oldRemoteURL
				})
				It("should download the remote registry and generate scan results", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL)
					// has a OR condition because we don't know if there would be violations or not
					Eventually(session, 10).Should(Or(gexec.Exit(helper.ExitCodeThree), gexec.Exit(helper.ExitCodeZero)))
				})
			})

			When("terraform registry remote url has a invalid version", func() {
				oldRemoteURL := remoteURL
				JustBeforeEach(func() {
					remoteURL = "terraform-aws-modules/vpc/aws:blah"
				})
				JustAfterEach(func() {
					remoteURL = oldRemoteURL
				})
				It("should error out and exit with status code 1", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanUtils.ScanCommand, "-r", "terraform-registry", "--remote-url", remoteURL)
					Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeOne))
				})
			})
		})
	})
})
