package tfv12

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/configs"
	"github.com/spf13/afero"
)

// ParseFile parses the given terraform file from the given file path
func (*TfV12) LoadIacFile(filePath string) (config interface{}, err error) {

	// new terraform config parser
	parser := configs.NewParser(afero.NewOsFs())

	config, diags := parser.LoadConfigFile(filePath)
	if diags != nil {
		log.Printf("failed to load config file '%s'. error:\n%v\n", diags)
		return config, fmt.Errorf("failed to load config file")
	}
	log.Printf("config:\n%+v\n", config)

	// successful
	return config, nil
}
