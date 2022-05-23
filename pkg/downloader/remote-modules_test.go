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
	"testing"
)

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
			addr: "git@github.com:tenable/terrascan.git",
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
			got := IsLocalSourceAddr(tt.addr)
			if got != tt.want {
				t.Errorf("got: '%v', want: '%v'", got, tt.want)
			}
		})
	}
}

func TestIsRegistrySourceAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid terraform registry source without host",
			args: args{
				addr: "terraform-aws-modules/eks/aws",
			},
			want: true,
		},
		{
			name: "valid terraform registry source with host",
			args: args{
				addr: "registry.terraform.io/terraform-aws-modules/eks/aws",
			},
			want: true,
		},
		{
			name: "invalid terraform registry source - 1",
			args: args{
				addr: "eks/azure",
			},
			want: false,
		},
		{
			name: "invalid terraform registry source - 2",
			args: args{
				addr: "test/terraform/invalid/eks/azure",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRegistrySourceAddr(tt.args.addr); got != tt.want {
				t.Errorf("IsRegistrySourceAddr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidRemoteType(t *testing.T) {
	type args struct {
		remoteType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid remote type - git",
			args: args{
				remoteType: "git",
			},
			want: true,
		},
		{
			name: "valid remote type - terraform-registry",
			args: args{
				remoteType: "terraform-registry",
			},
			want: true,
		},
		{
			name: "valid remote type - Git",
			args: args{
				remoteType: "Git",
			},
			want: true,
		},
		{
			name: "invalid remote type - test",
			args: args{
				remoteType: "test",
			},
		},
		{
			name: "invalid remote type - terraformRegistry",
			args: args{
				remoteType: "terraformRegistry",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRemoteType(tt.args.remoteType); got != tt.want {
				t.Errorf("IsValidRemoteType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRemoteTypeTerraformRegistry(t *testing.T) {
	type args struct {
		remoteType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid remote type terraform-registry",
			args: args{
				remoteType: "terraform-registry",
			},
			want: true,
		},
		{
			name: "valid remote type terraform-registry in mixed case and space",
			args: args{
				remoteType: " TerraForm-Registry  ",
			},
			want: true,
		},
		{
			name: "invalid remote type terraform-registry",
			args: args{
				remoteType: "terraformRegistry",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRemoteTypeTerraformRegistry(tt.args.remoteType); got != tt.want {
				t.Errorf("IsRemoteTypeTerraformRegistry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSourceAddrAndVersion(t *testing.T) {
	testSourceURL := "terraform-aws-modules/eks/aws"
	testVersion := "2.20.0"

	type args struct {
		sourceURL string
	}
	tests := []struct {
		name        string
		args        args
		wantSource  string
		wantVersion string
	}{
		{
			name: "source address with version",
			args: args{
				sourceURL: testSourceURL + ":" + testVersion,
			},
			wantSource:  testSourceURL,
			wantVersion: testVersion,
		},
		{
			name: "source address without version",
			args: args{
				sourceURL: testSourceURL,
			},
			wantSource: testSourceURL,
		},
		{
			name: "source address with empty version",
			args: args{
				sourceURL: testSourceURL + ": ",
			},
			wantSource: testSourceURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetSourceAddrAndVersion(tt.args.sourceURL)
			if got != tt.wantSource {
				t.Errorf("GetSourceAddrAndVersion() got source = %v, want %v", got, tt.wantSource)
			}
			if got1 != tt.wantVersion {
				t.Errorf("GetSourceAddrAndVersion() got version = %v, want %v", got1, tt.wantVersion)
			}
		})
	}
}
