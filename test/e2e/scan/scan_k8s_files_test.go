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
	scanUtils "github.com/tenable/terrascan/test/e2e/scan"
	"github.com/tenable/terrascan/test/helper"
)

var _ = Describe("Scan is run for k8s directories and files", func() {

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	var policyDir, iacDir string
	policyDir, err1 := filepath.Abs(policyRootRelPath)
	iacDir, err2 := filepath.Abs(filepath.Join(k8sIacRelPath, "kubernetes_ingress_violation"))

	It("should not error out while getting absolute path", func() {
		Expect(err1).NotTo(HaveOccurred())
		Expect(err2).NotTo(HaveOccurred())
	})

	Context("scan iac directories violating k8s policies", func() {
		Context("iac type k8s will be part of all iac", func() {
			When("k8s files are scanned but iac type is not specified", func() {
				It("should scan will all iac and display violations", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-d", iacDir}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					// exit code is 5 because iac files in directory has violations
					// and directory scan errors
					helper.ValidateExitCode(session, scanUtils.ScanTimeout, helper.ExitCodeFive)
				})
			})
		})

		k8sGoldenRelPath := filepath.Join("golden", "k8s_scans", "k8s", "kubernetes_ingress_violations")

		Context("iac type is specified as k8s", func() {
			It("should scan and display violations in human output format", func() {
				scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir}
				scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})

			When("-v flag is used for verbose output", func() {
				It("should display verbose output for human output format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-v"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_human_verbose.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is sarif", func() {
				It("should display violations in sarif format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "sarif"}
					scanUtils.RunScanAndAssertGoldenSarifOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_sarif.txt"), helper.ExitCodeThree, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_json.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_yaml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_xml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is junit-xml", func() {
				It("should display violations in junit-xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-d", iacDir, "-o", "junit-xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_junit_xml.txt"), helper.ExitCodeThree, true, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})

	Context("scan iac files violating k8s policies", func() {
		iacFile := filepath.Join(iacDir, "config.yaml")
		It("should not error out while getting absolute path", func() {
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
		})

		Context("iac type k8s will be part of all iac", func() {
			When("k8s files are scanned but iac type is not specified", func() {
				It("should scan will all iac and display violations", func() {
					scanArgs := []string{scanUtils.ScanCommand, "-f", iacFile}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, scanArgs...)
					// exit code is 1 because iac file is expected to be of terraform iac type by default, not k8s yaml
					helper.ValidateExitCode(session, scanUtils.ScanTimeout, helper.ExitCodeOne)
				})
			})
		})

		k8sGoldenRelPath := filepath.Join("golden", "k8s_scans", "k8s", "kubernetes_ingress_violations")

		Context("iac type is specified as k8s", func() {
			It("should scan and display violations in human output format", func() {
				scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile}
				scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
			})

			When("-v flag is used for verbose output", func() {
				It("should display verbose output for human output format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-v"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_human_verbose.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is sarif", func() {
				It("should display violations in sarif format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-o", "sarif"}
					scanUtils.RunScanAndAssertGoldenSarifOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_sarif.txt"), helper.ExitCodeThree, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-o", "json"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_json.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-o", "yaml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_yaml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-o", "xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_xml.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is junit-xml", func() {
				It("should display violations in junit-xml format", func() {
					scanArgs := []string{"-i", "k8s", "-p", policyDir, "-f", iacFile, "-o", "junit-xml"}
					scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_junit_xml.txt"), helper.ExitCodeThree, true, true, outWriter, errWriter, scanArgs...)
				})
			})
		})
	})
})
