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

package tfv12

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

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

// identicalAllResourceConfigs determines if a and b have identical contents
func identicalAllResourceConfigs(a, b output.AllResourceConfigs) (bool, error) {
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

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name    string
		dirPath string
		tfv12   TfV12
		want    output.AllResourceConfigs
		wantErr error
	}{
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			tfv12:   TfV12{},
			wantErr: errEmptyTFConfigDir,
		},
		{
			name:    "empty config",
			dirPath: "./testdata/testfile",
			tfv12:   TfV12{},
			wantErr: errEmptyTFConfigDir,
		},
		{
			name:    "incorrect module structure",
			dirPath: "./testdata/invalid-moduleconfigs",
			tfv12:   TfV12{},
			wantErr: errBuildTFConfigDir,
		},
		{
			name:    "load invalid config dir",
			dirPath: "./testdata",
			tfv12:   TfV12{},
			wantErr: errLoadConfigDir,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.tfv12.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

	table2 := []struct {
		name        string
		tfConfigDir string
		tfJSONFile  string
		tfv12       TfV12
		wantErr     error
	}{
		{
			name:        "config1",
			tfConfigDir: "./testdata/tfconfigs",
			tfJSONFile:  "testdata/tfjson/fullconfig.json",
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "module directory",
			tfConfigDir: "./testdata/moduleconfigs",
			tfJSONFile:  "./testdata/tfjson/moduleconfigs.json",
			tfv12:       TfV12{},
			wantErr:     nil,
		},
		{
			name:        "nested module directory",
			tfConfigDir: "./testdata/deep-modules",
			tfJSONFile:  "./testdata/tfjson/deep-modules.json",
			tfv12:       TfV12{},
			wantErr:     nil,
		},
	}

	for _, tt := range table2 {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.tfv12.LoadIacDir(tt.tfConfigDir)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			var want output.AllResourceConfigs

			// Read the expected value and unmarshal into want
			contents, _ := ioutil.ReadFile(tt.tfJSONFile)
			err := json.Unmarshal(contents, &want)
			if err != nil {
				t.Errorf("unexpected error unmarshalling want: %v", err)
			}

			match, err := identicalAllResourceConfigs(got, want)
			if err != nil {
				t.Errorf("unexpected error checking result: %v", err)
			}
			if !match {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(want, "", "  ")
				t.Errorf("got '%v', want: '%v'", string(g), string(w))
			}
		})
	}
}
