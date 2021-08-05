package kustomizev4

import (
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name      string
		filePath  string
		options   map[string]interface{}
		kustomize KustomizeV4
		typeOnly  bool
		want      output.AllResourceConfigs
		wantErr   error
	}{
		{
			name:      "load iac file is not supported for kustomize",
			filePath:  "/dummyfilepath.yaml",
			kustomize: KustomizeV4{},
			wantErr:   errLoadIacFileNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.kustomize.LoadIacFile(tt.filePath, tt.options)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			} else if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
