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

var terraformPublicRegistryURL = "registry.terraform.io"

// the service discovery endpoint
// source: https://www.terraform.io/docs/internals/remote-service-discovery.html#discovery-process
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

// returns the request url to fetch the source url of the remote module
func (tr TerraformRegistry) getReqURL(modulesPath, version string) string {
	// form https://{host}/{modulesPath}/{namespace}/{name}/{provider}/
	reqURL := "https://" + tr.Host + modulesPath + tr.Namespace + "/" + tr.Name + "/" + tr.Provider
	if version != "" {
		reqURL = reqURL + "/" + version
	}

	return reqURL
}

func (tr TerraformRegistry) getModulesPath() (string, error) {
	// get the modules path by hitting the discovery endpoint
	response, err := http.Get("https://" + tr.Host + discoveryEndpoint)
	if err != nil {
		return "", err
	}

	// if the response received is not 200, then return
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non 200 response for the request")
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	dr := DiscoveryResponse{}
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return "", err
	}
	return dr.Modules, nil
}

func (tr TerraformRegistry) getResourceURL(version string) (string, error) {
	// get the modules path by hitting the discovery endpoint
	modulesPath, err := tr.getModulesPath()
	if err != nil {
		return "", err
	}

	// check if the specified registry url exists
	response, err := http.Get(tr.getReqURL(modulesPath, version))
	if err != nil {
		return "", err
	}

	// if the response received is not 200, then return
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non 200 response for the request")
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	regRes := RegistryResponse{}
	err = json.Unmarshal(data, &regRes)
	if err != nil {
		return "", err
	}

	// form the source url
	return fmt.Sprintf("git::%s.git?ref=%s", regRes.Source, regRes.Tag), nil
}

// this function checks if the module source string is of format <HOST>/<NAMESPACE>/<NAME>/<PROVIDER>
// and returns a TerraforRegistry struct if the module string format is valid
func isRemoteModuleValid(source string) (*TerraformRegistry, bool) {
	URLParts := strings.Split(source, "/")
	partsLength := len(URLParts)

	// the module is of the form <HOST>/<NAMESPACE>/<NAME>/<PROVIDER>,
	// <HOST> value is present when the registry is not public terraform registry
	if partsLength < 3 || partsLength > 4 {
		return nil, false
	}

	// if the length is 3, host is terraform public registry
	if partsLength == 3 {
		// since the parts length is 3, the module has to belong to terraform public registry
		return NewTerraformRegistry(terraformPublicRegistryURL, URLParts[0], URLParts[1], URLParts[2]), true
	}

	// if the length is 4, the 1st element of the slice would be the hostname of the terraform registry
	// hit the discovery endpoint and check for 200 resopnse
	if partsLength == 4 {

		// every terraform registry has to provide a valid response for reqURL formed below,
		// if not, it is not a valid terraform registry.
		// source: https://www.terraform.io/docs/internals/remote-service-discovery.html#discovery-process

		reqURL := "https://" + URLParts[0] + discoveryEndpoint
		response, err := http.Get(reqURL)
		if err != nil {
			return nil, false
		}

		if response.StatusCode != http.StatusOK {
			return nil, false
		}

		return NewTerraformRegistry(URLParts[0], URLParts[1], URLParts[2], URLParts[3]), true
	}

	return nil, false
}
