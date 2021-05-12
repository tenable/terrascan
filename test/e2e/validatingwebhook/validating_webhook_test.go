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

package validatingwebhook_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/accurics/terrascan/test/e2e/validatingwebhook"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
)

const (
	certsFolder      = "certs"
	k8sWebhookAPIKey = "K8S_WEBHOOK_API_KEY"
	apiKeyValue      = "accurics"
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

		// sleep added so that the default serviceaccount is get created
		// this logic needs to be improved
		time.Sleep(20 * time.Second)
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

	Describe("terrascan server as validating webhook with various available config options", func() {

		When("validating webhook with default 'k8s-admission-control' config", func() {

			Context("by default validating webhook runs in blind mode", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string

				It("server should start running on port 9010", func() {
					configFileName = "config1.toml"
					// create a config file with default config values
					err := validatingwebhook.CreateConfigFile(configFileName, policyRootRelPath, nil)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9010"))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(webhookYamlRelPath)
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, "9010")
						Expect(err).NotTo(HaveOccurred())
					})

					When("pod creation addmission requested is sent to server", func() {
						It("server should get the addmission request to review", func() {
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
				err := validatingwebhook.CreateConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
				Expect(err).NotTo(HaveOccurred())

				os.Setenv(k8sWebhookAPIKey, apiKeyValue)
				args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
				Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9011"))
			})

			When("request is made to add server as a validating webhook", func() {
				It("should get registered with k8s cluster as validating webhook successfully", func() {

					webhookFileAbsPath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
					Expect(err).NotTo(HaveOccurred())

					webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFileAbsPath, certFileAbsPath, apiKeyValue, port)
					Expect(err).NotTo(HaveOccurred())
				})

				When("pod creation addmission requested is sent to server", func() {
					It("server should get the addmission request to review", func() {
						// remove the config file
						defer os.Remove(configFileName)

						createPod(session, webhookConfig.GetName())
					})
				})
			})
		})

		When("validating webhook config has 'denied-severity' specified", func() {

			Context("service to be created violates a policy which has 'MEDIUM' seveirty", func() {
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
					err := validatingwebhook.CreateConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9012"))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation addmission requested is sent to server", func() {
						It("server should get the addmission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName())
						})
					})
				})
			})
		})

		When("validating webhook config has 'denied-categories' specified", func() {
			Context("service to be created violates a policy which has 'Network Security' category", func() {
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
							Categories: []string{"Network Security"},
						},
					}
					err := validatingwebhook.CreateConfigFile(configFileName, policyRootRelPath, &terrascanConfig)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookAPIKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9013"))
				})

				When("request is made to add server as a validating webhook", func() {
					It("should get registered with k8s cluster as validating webhook successfully", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join(webhookYamlRelPath))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
						Expect(err).NotTo(HaveOccurred())
					})

					When("service creation addmission requested is sent to server", func() {
						It("server should get the addmission request to review and reject the request", func() {
							// remove the config file
							defer os.Remove(configFileName)

							createService(session, webhookConfig.GetName())
						})
					})
				})
			})
		})
	})
})

// createService creates a service and asserts for reject status,
// and deletes the resources
func createService(session *gexec.Session, webhookName string) {
	serviceYamlAbsPath, err := filepath.Abs(filepath.Join(serviceYamlPath))
	Expect(err).NotTo(HaveOccurred())

	service, err := kubeClient.CreateService(serviceYamlAbsPath)
	Eventually(session.Err, defaultTimeout).Should(gbytes.Say("handle: validating webhook request"))
	Expect(err).To(HaveOccurred())

	if e, ok := err.(*k8serr.StatusError); ok {
		Expect(e.Status().Code).To(BeNumerically("==", 403))
	} else {
		errMessage := fmt.Sprintf("expected error to be of type 'k8s.io/apimachinery/pkg/api/errors.StatusError', got of type %T", err)
		Fail(errMessage)
	}
	Expect(service).To(BeNil())

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
