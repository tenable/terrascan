package tfv12

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	hclConfigs "github.com/hashicorp/terraform/configs"
	"go.uber.org/zap"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// CreateResourceConfig creates output.ResourceConfig
func CreateResourceConfig(managedResource *hclConfigs.Resource) (resourceConfig output.ResourceConfig, err error) {

	// read source file
	fileBytes, err := ioutil.ReadFile(managedResource.DeclRange.Filename)
	if err != nil {
		zap.S().Errorf("failed to read terrafrom IaC file '%s'. error: '%v'", managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to read terraform file")
	}

	// convert resource config from hcl.Body to map[string]interface{}
	c := converter{bytes: fileBytes}
	hclBody := managedResource.Config.(*hclsyntax.Body)
	goOut, err := c.convertBody(hclBody)
	if err != nil {
		zap.S().Errorf("failed to convert hcl.Body to go struct; resource '%s', file: '%s'. error: '%v'",
			managedResource.Name, managedResource.DeclRange.Filename, err)
		return resourceConfig, fmt.Errorf("failed to convert hcl.Body to go struct")
	}

	// create a resource config
	resourceConfig = output.ResourceConfig{
		Name:   managedResource.Name,
		Type:   managedResource.Type,
		Source: managedResource.DeclRange.Filename,
		Config: goOut,
	}

	// successful
	zap.S().Debugf("successfully created resource config for resource '%s', file: '%s'", resourceConfig.Name, resourceConfig.Source)
	return resourceConfig, nil
}
