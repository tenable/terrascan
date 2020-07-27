package cloudprovider

import (
	"reflect"
	"testing"

	awsProvider "github.com/accurics/terrascan/pkg/cloud-providers/aws"
)

func TestNewCloudProvider(t *testing.T) {

	table := []struct {
		name      string
		cloudType supportedCloudType
		want      CloudProvider
		wantErr   error
	}{
		{
			name:      "aws provider",
			cloudType: aws,
			want:      &awsProvider.AWSProvider{},
			wantErr:   nil,
		},
		{
			name:      "not supported cloud type",
			cloudType: "not-supported",
			want:      nil,
			wantErr:   errCloudNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := NewCloudProvider(string(tt.cloudType))
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

func TestIsCloudSupported(t *testing.T) {

	table := []struct {
		name      string
		cloudType supportedCloudType
		want      bool
	}{
		{
			name:      "aws provider",
			cloudType: aws,
			want:      true,
		},
		{
			name:      "not supported cloud type",
			cloudType: "not-supported",
			want:      false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCloudSupported(string(tt.cloudType))
			if got != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}
