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
	"bytes"
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

var _ = Describe("Server Remote Scan", func() {
	// In case of adding new test case first push the test data and golden data files and then write test cases around that
	var session *gexec.Session
	var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
	port := "9011"

	Context("remote repo scan", func() {

		JustBeforeEach(func() {
			os.Setenv(terrascanServerPort, port)
		})
		JustAfterEach(func() {
			os.Setenv(terrascanServerPort, "")
		})

		// launches a server session on port 9011
		It("should start a new server session", func() {
			session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, serverUtils.ServerCommand)
			Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say(fmt.Sprintf("http server listening at port %s", port)))
		})

		Context("scan remote terraform git repo", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/terraform/v12/all/remote/dir/scan", host, port)

			awsAmiRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/aws/aws_ami_violation"
			When("remote repo violates aws_ami", func() {
				It("should report violations", func() {

					goldenFilePath, err := filepath.Abs(filepath.Join(awsAmiGoldenRelPath, "aws_ami_violation_json.txt"))
					Expect(err).NotTo(HaveOccurred())

					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = awsAmiRepoURL

					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
					serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
				})
			})

			When("remote repo violates aws_db_instance", func() {
				It("should report violations", func() {
					remoteRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/aws/aws_db_instance_violation"

					goldenFilePath, err := filepath.Abs(filepath.Join(awsDbInstanceGoldenRelPath, "aws_db_instance_json.txt"))
					Expect(err).NotTo(HaveOccurred())

					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = remoteRepoURL

					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
					serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
				})
			})

			Context("body attributes are present in the request", func() {
				When("config_only attribute is present", func() {
					It("should receive resource config response", func() {
						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = awsAmiRepoURL
						bodyAttrs["config_only"] = true

						responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						var responseResourceConfig output.AllResourceConfigs
						err := json.Unmarshal(responseBytes, &responseResourceConfig)
						Expect(err).NotTo(HaveOccurred())

						Expect(responseResourceConfig).To(HaveLen(1))
						// the iac file only contains aws_ami resource
						Expect(responseResourceConfig).To(HaveKey("aws_ami"))
					})
				})

				When("show_passed attribute is present", func() {

					It("should receive resource config response", func() {
						remoteRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/aws/aws_db_instance_violation"

						goldenFilePath, err := filepath.Abs(filepath.Join(awsDbInstanceGoldenRelPath, "aws_db_instance_json_show_passed.txt"))
						Expect(err).NotTo(HaveOccurred())

						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = remoteRepoURL
						bodyAttrs["show_passed"] = true

						responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
					})
				})

				Context("Unknown attribute is present in body", func() {
					Context("server would not read unknown attributes", func() {

						It("should receive 200 OK response", func() {

							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["unknown"] = "invalid"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						})
					})
				})

				When("non_recursive attribute is present in body", func() {
					Context("remote url contains terraform files", func() {
						It("should receive 200 OK response", func() {

							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["non_recursive"] = true

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						})
					})

					Context("remote url doesn't not contain terraform files", func() {
						It("should receive 400 bad request response", func() {

							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = "https://github.com/tenable/terrascan//test/e2e/test_data/iac/k8s/kubernetes_ingress_violation"
							bodyAttrs["non_recursive"] = true

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
						})
					})
				})

				Context("remote_type or remote_url is not present in the request", func() {
					errMessage := "remote url or destination dir path cannot be empty"

					When("both attributes are not present", func() {
						It("should receive a 400 bad request response", func() {
							serverUtils.MakeRemoteScanRequest(requestURL, nil, http.StatusBadRequest)
						})
					})

					When("when remote_type is not present", func() {
						It("should receive a 400 bad request response", func() {
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_url"] = awsAmiRepoURL

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("when remote_url is not present", func() {
						It("should receive a 400 bad request response", func() {
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})
				})

				Context("invalid values for known attributes", func() {

					When("remote_type is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "supplied remote type is not supported"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "invalid"
							bodyAttrs["remote_url"] = awsAmiRepoURL

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("scan_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.scan_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["scan_rules"] = "Rule.1"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("skip_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.skip_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["skip_rules"] = "Rule.1"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("config_only value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.config_only of type bool"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["config_only"] = "invalid"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("show_passed value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.show_passed of type bool"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["show_passed"] = "invalid"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("non_recursive value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.non_recursive of type bool"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["non_recursive"] = "invalid"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("severity value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal number into Go struct field scanRemoteRepoReq.severity of type string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["severity"] = 1

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("category value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal number into Go struct field scanRemoteRepoReq.categories of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["categories"] = 4

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("scan_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.scan_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["scan_rules"] = "Rule.1"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})

					When("skip_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.skip_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["skip_rules"] = "Rule.1"

							responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Expect(string(bytes.TrimSpace(responseBytes))).To(Equal(errMessage))
						})
					})
				})
			})
		})

		Context("scan remote k8s git repo", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/k8s/v1/all/remote/dir/scan", host, port)

			When("remote repo violates kubernetes_ingress", func() {
				It("should report violations", func() {
					remoteRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/k8s/kubernetes_ingress_violation"

					goldenFilePath, err := filepath.Abs(filepath.Join(k8sGoldenRelPath, "kubernetes_ingress_violations", "kubernetes_ingress_json.txt"))
					Expect(err).NotTo(HaveOccurred())

					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = remoteRepoURL

					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
					serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, responseBytes)
				})
			})
		})

		Context("scan remote with remote type s3", func() {
			It("should scan the s3 repo", func() {
				Skip("skipping until we have a s3 bucket is available")
			})
		})

		Context("scan remote with remote type gcp", func() {
			It("should scan the gcp repo", func() {
				Skip("skipping until we have a gcp resource is available")
			})
		})

		Context("scan remote with remote type http", func() {
			It("should scan the remote http repo", func() {
				Skip("skipping test for now")
			})
		})

		Context("scan remote with remote type terraform-registry", func() {
			It("should scan the terraform registry, return violations with 200 OK response", func() {
				requestURL := fmt.Sprintf("%s:%s/v1/terraform/v12/all/remote/dir/scan", host, port)
				bodyAttrs := make(map[string]interface{})
				bodyAttrs["remote_type"] = "terraform-registry"
				bodyAttrs["remote_url"] = "terraform-aws-modules/vpc/aws"

				serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
			})
		})

		Context("rules filtering options for remote scan", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/terraform/v14/all/remote/dir/scan", host, port)
			remoteRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/resource_skipping/terraform"

			When("scan_rules is used", func() {
				It("should receive violations and 200 OK response", func() {

					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = remoteRepoURL
					bodyAttrs["scan_rules"] = []string{"AWS.RDS.DS.High.1041"}
					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(1))
				})
			})

			When("skip_rules is used", func() {

				It("should receive violations and 200 OK response", func() {
					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = remoteRepoURL
					bodyAttrs["skip_rules"] = []string{"AWS.RDS.DataSecurity.High.0577"}
					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					// There are total 8 rules in the test policies directory, out of which 1 is skipped
					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(7))
				})
			})

			When("scan and skip rules is used", func() {
				It("should receive violations and 200 OK response", func() {
					bodyAttrs := make(map[string]interface{})
					bodyAttrs["remote_type"] = "git"
					bodyAttrs["remote_url"] = remoteRepoURL
					bodyAttrs["scan_rules"] = []string{"AWS.RDS.DS.High.1041", "AWS.AWS RDS.NS.High.0101", "AWS.RDS.DataSecurity.High.0577"}
					bodyAttrs["skip_rules"] = []string{"AWS.RDS.DataSecurity.High.0577"}
					responseBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)

					var responseEngineOutput policy.EngineOutput
					err := json.Unmarshal(responseBytes, &responseEngineOutput)
					Expect(err).NotTo(HaveOccurred())

					// Total rules to be validated would be (scan_rules count -  skip_rules count)
					Expect(responseEngineOutput.ViolationStore.Summary.TotalPolicies).To(BeIdenticalTo(2))
				})
			})

			When("severity is used", func() {
				When("severity is valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = remoteRepoURL
						bodyAttrs["severity"] = "HIGH"

						serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
					})
				})
			})

			When("categories are used", func() {
				When("categories are valid", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = remoteRepoURL
						bodyAttrs["categories"] = []string{"DATA PROTECTION", "compliance validation"}

						serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
					})
				})
			})

			Context("resource is skipped", func() {

				resourceSkipGoldenRelPath := filepath.Join(goldenFilesRelPath, "resource_skipping")

				When("tf file has resource skipped", func() {
					It("should receive violations result with 200 OK response", func() {
						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = remoteRepoURL

						goldenFilePath, err := filepath.Abs(filepath.Join(resourceSkipGoldenRelPath, "terraform_file_resource_skipping.txt"))
						Expect(err).NotTo(HaveOccurred())

						respBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, respBytes)
					})
				})

				When("k8s file has resource skipped", func() {
					It("should receive violations result with 200 OK response", func() {
						requestURL := fmt.Sprintf("%s:%s/v1/k8s/v1/all/remote/dir/scan", host, port)
						remoteRepoURL := "https://github.com/tenable/terrascan//test/e2e/test_data/iac/resource_skipping/kubernetes"
						bodyAttrs := make(map[string]interface{})
						bodyAttrs["remote_type"] = "git"
						bodyAttrs["remote_url"] = remoteRepoURL

						goldenFilePath, err := filepath.Abs(filepath.Join(resourceSkipGoldenRelPath, "kubernetes_file_resource_skipping.txt"))
						Expect(err).NotTo(HaveOccurred())

						respBytes := serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusOK)
						// assertion is required since result will Contain skipped_violations
						serverUtils.CompareResponseAndGoldenOutput(goldenFilePath, respBytes)
					})
				})
			})
		})
	})
})
