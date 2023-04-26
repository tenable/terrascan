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
package help_test

import (
	"io"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	helpUtils "github.com/tenable/terrascan/test/e2e/help"
	"github.com/tenable/terrascan/test/helper"
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
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_command.txt"), true)
		})
	})

	Describe("terrascan is run -h flag", func() {
		It("should print all supported commands and exit with status code 0", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "-h")
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_flag.txt"), true)
		})
	})

	Describe("terrascan is run with an unknown command", func() {
		It("should exit with status code 1 and display a error message", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "test")
			helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeOne, filepath.Join("golden", "incorrect_command.txt"), false)
		})
	})

	Describe("help is run", func() {
		Context("with no arguments", func() {
			It("should print the terrascan help and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand)
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_command.txt"), true)
			})
		})

		Context("for init command", func() {
			It("should print help for init and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "init")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_init.txt"), true)
			})
		})

		Context("for scan command", func() {
			It("should print help for scan and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "scan")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_scan.txt"), true)
			})
		})

		Context("for server command", func() {
			It("should print help for server and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "server")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_server.txt"), true)
			})
		})

		Context("for version command", func() {
			It("should print help for version and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "version")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_version.txt"), true)
			})
		})

		Context("for an unknown command", func() {
			It("should display that help topic is not available for entered command and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, helpCommand, "test")
				helpUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("golden", "help_unsupported_command.txt"), false)
			})
		})
	})
})
