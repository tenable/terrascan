package utils

import (
	"reflect"
	"testing"
)

func TestReadYamlFile(t *testing.T) {

	table := []struct {
		path  string
		empty bool
	}{
		{path: "./testdata/empty.yaml", empty: true},
		{path: "./testdata/pod.yaml"},
	}

	for _, tt := range table {
		t.Run(tt.path, func(t *testing.T) {
			_, gotErr := ReadYamlFile(tt.path)

			if !reflect.DeepEqual(gotErr, ErrYamlFileEmpty) && tt.empty {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, ErrYamlFileEmpty)
			}

		})
	}
}
