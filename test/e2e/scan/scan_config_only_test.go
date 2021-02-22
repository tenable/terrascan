package scan_test

import (
	"path/filepath"

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
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only")
						Eventually(session).Should(gexec.Exit(1))
						helper.ContainsErrorSubString(session, "please use yaml or json output format when using --config-only flag")
					})
				})
			})
		})

		When("output type is xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "xml")
					Eventually(session, 5).Should(gexec.Exit(1))
					helper.ContainsErrorSubString(session, "failed to write XML output. error: 'xml: unsupported type: output.AllResourceConfigs'")
				})
			})
		})

		When("output type is junit-xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "junit-xml")
					Eventually(session, 5).Should(gexec.Exit(1))
					helper.ContainsErrorSubString(session, "incorrect input for JunitXML writer, supported type is policy.EngineOutput")
				})
			})
		})
	})

	Describe("scan command is run using the --config-only flag for unsupported output types", func() {
		Context("for terraform files", func() {
			When("output type is json", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "json")
						Eventually(session, 5).Should(gexec.Exit(3))
						goldenFileAbsPath, err := filepath.Abs("golden/config_only/config_only_tf.json")
						Expect(err).NotTo(HaveOccurred())
						helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "yaml")
						Eventually(session, 5).Should(gexec.Exit(3))
						goldenFileAbsPath, err := filepath.Abs("golden/config_only/config_only_tf.yaml")
						Expect(err).NotTo(HaveOccurred())
						helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
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
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "json", "-i", "k8s")
						Eventually(session, 5).Should(gexec.Exit(3))
						goldenFileAbsPath, err := filepath.Abs("golden/config_only/config_only_k8s.json")
						Expect(err).NotTo(HaveOccurred())
						helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 3", func() {
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanComand, "-d", iacDir, "--config-only", "-o", "yaml", "-i", "k8s")
						Eventually(session, 5).Should(gexec.Exit(3))
						goldenFileAbsPath, err := filepath.Abs("golden/config_only/config_only_k8s.yaml")
						Expect(err).NotTo(HaveOccurred())
						helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
					})
				})
			})
		})
	})
})
