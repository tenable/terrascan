package init_test

import (
	"io"
	"os"
	"path/filepath"

	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"gopkg.in/src-d/go-git.v4"
)

const (
	initCommandTimeout = 60
)

var (
	defaultPolicyRepoPath  string = os.Getenv("HOME") + "/.terrascan"
	terrascanGitURL        string = "https://github.com/accurics/terrascan.git"
	terrascanDefaultBranch string = "master"
	terrascanConfigEnvName string = "TERRASCAN_CONFIG"
)

var _ = Describe("Init", func() {
	var session *gexec.Session
	var terrascanBinaryPath string

	var outWriter, errWriter io.Writer

	BeforeSuite(func() {
		os.RemoveAll(defaultPolicyRepoPath)
		terrascanBinaryPath = helper.GetTerrascanBinaryPath()
	})

	BeforeEach(func() {
		outWriter = gbytes.NewBuffer()
		errWriter = gbytes.NewBuffer()
	})

	AfterEach(func() {
		outWriter = nil
		errWriter = nil
	})

	Describe("terrascan init is run", func() {
		When("terrascan init is run without any flags", func() {
			It("should download policies and exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init")
				Eventually(session, initCommandTimeout).Should(gexec.Exit(0))
				Expect(outWriter).Should(gbytes.Say(""))
			})

			It("should download policies at TERRASCAN's base location", func() {
				BeADirectory().Match(defaultPolicyRepoPath)
			})

			Context("git repo should be validated", func() {
				var repo *git.Repository
				var err error
				It("should be a valid git repo", func() {
					repo, err = git.PlainOpen(defaultPolicyRepoPath)
					Expect(err).NotTo(HaveOccurred())
					Expect(repo).NotTo(BeNil())
				})
				It("should be terrascan git repo", func() {
					remote, err := repo.Remote("origin")
					Expect(err).NotTo(HaveOccurred())
					Expect(remote).NotTo(BeNil())
					remoteConfig := remote.Config()
					Expect(remoteConfig).NotTo(BeNil())
					err = remoteConfig.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(remoteConfig.URLs[0]).To(BeEquivalentTo(terrascanGitURL))
				})
				It("master branch should be present", func() {
					_, err = repo.Branch(terrascanDefaultBranch)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		When("terrascan init is run with -h flag", func() {
			It("should print help", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-h")
				goldenFileAbsPath, err := filepath.Abs("golden/init_help.txt")
				Expect(err).NotTo(HaveOccurred())
				helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
			})

			It("should exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-h")
				Eventually(session).Should(gexec.Exit(0))
			})
		})

		When("terrascan init command has typo. eg: inti", func() {
			It("should print command suggestion", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "inti")
				goldenFileAbsPath, err := filepath.Abs("golden/init_typo_help.txt")
				Expect(err).NotTo(HaveOccurred())
				helper.CompareActualWithGolden(session, goldenFileAbsPath, false)
			})

			It("should exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "inti")
				Eventually(session, 5).Should(gexec.Exit(1))
			})
		})

		When("terrascan init is run with -c flag", func() {
			Context("config file has valid policy config data", func() {
				It("should download policies as per the policy config in the config file", func() {
					Skip("skipping this test due to https://github.com/accurics/terrascan/issues/550, should be implemented when fixed")
				})
			})
		})
	})

	Describe("terrascan init is run when TERRASCAN_CONFIG is set", func() {
		When("the config file has invalid repo url", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, "config/invalid_repo.toml")
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init")
				Eventually(session, initCommandTimeout).Should(gexec.Exit(1))
				helper.ContainsErrorSubString(session, `failed to download policies. error: 'Get "https://repository/url/info/refs?service=git-upload-pack": dial tcp: lookup repository on 8.8.8.8:53: no such host'`)
			})
		})
		When("the config file has invalid branch name", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, "config/invalid_branch.toml")
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init")
				Eventually(session, initCommandTimeout).Should(gexec.Exit(1))
				helper.ContainsErrorSubString(session, `failed to checkout branch 'invalid-branch'. error: 'reference not found'`)
			})
		})
		When("the config file has invalid rego subdir", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, "config/invalid_rego_subdir.toml")
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init")
				Eventually(session, initCommandTimeout).Should(gexec.Exit(1))
				helper.ContainsErrorSubString(session, "invalid/path: no such file or directory")
			})
		})
		When("the config file has invalid path", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, "config/invalid_path.toml")
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				Skip("Skipping invalid path test until discussion with team")
			})
		})
	})
})
