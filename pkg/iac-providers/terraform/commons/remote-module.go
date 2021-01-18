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

package commons

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var modulesPath = "v1/modules"
var terraformPublicRegistryURL = "registry.terraform.io"
var discoveryEndpoint = "/.well-known/terraform.json"

// RegistryResponse is the response received from terraform registry
type RegistryResponse struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Provider  string `json:"provider"`
	Source    string `json:"source"`
	Tag       string `json:"tag"`
}

// DiscoveryResponse is the response received for discovery process
type DiscoveryResponse struct {
	Modules string `json:"modules.v1"`
}

// TerraformRegistry represents a terraform registry
type TerraformRegistry struct {
	Host      string
	Namespace string
	Name      string
	Provider  string
}

// NewTerraformRegistry returns a TerraformRegistry
func NewTerraformRegistry(host, namespace, name, provider string) *TerraformRegistry {
	tfRegistry := new(TerraformRegistry)
	tfRegistry.Host = host
	tfRegistry.Namespace = namespace
	tfRegistry.Name = name
	tfRegistry.Provider = provider
	return tfRegistry
}

func (tr TerraformRegistry) getReqURL() string {
	return "https://" + tr.Host + "/" + modulesPath + "/" + tr.Namespace + "/" + tr.Name + "/" + tr.Provider
}

func (tr TerraformRegistry) getResourceURL() (string, error) {
	res, err := http.Get(tr.getReqURL())
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non 200 response for the request")
	}

	data, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	regRes := RegistryResponse{}
	err = json.Unmarshal(data, &regRes)
	if err != nil {
		return "", err
	}

	// form the source url
	return fmt.Sprintf("git::%s.git?ref=%s", regRes.Source, regRes.Tag), nil
}

func isRemoteModuleValid(source string) (*TerraformRegistry, bool) {
	URLParts := strings.Split(source, "/")
	partsLength := len(URLParts)

	// the module is of the form <HOST>/<NAMESPACE>/<NAME>/<PROVIDER>,
	// <HOST> valud is present when the registry is not public terraform registry
	if partsLength < 3 || partsLength > 4 {
		return nil, false
	}

	// if the length is 3, host is terraform public registry
	if partsLength == 3 {
		reqURL := "https://" + terraformPublicRegistryURL + "/" + modulesPath + "/" + strings.Join(URLParts, "/")

		res, err := http.Get(reqURL)
		if err != nil {
			return nil, false
		}

		if res.StatusCode != http.StatusOK {
			return nil, false
		}

		return NewTerraformRegistry(terraformPublicRegistryURL, URLParts[0], URLParts[1], URLParts[2]), true
	}

	// if the length is 4, the 1st element of the slice would be the hostname
	// hit the discover endpoint and check for 200 resopnse
	if partsLength == 4 {
		reqURL := "https://" + URLParts[0] + discoveryEndpoint
		res, err := http.Get(reqURL)
		if err != nil {
			return nil, false
		}

		if res.StatusCode != http.StatusOK {
			return nil, false
		}

		return NewTerraformRegistry(URLParts[0], URLParts[1], URLParts[2], URLParts[3]), true
	}

	return nil, false
}
