package httputils

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

var (
	errNewRequest = fmt.Errorf("failed to create http request")
	errDoRequest  = fmt.Errorf("failed to make http request")
)

// default global http client
var client *http.Client

// init creates a http client which retries on errors like connection timeouts,
// server too slow respond etc.
func init() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	client = retryClient.StandardClient()
}

// SendRequest sends a http request on the given url
func SendRequest(method, url, token string, data []byte) (*http.Response, error) {

	var resp *http.Response

	// new http request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		zap.S().Errorf("failed to create http request; method: '%v', url: '%v'")
		return resp, errNewRequest
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer: '%s'", token))
	}

	// make request
	resp, err = client.Do(req)
	if err != nil {
		zap.S().Errorf("failed to make http request; method: '%v', url: '%v'")
		return resp, errDoRequest
	}

	return resp, err
}

// SendPOSTRequest sends a http POST request
func SendPOSTRequest(url, token string, data []byte) (*http.Response, error) {
	return SendRequest("POST", url, token, data)
}
