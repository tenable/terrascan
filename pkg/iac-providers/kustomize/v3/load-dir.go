package kustomizev3

import (
	"fmt"
	"path/filepath"

	k8sv1 "github.com/accurics/terrascan/pkg/iac-providers/kubernetes/v1"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
)

const (
	kustomizedirectory string = "kustomization"
)

var (
	errorKustomizeNotFound     = fmt.Errorf("kustomization.y(a)ml file not found in the directory")
	errorMultipleKustomizeFile = fmt.Errorf("multiple kustomization.y(a)ml found in the directory")
	errorFromKustomize         = fmt.Errorf("error from kustomization")
)

// LoadIacDir loads the kustomize directory and returns the ResourceConfig mapping which is evaluated by the policy engine
func (k *KustomizeV3) LoadIacDir(absRootDir string) (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	files, err := utils.FindFilesBySuffixInDir(absRootDir, KustomizeFileNames())
	if err != nil {
		zap.S().Error("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(files) == 0 {
		err = errorKustomizeNotFound
		zap.S().Error("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(files) > 1 {
		err = errorMultipleKustomizeFile
		zap.S().Error("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	kustomizeFileName := *files[0]
	yamlkustomizeobj, err := utils.ReadYamlFile(filepath.Join(absRootDir, kustomizeFileName))

	if err != nil {
		err = fmt.Errorf("unable to read the kustomization file in the directory : %v", err)
		zap.S().Error("error while reading the file", kustomizeFileName, zap.Error(err))
		return allResourcesConfig, err
	}

	// ResourceConfig representing the kustomization.y(a)ml file
	config := output.ResourceConfig{
		Name:   filepath.Dir(absRootDir),
		Type:   kustomizedirectory,
		Line:   1,
		ID:     kustomizedirectory + "." + filepath.Dir(absRootDir),
		Source: filepath.Join(absRootDir, kustomizeFileName),
		Config: yamlkustomizeobj,
	}

	allResourcesConfig[kustomizedirectory] = append(allResourcesConfig[kustomizedirectory], config)

	// obtaining list of IacDocuments from the target working directory
	iacDocuments, err := LoadKustomize(absRootDir, kustomizeFileName)
	if err != nil {
		zap.S().Error("error occurred while loading kustomize directory", zap.String("kustomize directory", absRootDir), zap.Error(err))
		return nil, err
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

	return allResourcesConfig, nil
}

// LoadKustomize loads up a 'kustomized' directory and returns a returns a list of IacDocuments
func LoadKustomize(basepath, filename string) ([]*utils.IacDocument, error) {
	fSys := filesys.MakeFsOnDisk()
	k := krusty.MakeKustomizer(fSys, krusty.MakeDefaultOptions())

	m, err := k.Run(basepath)
	if err != nil {
		return nil, errorFromKustomize
	}

	yaml, err := m.AsYaml()
	if err != nil {
		return nil, err
	}

	res, err := utils.LoadYAMLString(string(yaml), filename)
	if err != nil {
		return nil, err
	}

	return res, nil
}
