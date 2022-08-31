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

package output

import (
	"reflect"
	"testing"
)

func TestAllResourceConfigsGetResourceCount(t *testing.T) {
	tests := []struct {
		name      string
		a         AllResourceConfigs
		wantCount int
	}{
		{
			name:      "nil AllResourceConfigs",
			a:         nil,
			wantCount: 0,
		},
		{
			name: "non nil AllResourceConfigs",
			a: AllResourceConfigs{
				"key1": {
					ResourceConfig{},
					ResourceConfig{},
				},
				"key2": {
					ResourceConfig{},
					ResourceConfig{},
				},
				"key3": {
					ResourceConfig{},
					ResourceConfig{},
				},
				"key4": nil,
			},
			wantCount: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCount := tt.a.GetResourceCount(); gotCount != tt.wantCount {
				t.Errorf("AllResourceConfigs.GetResourceCount() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestAllResourceConfigsUpdateResourceConfigs(t *testing.T) {
	type args struct {
		resourceType string
		resources    []ResourceConfig
	}
	tests := []struct {
		name       string
		a          AllResourceConfigs
		args       args
		wantLength int
	}{
		{
			name: "empty resources",
			args: args{
				resourceType: "s3",
			},
			a:          nil,
			wantLength: 0,
		},
		{
			name: "resource present",
			args: args{
				resourceType: "s3",
				resources: []ResourceConfig{
					{
						Name:   "s3_bucket",
						Source: "s3_bucket.tf",
					},
				},
			},
			a: AllResourceConfigs{
				"s3": []ResourceConfig{
					{
						Name:   "s3_bucket",
						Source: "s3_bucket.tf",
					},
				},
			},
			wantLength: 1,
		},
		{
			name: "resource not present, but resource type",
			args: args{
				resourceType: "pod",
				resources: []ResourceConfig{
					{
						Name:   "terra_controller",
						Source: "terra_controller.yml",
					},
				},
			},
			a: AllResourceConfigs{
				"s3": []ResourceConfig{
					{
						Name:   "s3_bucket",
						Source: "s3_bucket.tf",
					},
				},
				"pod": []ResourceConfig{
					{
						Name:   "some_pod",
						Source: "some_pod.yml",
					},
				},
			},
			wantLength: 3,
		},
		{
			name: "resource and resource type both not present",
			args: args{
				resourceType: "job",
				resources: []ResourceConfig{
					{
						Name:   "cron_job",
						Source: "cron_job.yaml",
					},
				},
			},
			a: AllResourceConfigs{
				"s3": []ResourceConfig{
					{
						Name:   "s3_bucket",
						Source: "s3_bucket.tf",
					},
					{
						Name:   "zods_s3_bucket",
						Source: "zods_s3_bucket.tf",
					},
				},
				"pod": []ResourceConfig{
					{
						Name:   "some_pod",
						Source: "some_pod.yml",
					},
				},
			},
			wantLength: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.UpdateResourceConfigs(tt.args.resourceType, tt.args.resources)
			if tt.wantLength != tt.a.GetResourceCount() {
				t.Errorf("expected length of all resource config is %d, got %d", tt.wantLength, tt.a.GetResourceCount())
			}
		})
	}
}

func TestIsConfigPresent(t *testing.T) {
	testResourceConfig := ResourceConfig{
		Name:   "s3_bucket_hulk",
		Source: "test.tf",
		Config: map[string]interface{}{
			"key1": "smash",
		},
	}

	type args struct {
		resources      []ResourceConfig
		resourceConfig ResourceConfig
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty resources",
			args: args{
				resourceConfig: ResourceConfig{},
			},
			want: false,
		},
		{
			name: "resource name and source are not same",
			args: args{
				resources: []ResourceConfig{
					{
						Name:   "s3_bucket_thor",
						Source: "test.yaml",
					},
					{
						Name:   "ec2_instance_bruce",
						Source: "test.yaml",
					},
				},
				resourceConfig: testResourceConfig,
			},
			want: false,
		},
		{
			name: "resource name and source are not same",
			args: args{
				resources: []ResourceConfig{
					testResourceConfig,
					{
						Name:   "ec2_instance_bruce",
						Source: "test.yaml",
					},
				},
				resourceConfig: testResourceConfig,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsConfigPresent(tt.args.resources, tt.args.resourceConfig); got != tt.want {
				t.Errorf("IsConfigPresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllResourceConfigsFindAllResourcesByID(t *testing.T) {
	testS3ResourceConfig := ResourceConfig{
		ID: "s3.my_s3_bucket",
	}

	testS3LongIDResourceConfig := ResourceConfig{
		ID: "module.somemodule.s3.my_s3_bucket",
	}

	testResourceConfigList := []*ResourceConfig{&testS3ResourceConfig}

	type args struct {
		resourceID string
	}
	tests := []struct {
		name    string
		a       AllResourceConfigs
		args    args
		want    []*ResourceConfig
		wantErr bool
	}{
		{
			name:    "nil AllResourceConfigs",
			a:       nil,
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid resource id",
			a: AllResourceConfigs{
				"key": {},
			},
			args: args{
				resourceID: "id",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "one resource present in AllResourceConfigs",
			a: AllResourceConfigs{
				"s3": {
					testS3ResourceConfig,
				},
			},
			args: args{
				resourceID: "s3.my_s3_bucket",
			},
			want: testResourceConfigList,
		},
		{
			name: "multiple resources present in AllResourceConfigs",
			a: AllResourceConfigs{
				"s3": {
					testS3ResourceConfig,
					testS3ResourceConfig,
				},
				"ingress": {
					ResourceConfig{
						ID: "allow_ssh",
					},
				},
			},
			args: args{
				resourceID: "s3.my_s3_bucket",
			},
			want: []*ResourceConfig{
				&testS3ResourceConfig,
				&testS3ResourceConfig,
			},
		},
		{
			name: "resource not present in AllResourceConfigs",
			a: AllResourceConfigs{
				"s3": {
					testS3ResourceConfig,
				},
			},
			args: args{
				resourceID: "ec2.test_instance",
			},
			want: []*ResourceConfig{},
		},
		{
			name: "long resource ID",
			a: AllResourceConfigs{
				"s3": {
					testS3LongIDResourceConfig,
				},
			},
			args: args{
				resourceID: "module.somemodule.s3.my_s3_bucket",
			},
			want: []*ResourceConfig{
				&testS3LongIDResourceConfig,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.FindAllResourcesByID(tt.args.resourceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AllResourceConfigs.FindAllResourcesByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllResourceConfigs.FindAllResourcesByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
