package server_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	serverUtils "github.com/accurics/terrascan/test/e2e/server"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

const (
	host                   string = "http://localhost"
	defaultPort            int    = 9010
	terrascanConfigEnvName string = "TERRASCAN_CONFIG"
	terrascanServerPort    string = "TERRASCAN_SERVER_PORT"
	configFileName         string = "configFile.toml"
)

var (
	terrascanBinaryPath string
)

var _ = Describe("Server", func() {

	BeforeSuite(func() {
		terrascanBinaryPath = helper.GetTerrascanBinaryPath()
		createAndSetEnvConfigFile(configFileName)
	})

	AfterSuite(func() {
		gexec.Terminate()
		os.Remove(configFileName)
	})

	Describe("server command's help test", func() {

		Context("server is run with -h flag", func() {
			It("should print help for server and exit with status code 0", func() {
				outWriter, errWriter := gbytes.NewBuffer(), gbytes.NewBuffer()
				serverArgs := []string{serverUtils.ServerCommand, "-h"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, serverArgs...)
				serverUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeZero, filepath.Join("..", "help", "golden", "help_server.txt"), true)
			})
		})

		Context("server command is run with a typo. eg: servre", func() {
			It("should print server command suggestion and exit with status code 1", func() {
				outWriter, errWriter := gbytes.NewBuffer(), gbytes.NewBuffer()
				serverArgs := []string{"servre"}
				session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, serverArgs...)
				serverUtils.ValidateExitCodeAndOutput(session, helper.ExitCodeOne, filepath.Join("golden", "server_typo_help.txt"), false)
			})
		})
	})

	Describe("server is started and killed", func() {
		var session *gexec.Session
		var outWriter, errWriter io.Writer = gbytes.NewBuffer(), gbytes.NewBuffer()
		Context("server is started without any arguments", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanServerPort, "")
			})
			It("should start the api server session", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, serverUtils.ServerCommand)
			})
			Context("by default server is running at port 9010", func() {
				It("server should be accepting requests", func() {
					// logs are written in StdErr
					Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("Route GET - /health"))
					Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("Route POST - /v1/{iac}/{iacVersion}/{cloud}/local/file/scan"))
					Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("Route POST - /v1/{iac}/{iacVersion}/{cloud}/remote/dir/scan"))
					Eventually(session.Err, serverUtils.ServerCommandTimeout).Should(gbytes.Say("http server listening at port 9010"))
				})

				Context("request with no body on all handlers", func() {
					healthCheckURL := fmt.Sprintf("%s:%d/health", host, defaultPort)
					terraformV12LocalScanURL := fmt.Sprintf("%s:%d/v1/terraform/v12/all/local/file/scan", host, defaultPort)
					terrformV12RemoteScanURL := fmt.Sprintf("%s:%d/v1/terraform/v12/aws/remote/dir/scan", host, defaultPort)

					When("health check request is made", func() {
						It("should get 200 OK response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodGet, healthCheckURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusOK))
						})
					})

					When("GET request on file scan handler is made", func() {
						It("should receive method not allowed response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodGet, terraformV12LocalScanURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusMethodNotAllowed))
						})
					})

					When("GET request on remote scan handler is made", func() {
						It("should receive method not allowed response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodGet, terrformV12RemoteScanURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusMethodNotAllowed))
						})
					})

					When("POST request on file scan handler is made without body", func() {
						It("should receive internal server error response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodPost, terraformV12LocalScanURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusInternalServerError))
						})
					})

					When("POST request on remote scan handler is made without body", func() {
						It("should receive bad request response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodPost, terrformV12RemoteScanURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusBadRequest))
						})
					})

					When("POST request on health check", func() {
						It("should receive method not allowed response", func() {
							r, err := serverUtils.MakeHTTPRequest(http.MethodPost, healthCheckURL)
							Expect(err).NotTo(HaveOccurred())
							defer r.Body.Close()
							Expect(r).NotTo(BeNil())
							Expect(r.StatusCode).To(BeIdenticalTo(http.StatusMethodNotAllowed))
						})
					})

					Context("server is stopped", func() {
						It("should gracefully exit", func() {
							session.Terminate()
							_, err := serverUtils.MakeHTTPRequest(http.MethodGet, healthCheckURL)
							Expect(err).To(HaveOccurred())
						})
					})
				})
			})
		})
	})
})

// createAndSetEnvConfigFile creates a config file with test policy path
// and sets and env variable
func createAndSetEnvConfigFile(configFileName string) {
	var policyAbsPath, _ = filepath.Abs(filepath.Join("..", "test_data", "policies"))

	// contents of the config file
	configFileContents := fmt.Sprintf(`[policy]
path = "%s"
rego_subdir = "%s"`, policyAbsPath, policyAbsPath)

	// create config file in work directory
	file, err := os.Create(configFileName)
	Expect(err).NotTo(HaveOccurred())
	_, err = file.WriteString(configFileContents)
	Expect(err).NotTo(HaveOccurred())
	os.Setenv(terrascanConfigEnvName, configFileName)
}
