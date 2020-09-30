/*
    Copyright (C) 2020 Accurics, Inc.

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

package downloader

import (
	"reflect"
	"testing"
)

func TestNewGoGetter(t *testing.T) {
	t.Run("new GoGetter", func(t *testing.T) {
		var (
			want = &GoGetter{}
			got  = NewGoGetter()
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: '%v', want: '%v'", got, want)
		}
	})
}

func TestGetURLSubDir(t *testing.T) {

	table := []struct {
		name       string
		URL        string
		dest       string
		wantURL    string
		wantSubDir string
		wantErr    error
	}{
		{
			name:       "github url no subdir",
			URL:        "github.com/accurics/terrascan",
			dest:       "some-dest",
			wantURL:    "git::https://github.com/accurics/terrascan.git",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "github url with subdir",
			URL:        "github.com/accurics/terrascan//some-subdir",
			dest:       "some-dest",
			wantURL:    "git::https://github.com/accurics/terrascan.git",
			wantSubDir: "some-subdir",
			wantErr:    nil,
		},
		{
			name:       "github ssh url",
			URL:        "git@github.com:accurics/terrascan.git//some-subdir",
			dest:       "some-dest",
			wantURL:    "git::ssh://git@github.com/accurics/terrascan.git",
			wantSubDir: "some-subdir",
			wantErr:    nil,
		},
		{
			name:       "github basic auth",
			URL:        "git::ssh://username@example.com/storage.git//some-subdir",
			dest:       "some-dest",
			wantURL:    "git::ssh://username@example.com/storage.git",
			wantSubDir: "some-subdir",
			wantErr:    nil,
		},
		{
			name:       "github ref version",
			URL:        "git::https://example.com/vpc.git?ref=v1.2.0",
			dest:       "some-dest",
			wantURL:    "git::https://example.com/vpc.git?ref=v1.2.0",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "http url",
			URL:        "https://example.com/vpc-module.zip",
			dest:       "some-dest",
			wantURL:    "https://example.com/vpc-module.zip",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "http url with basic auth",
			URL:        "https://Aladdin:OpenSesame@www.example.com/index.html",
			dest:       "some-dest",
			wantURL:    "https://Aladdin:OpenSesame@www.example.com/index.html",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "s3 url",
			URL:        "s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip",
			dest:       "some-dest",
			wantURL:    "s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "gcs url",
			URL:        "gcs::https://www.googleapis.com/storage/v1/modules",
			dest:       "some-dest",
			wantURL:    "gcs::https://www.googleapis.com/storage/v1/modules",
			wantSubDir: "",
			wantErr:    nil,
		},
	}

	for _, tt := range table {
		g := NewGoGetter()
		gotURL, gotSubDir, gotErr := g.GetURLSubDir(tt.URL, tt.dest)
		if !reflect.DeepEqual(gotURL, tt.wantURL) {
			t.Errorf("url got: '%v', want: '%v'", gotURL, tt.wantURL)
		}
		if !reflect.DeepEqual(gotSubDir, tt.wantSubDir) {
			t.Errorf("url got: '%v', want: '%v'", gotSubDir, tt.wantSubDir)
		}
		if !reflect.DeepEqual(gotErr, tt.wantErr) {
			t.Errorf("url got: '%v', want: '%v'", gotErr, tt.wantErr)
		}
	}
}
