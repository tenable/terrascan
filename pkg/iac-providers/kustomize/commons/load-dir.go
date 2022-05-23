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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	k8sv1 "github.com/tenable/terrascan/pkg/iac-providers/kubernetes/v1"
	"github.com/tenable/terrascan/pkg/iac-providers/output"
	"github.com/tenable/terrascan/pkg/results"
	"github.com/tenable/terrascan/pkg/utils"
	"go.uber.org/zap"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

const (
	kustomizedirectory string = "kustomization"
)

var (
	kustomizeErrMessage = "error from kustomization. error : %v"
)

// NewKustomizeDirectoryLoader creates a new kustomizeDirectoryLoader
func NewKustomizeDirectoryLoader(absRootDir string, options map[string]interface{}, useKustomizeBinary bool, version string) KustomizeDirectoryLoader {
	kustomizeDirectoryLoader := KustomizeDirectoryLoader{
		absRootDir:         absRootDir,
		options:            options,
		useKustomizeBinary: useKustomizeBinary,
		version:            version,
	}

	return kustomizeDirectoryLoader
}

// LoadIacDir loads the kustomize directory and returns the ResourceConfig mapping which is evaluated by the policy engine
func (t KustomizeDirectoryLoader) LoadIacDir() (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	files, err := utils.FindFilesBySuffixInDir(t.absRootDir, KustomizeFileNames())
	if err != nil {
		return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "kustomize", Directory: t.absRootDir, ErrMessage: err.Error()})
	}

	if len(files) == 0 {
		errMsg := fmt.Sprintf("kustomization.y(a)ml file not found in the directory %s", t.absRootDir)
		return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "kustomize", Directory: t.absRootDir, ErrMessage: errMsg})
	}

	if len(files) > 1 {
		errMsg := fmt.Sprintf("multiple kustomization.y(a)ml found in the directory %s", t.absRootDir)
		return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "kustomize", Directory: t.absRootDir, ErrMessage: errMsg})
	}

	kustomizeFileName := *files[0]
	yamlkustomizeobj, err := utils.ReadYamlFile(filepath.Join(t.absRootDir, kustomizeFileName))

	if err != nil {
		err = fmt.Errorf("unable to read the kustomization file in the directory %s, error: %v", t.absRootDir, err)
		zap.S().Error("error while reading the file", kustomizeFileName, zap.Error(err))
		return allResourcesConfig, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "kustomize", Directory: t.absRootDir, ErrMessage: err.Error()})
	}

	// ResourceConfig representing the kustomization.y(a)ml file
	config := output.ResourceConfig{
		Name:   filepath.Dir(t.absRootDir),
		Type:   kustomizedirectory,
		Line:   1,
		ID:     kustomizedirectory + "." + filepath.Dir(t.absRootDir),
		Source: filepath.Join(t.absRootDir, kustomizeFileName),
		Config: yamlkustomizeobj,
	}

	allResourcesConfig[kustomizedirectory] = append(allResourcesConfig[kustomizedirectory], config)

	// obtaining list of IacDocuments from the target working directory
	iacDocuments, err := LoadKustomize(t.absRootDir, kustomizeFileName, t.version, t.useKustomizeBinary)
	if err != nil {
		errMsg := fmt.Sprintf("error occurred while loading kustomize directory '%s'. err: %v", t.absRootDir, err)
		zap.S().Error("error occurred while loading kustomize directory", zap.String("kustomize directory", t.absRootDir), zap.Error(err))
		return nil, multierror.Append(t.errIacLoadDirs, results.DirScanErr{IacType: "kustomize", Directory: t.absRootDir, ErrMessage: errMsg})
	}

	for _, doc := range iacDocuments {
		var k k8sv1.K8sV1
		var config *output.ResourceConfig

		config, err = k.Normalize(doc)
		if err != nil {
			zap.S().Warn("unable to normalize data", zap.Error(err), zap.String("file", doc.FilePath))
			continue
		}

		// TODO finding a better solution to detect accurate line number for tracing back the files causing violations
		config.Line = 1
		config.Source = doc.FilePath
		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}

	return allResourcesConfig, t.errIacLoadDirs
}

// LoadKustomize loads up a 'kustomized' directory and returns a returns a list of IacDocuments
func LoadKustomize(basepath, filename, version string, useKustomizeBinary bool) ([]*utils.IacDocument, error) {

	var yaml []byte
	if useKustomizeBinary {
		envVar := "KUSTOMIZE_" + version
		kustomizeBinaryPath, exists := os.LookupEnv(envVar)
		if !exists {
			return nil, fmt.Errorf(kustomizeErrMessage, "Environment variable "+envVar+" not set")
		}
		exe := kustomizeBinaryPath
		arg1 := "build"
		arg2 := basepath

		cmd := exec.Command(exe, arg1, arg2)

		var stdOut, stdErr bytes.Buffer
		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr

		err := cmd.Run()

		if err != nil {
			if stdErr.String() != "" {
				return nil, fmt.Errorf(kustomizeErrMessage, stdErr.String())
			}
			return nil, fmt.Errorf(kustomizeErrMessage, err)
		}

		yaml = stdOut.Bytes()
	} else {
		fSys := filesys.MakeFsOnDisk()
		k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())

		m, err := k.Run(fSys, basepath)
		if err != nil {
			return nil, fmt.Errorf(kustomizeErrMessage, err)
		}

		yaml, err = m.AsYaml()
		if err != nil {
			return nil, err
		}
	}

	res, err := utils.LoadYAMLString(string(yaml), filename)
	if err != nil {
		return nil, err
	}

	return res, nil
}
