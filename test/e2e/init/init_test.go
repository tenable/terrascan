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

package init_test

import (
	"io"
	"os"
	"path/filepath"
	"time"

	initUtil "github.com/accurics/terrascan/test/e2e/init"
	"github.com/accurics/terrascan/test/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"gopkg.in/src-d/go-git.v4"
)

var (
	initCommand            string = "init"
	defaultPolicyRepoPath  string = os.Getenv("HOME") + "/.terrascan"
	terrascanGitURL        string = "https://github.com/accurics/terrascan.git"
	terrascanDefaultBranch string = "master"
	terrascanConfigEnvName string = "TERRASCAN_CONFIG"
	kaiMoneyGitURL         string = "https://github.com/accurics/KaiMonkey.git"
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
				goldenFileAbsPath, err := filepath.Abs("golden/init_help.txt")
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
				goldenFileAbsPath, err := filepath.Abs("golden/init_typo_help.txt")
				Expect(err).NotTo(HaveOccurred())
				helper.CompareActualWithGolden(session, goldenFileAbsPath, false)
			})

			It("should exit with status code 1", func() {
				session = helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "inti")
				Eventually(session, 5).Should(gexec.Exit(helper.ExitCodeOne))
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
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeOne)
				helper.ContainsErrorSubString(session, `failed to download policies. error: 'Get "https://repository/url/info/refs?service=git-upload-pack": dial tcp:`)
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
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeOne)
				helper.ContainsErrorSubString(session, `failed to initialize terrascan. error : failed to checkout git branch 'invalid-branch'. error: 'reference not found'`)
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
				session = initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeOne)
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
			It("should download policies and exit with status code 0", func() {
				initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
			})
		})
		Context("the config file has valid data", func() {
			When("config file has different git repo and branch", func() {
				JustBeforeEach(func() {
					os.Setenv(terrascanConfigEnvName, "config/valid_config.toml")
				})
				JustAfterEach(func() {
					os.Setenv(terrascanConfigEnvName, "")
				})
				It("init should download the repo provided in the config file", func() {
					initUtil.RunInitCommand(terrascanBinaryPath, outWriter, errWriter, helper.ExitCodeZero)
				})
				Context("Kai Monkey git repo is downloaded", func() {
					It("should validate Kai Monkey repo in the policy path", func() {
						repo := initUtil.OpenGitRepo(defaultPolicyRepoPath)
						initUtil.ValidateGitRepo(repo, kaiMoneyGitURL)
					})
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
