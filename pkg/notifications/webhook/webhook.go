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

package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"

	httputils "github.com/accurics/terrascan/pkg/utils/http"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

var (
	errInitFailed = fmt.Errorf("failed to initialize webhook notifier")
)

// Init initalizes the webhook notifier, reads config file and configures the
// necessary parameters for webhook notifications to work
func (w *Webhook) Init(config interface{}) error {

	// config to *toml.Tree
	tomlConfig := config.(*toml.Tree)

	// initalize Webhook struct with url and token
	err := tomlConfig.Unmarshal(w)
	if err != nil {
		zap.S().Error(errInitFailed.Error())
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
	resp, err := httputils.SendPOSTRequest(w.URL, w.Token, dataBytes)
	if err != nil {
		zap.S().Errorf("failed to send webhook notification. error: '%v'", err)
		return err
	}

	// validate http response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		zap.S().Errorf("failed to webhook notification. Incorrect status code: '%v'", resp.StatusCode)
		return fmt.Errorf("webhook notification failed")
	}

	// successful
	zap.S().Debug("sent webhook notification")
	return nil
}
