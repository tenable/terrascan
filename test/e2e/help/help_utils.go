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

package help

import (
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/test/helper"
)

// ValidateExitCodeAndOutput validates the exit code and output of the command
func ValidateExitCodeAndOutput(session *gexec.Session, exitCode int, relFilePath string, isStdOut bool) {
	gomega.Eventually(session).Should(gexec.Exit(exitCode))
	goldenFileAbsPath, err := filepath.Abs(relFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	helper.CompareActualWithGolden(session, goldenFileAbsPath, isStdOut)
}
