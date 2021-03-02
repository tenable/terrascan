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
	"testing"

	"github.com/accurics/terrascan/pkg/downloader"
	hclConfigs "github.com/hashicorp/terraform/configs"
)

// test data
var (
	testLocalSourceAddr  = "./someModule"
	testRemoteSourceAddr = "terraform-aws-modules/eks/aws"

	testModuleReqA = &hclConfigs.ModuleRequest{
		SourceAddr: testLocalSourceAddr,
		Parent: &hclConfigs.Config{
			SourceAddr: "./eks/aws",
		},
	}

	testModuleReqB = &hclConfigs.ModuleRequest{
		SourceAddr: testLocalSourceAddr,
		Parent: &hclConfigs.Config{
			SourceAddr: testRemoteSourceAddr,
		},
	}
)

func TestProcessLocalSource(t *testing.T) {

	type args struct {
		req            *hclConfigs.ModuleRequest
		remoteModPaths map[string]string
		absRootDir     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no remote module",
			args: args{
				req:        testModuleReqA,
				absRootDir: "/home/somedir",
			},
			want: "/home/somedir/eks/aws/someModule",
		},
		{
			name: "with remote module",
			args: args{
				req: testModuleReqB,
				remoteModPaths: map[string]string{
					testRemoteSourceAddr: "/var/temp/testDir",
				},
			},
			want: "/var/temp/testDir/someModule",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processLocalSource(tt.args.req, tt.args.remoteModPaths, tt.args.absRootDir); got != tt.want {
				t.Errorf("processLocalSource() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestProcessTerraformRegistrySource(t *testing.T) {
	testTempDir := generateTempDir()

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
				tempDir:        generateTempDir(),
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
			got, err := processTerraformRegistrySource(tt.args.req, tt.args.remoteModPaths, tt.args.tempDir, tt.args.m)
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
