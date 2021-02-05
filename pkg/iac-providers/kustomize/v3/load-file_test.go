package kustomizev3

import (
	"reflect"
	"testing"

	iacloaderror "github.com/accurics/terrascan/pkg/iac-providers/iac-load-error"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacFile(t *testing.T) {

	table := []struct {
		name      string
		filePath  string
		kustomize KustomizeV3
		typeOnly  bool
		want      output.AllResourceConfigs
		wantErr   error
	}{
		{
			name:      "load iac file is not supported for kustomize",
			filePath:  "/dummyfilepath.yaml",
			kustomize: KustomizeV3{},
			wantErr:   errLoadIacFileNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.kustomize.LoadIacFile(tt.filePath)
			if e, ok := gotErr.(*iacloaderror.LoadError); !ok || e.Err != tt.wantErr {
				t.Errorf("TestLoadIacFile()= gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
			if tt.typeOnly && (reflect.TypeOf(gotErr)) != reflect.TypeOf(tt.wantErr) {
				t.Errorf("TestLoadIacFile()= gotErr: '%v', wantErr: '%v'", reflect.TypeOf(gotErr), reflect.TypeOf(tt.wantErr))
			}
		})
	}
}
