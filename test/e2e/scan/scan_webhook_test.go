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
	"path/filepath"

	scanUtils "github.com/accurics/terrascan/test/e2e/scan"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Scan Command using webhook args", func() {
	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	var policyDir, iacDir string
	var err error

	notificationURL := "https://httpbin.org/post"
	notificationToken := "token"

	Describe("terrascan scan command is run with --notification-webhook-url and --notification-webhook-token flag", func() {
		tfGoldenRelPath := filepath.Join("golden", "terraform_scans")
		tfAwsAmiGoldenRelPath := filepath.Join(tfGoldenRelPath, "aws", "aws_ami_violations")

		iacDir, err = filepath.Abs(filepath.Join(awsIacRelPath, "aws_ami_violation"))
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		policyDir, err = filepath.Abs(policyRootRelPath)
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("valid --notification-webhook-url and --notification-webhook-token flag is supplied", func() {
			It("should exit with status code 5", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "--notification-webhook-url", notificationURL, "--notification-webhook-token", notificationToken}
				scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)

			})
		})

		Context("only --notification-webhook-url flag is supplied", func() {
			It("should exit with status code 5", func() {
				scanArgs := []string{"-p", policyDir, "-i", "terraform", "-d", iacDir, "--notification-webhook-url", notificationURL}
				scanUtils.RunScanAndAssertGoldenOutputRegex(terrascanBinaryPath, filepath.Join(tfAwsAmiGoldenRelPath, "aws_ami_violation_human.txt"), helper.ExitCodeThree, false, true, outWriter, errWriter, scanArgs...)

			})
		})
	})

})
