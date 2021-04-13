package kustomizev3

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

const kustomizeErrPrefix = "error from kustomization."

var testDataDir = "testdata"
var multibasesDir = filepath.Join(testDataDir, "multibases")

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name          string
		dirPath       string
		kustomize     KustomizeV3
		want          output.AllResourceConfigs
		wantErr       error
		resourceCount int
	}{
		{
			name:          "invalid dirPath",
			dirPath:       "not-there",
			kustomize:     KustomizeV3{},
			wantErr:       &os.PathError{Err: syscall.ENOENT, Op: "open", Path: "not-there"},
			resourceCount: 0,
		},
		{
			name:          "simple-deployment",
			dirPath:       filepath.Join(testDataDir, "simple-deployment"),
			kustomize:     KustomizeV3{},
			resourceCount: 4,
		},
		{
			name:          "multibases",
			dirPath:       filepath.Join(multibasesDir, "base"),
			kustomize:     KustomizeV3{},
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       filepath.Join(multibasesDir, "dev"),
			kustomize:     KustomizeV3{},
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       filepath.Join(multibasesDir, "prod"),
			kustomize:     KustomizeV3{},
			resourceCount: 2,
		},

		{
			name:          "multibases",
			dirPath:       filepath.Join(multibasesDir, "stage"),
			kustomize:     KustomizeV3{},
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       multibasesDir,
			kustomize:     KustomizeV3{},
			resourceCount: 4,
		},
		{
			name:          "no-kustomize-directory",
			dirPath:       filepath.Join(testDataDir, "no-kustomizefile"),
			kustomize:     KustomizeV3{},
			wantErr:       fmt.Errorf("kustomization.y(a)ml file not found in the directory %s", filepath.Join(testDataDir, "no-kustomizefile")),
			resourceCount: 0,
		},
		{
			name:          "kustomize-file-empty",
			dirPath:       filepath.Join(testDataDir, "kustomize-file-empty"),
			kustomize:     KustomizeV3{},
			wantErr:       fmt.Errorf("unable to read the kustomization file in the directory %s, error: yaml file is empty", filepath.Join(testDataDir, "kustomize-file-empty")),
			resourceCount: 0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resourceMap, gotErr := tt.kustomize.LoadIacDir(tt.dirPath)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}

			resCount := utils.GetResourceCount(resourceMap)
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
		name        string
		basepath    string
		filename    string
		want        output.AllResourceConfigs
		wantErr     error
		checkPrefix bool
	}{
		{
			name:     "simple-deployment",
			basepath: filepath.Join(testDataDir, "simple-deployment"),
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases",
			basepath: multibasesDir,
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/base",
			basepath: filepath.Join(multibasesDir, "base"),
			filename: kustomizeYml,
			wantErr:  nil,
		},
		{
			name:     "multibases/dev",
			basepath: filepath.Join(multibasesDir, "dev"),
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/prod",
			basepath: filepath.Join(multibasesDir, "prod"),
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/stage",
			basepath: filepath.Join(multibasesDir, "stage"),
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/zero-violation-base",
			basepath: filepath.Join(multibasesDir, "zero-violation-base"),
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:        "erroneous-pod",
			basepath:    filepath.Join(testDataDir, "erroneous-pod"),
			filename:    kustomizeYaml,
			wantErr:     fmt.Errorf(kustomizeErrPrefix),
			checkPrefix: true,
		},
		{
			name:        "erroneous-deployment",
			basepath:    filepath.Join(testDataDir, "erroneous-deployment/"),
			filename:    kustomizeYaml,
			wantErr:     fmt.Errorf(kustomizeErrPrefix),
			checkPrefix: true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := LoadKustomize(tt.basepath, tt.filename)
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
