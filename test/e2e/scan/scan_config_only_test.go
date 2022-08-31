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
	iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))

	It("should not error out while getting absolute path", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("scan command is run using the --config-only flag for unsupported output types", func() {
		When("output type is human readable format", func() {
			Context("it doesn't support --config-only flag", func() {
				Context("human readable output format is the default output format", func() {
					It("should result in an error and exit with status code 1", func() {
						errString := "please use yaml or json output format when using --config-only or --config-with-error flags"
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only"}
						scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
					})
				})
			})
		})

		When("output type is xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					errString := "please use yaml or json output format when using --config-only or --config-with-error flags"
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "xml"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})

		When("output type is junit-xml", func() {
			Context("it doesn't support --config-only flag", func() {
				It("should result in an error and exit with status code 1", func() {
					errString := "please use yaml or json output format when using --config-only or --config-with-error flags"
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "junit-xml"}
					scanUtils.RunScanAndAssertErrorMessage(terrascanBinaryPath, helper.ExitCodeOne, scanUtils.ScanTimeout, errString, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Describe("scan command is run using the --config-only flag for supported output types", func() {
		Context("for terraform files", func() {
			When("output type is json", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 0", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "json"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 0", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "yaml"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
					})
				})
			})
		})

		Context("for yaml files", func() {
			JustBeforeEach(func() {
				iacDir, err = filepath.Abs(filepath.Join(k8sIacRelPath, "kubernetes_ingress_violation"))
			})
			When("output type is json", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 0", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "json", "-i", "k8s"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
					})
				})
			})

			When("output type is yaml", func() {
				Context("it supports --config-only flag", func() {
					It("should display config json and exit with status code 0", func() {
						scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir, "--config-only", "-o", "yaml", "-i", "k8s"}
						session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
						Eventually(session, scanUtils.ScanTimeout).Should(gexec.Exit(helper.ExitCodeZero))
					})
				})
			})
		})
	})
})
