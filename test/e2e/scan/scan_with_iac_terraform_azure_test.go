package scan_test

import (
	"path/filepath"

	scanUtils "github.com/accurics/terrascan/test/e2e/scan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScanWithIacTerraformAzure", func() {

	Describe("scan iac files violating azure policies against test policy set", func() {
		var policyDir, iacDir string
		var err error

		policyDir, err = filepath.Abs("../test_data/policies/")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		iacDir, err = filepath.Abs("../test_data/iac/azure/azurem_application_gateway_violation")
		It("should not error out while getting absolute path", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		Context("iac files violates azurem application gateway policy", func() {
			When("when output type is json", func() {
				It("should display violations in json format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "json"}
					scanUtils.RunScanCommandAndAssertJSONOutput(terrascanBinaryPath, "golden/terraform_scans/azure/azurem_application_gateway_violation/azurem_application_gateway_json.txt", 3, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is yaml", func() {
				It("should display violations in yaml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "yaml"}
					scanUtils.RunScanCommandAndAssertYAMLOutput(terrascanBinaryPath, "golden/terraform_scans/azure/azurem_application_gateway_violation/azurem_application_gateway_yaml.txt", 3, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			When("when output type is xml", func() {
				It("should display violations in xml format", func() {
					scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "xml"}
					scanUtils.RunScanCommandAndAssertXMLOutput(terrascanBinaryPath, "golden/terraform_scans/azure/azurem_application_gateway_violation/azurem_application_gateway_xml.txt", 3, false, true, outWriter, errWriter, scanArgs...)
				})
			})

			// When("when output type is junit-xml", func() {
			// 	It("should display violations in junit-xml format", func() {
			// 		scanArgs := []string{"-p", policyDir, "-d", iacDir, "-o", "junit-xml"}
			// 		scanUtils.RunScanCommandAndAssertTextOutput(terrascanBinaryPath, "golden/terraform_scans/azure/azurem_application_gateway_violation/azurem_application_gateway_junit_json.txt", 3, true, true, outWriter, errWriter, scanArgs...)
			// 	})
			// })
		})
	})
})
