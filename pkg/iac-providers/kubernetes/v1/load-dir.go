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
		zap.S().Debug("error while searching for iac files", zap.String("root dir", absRootDir), zap.Error(err))
		return allResourcesConfig, err
	}

	for fileDir, files := range fileMap {
		for i := range files {
			file := filepath.Join(fileDir, *files[i])

			var configData output.AllResourceConfigs
			if configData, err = k.LoadIacFile(file); err != nil {
				zap.S().Debug("error while loading iac files", zap.String("IAC file", file), zap.Error(err))
				continue
			}

			for key := range configData {
				// the source path formed for each resources is absolute, which should be relative
				resourceConfigs := configData[key]
				makeSourcePathRelative(absRootDir, resourceConfigs)

				allResourcesConfig[key] = append(allResourcesConfig[key], configData[key]...)
			}
		}
	}

	return allResourcesConfig, nil
}

// makeSourcePathRelative modifies the source path of each resource from absolute to relative path
func makeSourcePathRelative(absRootDir string, resourceConfigs []output.ResourceConfig) {
	for i := range resourceConfigs {
		r := &resourceConfigs[i]
		var err error

		oldSource := r.Source

		// update the source path
		r.Source, err = filepath.Rel(absRootDir, r.Source)

		// though this error should never occur, but, if occurs for some reason, assign the old value of source back
		if err != nil {
			r.Source = oldSource
			zap.S().Debug("error while getting the relative path for", zap.String("IAC file", oldSource), zap.Error(err))
			continue
		}
	}
}
