/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testDataDir = "testdata"

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
			want:    GetHomeDir(),
			wantErr: nil,
		},
		{
			name:    "dir in HOME dir",
			path:    "~/somedir",
			want:    filepath.Join(GetHomeDir(), "somedir"),
			wantErr: nil,
		},
		{
			name:    "testdata dir",
			path:    testDataDir,
			want:    filepath.Join(os.Getenv("PWD"), "testdata"),
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

func TestFindAllDirectories(t *testing.T) {

	table := []struct {
		name     string
		basePath string
		want     []string
		wantErr  error
	}{
		{
			name:     "happy path",
			basePath: testDataDir,
			want:     []string{testDataDir, filepath.Join(testDataDir, "emptydir"), filepath.Join(testDataDir, "testdir1"), filepath.Join(testDataDir, "testdir2")},
			wantErr:  nil,
		},
		{
			name:     "empty dir",
			basePath: filepath.Join(testDataDir, "emptydir"),
			want:     []string{filepath.Join(testDataDir, "emptydir")},
			wantErr:  nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := FindAllDirectories(tt.basePath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("gotErr: '%+v', wantErr: '%+v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}

	t.Run("invalid dir", func(t *testing.T) {
		basePath := filepath.Join(testDataDir, "nothere")
		_, gotErr := FindAllDirectories(basePath)
		if gotErr == nil {
			t.Errorf("got no error; error expected")
		}
	})
}

func TestGetFileURI(t *testing.T) {

	table := []struct {
		name    string
		path    string
		want    string
		wantWin string
		wantErr error
	}{
		{
			name:    "test Linux path",
			path:    "dir1/dir2/file.txt",
			want:    "file://dir1/dir2/file.txt",
			wantWin: "file:///dir1/dir2/file.txt",
			wantErr: nil,
		},
		{
			name:    "test Windows path",
			path:    "c:\\dir1\\dir2\\file.txt",
			want:    "file://c:/dir1/dir2/file.txt",
			wantWin: "file:///c:/dir1/dir2/file.txt",
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileURI(tt.path)
			if err != tt.wantErr {
				t.Errorf("unexpected error; got: '%v', want: '%v'", err, tt.wantErr)
			}
			if IsWindowsPlatform() {
				if got != tt.wantWin {
					t.Errorf("got: '%v', want: '%v'", got, tt.wantWin)
				}
			} else {
				if got != tt.want {
					t.Errorf("got: '%v', want: '%v'", got, tt.want)
				}
			}
		})
	}
}
