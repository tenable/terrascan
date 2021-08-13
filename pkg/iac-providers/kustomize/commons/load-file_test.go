package commons

import (
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name     string
		filePath string
		options  map[string]interface{}
		typeOnly bool
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "load iac file is not supported for kustomize",
			filePath: "/dummyfilepath.yaml",
			wantErr:  errLoadIacFileNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := LoadIacFile()
			if gotErr != tt.wantErr {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
