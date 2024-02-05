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

package init

import (
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/test/helper"
)

// InitCommandTimeout timeout
const InitCommandTimeout = 60

// RunInitCommand will execute the init command and verify exit code
func RunInitCommand(terrascanBinaryPath string, outWriter, errWriter io.Writer, exitCode int) *gexec.Session {
	session := helper.RunCommand(terrascanBinaryPath, outWriter, errWriter, "init")
	gomega.Eventually(session, InitCommandTimeout).Should(gexec.Exit(exitCode))
	return session
}

// OpenGitRepo checks if a directory is a git repo
func OpenGitRepo(repoPath string) *git.Repository {
	repo, err := git.PlainOpen(repoPath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(repo).NotTo(gomega.BeNil())
	return repo
}

// ValidateGitRepo validates a git repo and verifies the git url
func ValidateGitRepo(repo *git.Repository, gitURL string) {
	remote, err := repo.Remote("origin")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(remote).NotTo(gomega.BeNil())
	remoteConfig := remote.Config()
	gomega.Expect(remoteConfig).NotTo(gomega.BeNil())
	err = remoteConfig.Validate()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(remoteConfig.URLs[0]).To(gomega.BeEquivalentTo(gitURL))
}
