package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	serverUtils "github.com/accurics/terrascan/test/e2e/server"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Server Remote Scan", func() {

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

			awsAmiRepoURL := "https://github.com/accurics/terrascan//test/e2e/test_data/iac/aws/aws_ami_violation"
			When("remote repo violates aws_ami", func() {
				It("should report violations", func() {

					goldenFilePath, err := filepath.Abs(filepath.Join("..", "scan", "golden", "terraform_scans", "aws", "aws_ami_violations", "aws_ami_violation_json.txt"))
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
					remoteRepoURL := "https://github.com/accurics/terrascan//test/e2e/test_data/iac/aws/aws_db_instance_violation"

					goldenFilePath, err := filepath.Abs(filepath.Join("..", "scan", "golden", "terraform_scans", "aws", "aws_db_instance_violations", "aws_db_instance_json.txt"))
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
						remoteRepoURL := "https://github.com/accurics/terrascan//test/e2e/test_data/iac/aws/aws_db_instance_violation"

						goldenFilePath, err := filepath.Abs(filepath.Join("..", "scan", "golden", "terraform_scans", "aws", "aws_db_instance_violations", "aws_db_instance_json_show_passed.txt"))
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

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("scan_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.scan_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["scan_rules"] = "Rule.1"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("skip_rules value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.skip_rules of type []string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["skip_rules"] = "Rule.1"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("config_only value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.config_only of type bool"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["config_only"] = "invalid"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("show_passed value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal string into Go struct field scanRemoteRepoReq.show_passed of type bool"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["show_passed"] = "invalid"

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})

					When("severity value is invalid", func() {
						It("should receive a 400 bad request response", func() {
							errMessage := "json: cannot unmarshal number into Go struct field scanRemoteRepoReq.severity of type string"
							bodyAttrs := make(map[string]interface{})
							bodyAttrs["remote_type"] = "git"
							bodyAttrs["remote_url"] = awsAmiRepoURL
							bodyAttrs["severity"] = 1

							serverUtils.MakeRemoteScanRequest(requestURL, bodyAttrs, http.StatusBadRequest)
							Eventually(session.Err).Should(gbytes.Say(errMessage))
						})
					})
				})
			})
		})

		Context("scan remote k8s git repo", func() {
			requestURL := fmt.Sprintf("%s:%s/v1/k8s/v1/all/remote/dir/scan", host, port)

			When("remote repo violates kubernetes_ingress", func() {
				It("should report violations", func() {
					remoteRepoURL := "https://github.com/accurics/terrascan//test/e2e/test_data/iac/k8s/kubernetes_ingress_violation"

					goldenFilePath, err := filepath.Abs(filepath.Join("..", "scan", "golden", "k8s_scans", "k8s", "kubernetes_ingress_violations", "kubernetes_ingress_json.txt"))
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
	})
})
