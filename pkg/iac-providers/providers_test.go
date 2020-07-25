package iacprovider

import (
	"reflect"
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
