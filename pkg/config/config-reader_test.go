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

package config

import (
	"reflect"
	"testing"
)

func TestNewTerrascanConfigReader(t *testing.T) {
	testNotifier := Notifier{
		NotifierType: "webhook",
		NotifierConfig: map[string]interface{}{
			"url": "testurl1",
		},
	}
	testPolicy := Policy{
		BasePath: "custom-path",
		RepoPath: "rego-subdir",
		RepoURL:  "https://repository/url",
		Branch:   "branch-name",
	}
	testRules := Rules{
		ScanRules: []string{"rule.1", "rule.2", "rule.3", "rule.4", "rule.5"},
		SkipRules: []string{"rule.1"},
	}

	type args struct {
		fileName string
	}
	tests := []struct {
		name          string
		args          args
		want          *TerrascanConfigReader
		wantErr       bool
		assertGetters bool
		Policy
		notifications map[string]Notifier
		Rules
	}{
		{
			name: "empty config file",
			args: args{
				fileName: "",
			},
			want: &TerrascanConfigReader{},
		},
		{
			name: "non existent config file",
			args: args{
				fileName: "test",
			},
			wantErr: true,
			want:    &TerrascanConfigReader{},
		},
		{
			name: "invalid toml config file",
			args: args{
				fileName: "testdata/invalid.toml",
			},
			wantErr: true,
			want:    &TerrascanConfigReader{},
		},
		{
			name: "valid toml config file with partial fields",
			args: args{
				fileName: "testdata/terrascan-config.toml",
			},
			want: &TerrascanConfigReader{
				config: TerrascanConfig{
					Policy: testPolicy,
				},
			},
		},
		{
			name: "valid toml config file with all fields",
			args: args{
				fileName: "testdata/terrascan-config-all-fields.toml",
			},
			want: &TerrascanConfigReader{
				config: TerrascanConfig{
					Policy: testPolicy,
					Notifications: map[string]Notifier{
						"webhook1": testNotifier,
					},
					Rules: testRules,
				},
			},
			assertGetters: true,
			notifications: map[string]Notifier{
				"webhook1": testNotifier,
			},
			Policy: testPolicy,
			Rules:  testRules,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTerrascanConfigReader(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTerrascanConfigReader() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTerrascanConfigReader() = got %v, want %v", got, tt.want)
			}
			if tt.assertGetters {
				if !reflect.DeepEqual(got.GetPolicyConfig(), tt.Policy) {
					t.Errorf("NewTerrascanConfigReader() = got config: %v, want config: %v", got.GetPolicyConfig(), tt.Policy)
				}

				if !reflect.DeepEqual(got.GetNotifications(), tt.notifications) {
					t.Errorf("NewTerrascanConfigReader() = got notifications: %v, want notifications: %v", got.GetNotifications(), tt.notifications)
				}

				if !reflect.DeepEqual(got.GetRules(), tt.Rules) {
					t.Errorf("NewTerrascanConfigReader() = got rules: %v, want rules: %v", got.GetRules(), tt.Rules)
				}
			}
		})
	}
}
