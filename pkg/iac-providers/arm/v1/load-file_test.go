/*
    Copyright (C) 2022 Tenable, Inc.

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

package armv1

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var fileTestDataDir = filepath.Join(testDataDir, "file-test-data")

func TestLoadIacFile(t *testing.T) {
	table := []struct {
		wantErr  error
		want     output.AllResourceConfigs
		options  map[string]interface{}
		armv1    ARMV1
		name     string
		filePath string
		typeOnly bool
	}{
		{
			name:     "empty config file",
			filePath: filepath.Join(fileTestDataDir, "empty-file.json"),
			armv1:    ARMV1{},
			wantErr:  fmt.Errorf("unable to parse file testdata/file-test-data/empty-file.json"),
		},
		{
			name:     "key-vault",
			filePath: filepath.Join(fileTestDataDir, "key-vault.json"),
			armv1:    ARMV1{},
			wantErr:  nil,
		},
		{
			name:     "nonexistent file",
			filePath: "nonexistent.json",
			armv1:    ARMV1{},
			wantErr:  fmt.Errorf("unable to read file nonexistent.json"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.armv1.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}

func TestLinkedTemplateDownload(t *testing.T) {
	table := []struct {
		wantErr  error
		want     output.AllResourceConfigs
		options  map[string]interface{}
		armv1    ARMV1
		name     string
		filePath string
		typeOnly bool
	}{
		{
			wantErr: nil,
			want: output.AllResourceConfigs{
				"azurerm_storage_account": []output.ResourceConfig{{
					ID: "azurerm_storage_account.GEN-UNIQUE",
				}},
			},
			armv1:    ARMV1{},
			name:     "linked-template-download",
			filePath: filepath.Join(fileTestDataDir, "azuredeploy.json"),
			typeOnly: false,
		},
		{
			wantErr:  nil,
			armv1:    ARMV1{},
			name:     "linked-template-wrong-uri",
			filePath: filepath.Join(fileTestDataDir, "azuredeploy-wrong-uri.json"),
			typeOnly: false,
			want:     nil,
		},
		{
			wantErr:  nil,
			armv1:    ARMV1{},
			name:     "linked-template-nested",
			filePath: filepath.Join(fileTestDataDir, "azuredeploy-nested.json"),
			typeOnly: false,
			want: output.AllResourceConfigs{
				"azurerm_storage_account": []output.ResourceConfig{{
					ID: "azurerm_storage_account.GEN-UNIQUE",
				}},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			aRC, gotErr := tt.armv1.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
			if tt.want != nil {
				// check if the resource from downloaded template is present
				for resType, resources := range tt.want {
					if _, present := aRC[resType]; !present {
						t.Errorf("resources for type %v not found for file %v", resType, tt.filePath)
					}
					for _, resource := range resources {
						if !isIDPresent(resource.ID, aRC[resType]) {
							t.Errorf("resource ID %v not found for file %v", resource.ID, tt.filePath)
						}
					}
				}
			}
		})
	}
}

func isIDPresent(id string, resourceConfigs []output.ResourceConfig) bool {
	for _, r := range resourceConfigs {
		if id == r.ID {
			return true
		}
	}
	return false
}
