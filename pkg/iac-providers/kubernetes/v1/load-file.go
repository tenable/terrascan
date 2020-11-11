package k8sv1

import (
	"github.com/accurics/terrascan/pkg/utils"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"go.uber.org/zap"
)

// LoadIacFile loads the k8s file specified
// Note that a single k8s yaml file may contain multiple resource definitions
func (k *K8sV1) LoadIacFile(absRootPath string) (allResourcesConfig output.AllResourceConfigs, err error) {
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	var iacDocuments []*utils.IacDocument

	fileExt := k.getFileType(absRootPath)
	switch fileExt {
	case YAMLExtension:
		fallthrough
	case YAMLExtension2:
		iacDocuments, err = utils.LoadYAML(absRootPath)
	case JSONExtension:
		iacDocuments, err = utils.LoadJSON(absRootPath)
	default:
		zap.S().Error("unknown extension found", zap.String("extension", fileExt))
		return allResourcesConfig, err
	}
	if err != nil {
		zap.S().Info("failed to load file", zap.String("file", absRootPath))
		return allResourcesConfig, err
	}

	for _, doc := range iacDocuments {
		var config *output.ResourceConfig
		config, err = k.Normalize(doc)
		if err != nil {
			zap.S().Debug("unable to normalize data", zap.Error(err), zap.String("file", absRootPath))
			continue
		}

		config.Line = doc.StartLine
		config.Source = absRootPath

		allResourcesConfig[config.Type] = append(allResourcesConfig[config.Type], *config)
	}
	return allResourcesConfig, nil
}
