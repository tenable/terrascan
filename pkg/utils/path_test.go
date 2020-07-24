package utils

import (
	"os"
	"testing"
)

func TestGetAbsPath(t *testing.T) {

	table := []struct {
		name    string
		path    string
		want    string
		wantErr error
	}{
		{
			name:    "test PWD",
			path:    ".",
			want:    os.Getenv("PWD"),
			wantErr: nil,
		},
		{
			name:    "user HOME dir",
			path:    "~",
			want:    os.Getenv("HOME"),
			wantErr: nil,
		},
		{
			name:    "testdata dir",
			path:    "./testdata",
			want:    os.Getenv("PWD") + "/testdata",
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAbsPath(tt.path)
			if err != tt.wantErr {
				t.Errorf("unexpected error; got: '%v', want: '%v'", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}
