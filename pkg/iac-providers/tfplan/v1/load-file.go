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

package tfplan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

const (
	jqQuery = `[.planned_values.root_module | .. | select(.type? != null and .address? != null and .mode? == "managed") | {id: .address?, type: .type?, name: .name?, config: .values?, source: ""}]`
)

// LoadIacFile parses the given tfplan file from the given file path
func (t *TFPlan) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {

	zap.S().Debug("processing tfplan file")

	// validate if provide file is a valid tfplan file

	// read tfplan json file
	tfjson, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		errMsg := fmt.Sprintf("failed to read tfplan JSON file. error: '%v'", err)
		zap.S().Debug(errMsg)
		return allResourcesConfig, fmt.Errorf(errMsg)
	}

	// run jq query on tfplan json
	processed, err := utils.JQFilterWithQuery(jqQuery, tfjson)
	if err != nil {
		errMsg := fmt.Sprintf("failed to process tfplan JSON. error: '%v'", err)
		zap.S().Debug(errMsg)
		return allResourcesConfig, fmt.Errorf(errMsg)
	}

	// decode processed out into output.ResourceConfig
	var resourceConfigs []output.ResourceConfig
	if err := json.Unmarshal(processed, &resourceConfigs); err != nil {
		errMsg := fmt.Sprintf("failed to decode proceesed jq output. error: '%v'", err)
		zap.S().Debug(errMsg)
		return allResourcesConfig, fmt.Errorf(errMsg)
	}

	// create AllResourceConfigs from resourceConfigs
	allResourcesConfig = make(map[string][]output.ResourceConfig)
	for _, r := range resourceConfigs {
		r.ID = getTFID(r.ID)
		if _, present := allResourcesConfig[r.Type]; !present {
			allResourcesConfig[r.Type] = []output.ResourceConfig{r}
		} else {
			allResourcesConfig[r.Type] = append(allResourcesConfig[r.Type], r)
		}
	}

	// return output
	return allResourcesConfig, nil
}

// getTFID returns a valid resource ID for terraform
func getTFID(id string) string {
	split := strings.Split(id, ".")
	if len(split) <= 2 {
		return strings.Join(split, ".")
	}
	return strings.Join(split[len(split)-2:], ".")
}
