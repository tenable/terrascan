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

package tfv12

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/downloader"
)

var (
	someDest = "some-dest"
	someAddr = "some-address"
)

// MockDownloader mocks the downloader.Downloader interface
type MockDownloader struct {
	output          string
	subDir          string
	errDownload     error
	errGetURLSubDir error
	errSubDirGlob   error
}

// Download mock method
func (m *MockDownloader) Download(url, destDir string) (string, error) {
	return m.output, m.errDownload
}

// DownloadWithType mock method
func (m *MockDownloader) DownloadWithType(remoteType, url, destDir string) (string, error) {
	return m.output, m.errDownload
}

// GetURLSubDir mock method
func (m *MockDownloader) GetURLSubDir(url, destDir string) (string, string, error) {
	return m.output, m.subDir, m.errGetURLSubDir
}

// SubDirGlob mock method
func (m *MockDownloader) SubDirGlob(destDir, subDir string) (string, error) {
	return m.output, m.errSubDirGlob
}

// ---------------------- unit tests -------------------------------- //
func TestIsLocalSourceAddr(t *testing.T) {

	table := []struct {
		name string
		addr string
		want bool
	}{
		{
			name: "local dir with ./",
			addr: "./somedir",
			want: true,
		},
		{
			name: "local dir with ../",
			addr: "../somedir",
			want: true,
		},
		{
			name: "local dir with .\\",
			addr: ".\\somedir",
			want: true,
		},
		{
			name: "local dir with ..\\",
			addr: "..\\somedir",
			want: true,
		},
		{
			name: "git repo",
			addr: "git@github.com:accurics/terrascan.git",
			want: false,
		},
		{
			name: "http url",
			addr: "http://i-am-not-there.com",
			want: false,
		},
		{
			name: "https url",
			addr: "https://i-am-not-there.com",
			want: false,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := isLocalSourceAddr(tt.addr)
			if got != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

func TestNewRemoteModuleInstaller(t *testing.T) {
	var (
		got            = NewRemoteModuleInstaller()
		wantDownloader = downloader.NewDownloader()
	)
	if !reflect.DeepEqual(got.downloader, wantDownloader) {
		t.Errorf("downloader got: '%v', want: '%v'", got.downloader, wantDownloader)
	}
}

func TestDownloadModule(t *testing.T) {

	var (
		wantOutput string = "want-output"
		wantSubDir string = "want-subdir"
		wantErr    error  = fmt.Errorf("want-error")
		wantNoErr  error  = nil
	)

	table := []struct {
		name       string
		addr       string
		dest       string
		wantOutput string
		wantSubDir string
		wantErr    error
		r          *RemoteModuleInstaller
	}{
		{
			name:       "GetURLSubDir error",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: "",
			wantErr:    wantErr,
			r: &RemoteModuleInstaller{
				cache: make(map[string]string),
				downloader: &MockDownloader{
					output:          wantOutput,
					subDir:          wantSubDir,
					errGetURLSubDir: wantErr,
				},
			},
		},
		{
			name:       "Download error",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: "",
			wantErr:    wantErr,
			r: &RemoteModuleInstaller{
				cache: make(map[string]string),
				downloader: &MockDownloader{
					output:      wantOutput,
					subDir:      wantSubDir,
					errDownload: wantErr,
				},
			},
		},
		{
			name:       "SubDirGlob error",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: "",
			wantErr:    wantErr,
			r: &RemoteModuleInstaller{
				cache: make(map[string]string),
				downloader: &MockDownloader{
					output:        wantOutput,
					subDir:        wantSubDir,
					errSubDirGlob: wantErr,
				},
			},
		},
		{
			name:       "no error",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: wantOutput,
			wantErr:    wantNoErr,
			r: &RemoteModuleInstaller{
				cache: make(map[string]string),
				downloader: &MockDownloader{
					output: wantOutput,
					subDir: wantSubDir,
				},
			},
		},
		{
			name:       "cache test",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: wantOutput,
			wantErr:    wantNoErr,
			r: &RemoteModuleInstaller{
				cache: map[string]string{wantOutput: "some-cache"},
				downloader: &MockDownloader{
					output: wantOutput,
					subDir: wantSubDir,
				},
			},
		},
		{
			name:       "no subdir",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: someDest,
			wantErr:    wantNoErr,
			r: &RemoteModuleInstaller{
				cache: make(map[string]string),
				downloader: &MockDownloader{
					output: wantOutput,
				},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, gotErr := tt.r.DownloadModule(tt.addr, tt.dest)
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("output got: '%v', want: '%v'", gotOutput, tt.wantOutput)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("error got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func TestCleanUp(t *testing.T) {

	t.Run("test cache clean up", func(t *testing.T) {

		// create temp dir
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Errorf("failed to create temp dir. error: '%v'", err)
		}

		// store temp dir into the installer cache
		r := &RemoteModuleInstaller{
			cache: map[string]string{"some-module": tempDir},
		}

		// run clean up, expect clean up to delete the temp dir
		r.CleanUp()

		// check if temp dir is deleted
		_, err = os.Stat(tempDir)
		if err == nil {
			t.Errorf("clean up failed")
		}
	})
}
