package iacprovider

import (
	"reflect"

	tfv14 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v14"
)

// terraform specific constants
const (
	terraform               supportedIacType    = "terraform"
	terraformV14            supportedIacVersion = "v14"
	terraformDefaultVersion                     = terraformV14
)

// register terraform as an IaC provider with terrascan
func init() {

	// register iac provider
	RegisterIacProvider(terraform, terraformV14, terraformDefaultVersion, reflect.TypeOf(tfv14.TfV14{}))
}
