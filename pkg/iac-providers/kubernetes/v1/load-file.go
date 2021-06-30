package k8sv1

import (
	"fmt"
	"path/filepath"

	"github.com/accurics/terrascan/pkg/utils"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

// LoadIacFile loads the k8s file specified
// Note that a single k8s yaml file may contain multiple resource definitions
func (k *K8sV1) LoadIacFile(absFilePath string, options map[string]interface{}) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	var iacDocuments []*utils.IacDocument

	fileExt := k.getFileType(absFilePath)
	switch fileExt {
	case YAMLExtension:
		fallthrough
	case YAMLExtension2:
		iacDocuments, err = utils.LoadYAML(absFilePath)
	case JSONExtension:
		iacDocuments, err = utils.LoadJSON(absFilePath)
	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return allResourcesConfig, fmt.Errorf("unknown file extension for file %s", absFilePath)
	}
	if err != nil {
		zap.S().Debug("failed to load file", zap.String("file", absFilePath))
		return allResourcesConfig, err
	}

	for _, doc := range iacDocuments {
		var config *output.ResourceConfig
		config, err = k.Normalize(doc)
		if err != nil {
			zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absFilePath))
			continue
		}

		config.Line = doc.StartLine
		config.Source = k.getSourceRelativePath(absFilePath)

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}
	return allResourcesConfig, nil
}

// getSourceRelativePath fetches the relative path of file being loaded
func (k *K8sV1) getSourceRelativePath(sourceFile string) string {

	// rootDir should be empty when file scan was initiated by user
	if k.absRootDir == "" {
		return filepath.Base(sourceFile)
	}
	relPath, err := filepath.Rel(k.absRootDir, sourceFile)
	if err != nil {
		zap.S().Debug("error while getting the relative path for", zap.String("IAC file", sourceFile), zap.Error(err))
		return sourceFile
	}
	return relPath
}
