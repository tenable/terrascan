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

package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	httputils "github.com/tenable/terrascan/pkg/utils/http"
	"go.uber.org/zap"
)

var (
	errInitFailed = fmt.Errorf("failed to initialize webhook notifier")
	// ErrNilConfigData will be returned when config is nil
	ErrNilConfigData = fmt.Errorf("config data is nil")
)

// Init initializes the webhook notifier, reads config file and configures the
// necessary parameters for webhook notifications to work
func (w *Webhook) Init(config interface{}) error {
	// return error if config data is not present
	if config == nil {
		return ErrNilConfigData
	}

	// config to *toml.Tree
	webhookConfig, ok := config.(map[string]interface{})
	if !ok {
		zap.S().Errorf("error type casting webhook config data")
		return errInitFailed
	}

	// initialize Webhook struct with url and token

	jsonData, err := json.Marshal(webhookConfig)
	if err != nil {
		zap.S().Error("error while marshalling webhook config data", zap.Error(err))
		return errInitFailed
	}

	if err = json.Unmarshal(jsonData, w); err != nil {
		zap.S().Error("error while un-marshalling webhook config data", zap.Error(err))
		return errInitFailed
	}

	// succesful
	zap.S().Debug("initialized webhook notifier")
	return nil
}

// SendNotification sends webhook notification i.e sends a http POST request
// to the configured URL
func (w *Webhook) SendNotification(data interface{}) error {

	// convert data to json
	dataBytes, _ := json.Marshal(data)

	// make http POST request
	resp, err := httputils.SendPOSTRequest(w.URL, w.Token, dataBytes, http.Header{})
	if err != nil {
		zap.S().Errorf("failed to send webhook notification. error: '%v'", err)
		return err
	}

	// validate http response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		zap.S().Errorf("failed to send webhook notification, status code: '%v'", resp.StatusCode)
		return fmt.Errorf("webhook notification failed")
	}

	// successful
	zap.S().Debug("sent webhook notification")
	return nil
}
