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

package commons

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
)

// test data
var (
	testLocalSourceAddr  = "./someModule"
	testRemoteSourceAddr = "terraform-aws-modules/eks/aws"
	testDirPath          = filepath.Join("root", "test")
	testFileNamePath     = filepath.Join(testDirPath, "main.tf")

	testModuleReqA = &hclConfigs.ModuleRequest{
		SourceAddr: testLocalSourceAddr,
		CallRange:  hcl.Range{Filename: testFileNamePath},
	}
)

func TestProcessLocalSource(t *testing.T) {

	type args struct {
		req *hclConfigs.ModuleRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no remote module",
			args: args{
				req: testModuleReqA,
			},
			want: filepath.Join(testDirPath, "someModule"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := NewTerraformDirectoryLoader("", false)
			if got := dl.processLocalSource(tt.args.req); got != tt.want {
				t.Errorf("processLocalSource() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestProcessTerraformRegistrySource(t *testing.T) {
	testTempDir := utils.GenerateTempDir()

	type args struct {
		req            *hclConfigs.ModuleRequest
		remoteModPaths map[string]string
		tempDir        string
		m              downloader.ModuleDownloader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid registry host",
			args: args{
				req: &hclConfigs.ModuleRequest{
					SourceAddr: "test.com/test/eks/aws",
				},
				remoteModPaths: make(map[string]string),
				tempDir:        utils.GenerateTempDir(),
				m:              downloader.NewRemoteDownloader(),
			},
			wantErr: true,
		},
		{
			name: "valid registry source",
			args: args{
				req: &hclConfigs.ModuleRequest{
					SourceAddr: testRemoteSourceAddr,
				},
				remoteModPaths: make(map[string]string),
				tempDir:        testTempDir,
				m:              downloader.NewRemoteDownloader(),
			},
			want: testTempDir,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.args.tempDir)
			dl := NewTerraformDirectoryLoader("", false)
			got, err := dl.processTerraformRegistrySource(tt.args.req, tt.args.tempDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTerraformRegistrySource() got error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("processTerraformRegistrySource() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
