package runtime

import (
	"reflect"
	"testing"
)

func TestValidateInputs(t *testing.T) {

	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "valid filePath",
			executor: Executor{
				filePath:   "./testdata/testfile",
				dirPath:    "",
				cloudType:  "aws",
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr: nil,
		},
		{
			name: "valid dirPath",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  "aws",
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr: nil,
		},
		{
			name: "empty iac path",
			executor: Executor{
				filePath: "",
				dirPath:  "",
			},
			wantErr: errEmptyIacPath,
		},
		{
			name: "incorrect iac path",
			executor: Executor{
				filePath: "./testdata/testfile",
				dirPath:  "./testdata/testdir",
			},
			wantErr: errIncorrectIacPath,
		},
		{
			name: "filepath does not exist",
			executor: Executor{
				filePath: "./testdata/notthere",
			},
			wantErr: errFileNotExists,
		},
		{
			name: "directory does not exist",
			executor: Executor{
				dirPath: "./testdata/notthere",
			},
			wantErr: errDirNotExists,
		},
		{
			name: "invalid cloud",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  "nothere",
				iacType:    "terraform",
				iacVersion: "v12",
			},
			wantErr: errCloudNotSupported,
		},
		{
			name: "invalid iac type",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  "aws",
				iacType:    "notthere",
				iacVersion: "v12",
			},
			wantErr: errIacNotSupported,
		},
		{
			name: "invalid iac version",
			executor: Executor{
				filePath:   "",
				dirPath:    "./testdata/testdir",
				cloudType:  "aws",
				iacType:    "terraform",
				iacVersion: "notthere",
			},
			wantErr: errIacNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.ValidateInputs()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error, gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
