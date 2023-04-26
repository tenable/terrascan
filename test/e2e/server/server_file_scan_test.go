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

package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/policy"
	serverUtils "github.com/tenable/terrascan/test/e2e/server"
	"github.com/tenable/terrascan/test/helper"
)

var _ = Describe("Server File Scan", func() {

	var session *gexec.Session
	var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
	port := "9012"

	Context("file scan tests", func() {

		JustBeforeEach(func() {
			os.Setenv(terrascanServerPort, port)
		})
		JustAfterEach(func() {
			os.Setenv(terrascanServerPort, "")
		})

		// launches a server session on port 9012
		It("should start a new server session", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, serverUtils.ServerCommand)
			Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say(fmt.Sprintf("http server listening at port %s", port)))
		})

		awsDbInsViolationRelPath := filepath.Join(awsIacRelPath, "aws_db_instance_violation")
		k8sIngressViolationRelPath := filepath.Join(k8sIacRelPath, "kubernetes_ingress_violation")
		awsAmiViolationRelPath := filepath.Join(awsIacRelPath, "aws_ami_violation")

		Context("terraform file scan", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/terraform/v12/all/local/file/scan", host, port)
			When("iac file violates aws_db_instance_violation", func() {
				It("should return violations successfully", func() {
					iacFilePath, err := filepath.Abs(filepath.Join(awsDbInsViolationRelPath, "main.tf"))
					Expect(err).NotTo(HaveOccurred())

					goldenFilePath, err := filepath.Abs(filepath.Join(awsDbInstanceGoldenRelPath, "aws_db_instance_json.txt"))
					Expect(err).NotTo(HaveOccurred())

					responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusOK)
					serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
				})
			})
			When("iac file violates aws_ami", func() {
				It("should return violations successfully", func() {
					iacFilePath, err := filepath.Abs(filepath.Join(awsAmiViolationRelPath, "main.tf"))
					Expect(err).NotTo(HaveOccurred())

					goldenFilePath, err := filepath.Abs(filepath.Join(awsAmiGoldenRelPath, "aws_ami_violation_json.txt"))
					Expect(err).NotTo(HaveOccurred())

					responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusOK)
					serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
				})
			})

			Context("iac provider or iac version is invalid", func() {
				errMessage := "iac type or version not supported"
				iacFilePath, err := filepath.Abs(filepath.Join(awsDbInsViolationRelPath, "main.tf"))

				When("iac provider is invalid", func() {
					It("should get and error response", func() {
						requestURL := fmt.Sprintf("%s:%s/v1/%s/v12/all/local/file/scan", host, port, "terra")
						Expect(err).NotTo(HaveOccurred())

						serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusBadRequest)
						Eventually(session.Err).Should(gbytes.Say(errMessage))
					})
				})
				When("iac version is invalid", func() {
					It("should get and error response", func() {
						requestURL := fmt.Sprintf("%s:%s/v1/terraform/%s/all/local/file/scan", host, port, "v2")
						Expect(err).NotTo(HaveOccurred())

						serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusBadRequest)
						Eventually(session.Err).Should(gbytes.Say(errMessage))
					})
				})
			})

			Context("body attributes present in the request", func() {
				awsAmiIacFilePath, _ := filepath.Abs(filepath.Join(awsAmiViolationRelPath, "main.tf"))

				When("config_only attribute is set", func() {

					It("should receive resource config in response", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["config_only"] = "true"

						responseBytes := serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)

						var responseResourceConfig output.AllResourceConfigs
						err := json.Unmarshal(responseBytes, &responseResourceConfig)
						Expect(err).NotTo(HaveOccurred())

						Expect(responseResourceConfig).To(HaveLen(1))
						// the iac file only contains aws_ami resource
						Expect(responseResourceConfig).To(HaveKey("aws_ami"))
					})
				})

				When("show_passed attribute is set", func() {
					It("should receive passed rules along with violations", func() {
						iacFilePath, err := filepath.Abs(filepath.Join(awsDbInsViolationRelPath, "main.tf"))
						Expect(err).NotTo(HaveOccurred())

						goldenFilePath, err := filepath.Abs(filepath.Join(awsDbInstanceGoldenRelPath, "aws_db_instance_json_show_passed.txt"))
						Expect(err).NotTo(HaveOccurred())

						bodyAttrs := make(map[string]string)
						bodyAttrs["show_passed"] = "true"
						responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, bodyAttrs, http.StatusOK)
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
					})
				})

				Context("unknown body attributes are present", func() {
					Context("api server ignores unknown attributes", func() {
						It("should receive violations and 200 OK response", func() {
							bodyAttrs := make(map[string]string)
							bodyAttrs["unknown_attribute"] = "someValue"

							serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)
						})
					})
				})

				Context("body attributes have invalid value", func() {
					When("config_only has invalid value", func() {
						It("should receive an error", func() {
							bodyAttrs := make(map[string]string)
							bodyAttrs["config_only"] = "invalid"

							serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say(`error while reading 'config_only' value. error: 'strconv.ParseBool: parsing "invalid": invalid syntax'`))
						})
					})

					When("show_passed has invalid value", func() {
						It("should receive an error", func() {
							bodyAttrs := make(map[string]string)
							bodyAttrs["show_passed"] = "invalid"

							serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say(`error while reading 'show_passed' value. error: 'strconv.ParseBool: parsing "invalid": invalid syntax'`))
						})
					})
				})
			})
		})

		Context("k8s file scan", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/k8s/v1/all/local/file/scan", host, port)
			It("should return violations successfully", func() {
				iacFilePath, err := filepath.Abs(filepath.Join(k8sIngressViolationRelPath, "config.yaml"))
				Expect(err).NotTo(HaveOccurred())

				goldenFilePath, err := filepath.Abs(filepath.Join(k8sIngressViolationGoldenRelPath, "kubernetes_ingress_json.txt"))
				Expect(err).NotTo(HaveOccurred())

				responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusOK)
				serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
			})
		})

		Context("rules filtering options for file scan", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/terraform/v12/all/local/file/scan", host, port)
			iacFilePath, _ := filepath.Abs(filepath.Join(awsDbInsViolationRelPath, "main.tf"))

			When("scan_rules is used", func() {
				It("should receive violations and 200 OK response", func() {

					bodyAttrs := make(map[string]string)
					bodyAttrs["scan_rules"] = "AWS.RDS.DS.High.1041"
					responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(1))
				})
			})

			When("skip_rules is used", func() {

				It("should receive violations and 200 OK response", func() {
					bodyAttrs := make(map[string]string)
					bodyAttrs["skip_rules"] = "AWS.RDS.DataSecurity.High.0577"
					responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					// There are total 8 rules in the test policies directory, out of which 1 is skipped
					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(7))
				})
			})

			When("scan and skip rules is used", func() {
				It("should receive violations and 200 OK response", func() {
					bodyAttrs := make(map[string]string)
					bodyAttrs["scan_rules"] = "AWS.RDS.DS.High.1041,AWS.AWS RDS.NS.High.0101,AWS.RDS.DataSecurity.High.0577"
					bodyAttrs["skip_rules"] = "AWS.RDS.DataSecurity.High.0577"
					responseBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					// Total rules to be validated would be (scan_rules count -  skip_rules count)
					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(2))
				})
			})

			When("severity is used", func() {
				awsAmiIacFilePath, _ := filepath.Abs(filepath.Join(awsAmiViolationRelPath, "main.tf"))

				When("severity is invalid", func() {
					It("should receive a 400 bad request", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["severity"] = "1"

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusBadRequest)
						Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("severity level not supported"))
					})
				})

				When("severity is valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["severity"] = "HIGH"

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)
					})
				})
			})

			When("category is used", func() {
				awsAmiIacFilePath, _ := filepath.Abs(filepath.Join(awsAmiViolationRelPath, "main.tf"))

				When("category is invalid", func() {
					It("should receive a 400 bad request", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["categories"] = "Wrong caTeGory "

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusBadRequest)
						Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("category not supported"))
					})
				})

				When("multiple categories are sent but some of them are invalid", func() {
					It("should receive a 400 bad request", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["categories"] = " dATa pROtECtION, IDENTITY is Access Management "

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusBadRequest)
						Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("category not supported"))
					})
				})

				When("category is valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["categories"] = "DATA PROTECTION"

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)
					})
				})

				When("category is wrongly formatted but valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["categories"] = " dATa pROtECtION "

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)
					})
				})

				When("multiple categories are sent and all of them are valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]string)
						bodyAttrs["categories"] = " dATa pROtECtION, IDENTITY and Access Management "

						serverUtils.MakeFileScanRequest(awsAmiIacFilePath, requestURL, bodyAttrs, http.StatusOK)
					})
				})
			})

			Context("resource is skipped", func() {
				resourceSkipIacRelPath := filepath.Join(iacRootRelPath, "resource_skipping")

				When("tf file has resource skipped", func() {
					It("should receive violations result with 200 OK response", func() {
						iacFilePath, err := filepath.Abs(filepath.Join(resourceSkipIacRelPath, "terraform", "main.tf"))
						Expect(err).NotTo(HaveOccurred())

						goldenFilePath, err := filepath.Abs(filepath.Join(goldenFilesRelPath, "resource_skipping", "terraform_file_resource_skipping.txt"))
						Expect(err).NotTo(HaveOccurred())

						respBytes := serverUtils.MakeFileScanRequest(iacFilePath, requestURL, nil, http.StatusOK)
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, respBytes)
					})
				})

				When("k8s file has resource skipped", func() {
					k8sRequestURL := fmt.Sprintf("%s:%s/v1/k8s/v1/all/local/file/scan", host, port)
					It("should receive violations result with 200 OK response", func() {
						iacFilePath, err := filepath.Abs(filepath.Join(resourceSkipIacRelPath, "kubernetes", "config.yaml"))
						Expect(err).NotTo(HaveOccurred())

						goldenFilePath, err := filepath.Abs(filepath.Join(goldenFilesRelPath, "resource_skipping", "kubernetes_file_resource_skipping.txt"))
						Expect(err).NotTo(HaveOccurred())

						respBytes := serverUtils.MakeFileScanRequest(iacFilePath, k8sRequestURL, nil, http.StatusOK)
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, respBytes)
					})
				})
			})
		})
	})
})
