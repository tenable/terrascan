package test

import (
	"encoding/json"
	"reflect"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// prepareAllResourceConfigs prepares a
// map[string]map[string]output.ResourceConfig
// from the output.AllResourceConfigs, which is a
// map[string][]output.ResourceConfig
//
// The goal is to put the [] into a map[string] so that we don't rely on the
// implicit order of the [], but can use the keys for ordering.
// The key is computed from the source and id, which should be globally unique.
func prepareAllResourceConfigs(v output.AllResourceConfigs) ([]byte, error) {

	newval := make(map[string]map[string]output.ResourceConfig, len(v))
	for key, val := range v {
		newval[key] = make(map[string]output.ResourceConfig, len(val))
		for _, item := range val {
			newkey := item.Source + "##" + item.ID
			newval[key][newkey] = item
		}
	}

	contents, err := json.Marshal(newval)
	if err != nil {
		return []byte{}, err
	}

	return contents, nil
}

// IdenticalAllResourceConfigs determines if a and b have identical contents
func IdenticalAllResourceConfigs(a, b output.AllResourceConfigs) (bool, error) {
	value1, err := prepareAllResourceConfigs(a)
	if err != nil {
		return false, err
	}
	value2, err := prepareAllResourceConfigs(b)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(value1, value2), nil
}
