package tfv12

import (
	"fmt"
	"log"

	hclConfigs "github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
)

// LoadIacFile parses the given terraform file from the given file path
func (*TfV12) LoadIacFile(filePath string) (allResourcesConfig output.AllResourceConfigs, err error) {

	// get absolute path
	absFilePath, err := utils.GetAbsPath(filePath)
	if err != nil {
		return allResourcesConfig, err
	}

	// new terraform config parser
	parser := hclConfigs.NewParser(afero.NewOsFs())

	hclFile, diags := parser.LoadConfigFile(absFilePath)
	if diags != nil {
		log.Printf("failed to load config file '%s'. error:\n%v\n", diags)
		return allResourcesConfig, fmt.Errorf("failed to load config file")
	}
	if hclFile == nil && diags.HasErrors() {
		log.Printf("error occured while loading config file. error:\n%v\n", diags)
		return allResourcesConfig, fmt.Errorf("failed to load config file")
	}

	// traverse through all current's resources
	for _, managedResource := range hclFile.ManagedResources {

		// create output.ResourceConfig from hclConfigs.Resource
		resourceConfig, err := CreateResourceConfig(managedResource)
		if err != nil {
			return allResourcesConfig, fmt.Errorf("failed to create ResourceConfig")
		}

		// append resource config to list of all resources
		allResourcesConfig = append(allResourcesConfig, resourceConfig)
	}

	// successful
	return allResourcesConfig, nil
}
