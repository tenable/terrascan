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

package server

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/tenable/terrascan/pkg/policy"
	"github.com/tenable/terrascan/test/helper"
)

const (
	// ServerCommandTimeout is the default time for server command
	ServerCommandTimeout int = 10

	// ServerCommand is terrascan's server command
	ServerCommand string = "server"
)

// ValidateExitCodeAndOutput validates the exit code and output of the command
func ValidateExitCodeAndOutput(session *gexec.Session, exitCode int, relFilePath string, isStdOut bool) {
	gomega.Eventually(session).Should(gexec.Exit(exitCode))
	goldenFileAbsPath, err := filepath.Abs(relFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	helper.CompareActualWithGolden(session, goldenFileAbsPath, isStdOut)
}

// MakeHTTPRequest calls health handler of the api server
func MakeHTTPRequest(method, URL string) (*http.Response, error) {
	r, err := http.NewRequest(method, URL, nil)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	resp, err := http.DefaultClient.Do(r)
	return resp, err
}

// MakeFileScanRequest calls the file scan handler of api server
func MakeFileScanRequest(iacFilePath, URL string, bodyAttributes map[string]string, expectedStatusCode int) []byte {
	// open the iac file to upload
	iacFile, err := os.Open(iacFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	defer iacFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// create a multiple part form file
	part, err := writer.CreateFormFile("file", filepath.Base(iacFilePath))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	_, err = io.Copy(part, iacFile)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// add other body attributes to the request
	for k, v := range bodyAttributes {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	req, err := http.NewRequest(http.MethodPost, URL, body)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// make the http request
	resp, err := http.DefaultClient.Do(req)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(resp).NotTo(gomega.BeNil())

	// assert that the received response code matches expected status code
	gomega.Expect(resp.StatusCode).To(gomega.BeIdenticalTo(expectedStatusCode))
	defer resp.Body.Close()

	// read response body
	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return respBody.Bytes()
}

// MakeRemoteScanRequest calls the file scan handler of api server
func MakeRemoteScanRequest(URL string, bodyAttributes map[string]interface{}, expectedStatusCode int) []byte {
	// add body params
	reqBody, err := json.Marshal(bodyAttributes)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// create http post request with body params
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(reqBody))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	req.Header.Set("Content-Type", "application/json")

	// make http request
	resp, err := http.DefaultClient.Do(req)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(resp).NotTo(gomega.BeNil())

	// assert that the received response code matches expected status code
	gomega.Expect(resp.StatusCode).To(gomega.BeEquivalentTo(expectedStatusCode))
	defer resp.Body.Close()

	// read response body
	respBody := &bytes.Buffer{}
	_, err = respBody.ReadFrom(resp.Body)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return respBody.Bytes()
}

// CompareResponseAndGoldenOutput compares the json response and golden json output
func CompareResponseAndGoldenOutput(goldenFilePath string, responseBytes []byte) {
	var responseEngineOutput, fileDataEngineOutput policy.EngineOutput
	fileBytes, err := os.ReadFile(goldenFilePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	err = json.Unmarshal(responseBytes, &responseEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = json.Unmarshal(fileBytes, &fileDataEngineOutput)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	helper.CompareSummaryAndViolations(responseEngineOutput, fileDataEngineOutput)
}
