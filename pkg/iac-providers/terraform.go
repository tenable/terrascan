package iacprovider

import (
	"reflect"

	tfv12 "github.com/accurics/terrascan/pkg/iac-providers/terraform/v12"
)

// terraform specific constants
const (
	terraform               supportedIacType    = "terraform"
	terraformV12            supportedIacVersion = "v12"
	terraformDefaultVersion                     = terraformV12
)

// register terraform as an IaC provider with terrascan
func init() {

	// register iac provider
	RegisterIacProvider(terraform, terraformV12, terraformDefaultVersion, reflect.TypeOf(tfv12.TfV12{}))
}
