package help_test

import (
	"io"

	helpUtils "github.com/accurics/terrascan/test/e2e/help"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	helpCommand string = "help"
)

var _ = Describe("Help", func() {

	var session *gexec.Session
	var terrascanBinaryPath string

	var outWriter, errWriter io.Writer

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

	Describe("terrascan is run without any command", func() {
		It("should print all supported commands and exit with status code 0", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter)
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_command.txt", true)
		})
	})

	Describe("terrascan is run -h flag", func() {
		It("should print all supported commands and exit with status code 0", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "-h")
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_flag.txt", true)
		})
	})

	Describe("terrascan is run with an unkonwn command", func() {
		It("should exit with status code 1 and display a error message", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "test")
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeOne, "golden/incorrect_command.txt", false)
		})
	})

	Describe("help is run", func() {
		Context("with no arguments", func() {
			It("should print the terrascan help and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand)
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_command.txt", true)
			})
		})

		Context("for init command", func() {
			It("should print help for init and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "init")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_init.txt", true)
			})
		})

		Context("for scan command", func() {
			It("should print help for init and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "scan")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_scan.txt", true)
			})
		})

		Context("for server command", func() {
			It("should print help for init and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "server")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_server.txt", true)
			})
		})

		Context("for version command", func() {
			It("should print help for init and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "version")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_version.txt", true)
			})
		})

		Context("for an unkonwn command", func() {
			It("should display that help topic is not available for entered command and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "test")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, "golden/help_unsupported_command.txt", false)
			})
		})
	})
})
