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

package armv1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/utils"
)

const (
	repoURL     = "https://github.com/tenable/kaimonkey.git"
	branch      = "master"
	artifacts   = "artifacts"
	provider    = "arm"
	testDataDir = "testdata"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestLoadIacDir(t *testing.T) {
	invalidDirErr := &os.PathError{Err: syscall.ENOENT, Op: "lstat", Path: "not-there"}
	if utils.IsWindowsPlatform() {
		invalidDirErr = &os.PathError{Err: syscall.ENOENT, Op: "CreateFile", Path: "not-there"}
	}

	var linkedResConf output.ResourceConfig
	if templateData, err := os.ReadFile(filepath.Join(testDataDir, "linked", "output.json")); err == nil {
		err := json.Unmarshal(templateData, &linkedResConf)
		if err != nil {
			t.Errorf("output file not found for linked template test, got %T", err)
		}
	}

	table := []struct {
		wantErr error
		want    output.AllResourceConfigs
		armv1   ARMV1
		name    string
		dirPath string
		options map[string]interface{}
	}{
		{
			name:    "empty config",
			dirPath: filepath.Join(testDataDir, "testfile"),
			armv1:   ARMV1{},
			wantErr: multierror.Append(fmt.Errorf("no directories found for path %s", filepath.Join(testDataDir, "testfile"))),
		},
		{
			name:    "invalid dirPath",
			dirPath: "not-there",
			armv1:   ARMV1{},
			wantErr: multierror.Append(invalidDirErr),
		},
		{
			name:    "key-vault",
			dirPath: filepath.Join(testDataDir, "key-vault"),
			armv1:   ARMV1{},
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			aRC, gotErr := tt.armv1.LoadIacDir(tt.dirPath, tt.options)
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
			if tt.want != nil {
				assert.Equal(t, tt.want, aRC)
			}
		})
	}
}

func TestARMMapper(t *testing.T) {
	root := filepath.Join(artifacts, provider)
	dirList := make([]string, 0)
	err := filepath.Walk(root, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo != nil && fileInfo.IsDir() {
			dirList = append(dirList, filePath)
		}
		return err
	})

	if err != nil {
		t.Error(err)
	}

	options := make(map[string]interface{})

	armv1 := ARMV1{}

	// get output json to verify
	var testArc output.AllResourceConfigs
	outputData, err := os.ReadFile(filepath.Join(root, "output.json"))
	if err != nil {
		t.Errorf("error reading output.json ResourceConfig, %T", err)
	}

	err = json.Unmarshal(outputData, &testArc)
	if err != nil {
		t.Errorf("error loading output.json ResourceConfig, %T", err)
	}

	t.Run(root, func(t *testing.T) {

		allResourceConfigs, gotErr := armv1.LoadIacDir(root, options)
		_, ok := gotErr.(*multierror.Error)
		if !ok {
			t.Errorf("expected multierror.Error, got %T", gotErr)
		}

		// check if resource count is as expected
		for resType := range testArc {
			if allResourceConfigs[resType] == nil {
				t.Errorf("resource Type %s from test data", resType)
			}
			assert.Equal(t, len(allResourceConfigs[resType]), len(testArc[resType]))
		}

	})

}

func setup() {
	err := downloadArtifacts()
	if err != nil {
		zap.S().Fatal(err)
	}
}

func shutdown() {
	os.RemoveAll(artifacts)
}

func downloadArtifacts() error {
	os.RemoveAll(artifacts)

	r, err := git.PlainClone(artifacts, false, &git.CloneOptions{
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
