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

package validatingwebhook_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/utils"
	"github.com/tenable/terrascan/test/e2e/server"
	"github.com/tenable/terrascan/test/e2e/validatingwebhook"
	"github.com/tenable/terrascan/test/helper"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
)

const (
	certsFolder      = "certs"
	k8sWebhookAPIKey = "K8S_WEBHOOK_API_KEY"
	apiKeyValue      = "tenable"
	defaultTimeout   = 10
)

var (
	kubeClient          *validatingwebhook.KubernetesClient
	terrascanBinaryPath string
	certFileAbsPath     string
	privKeyFileAbsPath  string
	policyRootRelPath   = filepath.Join("..", "test_data", "policies")
	webhookYamlRelPath  = filepath.Join("test-data", "yamls", "webhook.yaml")
	podYamlRelPath      = filepath.Join("test-data", "yamls", "pod.yaml")
	serviceYamlPath     = filepath.Join("test-data", "yamls", "service.yaml")
)

var _ = Describe("ValidatingWebhook", func() {

	BeforeSuite(func() {

		// delete the default cluster if it is already running
		err := validatingwebhook.DeleteDefaultKindCluster()
		if err != nil {
			message := fmt.Sprintf("error while deleting cluster. err: %v", err)
			Fail(message)
		}

		// create a new cluster
		err = validatingwebhook.CreateDefaultKindCluster()
		if err != nil {
			message := fmt.Sprintf("error while creating cluster. err: %v", err)
			Fail(message)
		}

		// get k8s client
		kubeClient, err = validatingwebhook.NewKubernetesClient()
		if err != nil {
			errMessage := fmt.Sprintf("failed to connected to default k8s cluster, error: %s", err.Error())
			Fail(errMessage)
		}

		// get terrascan binary path
		terrascanBinaryPath = helper.GetTerrascanBinaryPath()

		// create tls certificates for server
		certFileAbsPath, privKeyFileAbsPath, err = validatingwebhook.CreateCertificate(certsFolder, "server.crt", "priv.key")
		if err != nil {
			errMessage := fmt.Sprintf("failed to create certificates, error: %s", err.Error())
			Fail(errMessage)
		}

		// wait 1 minute max for the service account to get created
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		err = kubeClient.WaitForServiceAccount(ctx)
		if err != nil {
			errMessage := fmt.Sprintf("service account for default namespace not created after 1 minute, error: %s", err.Error())
			Fail(errMessage)
		}
	})

	AfterSuite(func() {
		if utils.IsWindowsPlatform() {
			gexec.Kill()
		} else {
			gexec.Terminate()
		}

		// delete the cluster
		err := validatingwebhook.DeleteDefaultKindCluster()
		if err != nil {
			message := fmt.Sprintf("error while deleting cluster. err: %v", err)
			Fail(message)
		}

		os.RemoveAll(certsFolder)
	})

	// log message to be asserted when sever starts
	assertLogMessage := "https server listening at port %s"

	Describe("terrascan server as validating webhook with various available config options", func() {

		When("validating webhook with default 'k8s-admission-control' config", func() {

			Context("by default validating webhook runs in blind mode", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string
				port := "9010"

				It("server should start running on port 9010", func() {
					configFileName = "config1.toml"
					// create a config file with default config values
					err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, nil)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-l", "debug"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
				})

				Context("in blind mode, log end points return error response", func() {
					myIP, err := validatingwebhook.GetIP()
					When("request is made to the get a single log endpoint", func() {
						It("should return a 400 bad request response", func() {
							Expect(err).NotTo(HaveOccurred())

							requestURL := fmt.Sprintf("http://%s:%s/k8s/webhooks/logs/%s", myIP.To4(), port, apiKeyValue)
							resp, err := server.MakeHTTPRequest("GET", requestURL)
							Expect(err).NotTo(HaveOccurred())
							Expect(resp).NotTo(BeNil())
							Expect(resp.StatusCode).To(BeIdenticalTo(http.StatusBadRequest))
						})
					})

					When("request is made to the get a single log endpoint", func() {
						It("should return a 400 bad request response", func() {
							Expect(err).NotTo(HaveOccurred())

							requestURL := fmt.Sprintf("http://%s:%s/k8s/webhooks/%s/logs", myIP.To4(), port, apiKeyValue)
							resp, err := server.MakeHTTPRequest("GET", requestURL)
							Expect(err).NotTo(HaveOccurred())
							Expect(resp).NotTo(BeNil())
							Expect(resp.StatusCode).To(BeIdenticalTo(http.StatusBadRequest))
						})
					})
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(webhookYamlRelPath)
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, "9010")
						Expect(err).NotTo(HaveOccurred())
					})

					When("pod creation admission requested is sent to server", func() {
						It("server should get the admission request to review", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createPod(session, webhookConfig.GetName())
						})
					})
				})
			})
		})

		When("validating webhook config has 'dashboard' and 'save-requests' is enabled", func() {
			var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
			var session *gexec.Session
			var webhookConfig *admissionv1.ValidatingWebhookConfiguration
			var configFileName string
			var port string

			It("server should start running on port 9011", func() {
				port = "9011"
				configFileName = "config2.toml"

				// create a config file with 'dashboard' set to true
				terrascanConfig := config.TerrascanConfig{
					K8sAdmissionControl: config.K8sAdmissionControl{
						Dashboard:    true,
						SaveRequests: true,
					},
				}
				err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
				Expect(err).NotTo(HaveOccurred())

				os.Setenv(k8sWebhookAPIKey, apiKeyValue)
				args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port, "-l", "debug"}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
				Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
			})

			When("request is made to add server as a validating webhook", func() {
				It("should get registered with k8s cluster as validating webhook successfully", func() {

					webhookFileAbsPath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
					Expect(err).NotTo(HaveOccurred())

					webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFileAbsPath, certFileAbsPath, apiKeyValue, port)
					Expect(err).NotTo(HaveOccurred())
				})

				When("pod creation admission requested is sent to server", func() {
					It("server should get the admission request to review", func() {
						// remove the config file
						defer os.Remove(configFileName)

						createPod(session, webhookConfig.GetName())
					})
				})
			})
		})

		When("validating webhook config has 'denied-severity' specified", func() {

			Context("service to be created violates a policy with specified denied severity", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string
				var port string

				It("server should start running on port 9012", func() {
					port = "9012"
					configFileName = "config3.toml"

					// create a config file with desired severity specified
					terrascanConfig := config.TerrascanConfig{
						K8sAdmissionControl: config.K8sAdmissionControl{
							DeniedSeverity: "MEDIUM",
						},
					}
					err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port, "-l", "debug"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation admission requested is sent to server", func() {
						It("server should get the admission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName(), true)
						})
					})
				})
			})

			Context("service to be created violates a policy which doesn't have the desired severity", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string
				var port string

				It("server should start running on port 9013", func() {
					port = "9013"
					configFileName = "config4.toml"

					// create a config file with desired severity specified
					terrascanConfig := config.TerrascanConfig{
						K8sAdmissionControl: config.K8sAdmissionControl{
							DeniedSeverity: "HIGH",
						},
					}
					err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port, "-l", "debug"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation admission requested is sent to server", func() {
						It("server should get the admission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName(), false)
						})
					})
				})
			})
		})

		When("validating webhook config has 'denied-categories' specified", func() {
			Context("service to be created violates the denied category", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string
				var port string

				It("server should start running on port 9014", func() {
					port = "9014"
					configFileName = "config5.toml"

					// create a config file with desired severity specified
					terrascanConfig := config.TerrascanConfig{
						K8sAdmissionControl: config.K8sAdmissionControl{
							Categories: []string{"Network Security"},
						},
					}
					err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port, "-l", "debug"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation admission requested is sent to server", func() {
						It("server should get the admission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName(), true)
						})
					})
				})
			})

			Context("service to be created does not violate the denied category", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string
				var port string

				It("server should start running on port 9015", func() {
					port = "9015"
					configFileName = "config6.toml"

					// create a config file with desired severity specified
					terrascanConfig := config.TerrascanConfig{
						K8sAdmissionControl: config.K8sAdmissionControl{
							Categories: []string{"Doesn't Exist"},
						},
					}
					err := validatingwebhook.CreateTerrascanConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port, "-l", "debug"}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say(fmt.Sprintf(assertLogMessage, port)))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation admission requested is sent to server", func() {
						It("server should get the admission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName(), false)
						})
					})
				})
			})
		})
	})
})

