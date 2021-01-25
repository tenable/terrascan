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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-version"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/hashicorp/terraform/registry/response"
)

func TestCleanUp(t *testing.T) {

	t.Run("test cache clean up", func(t *testing.T) {

		// create temp dir
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Errorf("failed to create temp dir. error: '%v'", err)
		}

		// store temp dir into the installer cache
		r := &GoGetter{
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

func TestDownloadRemoteModule(t *testing.T) {

	// disable terraform logs when TF_LOG env variable is not set
	if os.Getenv("TF_LOG") == "" {
		log.SetOutput(ioutil.Discard)
	}

	testConstraintsSecurityGroup, _ := version.NewConstraint("3.17.0")
	testModuleSecurityGroup, _ := regsrc.ParseModuleSource("terraform-aws-modules/security-group/aws")
	testModuleInvalidProvider, _ := regsrc.ParseModuleSource("terraform-aws-modules/testgroup/testprovider")
	testModuleRdsWithRawSubModule, _ := regsrc.ParseModuleSource("terraform-aws-modules/security-group/aws//db_subnet_group")
	testRawSubModuleName := "db_subnet_group"

	type fields struct {
		downloader Downloader
	}
	type args struct {
		requiredVersion hclConfigs.VersionConstraint
		module          *regsrc.Module
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            string
		wantErr         bool
		hasRawSubModule bool
		rawSubModule    string
	}{
		{
			name: "remote module download with valid module and version",
			fields: fields{
				downloader: NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{
					Required: testConstraintsSecurityGroup,
				},
				module: testModuleSecurityGroup,
			},
		},
		{
			name: "remote module download with valid module, without version",
			fields: fields{
				downloader: NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
		},
		{
			name: "remote module download with invalid module source",
			fields: fields{
				downloader: NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleInvalidProvider,
			},
			wantErr: true,
		},
		{
			name: "remote module download with raw sub module",
			fields: fields{
				downloader: NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleRdsWithRawSubModule,
			},
			hasRawSubModule: true,
			rawSubModule:    testRawSubModuleName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
			tt.want = testDir
			if tt.wantErr {
				tt.want = ""
			}

			if tt.hasRawSubModule {
				tt.want = tt.want + string(os.PathSeparator) + tt.rawSubModule
			}

			defer os.RemoveAll(testDir)
			got, err := tt.fields.downloader.DownloadRemoteModule(tt.args.requiredVersion, testDir, tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteModuleInstaller.DownloadRemoteModule() got error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RemoteModuleInstaller.DownloadRemoteModule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVersionToDownload(t *testing.T) {
	source := "terraform-aws-modules/security-group/aws"
	testModule, _ := regsrc.ParseModuleSource(source)
	testVersionConstraint, _ := version.NewConstraint("3.17.0")
	testVersionConstraintMatching, _ := version.NewConstraint("1.2.0")
	testVersion, _ := version.NewVersion("1.2.0")
	testVersions := []*response.ModuleVersion{
		{
			Version: "1.0.0",
		},
		{
			Version: "1.1.0",
		},
		{
			Version: "1.2.0",
		},
	}

	type args struct {
		moduleVersions  *response.ModuleVersions
		requiredVersion hclConfigs.VersionConstraint
		module          *regsrc.Module
	}
	tests := []struct {
		name    string
		args    args
		want    *version.Version
		wantErr bool
	}{
		{
			name: "invalid version",
			args: args{
				moduleVersions: &response.ModuleVersions{
					Modules: []*response.ModuleProviderVersions{
						{
							Source: "",
							Versions: []*response.ModuleVersion{
								{},
							},
						},
					},
				},
				module: testModule,
				requiredVersion: hclConfigs.VersionConstraint{
					Required: testVersionConstraint,
				},
			},
			wantErr: true,
		},
		{
			name: "no matching versions",
			args: args{
				moduleVersions: &response.ModuleVersions{
					Modules: []*response.ModuleProviderVersions{
						{
							Source:   "source",
							Versions: testVersions,
						},
					},
				},
				module: testModule,
				requiredVersion: hclConfigs.VersionConstraint{
					Required: testVersionConstraint,
				},
			},
			wantErr: true,
		},
		{
			name: "no required version specified",
			args: args{
				moduleVersions: &response.ModuleVersions{
					Modules: []*response.ModuleProviderVersions{
						{
							Source:   "source",
							Versions: testVersions,
						},
					},
				},
				module: testModule,
			},
			want: testVersion,
		},
		{
			name: "required version specified",
			args: args{
				moduleVersions: &response.ModuleVersions{
					Modules: []*response.ModuleProviderVersions{
						{
							Source:   "source",
							Versions: testVersions,
						},
					},
				},
				module: testModule,
				requiredVersion: hclConfigs.VersionConstraint{
					Required: testVersionConstraintMatching,
				},
			},
			want: testVersion,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getVersionToDownload(tt.args.moduleVersions, tt.args.requiredVersion, tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVersionToDownload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVersionToDownload() = %v, want %v", got, tt.want)
			}
		})
	}
}
