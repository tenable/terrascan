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

package downloader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/registry/regsrc"
	"github.com/hashicorp/terraform/registry/response"
	"github.com/tenable/terrascan/pkg/utils"
)

var (
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

// terraformRegistryClientMocks is a mock implementation of terraform registry client
type MockTerraformRegistryClient struct {
	location       string
	moduleVersions *response.ModuleVersions
	versionError   error
	locationError  error
}

// ModuleVersions mock method
func (t *MockTerraformRegistryClient) ModuleVersions(module *regsrc.Module) (*response.ModuleVersions, error) {
	if t.versionError != nil {
		return nil, t.versionError
	}
	return t.moduleVersions, nil
}

// ModuleLocation mock method
func (t *MockTerraformRegistryClient) ModuleLocation(module *regsrc.Module, version string) (string, error) {
	if t.locationError != nil {
		return "", t.locationError
	}
	return t.location, nil
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
		r          *remoteModuleInstaller
	}{
		{
			name:       "GetURLSubDir error",
			addr:       someAddr,
			dest:       someDest,
			wantOutput: "",
			wantErr:    wantErr,
			r: &remoteModuleInstaller{
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
			r: &remoteModuleInstaller{
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
			r: &remoteModuleInstaller{
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
			r: &remoteModuleInstaller{
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
			r: &remoteModuleInstaller{
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
			r: &remoteModuleInstaller{
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
		tempDir, err := os.MkdirTemp("", "")
		if err != nil {
			t.Errorf("failed to create temp dir. error: '%v'", err)
		}

		// store temp dir into the installer cache
		r := &remoteModuleInstaller{
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
		log.SetOutput(io.Discard)
	}

	testConstraintsSecurityGroup, _ := version.NewConstraint("3.17.0")
	testModuleSecurityGroup, _ := regsrc.ParseModuleSource("terraform-aws-modules/security-group/aws")
	testModuleInvalidProvider, _ := regsrc.ParseModuleSource("terraform-aws-modules/testgroup/testprovider")
	testModuleRdsWithRawSubModule, _ := regsrc.ParseModuleSource("terraform-aws-modules/security-group/aws//db_subnet_group")
	testRawSubModuleName := "db_subnet_group"
	testModuleVersions := []*response.ModuleProviderVersions{
		{
			Source: "source",
			Versions: []*response.ModuleVersion{
				{
					Version: "12.3.1",
				},
			},
		},
	}
	testNoModuleVersions := []*response.ModuleProviderVersions{
		{
			Source:   "no version",
			Versions: []*response.ModuleVersion{},
		},
	}

	type fields struct {
		moduleDownloader ModuleDownloader
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
				moduleDownloader: NewRemoteDownloader(),
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
				moduleDownloader: NewRemoteDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
		},
		{
			name: "remote module download with invalid module source",
			fields: fields{
				moduleDownloader: NewRemoteDownloader(),
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
				moduleDownloader: NewRemoteDownloader(),
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleRdsWithRawSubModule,
			},
			hasRawSubModule: true,
			rawSubModule:    testRawSubModuleName,
		},
		{
			name: "error while getting module versions",
			fields: fields{
				moduleDownloader: &remoteModuleInstaller{
					cache:      make(map[string]string),
					downloader: &MockDownloader{},
					terraformRegistryClient: &MockTerraformRegistryClient{
						versionError: fmt.Errorf("Error while getting versions"),
					},
				},
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
			wantErr: true,
		},
		{
			name: "no versions available to download",
			fields: fields{
				moduleDownloader: &remoteModuleInstaller{
					cache:      make(map[string]string),
					downloader: &MockDownloader{},
					terraformRegistryClient: &MockTerraformRegistryClient{
						moduleVersions: &response.ModuleVersions{
							Modules: testNoModuleVersions,
						},
					},
				},
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
			wantErr: true,
		},
		{
			name: "error while getting module location",
			fields: fields{
				moduleDownloader: &remoteModuleInstaller{
					cache:      make(map[string]string),
					downloader: &MockDownloader{},
					terraformRegistryClient: &MockTerraformRegistryClient{
						moduleVersions: &response.ModuleVersions{
							Modules: testModuleVersions,
						},
						locationError: fmt.Errorf("Error getting module location"),
					},
				},
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
			wantErr: true,
		},
		{
			name: "error while downloading module",
			fields: fields{
				moduleDownloader: &remoteModuleInstaller{
					cache: make(map[string]string),
					downloader: &MockDownloader{
						errGetURLSubDir: fmt.Errorf("Error in download"),
					},
					terraformRegistryClient: &MockTerraformRegistryClient{
						moduleVersions: &response.ModuleVersions{
							Modules: testModuleVersions,
						},
						location: "test",
					},
				},
			},
			args: args{
				requiredVersion: hclConfigs.VersionConstraint{},
				module:          testModuleSecurityGroup,
			},
			wantErr: true,
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
			got, err := tt.fields.moduleDownloader.DownloadRemoteModule(tt.args.requiredVersion, testDir, tt.args.module)
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

func TestAuthenticatedRegistryClient(t *testing.T) {
	testDirPath := filepath.Join("testdata", "run-test")

	tests := []struct {
		name     string
		filename string
	}{
		{
			name:     "valid terraformrc file",
			filename: "terraformrc",
		},
		{
			name:     "invalid terraformrc file",
			filename: "badterraformrc",
		},
		{
			name:     "invalid terraformrc file",
			filename: "nonexistentfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRCFilePath := filepath.Join(testDirPath, tt.filename)
			got := NewAuthenticatedRegistryClient(testRCFilePath)
			if got == nil {
				t.Errorf("NewAuthenticatedRegistryClient(%s) = nil, expected terraformRegistryClient", testRCFilePath)
			}
		})
	}
}

func TestBuildDiscoServices(t *testing.T) {
	testDirPath := "testdata"

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid terraformrc file to parse",
			filename: "terraformrc",
			wantErr:  false,
		},
		{
			name:     "invalid terraformrc file to parse",
			filename: "badterraformrc",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRCFilePath := filepath.Join(testDirPath, tt.filename)
			b, err := os.ReadFile(testRCFilePath)
			if err != nil {
				t.Errorf("Error reading %s: %s", testRCFilePath, err)
				return
			}

			_, err = buildDiscoServices(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildDiscoServices() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestConvertCredentialMapToHostMap(t *testing.T) {
	tests := []struct {
		name          string
		credentialMap map[string]map[string]interface{}
		wantNilResult bool
	}{
		{
			name:          "nil credentialMap",
			credentialMap: nil,
			wantNilResult: true,
		},
		{
			name:          "empty credentialMap",
			credentialMap: make(map[string]map[string]interface{}),
			wantNilResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostMap := convertCredentialMapToHostMap(tt.credentialMap)
			if (hostMap == nil) != tt.wantNilResult {
				t.Errorf("convertCredentialMapToHostMap() error: got nil result %v, wantNilResult %v", hostMap == nil, tt.wantNilResult)
				return
			}
		})
	}

}
