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

package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// FindResourceByID Finds a given resource within the resource map and returns a reference to that resource
func FindResourceByID(resourceID string, normalizedResources *output.AllResourceConfigs) (*output.ResourceConfig, error) {
	resTypeName := strings.Split(resourceID, ".")
	if len(resTypeName) < 2 {
		return nil, fmt.Errorf("resource ID has an invalid format %s", resourceID)
	}

	resourceType := resTypeName[0]

	found := false
	var resource output.ResourceConfig
	resourceTypeList := (*normalizedResources)[resourceType]
	for i := range resourceTypeList {
		if resourceTypeList[i].ID == resourceID {
			resource = resourceTypeList[i]
			found = true
			break
		}
	}

	if !found {
		return nil, nil
	}

	return &resource, nil
}

// GetResourceCount gives out the total number of resources present in a output.ResourceConfig object.
// Since the ResourceConfig mapping stores resources in lists which can be located resourceMapping[Type],
// `len(resourceMapping)` does not give the count of the resources but only gives out the total number of
// the type of resources inside the object.
func GetResourceCount(resourceMapping map[string][]output.ResourceConfig) (count int) {
	count = 0
	for _, list := range resourceMapping {
		count = count + len(list)
	}
	return
}

// UpdateResourceConfigs adds a resource of given type if it is not present in allResources
func UpdateResourceConfigs(resourceType string, resources []output.ResourceConfig, allResources output.AllResourceConfigs) {
	for _, res := range resources {
		if !IsConfigPresent(allResources[resourceType], res) {
			allResources[resourceType] = append(allResources[resourceType], res)
		}
	}
}

// IsConfigPresent checks whether a resource is already present in the list of configs or not.
// The equality of a resource is based on name, source and config of the resource.
func IsConfigPresent(resources []output.ResourceConfig, resourceConfig output.ResourceConfig) bool {
	for _, resource := range resources {
		if resource.Name == resourceConfig.Name && resource.Source == resourceConfig.Source {
			if reflect.DeepEqual(resource.Config, resourceConfig.Config) {
				return true
			}
		}
	}
	return false
}
