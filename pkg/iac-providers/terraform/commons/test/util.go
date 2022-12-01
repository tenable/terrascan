package test

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

// prepareAllResourceConfigs prepares a
// map[string]map[string]output.ResourceConfig
// from the output.AllResourceConfigs, which is a
// map[string][]output.ResourceConfig
//
// The goal is to put the [] into a map[string] so that we don't rely on the
// implicit order of the [], but can use the keys for ordering.
// The key is computed from the source and id, which should be globally unique.
func prepareAllResourceConfigs(v output.AllResourceConfigs, modifySourceVal bool) ([]byte, error) {

	newval := make(map[string]map[string]output.ResourceConfig, len(v))
	for key, val := range v {
		newval[key] = make(map[string]output.ResourceConfig, len(val))
		for _, item := range val {
			if modifySourceVal {
				item.Source = filepath.Join(strings.Split(item.Source, "/")...)
			}
			// we pull latest version available for provider version hence is subject to change
			item.ProviderVersion = ""

			newkey := item.Source + "##" + item.ID
			newval[key][newkey] = item
		}
	}

	contents, err := json.Marshal(newval)
	if !modifySourceVal {
		contents = utils.ReplaceCarriageReturnBytes(contents)
	}

	if err != nil {
		return []byte{}, err
	}

	return contents, nil
}

// IdenticalAllResourceConfigs determines if a and b have identical contents
func IdenticalAllResourceConfigs(a, b output.AllResourceConfigs) (bool, error) {
	value1, err := prepareAllResourceConfigs(a, false)
	if err != nil {
		return false, err
	}
	value2, err := prepareAllResourceConfigs(b, true)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(value1, value2), nil
}
