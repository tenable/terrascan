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

package helmv3

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name          string
		dirPath       string
		helmv3        HelmV3
		want          output.AllResourceConfigs
		wantErr       error
		resourceCount int
	}{
		{
			name:          "happy path (credit to madhuakula/kubernetes-goat)",
			dirPath:       "./testdata/happy-path",
			helmv3:        HelmV3{},
			wantErr:       nil,
			resourceCount: 3,
		},
		{
			name:          "happy path with subchart (credit to madhuakula/kubernetes-goat)",
			dirPath:       "./testdata/happy-path-with-subchart",
			helmv3:        HelmV3{},
			wantErr:       nil,
			resourceCount: 5,
		},
		{
			name:          "bad directory",
			dirPath:       "./testdata/bad-dir",
			helmv3:        HelmV3{},
			wantErr:       &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "./testdata/bad-dir"},
			resourceCount: 0,
		},
		{
			name:          "no helm charts in directory",
			dirPath:       "./testdata/no-helm-charts",
			helmv3:        HelmV3{},
			wantErr:       errNoHelmChartsFound,
			resourceCount: 0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resources, gotErr := tt.helmv3.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			resCount := len(resources)
			if resCount != tt.resourceCount {
				t.Errorf("resource count (%d) does not match expected (%d)", resCount, tt.resourceCount)
			}
		})
	}

}

func TestLoadChart(t *testing.T) {

	table := []struct {
		name      string
		chartPath string
		helmv3    HelmV3
		want      output.AllResourceConfigs
		wantErr   error
	}{
		{
			name:      "happy path (credit to madhuakula/kubernetes-goat)",
			chartPath: "./testdata/happy-path/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   nil,
		},
		{
			name:      "unreadable chart file",
			chartPath: "./testdata/bad-chart-file",
			helmv3:    HelmV3{},
			wantErr:   &os.PathError{Err: syscall.EISDIR, Op: "read", Path: "./testdata/bad-chart-file"},
		},
		{
			name:      "unmarshal bad chart",
			chartPath: "./testdata/bad-chart-file/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &yaml.TypeError{Errors: []string{"line 1: cannot unmarshal !!str `:bad ba...` into helmv3.helmChartData"}},
		},
		{
			name:      "chart path with no values.yaml",
			chartPath: "./testdata/chart-no-values/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &os.PathError{Err: syscall.ENOENT, Op: "stat", Path: "testdata/chart-no-values/values.yaml"},
		},
		{
			name:      "chart path with unreadable values.yaml",
			chartPath: "./testdata/chart-unreadable-values/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &os.PathError{Err: syscall.EISDIR, Op: "read", Path: "testdata/chart-unreadable-values/values.yaml"},
		},
		{
			name:      "chart path with unreadable values.yaml",
			chartPath: "./testdata/chart-bad-values/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &yaml.TypeError{Errors: []string{"line 1: cannot unmarshal !!str `:bad <bad` into map[string]interface {}"}},
		},
		{
			name:      "chart path no template dir",
			chartPath: "./testdata/chart-no-template-dir/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "testdata/chart-no-template-dir/templates"},
		},
		//{
		//	name:      "chart path skip test dir",
		//	chartPath: "./testdata/chart-skip-test-dir/Chart.yaml",
		//	helmv3:    HelmV3{},
		//	wantErr:   errSkipTestDir,
		//},
		{
			name:      "chart path bad template file",
			chartPath: "./testdata/chart-bad-template-file/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   &os.PathError{Err: syscall.EISDIR, Op: "read", Path: "testdata/chart-bad-template-file/templates/service.yaml"},
		},
		{
			name:      "chart path bad chart name",
			chartPath: "./testdata/chart-bad-name/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   errBadChartName,
		},
		{
			name:      "chart path bad chart version",
			chartPath: "./testdata/chart-bad-version/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   errBadChartVersion,
		},
		{
			name:      "chart path rendering error",
			chartPath: "./testdata/chart-rendering-error/Chart.yaml",
			helmv3:    HelmV3{},
			wantErr:   fmt.Errorf("parse error at (metadata-db/templates/ingress.yaml:40): unexpected {{end}}"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, _, gotErr := tt.helmv3.loadChart(tt.chartPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}

}
