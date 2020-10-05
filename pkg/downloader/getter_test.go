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
	"fmt"
	"reflect"
	"testing"
)

var (
	someDest   = "some-dest"
	someType   = "some-type"
	someSubdir = "some-subdir"
	someURL    = "some-url"
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
			dest:       someDest,
			wantURL:    "git::https://github.com/accurics/terrascan.git",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "github url with subdir",
			URL:        "github.com/accurics/terrascan//some-subdir",
			dest:       someDest,
			wantURL:    "git::https://github.com/accurics/terrascan.git",
			wantSubDir: someSubdir,
			wantErr:    nil,
		},
		{
			name:       "github ssh url",
			URL:        "git@github.com:accurics/terrascan.git//some-subdir",
			dest:       someDest,
			wantURL:    "git::ssh://git@github.com/accurics/terrascan.git",
			wantSubDir: someSubdir,
			wantErr:    nil,
		},
		{
			name:       "github basic auth",
			URL:        "git::ssh://username@example.com/storage.git//some-subdir",
			dest:       someDest,
			wantURL:    "git::ssh://username@example.com/storage.git",
			wantSubDir: someSubdir,
			wantErr:    nil,
		},
		{
			name:       "github ref version",
			URL:        "git::https://example.com/vpc.git?ref=v1.2.0",
			dest:       someDest,
			wantURL:    "git::https://example.com/vpc.git?ref=v1.2.0",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "http url",
			URL:        "https://example.com/vpc-module.zip",
			dest:       someDest,
			wantURL:    "https://example.com/vpc-module.zip",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "http url with basic auth",
			URL:        "https://Aladdin:OpenSesame@www.example.com/index.html",
			dest:       someDest,
			wantURL:    "https://Aladdin:OpenSesame@www.example.com/index.html",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "s3 url",
			URL:        "s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip",
			dest:       someDest,
			wantURL:    "s3::https://s3-eu-west-1.amazonaws.com/examplecorp-terraform-modules/vpc.zip",
			wantSubDir: "",
			wantErr:    nil,
		},
		{
			name:       "gcs url",
			URL:        "gcs::https://www.googleapis.com/storage/v1/modules",
			dest:       someDest,
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
			t.Errorf("subdir got: '%v', want: '%v'", gotSubDir, tt.wantSubDir)
		}
		if !reflect.DeepEqual(gotErr, tt.wantErr) {
			t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
		}
	}
}

func TestDownload(t *testing.T) {

	table := []struct {
		name     string
		URL      string
		dest     string
		wantDest string
		wantErr  error
	}{
		{
			name:     "empty URL",
			URL:      "",
			dest:     someDest,
			wantDest: "",
			wantErr:  ErrEmptyURLDest,
		},
		{
			name:     "empty destination",
			URL:      someURL,
			dest:     "",
			wantDest: "",
			wantErr:  ErrEmptyURLDest,
		},
		{
			name:     "invalid url",
			URL:      "github.com/some-repo",
			dest:     someDest,
			wantDest: "",
			wantErr:  fmt.Errorf("GitHub URLs should be github.com/username/repo"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoGetter()
			gotDest, gotErr := g.Download(tt.URL, tt.dest)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(gotDest, tt.wantDest) {
				t.Errorf("dest got: '%v', want: '%v'", gotDest, tt.wantDest)
			}
		})
	}
}

func TestDownloadWithType(t *testing.T) {

	table := []struct {
		name     string
		Type     string
		URL      string
		dest     string
		wantDest string
		wantErr  error
	}{
		{
			name:     "empty URL and Type",
			Type:     "",
			URL:      "",
			dest:     someDest,
			wantDest: "",
			wantErr:  ErrEmptyURLType,
		},
		{
			name:     "empty URL",
			Type:     someType,
			URL:      "",
			dest:     someDest,
			wantDest: "",
			wantErr:  ErrEmptyURLDest,
		},
		{
			name:     "empty Type",
			Type:     "",
			URL:      someURL,
			dest:     someDest,
			wantDest: "",
			wantErr:  ErrEmptyURLDest,
		},
		{
			name:     "empty dest",
			Type:     someType,
			URL:      someURL,
			dest:     "",
			wantDest: "",
			wantErr:  ErrEmptyURLDest,
		},
		{
			name:     "invalid url",
			Type:     someType,
			URL:      "github.com/some-url",
			dest:     someDest,
			wantDest: "",
			wantErr:  fmt.Errorf("download not supported for scheme 'some-type'"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGoGetter()
			gotDest, gotErr := g.DownloadWithType(tt.Type, tt.URL, tt.dest)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(gotDest, tt.wantDest) {
				t.Errorf("dest got: '%v', want: '%v'", gotDest, tt.wantDest)
			}
		})
	}
}
