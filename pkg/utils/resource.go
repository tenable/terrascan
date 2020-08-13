package utils

import (
	"fmt"
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
