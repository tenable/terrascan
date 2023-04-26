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

package commons

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
	"github.com/tenable/terrascan/pkg/downloader"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/iac-providers/terraform/commons/test"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
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

	invalidDirErrStringTemplate = "directory '%s' has no terraform config files"
	testDataDir                 = "testdata"
	tfJSONDir                   = filepath.Join(testDataDir, "tfjson")
)

func TestProcessLocalSource(t *testing.T) {

	type args struct {
		req *hclConfigs.ModuleRequest
	}
	tests := []struct {
		name             string
		args             args
		want             string
		options          map[string]interface{}
		terraformVersion string
	}{
		{
			name: "no remote module",
			args: args{
				req: testModuleReqA,
			},
			terraformVersion: "0.15.0",
			want:             filepath.Join(testDirPath, "someModule"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := NewTerraformDirectoryLoader("", tt.terraformVersion, tt.options)
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
		name             string
		args             args
		want             string
		wantErr          bool
		options          map[string]interface{}
		terraformVersion string
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
			wantErr:          true,
			terraformVersion: "0.15.0",
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
			dl := NewTerraformDirectoryLoader("", tt.terraformVersion, tt.options)
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

func TestGetRemoteLocation(t *testing.T) {
	type args struct {
		cache        map[string]string
		resourcePath string
	}
	tests := []struct {
		name          string
		args          args
		wantRemoteURL string
		wantTmpDir    string
	}{
		{
			name: "empty cache",
			args: args{
				resourcePath: "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
			},
			wantRemoteURL: "",
			wantTmpDir:    "",
		},
		{
			name: "resource is local",
			args: args{
				cache:        map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/"},
				resourcePath: "modules/db_parameter_group/main.tf",
			},
			wantRemoteURL: "",
			wantTmpDir:    "",
		},
		{
			name: "resource is local and in same scan dir",
			args: args{
				cache:        map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/"},
				resourcePath: "main.tf",
			},
			wantRemoteURL: "",
			wantTmpDir:    "",
		},
		{
			name: "tempdir is empty",
			args: args{
				cache:        map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": ""},
				resourcePath: "modules/db_parameter_group/main.tf",
			},
			wantRemoteURL: "",
			wantTmpDir:    "",
		},
		{
			name: "tempdir mapping is present cache",
			args: args{
				cache:        map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/", "git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.10.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791fcs/"},
				resourcePath: "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
			},
			wantRemoteURL: "git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0",
			wantTmpDir:    "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/",
		},
		{
			name: "source path is local and length of path is greater than tempDirs",
			args: args{
				cache:        map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/", "git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.10.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791fcs/"},
				resourcePath: "/user/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
			},
			wantRemoteURL: "",
			wantTmpDir:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemoteURL, gotTmpDir := GetRemoteLocation(tt.args.cache, tt.args.resourcePath)
			if gotRemoteURL != tt.wantRemoteURL {
				t.Errorf("GetRemoteLocation() gotRemoteURL = %v, want %v", gotRemoteURL, tt.wantRemoteURL)
			}
			if gotTmpDir != tt.wantTmpDir {
				t.Errorf("GetRemoteLocation() gotTmpDir = %v, want %v", gotTmpDir, tt.wantTmpDir)
			}
		})
	}
}

func TestGetConfigSource(t *testing.T) {
	type args struct {
		remoteURLMapping map[string]string
		resourceConfig   output.ResourceConfig
		absRootDir       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "remote module resource",
			args: args{
				remoteURLMapping: map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/"},
				resourceConfig: output.ResourceConfig{
					ID:     "azurerm_virtual_network.vnet",
					Name:   "vnet",
					Source: "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
				},
			},
			want:    "git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0/modules/db_parameter_group/main.tf",
			wantErr: false,
		},
		{
			name: "local module resource",
			args: args{
				remoteURLMapping: map[string]string{"git::https:/github.com/terraform-aws-modules/terraform-aws-rds?ref=v2.20.0": "/var/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/"},
				resourceConfig: output.ResourceConfig{
					ID:     "azurerm_virtual_network.vnet",
					Name:   "vnet",
					Source: "/user/folders/y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
				},
				absRootDir: "/user/folders/",
			},
			want:    "y5/y1qlrpl90rs_3n06z_qgjwv00000gn/T/791rns/modules/db_parameter_group/main.tf",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := GetConfigSource(tt.args.remoteURLMapping, tt.args.resourceConfig, tt.args.absRootDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConfigSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRemoteModuleIfPresentInTerraformSrc(t *testing.T) {
	absRootDir, err := filepath.Abs(filepath.Dir(filepath.Join("testdata", "terraform_cache_use_in_scan", "remote-module.tf")))
	if err != nil {
		zap.S().Error("error finding working directory", err)
	}
	terraformInitRegs := filepath.Join(absRootDir, terraformModuleInstallDir, "network")
	type fields struct {
		Cache map[string]TerraformModuleManifest
	}
	type args struct {
		req *hclConfigs.ModuleRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantSrc      string
		wantDestpath string
	}{
		{
			name: "module present in terraform cache",
			fields: fields{
				Cache: make(map[string]TerraformModuleManifest),
			},
			args: args{
				req: &hclConfigs.ModuleRequest{
					SourceAddr:      "Azure/network/azurerm",
					SourceAddrRange: hcl.Range{Filename: filepath.Join("testdata", "terraform_cache_use_in_scan", "remote-module.tf")},
				},
			},
			wantSrc:      "Azure/network/azurerm",
			wantDestpath: terraformInitRegs,
		},
		{
			name: "module not present in terraform cache",
			fields: fields{
				Cache: make(map[string]TerraformModuleManifest),
			},
			args: args{
				req: &hclConfigs.ModuleRequest{
					SourceAddr:      "Azure/network/azurermtest",
					SourceAddrRange: hcl.Range{Filename: filepath.Join("testdata", "terraform_cache_use_in_scan", "remote-module.tf")},
				},
			},
			wantSrc:      "",
			wantDestpath: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TerraformDirectoryLoader{
				absRootDir:               absRootDir,
				terraformInitModuleCache: tt.fields.Cache,
			}
			gotSrc, gotDestpath := tr.GetRemoteModuleIfPresentInTerraformSrc(tt.args.req)
			if gotSrc != tt.wantSrc {
				t.Errorf("TerraformModuleManifestCache.GetRemoteModuleIfPresentInTerraformSrc() gotSrc = %v, want %v", gotSrc, tt.wantSrc)
			}
			if gotDestpath != tt.wantDestpath {
				t.Errorf("TerraformModuleManifestCache.GetRemoteModuleIfPresentInTerraformSrc() gotDestpath = %v, want %v", gotDestpath, tt.wantDestpath)
			}
		})
	}
}

func TestTerraformDirectoryLoaderLoadIacDir(t *testing.T) {
	var nilMultiErr *multierror.Error = nil
	tests := []struct {
		name             string
		tfConfigDir      string
		tfJSONFile       string
		options          map[string]interface{}
		wantErr          error
		terraformVersion string
	}{
		{
			name:        "directory with resources having container defined",
			tfConfigDir: filepath.Join(testDataDir, "terraform-container-extraction"),
			tfJSONFile:  filepath.Join(tfJSONDir, "output-with-containers.json"),
			wantErr: multierror.Append(
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "terraform-container-extraction")),
				fmt.Errorf(invalidDirErrStringTemplate, filepath.Join(testDataDir, "terraform-container-extraction/terraform-aws-provider/task-definitions")),
			),
			terraformVersion: "0.15.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TerraformDirectoryLoader{
				absRootDir:               tt.tfConfigDir,
				remoteDownloader:         downloader.NewRemoteDownloader(),
				parser:                   hclConfigs.NewParser(afero.NewOsFs()),
				terraformInitModuleCache: make(map[string]TerraformModuleManifest),
				terraformVersion:         tt.terraformVersion,
			}
			got, gotErr := tr.LoadIacDir()
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.wantErr == nilMultiErr {
				if err := me.ErrorOrNil(); err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else if me.Error() != tt.wantErr.Error() {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			var want output.AllResourceConfigs

			// Read the expected value and unmarshal into want
			contents, _ := os.ReadFile(tt.tfJSONFile)
			if utils.IsWindowsPlatform() {
				contents = utils.ReplaceWinNewLineBytes(contents)
			}

			err := json.Unmarshal(contents, &want)
			if err != nil {
				t.Errorf("unexpected error unmarshalling want: %v", err)
			}

			match, err := test.IdenticalAllResourceConfigs(got, want)
			if err != nil {
				t.Errorf("unexpected error checking result: %v", err)
			}
			if !match {
				g, _ := json.MarshalIndent(got, "", "  ")
				w, _ := json.MarshalIndent(want, "", "  ")
				t.Errorf("got '%v', want: '%v'", string(g), string(w))
			}
		})
	}
}
