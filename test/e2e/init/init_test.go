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

package init_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/utils"
	initUtil "github.com/tenable/terrascan/test/e2e/init"
	"github.com/tenable/terrascan/test/helper"
)

var (
	initCommand            = "init"
	defaultPolicyRepoPath  = filepath.Join(utils.GetHomeDir(), ".terrascan")
	terrascanGitURL        = "https://github.com/tenable/terrascan.git"
	terrascanDefaultBranch = "master"
	terrascanConfigEnvName = "TERRASCAN_CONFIG"
	kaiMoneyGitURL         = "https://github.com/tenable/KaiMonkey.git"

	testPolicyRepoPath = filepath.Join(utils.GetHomeDir(), ".terrascan-test")
	testRegoSubDirPath = filepath.Join(testPolicyRepoPath, "pkg", "policies", "opa", "rego")
	warnNoBasePath     = fmt.Sprintf("policy rego_subdir specified in configfile '%s', but base path not specified. applying default base path value", filepath.Join("config", "relative_rego_subdir.toml"))
	warnNoSubDirPath   = fmt.Sprintf("policy base path specified in configfile '%s', but rego_subdir path not specified. applying default rego_subdir value", filepath.Join("config", "home_prefixed_path.toml"))
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
		When("without any flags", func() {

			It("should download policies and exit with status code 0", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
				Expect(outWriter).Should(gbytes.Say(""))
			})

			It("should download policies at TERRASCAN's base location", func() {
				BeADirectory().Match(defaultPolicyRepoPath)
			})

			Context("git repo should be validated", func() {
				var repo *git.Repository
				var err error
				It("should be a valid git repo", func() {
					repo = initUtil.OpenGitRepo(defaultPolicyRepoPath)
				})
				It("should be terrascan git repo", func() {
					initUtil.ValidateGitRepo(repo, terrascanGitURL)
				})
				It("master branch should be present", func() {
					_, err = repo.Branch(terrascanDefaultBranch)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		When("terrascan init is run with -h flag", func() {
			It("should print help", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, initCommand, "-h")
				goldenFileAbsPath, err := filepath.Abs(filepath.Join("..", "help", "golden", "help_init.txt"))
				Expect(err).NotTo(HaveOccurred())
				helper.CompareActualWithGolden(session, goldenFileAbsPath, true)
			})

			It("should exit with status code 0", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, initCommand, "-h")
				Eventually(session).Should(gexec.Exit(helper.ExitCodeZero))
			})
		})

		When("terrascan init command has typo. eg: inti", func() {
			It("should print command suggestion", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "inti")
				goldenFileAbsPath, err := filepath.Abs(filepath.Join("golden", "init_typo_help.txt"))
				Expect(err).NotTo(HaveOccurred())
				helper.CompareActualWithGolden(session, goldenFileAbsPath, false)
			})

			It("should exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "inti")
				Eventually(session, 5).Should(gexec.Exit(helper.ExitCodeOne))
			})
		})
	})

	Describe("terrascan init is run with -c flag", func() {

		Context("config file has valid policy repo and branch data", func() {
			It("should download policies as per the policy config in the config file", func() {
				configFile := filepath.Join("config", "valid_repo.toml")
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-c", configFile)
				helper.ValidateExitCode(session, initUtil.InitCommandTimeout, helper.ExitCodeZero)
			})

			Context("KaiMonkey git repo is downloaded", func() {
				It("should validate KaiMonkey repo in the policy path", func() {
					repo := initUtil.OpenGitRepo(defaultPolicyRepoPath)
					initUtil.ValidateGitRepo(repo, kaiMoneyGitURL)
				})
				os.RemoveAll(defaultPolicyRepoPath)
			})

		})

		Context("config file has valid policy path and rego_subdir data", func() {
			It("should download policies as per the policy config in the config file", func() {
				configFile := filepath.Join("config", "valid_paths.toml")
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-c", configFile)
				helper.ValidateExitCode(session, initUtil.InitCommandTimeout, helper.ExitCodeZero)
			})

			It("should validate terrascan repo in the policy path", func() {
				repo := initUtil.OpenGitRepo(testPolicyRepoPath)
				initUtil.ValidateGitRepo(repo, terrascanGitURL)
			})

			os.RemoveAll(testPolicyRepoPath)
		})

		Context("config file has all valid policy config data", func() {

			It("should download policies as per the policy config in the config file", func() {
				configFile := filepath.Join("config", "valid_config.toml")
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-c", configFile)
				helper.ValidateExitCode(session, initUtil.InitCommandTimeout, helper.ExitCodeZero)
			})

			Context("terrascan git repo is downloaded", func() {
				It("should validate terrascan repo in the policy path", func() {
					repo := initUtil.OpenGitRepo(testPolicyRepoPath)
					initUtil.ValidateGitRepo(repo, terrascanGitURL)
					helper.ValidateDirectoryExists(testRegoSubDirPath)
				})
			})

			os.RemoveAll(testPolicyRepoPath)
		})

		Context("config file has all valid policy paths with ~ prefix base path", func() {

			It("should download policies as per the policy config in the config file", func() {
				configFile := filepath.Join("config", "home_prefix_path_config.toml")
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init", "-c", configFile)
				helper.ValidateExitCode(session, initUtil.InitCommandTimeout, helper.ExitCodeZero)
			})

			Context("terrascan git repo is downloaded", func() {
				It("should validate terrascan repo in the policy path", func() {
					repo := initUtil.OpenGitRepo(testPolicyRepoPath)
					initUtil.ValidateGitRepo(repo, terrascanGitURL)
					subpath := testRegoSubDirPath //filepath.Join(path, "pkg/policies/opa/rego")
					helper.ValidateDirectoryExists(subpath)
				})
			})

			os.RemoveAll(testPolicyRepoPath)
		})

	})

	Describe("terrascan init is run when TERRASCAN_CONFIG is set", func() {
		When("the config file has invalid repo url", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "invalid_repo.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeOne)
				helper.ContainsErrorSubString(session, `failed to download policies. error: 'Get "https://repository/url/info/refs?service=git-upload-pack": dial tcp:`)
			})
		})
		When("the config file has invalid branch name", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "invalid_branch.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should error out and exit with status code 1", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeOne)
				helper.ContainsErrorSubString(session, `failed to initialize terrascan. error : failed to checkout git branch 'invalid-branch'. error: 'reference not found'`)
			})
		})
		When("the config file has relative rego subdir", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "relative_rego_subdir.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})

			// The current behavior of terrascan is that, in case of init command, even if the value of
			// rego_subdir is an invalid/non-existant directory, the init is successful and repoURL will be
			// cloned at the base path (either default or based on config file)
			It("should log a warning & download the policies at default base path", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
				repo := initUtil.OpenGitRepo(defaultPolicyRepoPath)
				initUtil.ValidateGitRepo(repo, terrascanGitURL)
			})

			It("should log a warning stating no base path specified", func() {
				helper.ContainsErrorSubString(session, warnNoBasePath)
			})
		})

		When("the config file has relative path", func() {
			path, err := utils.GetAbsPath("policy/base_path")

			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "relative_path.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})

			// The current behavior of terrascan is that, when init command is being run with an invalid/
			// non-existant base path, the specified path gets created and repoURL is cloned at that location
			It("should work fine and give out exit code zero", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
			})
			It("should download the policy repo at the specified path (relative to the cwd)", func() {
				Expect(err).ToNot(HaveOccurred())
				repo := initUtil.OpenGitRepo(path)
				initUtil.ValidateGitRepo(repo, terrascanGitURL)
				os.RemoveAll(path)
			})
		})

		When("the config file has relative path with kai monkey repository specified", func() {
			path, err := utils.GetAbsPath("policy/base_path")

			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "kai_monkey_relative_path.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})
			It("should work fine and give out exit code zero", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
			})
			It("should download the policy repo at the specified path (relative to the cwd)", func() {
				Expect(err).ToNot(HaveOccurred())
				repo := initUtil.OpenGitRepo(path)
				initUtil.ValidateGitRepo(repo, kaiMoneyGitURL)
				os.RemoveAll(path)
			})

		})

		When("the config file has a ~ prefixed path and no rego_subdir", func() {
			JustBeforeEach(func() {
				os.Setenv(terrascanConfigEnvName, filepath.Join("config", "home_prefixed_path.toml"))
			})
			JustAfterEach(func() {
				os.Setenv(terrascanConfigEnvName, "")
			})

			It("should download the policies at $HOME/<path>", func() {
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
				repo := initUtil.OpenGitRepo(testPolicyRepoPath)
				initUtil.ValidateGitRepo(repo, terrascanGitURL)
			})

			It("should log a warning stating no rego_subdir specified", func() {
				helper.ContainsErrorSubString(session, warnNoSubDirPath)
			})

			os.RemoveAll(testPolicyRepoPath)
		})

		Context("the config file has valid data", func() {
			When("config file has different git repo and branch", func() {
				JustBeforeEach(func() {
					os.Setenv(terrascanConfigEnvName, filepath.Join("config", "valid_repo.toml"))
				})
				JustAfterEach(func() {
					os.Setenv(terrascanConfigEnvName, "")
				})
				It("init should download the repo provided in the config file", func() {
					initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
					// Kai Monkey git repo is downloaded
					// validate Kai Monkey repo in the repo path
					repo := initUtil.OpenGitRepo(defaultPolicyRepoPath)
					initUtil.ValidateGitRepo(repo, kaiMoneyGitURL)
				})
			})
		})

	})

	Describe("terrascan init is run multiple times", func() {
		Context("init clones the git repo to a temp dir, deletes policy path and renames tempdir to policy path", func() {
			Context("running init the first time", func() {
				var modifiedTime time.Time
				It("should download policies at the default policy path", func() {
					initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
					fi, err := os.Stat(defaultPolicyRepoPath)
					Expect(err).ToNot(HaveOccurred())
					modifiedTime = fi.ModTime()
				})
				Context("running init the second time", func() {
					It("should download policies again at the default policy path", func() {
						initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
						fi, err := os.Stat(defaultPolicyRepoPath)
						Expect(err).ToNot(HaveOccurred())
						Expect(fi.ModTime()).To(BeTemporally(">", modifiedTime))
					})
				})
			})
		})
	})
})
