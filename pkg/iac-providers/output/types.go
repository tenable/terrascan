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

package output

import (
	"fmt"
	"reflect"
	"strings"
)

// ResourceConfig describes a resource present in IaC
type ResourceConfig struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	ModuleName string      `json:"module_name,omitempty" yaml:"module_name,omitempty"`
	Source     string      `json:"source"`
	PlanRoot   string      `json:"plan_root,omitempty" yaml:"plan_root,omitempty" `
	Line       int         `json:"line"`
	Type       string      `json:"type"`
	Config     interface{} `json:"config"`
	// SkipRules will hold the rules to be skipped for the resource.
	// Each iac provider should append the rules to be skipped for a resource,
	// while extracting resource from the iac files
	SkipRules []SkipRule `json:"skip_rules" yaml:"skip_rules"`
}

// SkipRule struct will hold the skipped rule and any comment for the skipped rule
type SkipRule struct {
	Rule    string `json:"rule"`
	Comment string `json:"comment"`
}

// AllResourceConfigs is a list/slice of resource configs present in IaC
type AllResourceConfigs map[string][]ResourceConfig

// FindResourceByID Finds a given resource within the resource map and returns a reference to that resource
func (a AllResourceConfigs) FindResourceByID(resourceID string) (*ResourceConfig, error) {
	if len(a) == 0 {
		return nil, fmt.Errorf("AllResourceConfigs is nil or doesn't contain any resource type")
	}
	resTypeName := strings.Split(resourceID, ".")
	if len(resTypeName) < 2 {
		return nil, fmt.Errorf("resource ID has an invalid format %s", resourceID)
	}

	resourceType := resTypeName[0]

	found := false
	var resource ResourceConfig
	resourceTypeList := a[resourceType]
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

// FindAllResourcesByID Finds all resources within the resource map
func (a AllResourceConfigs) FindAllResourcesByID(resourceID string) ([]*ResourceConfig, error) {
	if len(a) == 0 {
		return nil, fmt.Errorf("AllResourceConfigs is nil or doesn't contain any resource type")
	}
	resTypeName := strings.Split(resourceID, ".")
	if len(resTypeName) < 2 {
		return nil, fmt.Errorf("resource ID has an invalid format %s", resourceID)
	}

	resourceType := resTypeName[0]

	resources := make([]*ResourceConfig, 0)
	resourceTypeList := a[resourceType]
	for i := range resourceTypeList {
		if resourceTypeList[i].ID == resourceID {
			resources = append(resources, &resourceTypeList[i])
		}
	}

	return resources, nil
}

// GetResourceCount gives out the total number of resources present in a output.ResourceConfig object.
// Since the ResourceConfig mapping stores resources in lists which can be located resourceMapping[Type],
// `len(resourceMapping)` does not give the count of the resources but only gives out the total number of
// the type of resources inside the object.
func (a AllResourceConfigs) GetResourceCount() (count int) {
	// handles nil map
	if len(a) == 0 {
		return 0
	}
	count = 0
	for _, list := range a {
		count = count + len(list)
	}
	return
}

// UpdateResourceConfigs adds a resource of given type if it is not present in allResources
func (a AllResourceConfigs) UpdateResourceConfigs(resourceType string, resources []ResourceConfig) {
	if _, ok := a[resourceType]; !ok {
		if len(a) == 0 {
			a = make(AllResourceConfigs)
		}
		a[resourceType] = resources
		return
	}
	for _, res := range resources {
		if !IsConfigPresent(a[resourceType], res) {
			a[resourceType] = append(a[resourceType], res)
		}
	}
}

// IsConfigPresent checks whether a resource is already present in the list of configs or not.
// The equality of a resource is based on name, source and config of the resource.
func IsConfigPresent(resources []ResourceConfig, resourceConfig ResourceConfig) bool {
	for _, resource := range resources {
		if resource.Name == resourceConfig.Name && resource.Source == resourceConfig.Source {
			if reflect.DeepEqual(resource.Config, resourceConfig.Config) {
				return true
			}
		}
	}
	return false
}
