package kustomizev3

import (
	"errors"
	"path/filepath"

	k8sv1 "github.com/accurics/terrascan/pkg/iac-providers/kubernetes/v1"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
)

const (
	kustomizedirectory string = "kustomize_directory"
)

// LoadIacDir loads the kustomize directory
func (k *KustomizeV3) LoadIacDir(absRootDir string) (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	files, err := utils.FindFilesBySuffixInCurrentDir(absRootDir, KustomizeFileNames())
	if err != nil {
		zap.S().Warn("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(files) == 0 {
		err := errors.New("could not find a kustomization.yaml/yml file in the directory")
		zap.S().Warn("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	if len(files) > 1 {
		err := errors.New("a directory cannot have more than 1 kustomization.yaml/yml file")
		zap.S().Warn("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	var config output.ResourceConfig
	config.Type = kustomizedirectory
	config.Name = filepath.Dir(absRootDir)
	config.Line = 0
	config.ID = config.Type + "." + config.Name

	var yamlkustomizeobj map[string]interface{}
	var kustomizeFileName string
	for _, filename := range KustomizeFileNames() {
		yamlkustomizeobj, err = utils.ReadYamlFile(filepath.Join(absRootDir, filename))
		if err == nil {
			kustomizeFileName = filename
			break
		}
	}

	if len(yamlkustomizeobj) == 0 {
		err := errors.New("unable to read any kustomization file in the directory")
		zap.S().Warn("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	config.Source = filepath.Join(absRootDir, kustomizeFileName)
	config.Config = yamlkustomizeobj

	allResourcesConfig[kustomizedirectory] = append(allResourcesConfig[kustomizedirectory], config)

	iacDocumentMap := make(map[string][]*utils.IacDocument)
	var iacDocuments []*utils.IacDocument

	iacDocuments, err = loadKustomize(absRootDir, kustomizeFileName)
	if err != nil {
		zap.S().Warn("error occurred while loading kustomize directory", zap.String("kustomize directory", absRootDir), zap.Error(err))
		return nil, err
	}

	iacDocumentMap[absRootDir] = iacDocuments

	for _, iacDocuments := range iacDocumentMap {
		for _, doc := range iacDocuments {
			// @TODO add k8s version check
			var k k8sv1.K8sV1
			var config *output.ResourceConfig

			config, err = k.Normalize(doc)
			if err != nil {
				zap.S().Warn("unable to normalize data", zap.Error(err), zap.String("file", doc.FilePath))
				continue
			}

			config.Line = 1
			config.Source = doc.FilePath

			allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
		}
	}

	return allResourcesConfig, nil
}

func loadKustomize(basepath, filename string) ([]*utils.IacDocument, error) {
	fSys := filesys.MakeFsOnDisk()
	k := krusty.MakeKustomizer(fSys, krusty.MakeDefaultOptions())

	m, err := k.Run(basepath)
	if err != nil {
		return nil, err
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
