package k8sv1

import (
	"fmt"

	"github.com/accurics/terrascan/pkg/utils"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

// LoadIacFile loads the k8s file specified
// Note that a single k8s yaml file may contain multiple resource definitions
func (k *K8sV1) LoadIacFile(filePath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	var iacDocuments []*utils.IacDocument

	fileExt := k.getFileType(filePath)
	switch fileExt {
	case YAMLExtension:
		fallthrough
	case YAMLExtension2:
		iacDocuments, err = utils.LoadYAML(filePath)
	case JSONExtension:
		iacDocuments, err = utils.LoadJSON(filePath)
	default:
		zap.S().Debug("unknown extension found", zap.String("extension", fileExt))
		return allResourcesConfig, fmt.Errorf("unknown file extension for file %s", filePath)
	}
	if err != nil {
		zap.S().Debug("failed to load file", zap.String("file", filePath))
		return allResourcesConfig, err
	}

	for _, doc := range iacDocuments {
		var config *output.ResourceConfig
		config, err = k.Normalize(doc)
		if err != nil {
			zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", filePath))
			continue
		}

		config.Line = doc.StartLine
		config.Source = filePath

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}
	return allResourcesConfig, nil
}
