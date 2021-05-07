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
)

const (
	certsFolder      = "certs"
	k8sWebhookApiKey = "K8S_WEBHOOK_API_KEY"
	apiKeyValue      = "accurics"
	defaultTimeout   = 10
)

var (
	kubeClient          *validatingwebhook.KubernetesClient
	terrascanBinaryPath string
	certFileAbsPath     string
	privKeyFileAbsPath  string
	policyRootRelPath   = filepath.Join("..", "test_data", "policies")
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

		os.RemoveAll("certs")
	})

	Describe("terrascan server as validating webhook with various available config options", func() {

		When("validating webhook with default 'k8s-admission-control' config", func() {

			Context("by default validating webhook runs in blind mode", func() {
				var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
				var session *gexec.Session
				var webhookConfig *admissionv1.ValidatingWebhookConfiguration
				var configFileName string

				It("should run server successfully", func() {
					configFileName = "config1.toml"
					// create a config file with default config values
					err := validatingwebhook.CreateConfigFile(configFileName, policyRootRelPath, nil)
					Expect(err).NotTo(HaveOccurred())

					os.Setenv(k8sWebhookApiKey, apiKeyValue)
					args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath}
					session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route GET - /health"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route POST - /v1/{iac}/{iacVersion}/{cloud}/local/file/scan"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route POST - /v1/{iac}/{iacVersion}/{cloud}/remote/dir/scan"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route GET - /k8s/webhooks/{apiKey}/logs"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route GET - /k8s/webhooks/logs/{uid}"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("Route POST - /v1/k8s/webhooks/{apiKey}/scan/validate"))
					Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9010"))
				})

				When("server is added as a validating webhook and resource is created", func() {
					It("should get registered with k8s cluster as validating webhook", func() {

						webhookFilePath, err := filepath.Abs(filepath.Join("test-data", "yamls", "webhook.yaml"))
						Expect(err).NotTo(HaveOccurred())

						webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, "9010")
						Expect(err).NotTo(HaveOccurred())
					})

					When("resource creation is requested, it should be handled by terrascan server", func() {
						It("should be handled by terrascan server", func() {
							// remove the config file
							defer os.Remove(configFileName)

							podYamlAbsPath, err := filepath.Abs(filepath.Join("test-data", "yamls", "pod.yaml"))
							Expect(err).NotTo(HaveOccurred())

							pod, err := kubeClient.CreatePod(podYamlAbsPath)
							Eventually(session.Err, 10).Should(gbytes.Say("handle: validating webhook request"))
							Expect(err).NotTo(HaveOccurred())
							Expect(pod).NotTo(BeNil())

							// delete pod
							err = kubeClient.DeletePod(pod.GetName())
							Expect(err).NotTo(HaveOccurred())

							// delete validating webhook configuration
							err = kubeClient.DeleteValidatingWebhookConfiguration(webhookConfig.GetName())
							Expect(err).NotTo(HaveOccurred())

							if utils.IsWindowsPlatform() {
								session.Kill()
							} else {
								session.Terminate()
							}
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

			It("should run server successfully", func() {
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

				os.Setenv(k8sWebhookApiKey, apiKeyValue)
				args := []string{"server", "-c", configFileName, "--cert-path", certFileAbsPath, "--key-path", privKeyFileAbsPath, "-p", port}
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, args...)
				Eventually(session.Err, defaultTimeout).Should(gbytes.Say("http server listening at port 9011"))
			})

			When("server is added as a validating webhook and resource is created", func() {
				It("should get registered with k8s cluster as validating webhook", func() {

					webhookFilePath, err := filepath.Abs(filepath.Join("test-data", "yamls", "webhook.yaml"))
					Expect(err).NotTo(HaveOccurred())

					webhookConfig, err = kubeClient.CreateValidatingWebhookConfiguration(webhookFilePath, certFileAbsPath, apiKeyValue, port)
					Expect(err).NotTo(HaveOccurred())
				})

				When("resource creation is requested, it should be handled by terrascan server", func() {
					It("should be handled by terrascan server", func() {
						// remove the config file
						defer os.Remove(configFileName)

						podYamlAbsPath, err := filepath.Abs(filepath.Join("test-data", "yamls", "pod.yaml"))
						Expect(err).NotTo(HaveOccurred())

						pod, err := kubeClient.CreatePod(podYamlAbsPath)
						Eventually(session.Err, 10).Should(gbytes.Say("handle: validating webhook request"))
						Expect(err).NotTo(HaveOccurred())
						Expect(pod).NotTo(BeNil())

						// delete pod
						err = kubeClient.DeletePod(pod.GetName())
						Expect(err).NotTo(HaveOccurred())

						// delete validating webhook configuration
						err = kubeClient.DeleteValidatingWebhookConfiguration(webhookConfig.GetName())
						Expect(err).NotTo(HaveOccurred())

						if utils.IsWindowsPlatform() {
							session.Kill()
						} else {
							session.Terminate()
						}
					})
				})
			})
		})

		When("validating webhook config has 'denied-severity' specified", func() {

		})

		When("validating webhook config has 'denied-categories' specified", func() {

		})
	})
})
