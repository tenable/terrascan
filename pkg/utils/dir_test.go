package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsDirExists(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "directory doesnot exists",
			args: args{dir: "test"},
			want: false,
		},
		{
			name: "directory  exists",
			args: args{dir: filepath.Join(os.Getenv("PWD"), "testdata")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDirExists(tt.args.dir); got != tt.want {
				t.Errorf("IsDirExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
