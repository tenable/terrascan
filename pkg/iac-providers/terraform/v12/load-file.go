package tfv12

import (
	"fmt"

	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

var (
	errFailedLoadConfigFile = fmt.Errorf("failed to load config file")
)

// LoadIacFile parses the given terraform file from the given file path
func (*TfV12) LoadIacFile(absFilePath string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// new terraform config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	hclFile, diags := parser.LoadConfigFile(absFilePath)
	if diags != nil {
		zap.S().Errorf("failed to load config file '%s'. error:\n%v\n", diags)
		return allResourcesConfig, errFailedLoadConfigFile
	}
	if hclFile == nil && diags.HasErrors() {
		zap.S().Errorf("error occured while loading config file. error:\n%v\n", diags)
		return allResourcesConfig, errFailedLoadConfigFile
	}

	// initialize normalized output
	allResourcesConfig = make(map[string][]output.ResourceConfig)

	// traverse through all current's resources
	for _, managedResource := range hclFile.ManagedResources {

		// create output.ResourceConfig from hclConfigs.Resource
		resourceConfig, err := CreateResourceConfig(managedResource)
		if err != nil {
			return allResourcesConfig, fmt.Errorf("failed to create ResourceConfig")
		}

		// append resource config to list of all resources
		// allResourcesConfig = append(allResourcesConfig, resourceConfig)

		// append to normalized output
		if _, present := allResourcesConfig[resourceConfig.Type]; !present {
			allResourcesConfig[resourceConfig.Type] = []output.ResourceConfig{resourceConfig}
		} else {
			allResourcesConfig[resourceConfig.Type] = append(allResourcesConfig[resourceConfig.Type], resourceConfig)
		}
	}

	// successful
	return allResourcesConfig, nil
}