// createService creates a service and asserts for reject status,
// and deletes the resources
func createService(session *gexec.Session, webhookName string, shouldBeDenied bool) {
	serviceYamlAbsPath, err := filepath.Abs(filepath.Join(serviceYamlPath))
	Expect(err).NotTo(HaveOccurred())

	service, err := kubeClient.CreateService(serviceYamlAbsPath)
	Eventually(session.Err, defaultTimeout).Should(gbytes.Say("handle: validating webhook request"))

	if shouldBeDenied {
		Expect(err).To(HaveOccurred())
		if e, ok := err.(*k8serr.StatusError); ok {
			Expect(e.Status().Code).To(BeNumerically("==", 403))
		} else {
			errMessage := fmt.Sprintf("expected error to be of type 'k8s.io/apimachinery/pkg/api/errors.StatusError', got of type %T", err)
			Fail(errMessage)
		}
		Expect(service).To(BeNil())
	} else {
		Expect(err).NotTo(HaveOccurred())

		err = kubeClient.DeleteService(service.GetName())
		Expect(err).NotTo(HaveOccurred())
	}

	// delete validating webhook configuration
	err = kubeClient.DeleteValidatingWebhookConfiguration(webhookName)
	Expect(err).NotTo(HaveOccurred())

	if utils.IsWindowsPlatform() {
		session.Kill()
	} else {
		session.Terminate()
	}
}

// createPod creates a pod and asserts for reject status,
// and deletes the resources
func createPod(session *gexec.Session, webhookName string) {
	podYamlAbsPath, err := filepath.Abs(filepath.Join(podYamlRelPath))
	Expect(err).NotTo(HaveOccurred())

	pod, err := kubeClient.CreatePod(podYamlAbsPath)
	Eventually(session.Err, defaultTimeout).Should(gbytes.Say("handle: validating webhook request"))
	Expect(err).NotTo(HaveOccurred())
	Expect(pod).NotTo(BeNil())

	// delete pod
	err = kubeClient.DeletePod(pod.GetName())
	Expect(err).NotTo(HaveOccurred())

	// delete validating webhook configuration
	err = kubeClient.DeleteValidatingWebhookConfiguration(webhookName)
	Expect(err).NotTo(HaveOccurred())

	if utils.IsWindowsPlatform() {
		session.Kill()
	} else {
		session.Terminate()
	}
}
