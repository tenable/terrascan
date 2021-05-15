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

package mapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	repoURL  = "https://github.com/accurics/KaiMonkey.git"
	branch   = "master"
	basePath = "test_data"
	provider = "arm"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestARMMapper(t *testing.T) {
	root := filepath.Join(basePath, provider)
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

	m := NewMapper(provider)
	for i := 1; i < len(dirList); i++ {
		dir := dirList[i]
		fileInfo, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Error(err)
			continue
		}

		t.Run(dir, func(t *testing.T) {
			doc, err := iacDocumentFromFile(filepath.Join(dir, fileInfo[0].Name()))
			if err != nil {
				t.Error(err)
			}

			params, err := parametersFromFile(filepath.Join(dir, fileInfo[1].Name()))
			if err != nil {
				t.Error(err)
			}
			_, err = m.Map(doc, params)
			if err != nil {
				t.Error(err)
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
	os.RemoveAll(basePath)
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

func iacDocumentFromFile(name string) (*utils.IacDocument, error) {
	data, err := readFile(name)
	if err != nil {
		return nil, err
	}

	return &utils.IacDocument{
		Type:      utils.JSONDoc,
		StartLine: 0,
		EndLine:   183,
		FilePath:  filepath.Join("test_data", name),
		Data:      data,
	}, nil
}

func readFile(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func parametersFromFile(name string) (map[string]interface{}, error) {
	data, err := readFile(name)
	if err != nil {
		return nil, err
	}

	var params map[string]interface{}
	err = json.Unmarshal(data, &params)
	if err != nil {
		return nil, err
	}

	res, err := extractParameterValues(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func extractParameterValues(params map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(params["parameters"])
	if err != nil {
		return nil, err
	}
	var npm map[string]struct {
		Value interface{} `json:"value"`
	}
	err = json.Unmarshal(data, &npm)
	if err != nil {
		return nil, err
	}

	finalParams := map[string]interface{}{}
	for key, value := range npm {
		finalParams[key] = value.Value
	}
	return finalParams, nil
}
