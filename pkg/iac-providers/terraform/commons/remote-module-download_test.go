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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/downloader"
	"github.com/accurics/terrascan/pkg/utils"
	version "github.com/hashicorp/go-version"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry/regsrc"
)

func TestRemoteModuleInstallerDownloadRemoteModule(t *testing.T) {

	// disable terraform logs when TF_LOG env variable is not set
	if os.Getenv("TF_LOG") == "" {
		log.SetOutput(ioutil.Discard)
	}

	testConstraintsSecurityGroup, _ := version.NewConstraint("3.17.0")
	testModuleSecurityGroup, _ := regsrc.ParseModuleSource("terraform-aws-modules/security-group/aws")
	testModuleInvalidProvider, _ := regsrc.ParseModuleSource("terraform-aws-modules/testgroup/testprovider")

	type fields struct {
		cache      InstalledCache
		downloader downloader.Downloader
	}
	type args struct {
		requiredVersion hclConfigs.VersionConstraint
		module          *regsrc.Module
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "remote module download with valid module and version",
			fields: fields{
				cache:      make(map[string]string),
				downloader: downloader.NewDownloader(),
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
				cache:      make(map[string]string),
				downloader: downloader.NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
		},
		{
			name: "remote module download with invalid module source",
			fields: fields{
				cache:      make(map[string]string),
				downloader: downloader.NewDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleInvalidProvider,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RemoteModuleInstaller{
				cache:      tt.fields.cache,
				downloader: tt.fields.downloader,
			}
			testDir := filepath.Join(os.TempDir(), utils.GenRandomString(6))
			tt.want = testDir
			if tt.wantErr {
				tt.want = ""
			}

			defer os.RemoveAll(testDir)
			got, err := r.DownloadRemoteModule(tt.args.requiredVersion, testDir, tt.args.module)
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
