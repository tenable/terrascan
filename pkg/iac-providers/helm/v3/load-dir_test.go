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
	"path"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"github.com/hashicorp/go-multierror"
)

var testDataDir = "testdata"

func TestLoadIacDir(t *testing.T) {

	invalidDirErr := &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: filepath.Join(testDataDir, "bad-dir")}
	if utils.IsWindowsPlatform() {
		invalidDirErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: filepath.Join(testDataDir, "bad-dir")}
	}

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
			dirPath:       filepath.Join(testDataDir, "happy-path"),
			helmv3:        HelmV3{},
			resourceCount: 3,
		},
		{
			name:          "happy path with subchart (credit to madhuakula/kubernetes-goat)",
			dirPath:       filepath.Join(testDataDir, "happy-path-with-subchart"),
			helmv3:        HelmV3{},
			resourceCount: 5,
		},
		{
			name:          "bad directory",
			dirPath:       filepath.Join(testDataDir, "bad-dir"),
			helmv3:        HelmV3{},
			wantErr:       multierror.Append(invalidDirErr),
			resourceCount: 0,
		},
		{
			name:          "no helm charts in directory",
			dirPath:       filepath.Join(testDataDir, "no-helm-charts"),
			helmv3:        HelmV3{},
			wantErr:       multierror.Append(fmt.Errorf("no helm charts found in directory %s", filepath.Join(testDataDir, "no-helm-charts"))),
			resourceCount: 0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resources, gotErr := tt.helmv3.LoadIacDir(tt.dirPath, false)
			me, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
			if tt.wantErr == nil {
				if err := me.ErrorOrNil(); err != nil {
					t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else if me.Error() != tt.wantErr.Error() {
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

	chartPathNoValuesYAMLErr := &os.PathError{Err: syscall.ENOENT, Op: "stat", Path: filepath.Join(testDataDir, "chart-no-values", "values.yaml")}
	chartPathNoTemplateDirErr := &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: filepath.Join(testDataDir, "chart-no-template-dir", "templates")}
	unreadableChartFileErr := &os.PathError{Err: syscall.EISDIR, Op: "read", Path: filepath.Join(testDataDir, "bad-chart-file")}
	chartPathUnreadableValuesErr := &os.PathError{Err: syscall.EISDIR, Op: "read", Path: filepath.Join(testDataDir, "chart-unreadable-values", "values.yaml")}
	chartPathBadTemplateErr := &os.PathError{Err: syscall.EISDIR, Op: "read", Path: filepath.Join(testDataDir, "chart-bad-template-file", "templates", "service.yaml")}

	if utils.IsWindowsPlatform() {
		chartPathNoValuesYAMLErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: filepath.Join(testDataDir, "chart-no-values", "values.yaml")}
		chartPathNoTemplateDirErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: filepath.Join(testDataDir, "chart-no-template-dir", "templates")}
		unreadableChartFileErr = &os.PathError{Err: syscall.Errno(6), Op: "read", Path: filepath.Join(testDataDir, "bad-chart-file")}
		chartPathUnreadableValuesErr = &os.PathError{Err: syscall.Errno(6), Op: "read", Path: filepath.Join(testDataDir, "chart-unreadable-values", "values.yaml")}
		chartPathBadTemplateErr = &os.PathError{Err: syscall.Errno(6), Op: "read", Path: filepath.Join(testDataDir, "chart-bad-template-file", "templates", "service.yaml")}
	}

	table := []struct {
		name      string
		chartPath string
		helmv3    HelmV3
		want      output.AllResourceConfigs
		wantErr   error
	}{
		{
			name:      "happy path (credit to madhuakula/kubernetes-goat)",
			chartPath: filepath.Join(testDataDir, "happy-path", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   nil,
		},
		{
			name:      "unreadable chart file",
			chartPath: filepath.Join(testDataDir, "bad-chart-file"),
			helmv3:    HelmV3{},
			wantErr:   unreadableChartFileErr,
		},
		{
			name:      "unmarshal bad chart",
			chartPath: filepath.Join(testDataDir, "bad-chart-file", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   &yaml.TypeError{Errors: []string{"line 1: cannot unmarshal !!str `:bad ba...` into helmv3.helmChartData"}},
		},
		{
			name:      "chart path with no values.yaml",
			chartPath: filepath.Join(testDataDir, "chart-no-values", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   chartPathNoValuesYAMLErr,
		},
		{
			name:      "chart path with unreadable values.yaml",
			chartPath: filepath.Join(testDataDir, "chart-unreadable-values", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   chartPathUnreadableValuesErr,
		},
		{
			name:      "chart path with unreadable values.yaml",
			chartPath: filepath.Join(testDataDir, "chart-bad-values", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   &yaml.TypeError{Errors: []string{"line 1: cannot unmarshal !!str `:bad <bad` into map[string]interface {}"}},
		},
		{
			name:      "chart path no template dir",
			chartPath: filepath.Join(testDataDir, "chart-no-template-dir", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   chartPathNoTemplateDirErr,
		},
		//{
		//	name:      "chart path skip test dir",
		//	chartPath: "./testdata/chart-skip-test-dir/Chart.yaml",
		//	helmv3:    HelmV3{},
		//	wantErr:   errSkipTestDir,
		//},
		{
			name:      "chart path bad template file",
			chartPath: filepath.Join(testDataDir, "chart-bad-template-file", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   chartPathBadTemplateErr,
		},
		{
			name:      "chart path bad chart name",
			chartPath: filepath.Join(testDataDir, "chart-bad-name", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   errBadChartName,
		},
		{
			name:      "chart path bad chart version",
			chartPath: filepath.Join(testDataDir, "chart-bad-version", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   errBadChartVersion,
		},
		{
			name:      "chart path rendering error",
			chartPath: filepath.Join(testDataDir, "chart-rendering-error", "Chart.yaml"),
			helmv3:    HelmV3{},
			wantErr:   fmt.Errorf("parse error at (%s:40): unexpected {{end}}", path.Join("metadata-db", filepath.Join("templates", "ingress.yaml"))),
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
