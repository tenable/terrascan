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

package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	iacProvider "github.com/tenable/terrascan/pkg/iac-providers"
	"go.uber.org/zap"
)

// IacProvider contains response body for iac providers
type IacProvider struct {
	Type           string   `json:"type"`
	Versions       []string `json:"versions"`
	DefaultVersion string   `json:"defaultVersion"`
}

// iacProviders returns list of iac providers
func (g *APIHandler) iacProviders(w http.ResponseWriter, r *http.Request) {
	var providers = []IacProvider{}
	for _, provider := range iacProvider.SupportedIacProviders() {
		providers = append(providers, IacProvider{
			Type:           string(provider),
			Versions:       iacProvider.GetProviderIacVersions(provider),
			DefaultVersion: iacProvider.GetDefaultIacVersion(provider),
		})
	}

	response, err := json.MarshalIndent(providers, "", "  ")
	if err != nil {
		errMsg := fmt.Sprintf("failed to create JSON. error: '%v'", err)
		zap.S().Error(errMsg)
		apiErrorResponse(w, errMsg, http.StatusInternalServerError)
		return
	}

	apiResponse(w, string(response), http.StatusOK)
}
