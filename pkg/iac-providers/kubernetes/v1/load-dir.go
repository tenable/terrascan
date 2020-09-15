package k8sv1

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

func (*K8sV1) getFileType(file string) string {
	if strings.HasSuffix(file, YAMLExtension) {
		return YAMLExtension
	} else if strings.HasSuffix(file, YAMLExtension2) {
		return YAMLExtension2
	} else if strings.HasSuffix(file, JSONExtension) {
		return JSONExtension
	}
	return UnknownExtension
}

// LoadIacDir loads all k8s files in the current directory
func (k *K8sV1) LoadIacDir(absRootDir string) (output.AllResourceConfigs, error) {

	allResourcesConfig := make(map[string][]output.ResourceConfig)

	fileMap, err := utils.FindFilesBySuffix(absRootDir, K8sFileExtensions())
	if err != nil {
		zap.S().Warn("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])

			var configData output.AllResourceConfigs
			if configData, err = k.LoadIacFile(file); err != nil {
				continue
			}

			for key := range configData {
				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, nil
}
