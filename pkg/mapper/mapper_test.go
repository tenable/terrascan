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

package mapper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestARMMapper(t *testing.T) {
	tests := []struct {
		name          string
		template      string
		parameters    string
		expectedError bool
	}{
		{
			name:          "test-for-valid-json",
			template:      "key-vault/azuredeploy.json",
			parameters:    "key-vault/azuredeploy.parameters.json",
			expectedError: false,
		},
	}

	m := NewMapper("arm")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, err := iacDocumentFromFile(test.template)
			if err != nil {
				t.Error(err)
			}

			params, err := parametersFromFile(test.parameters)
			if err != nil {
				t.Error(err)
			}
			_, err = m.Map(doc, params)
			if err != nil {
				assert.True(t, test.expectedError)
				return
			}
			assert.False(t, test.expectedError)
		})
	}
}

func iacDocumentFromFile(name string) (*utils.IacDocument, error) {
	data, err := readFile(name)
	if err != nil {
		return nil, err
	}

	return &utils.IacDocument{
		Type:      utils.JSONDoc,
		StartLine: 0,
		EndLine:   183,
		FilePath:  filepath.Join("test_data", name),
		Data:      data,
	}, nil
}

func readFile(name string) ([]byte, error) {
	const testData = "test_data"
	f, err := os.Open(filepath.Join(testData, name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func parametersFromFile(name string) (map[string]interface{}, error) {
	data, err := readFile(name)
	if err != nil {
		return nil, err
	}

	var params map[string]interface{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, err
	}

	res, err := extractParameterValues(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func extractParameterValues(params map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(params["parameters"])
	if err != nil {
		return nil, err
	}
	var npm map[string]struct {
		Value interface{} `json:"value"`
	}
	err = json.Unmarshal(data, &npm)
	if err != nil {
		return nil, err
	}

	finalParams := map[string]interface{}{}
	for key, value := range npm {
		finalParams[key] = value.Value
	}
	return finalParams, nil
}
