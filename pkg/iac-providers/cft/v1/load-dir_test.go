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

package cftv1

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"go.uber.org/zap"

	"github.com/hashicorp/go-multierror"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

const (
	repoURL  = "https://github.com/tenable/kaimonkey.git"
	branch   = "master"
	basePath = "artifacts"
	provider = "cft"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestLoadIacDir(t *testing.T) {
	testDataDir := "testdata"

	pathErr := &os.PathError{Op: "lstat", Path: "not-there", Err: syscall.ENOENT}
	if utils.IsWindowsPlatform() {
		pathErr = &os.PathError{Op: "CreateFile", Path: "not-there", Err: syscall.ENOENT}
	}

	table := []struct {
		wantErr error
		want    output.AllResourceConfigs
		cftv1   CFTV1
		name    string
		dirPath string
		options map[string]interface{}
	}{
		{
			name:    "empty config",
			dirPath: filepath.Join(testDataDir, "testfile"),
			cftv1:   CFTV1{},
			wantErr: multierror.Append(fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile"))),
		},
		{
			name:    "load config dir with sub directories",
			dirPath: filepath.Join(testDataDir, "templates"),
			cftv1:   CFTV1{},
			want:    map[string][]output.ResourceConfig{},
			wantErr: nil,
		},
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			cftv1:   CFTV1{},
			wantErr: multierror.Append(pathErr),
		},
		{
			name:    "load valid dir",
			dirPath: filepath.Join(testDataDir, "templates", "s3"),
			cftv1:   CFTV1{},
			want:    map[string][]output.ResourceConfig{},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := tt.cftv1.LoadIacDir(tt.dirPath, tt.options)
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
		})
	}
}

func TestCFTMapper(t *testing.T) {
	root := filepath.Join(basePath, provider)
	dirList, err := os.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}
	options := make(map[string]interface{})
	cftv1 := CFTV1{}
	for _, dir := range dirList {
		resourceDir := filepath.Join(root, dir.Name())
		t.Run(resourceDir, func(t *testing.T) {
			allResourceConfigs, gotErr := cftv1.LoadIacDir(resourceDir, options)

			// load expected output.json from test artifacts
			var testArc output.AllResourceConfigs
			outputData, err := os.ReadFile(filepath.Join(resourceDir, "output.json"))
			if err != nil {
				t.Errorf("error reading output.json ResourceConfig, %T", err)
			}

			err = json.Unmarshal(outputData, &testArc)
			if err != nil {
				t.Errorf("error loading output.json ResourceConfig, %T", err)
			}

			// check if resourcetype and resources are present
			for name, resources := range testArc {
				if allResourceConfigs[name] == nil {
					t.Errorf("resource Type %s from test data %s not found", name, resourceDir)
				}
				for _, testResource := range resources {
					var found bool
					for _, resource := range allResourceConfigs[name] {
						if resource.ID == testResource.ID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("resource %s of type %s from test data not found", testResource.ID, name)
					}
				}
			}

			_, ok := gotErr.(*multierror.Error)
			if !ok {
				t.Errorf("expected multierror.Error, got %T", gotErr)
			}
		})
	}
}

func setup() {
	err := downloadArtifacts()
	if err != nil {
		zap.S().Fatal(err)
	}
}

func shutdown() {
	_ = os.RemoveAll(basePath)
}

func downloadArtifacts() error {
	os.RemoveAll(basePath)

	r, err := git.PlainClone(basePath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	if err != nil {
		return err
	}
	return nil
}
