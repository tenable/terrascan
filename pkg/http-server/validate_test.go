package httpserver

import (
	"fmt"
	"reflect"
	"testing"
)

func TestValidateFiles(t *testing.T) {
	server := APIServer{}
	table := []struct {
		name           string
		privateKeyFile string
		certFile       string
		wantOutput     interface{}
		wantErr        error
	}{
		{
			name:           "normal file names",
			privateKeyFile: "key",
			certFile:       "cert",
			wantErr:        nil,
		},
		{
			name:           "error in both privatekey and certfile filenames",
			privateKeyFile: "",
			certFile:       "server.crt",
			wantErr:        fmt.Errorf("certificate file provided but private key file missing"),
		},
		{
			name:           "error in privatekey filename",
			privateKeyFile: "",
			certFile:       "",
			wantErr:        nil,
		},
		{
			name:           "error in certfile filename",
			privateKeyFile: "keyfile",
			certFile:       "",
			wantErr:        fmt.Errorf("private key file provided but certficate file missing"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := server.validateFiles(tt.privateKeyFile, tt.certFile)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				if tt.wantErr != nil && gotErr != nil && tt.wantErr.Error() != gotErr.Error() {
					t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
				}
			}
		})
	}
}
