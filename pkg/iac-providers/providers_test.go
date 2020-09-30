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

package iacprovider

import (
	"reflect"
	"sort"
	"testing"

	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
)

func TestNewIacProvider(t *testing.T) {

	table := []struct {
		name       string
		iacType    supportedIacType
		iacVersion supportedIacVersion
		want       IacProvider
		wantErr    error
	}{
		{
			name:       "terraform v12",
			iacType:    terraform,
			iacVersion: terraformV12,
			want:       &tfv12.TfV12{},
			wantErr:    nil,
		},
		{
			name:       "not supported iac type",
			iacType:    "not-supported",
			iacVersion: terraformV12,
			want:       nil,
			wantErr:    errIacNotSupported,
		},
		{
			name:       "not supported iac version",
			iacType:    terraform,
			iacVersion: "not-supported",
			want:       nil,
			wantErr:    errIacNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewIacProvider(string(tt.iacType), string(tt.iacVersion))
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

func TestIsIacSupported(t *testing.T) {

	table := []struct {
		name       string
		iacType    supportedIacType
		iacVersion supportedIacVersion
		want       bool
	}{
		{
			name:       "terraform v12",
			iacType:    terraform,
			iacVersion: terraformV12,
			want:       true,
		},
		{
			name:       "not supported iac type",
			iacType:    "not-supported",
			iacVersion: terraformV12,
			want:       false,
		},
		{
			name:       "not supported iac version",
			iacType:    terraform,
			iacVersion: "not-supported",
			want:       false,
		},
		{
			name:       "not supported iac type and version",
			iacType:    "not-supported-type",
			iacVersion: "not-supported-version",
			want:       false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIacSupported(string(tt.iacType), string(tt.iacVersion))
			if got != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

func TestSupportedIacProviders(t *testing.T) {
	t.Run("supported iac providers", func(t *testing.T) {
		var want []string
		for k := range supportedIacProviders {
			want = append(want, string(k))
		}
		sort.Strings(want)
		got := SupportedIacProviders()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}
