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
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
)

var testDataDir = "testdata"
var multibasesDir = filepath.Join(testDataDir, "multibases")

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name          string
		kustomize     KustomizeDirectoryLoader
		want          output.AllResourceConfigs
		wantErr       error
		resourceCount int
		options       map[string]interface{}
	}{
		{
			name: "invalid dirPath",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: "not-there",
			},
			wantErr:       multierror.Append(&os.PathError{Err: syscall.ENOENT, Op: "open", Path: "not-there"}),
			resourceCount: 0,
		},
		{
			name: "simple-deployment",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(testDataDir, "simple-deployment"),
			},
			resourceCount: 4,
		},
		{
			name: "multibases",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(multibasesDir, "base"),
			},
			resourceCount: 2,
		},
		{
			name: "multibases",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(multibasesDir, "dev"),
			},
			resourceCount: 2,
		},
		{
			name: "multibases",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(multibasesDir, "prod"),
			},
			resourceCount: 2,
		},

		{
			name: "multibases",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(multibasesDir, "stage"),
			},
			resourceCount: 2,
		},
		{
			name: "multibases",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: multibasesDir,
			},
			resourceCount: 4,
		},
		{
			name: "no-kustomize-directory",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(testDataDir, "no-kustomizefile"),
			},
			wantErr:       multierror.Append(fmt.Errorf("kustomization.y(a)ml file not found in the directory %s", filepath.Join(testDataDir, "no-kustomizefile"))),
			resourceCount: 0,
		},
		{
			name: "kustomize-file-empty",
			kustomize: KustomizeDirectoryLoader{
				absRootDir: filepath.Join(testDataDir, "kustomize-file-empty"),
			},
			wantErr:       multierror.Append(fmt.Errorf("unable to read the kustomization file in the directory %s, error: yaml file is empty", filepath.Join(testDataDir, "kustomize-file-empty"))),
			resourceCount: 0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resourceMap, gotErr := tt.kustomize.LoadIacDir()
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

			resCount := resourceMap.GetResourceCount()
			if resCount != tt.resourceCount {
				t.Errorf("resource count (%d) does not match expected (%d)", resCount, tt.resourceCount)
			}
		})
	}

}

func TestLoadKustomize(t *testing.T) {
	kustomizeYaml := "kustomization.yaml"
	kustomizeYml := "kustomization.yml"

	table := []struct {
		name            string
		basepath        string
		filename        string
		want            output.AllResourceConfigs
		wantErr         error
		checkPrefix     bool
		options         map[string]interface{}
		kustomizeBinary bool
		exe             string
	}{
		{
			name:            "simple-deployment",
			basepath:        filepath.Join(testDataDir, "simple-deployment"),
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases",
			basepath:        multibasesDir,
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases/base",
			basepath:        filepath.Join(multibasesDir, "base"),
			filename:        kustomizeYml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases/dev",
			basepath:        filepath.Join(multibasesDir, "dev"),
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases/prod",
			basepath:        filepath.Join(multibasesDir, "prod"),
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases/stage",
			basepath:        filepath.Join(multibasesDir, "stage"),
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "multibases/zero-violation-base",
			basepath:        filepath.Join(multibasesDir, "zero-violation-base"),
			filename:        kustomizeYaml,
			wantErr:         nil,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "erroneous-pod",
			basepath:        filepath.Join(testDataDir, "erroneous-pod"),
			filename:        kustomizeYaml,
			wantErr:         fmt.Errorf("error from kustomization"),
			checkPrefix:     true,
			kustomizeBinary: false,
			exe:             "",
		},
		{
			name:            "erroneous-deployment",
			basepath:        filepath.Join(testDataDir, "erroneous-deployment/"),
			filename:        kustomizeYaml,
			wantErr:         fmt.Errorf("error from kustomization"),
			checkPrefix:     true,
			kustomizeBinary: false,
			exe:             "",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := LoadKustomize(tt.basepath, tt.filename, tt.exe, tt.kustomizeBinary)
			if tt.checkPrefix {
				if tt.wantErr != nil && !strings.HasPrefix(gotErr.Error(), tt.wantErr.Error()) {
					t.Errorf("unexpected error; gotErr: '%v', expected prefix: '%v'", gotErr, tt.wantErr.Error())
				}
			} else if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
