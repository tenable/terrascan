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

package k8sv1

import (
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/iac-providers/kubernetes/v1/testdata"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

func TestK8sV1ExtractContainerImages(t *testing.T) {
	type args struct {
		doc  *utils.IacDocument
		kind string
	}
	tests := []struct {
		name                       string
		k                          *K8sV1
		args                       args
		wantContainerImageList     []output.ContainerDetails
		wantInitContainerImageList []output.ContainerDetails
		wantErr                    bool
	}{
		{
			name: "empty document object",
			args: args{
				doc:  &utils.IacDocument{},
				kind: "CRD",
			},
			wantErr:                    true,
			wantContainerImageList:     []output.ContainerDetails{},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "pod json document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "json",
					Data: testdata.PodJSONTemplate,
				},
				kind: "Pod",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "healthz", Image: "k8s.gcr.io/exechealthz-amd64:1.2"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "pod yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.PodYAMLTemplate,
				},
				kind: "Pod",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "myapp-container", Image: "nginx"}},
			wantInitContainerImageList: []output.ContainerDetails{{Name: "myapp-container", Image: "busybox"}},
		},
		{
			name: "cronjob yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.CronJobYAMLTemplate,
				},
				kind: "CronJob",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "hello", Image: "busybox"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "job yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.JobYAMLTemplate,
				},
				kind: "Job",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "c", Image: "gcr.io/terrascan/job-wq-1"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "deployment yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.DeploymentYAMLTemplate,
				},
				kind: "Deployment",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "nginx", Image: "nginx:1.14.2"}},
			wantInitContainerImageList: []output.ContainerDetails{{Name: "init", Image: "busybox"}},
		},
		{
			name: "daemonset yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.DaemonSetYAMLTemplate,
				},
				kind: "DaemonSet",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "fluentd-elasticsearch", Image: "quay.io/fluentd_elasticsearch/fluentd:v2.5.2"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "replicaset yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.ReplicaSetYAMLTemplate,
				},
				kind: "ReplicaSet",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "php-redis", Image: "gcr.io/google_samples/gb-frontend:v3"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
		{
			name: "replicationcontroller yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.ReplicationControllerTemplate,
				},
				kind: "ReplicationController",
			},
			wantContainerImageList: []output.ContainerDetails{
				{Name: "nginx", Image: "nginx:latest"},
				{Name: "sidecar1", Image: "sidecar-image-1"},
				{Name: "sidecar2", Image: "sidecar-image-2"}},
			wantInitContainerImageList: []output.ContainerDetails{
				{Name: "init1", Image: "init-image-1"},
				{Name: "init2", Image: "init-image-2"},
				{Name: "init3", Image: "init-image-3"}},
		},
		{
			name: "statefulSet yaml document object",
			args: args{
				doc: &utils.IacDocument{
					Type: "yaml",
					Data: testdata.StatefulSetTemplate,
				},
				kind: "StatefulSet",
			},
			wantContainerImageList:     []output.ContainerDetails{{Name: "nginx", Image: "k8s.gcr.io/nginx-slim:0.8"}},
			wantInitContainerImageList: []output.ContainerDetails{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &K8sV1{}
			gotContainerImageList, gotInitContainerImageList, err := k.extractContainerImages(tt.args.kind, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sV1.extractContainerImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContainerImageList, tt.wantContainerImageList) {
				t.Errorf("K8sV1.extractResource() got = %v, want %v", gotContainerImageList, tt.wantContainerImageList)
			}
			if !reflect.DeepEqual(gotInitContainerImageList, tt.wantInitContainerImageList) {
				t.Errorf("K8sV1.extractResource() got InitContainerImageList = %v, want %v", gotInitContainerImageList, tt.wantInitContainerImageList)
			}
		})
	}
}
